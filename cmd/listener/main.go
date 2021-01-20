package main

import (
	"combotest/internal/app/acceptor"
	"combotest/internal/app/auth"
	"combotest/internal/app/loader"
	"combotest/internal/app/pool"
	delivery "combotest/internal/delivery/http"
	"combotest/internal/repository"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

const (
	CONN_ACCEPTOR_HOST = "localhost"
	CONN_ACCEPTOR_PORT = "4000"
	CONN_ACCEPTOR_TYPE = "tcp"

	POOL_SIZE   = 1024
	WORKERS_NUM = 3

	MONGODB_URI = "mongodb://localhost:27017"
	DB_NAME     = "test"

	CHUNCK_SIZE uint32 = 2
	LOADERS_NUM uint32 = 2

	HASH_SALT = "salt"

	API_ADDRESS = ":8080"
)

var (
	EXPIRE_TIME        = 1 * time.Hour
	WAIT_EVENT_TIMEOUT = 5 * time.Second
	SECRET_KEY         *ecdsa.PrivateKey
)

var log = &logrus.Logger{
	Out:       os.Stderr,
	Formatter: new(logrus.TextFormatter),
	Hooks:     make(logrus.LevelHooks),
	Level:     logrus.DebugLevel,
}

func main() {
	SECRET_KEY, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	// TODO: parse congig from file

	toEncryptEventsPool := pool.NewChanEventPool(POOL_SIZE) // пул для событий требующих шифровку данных
	encryptedEventsPool := pool.NewChanEventPool(POOL_SIZE) // пул событий с зашифрованными данными
	// пул воркеров, которые берут из пула события, шифруют данные и отправляют дальше в другой пул
	workersPool := pool.NewEncryptWorkersPool(log)
	workersPool.Init(toEncryptEventsPool, encryptedEventsPool, WORKERS_NUM)

	acfg := acceptor.Config{
		ConnType: CONN_ACCEPTOR_TYPE,
		Host:     CONN_ACCEPTOR_HOST,
		Port:     CONN_ACCEPTOR_PORT,
		Timeout:  WAIT_EVENT_TIMEOUT, // так как открытие tcp соединения требует накладных расходов,
		// то обработчик отсавляет tcp соединение открытым. Если агент по истечении времени = WAIT_EVENT_TIMEOUT не прислал ничего
		// соединение разрывается.
	}
	acceptr := acceptor.NewAcceptor(acfg, toEncryptEventsPool, log) // обработчик входящих событий

	// Data base
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MONGODB_URI))
	if err != nil {
		log.WithError(err).Fatal("connect to db")
	}
	defer func() {
		if e := client.Disconnect(ctx); e != nil {
			log.WithError(err).Error("disconect")
		}
	}()

	db := client.Database(DB_NAME)

	// Repository
	er := repository.NewEventsRepository(db)
	ur := repository.NewUsersRepository(db)

	// Usecases
	authz := auth.NewAuthorizer(ur, HASH_SALT, SECRET_KEY, EXPIRE_TIME)

	// Delivery
	router := mux.NewRouter()

	ah := delivery.AuthHandler{
		Auth: authz,
		Log:  log,
	}
	ah.SetAuthHandler(router)

	server := &http.Server{
		Addr:    API_ADDRESS,
		Handler: router,
	}

	// eh := delivery.EventsHandler{}
	// eh.SetEventsHandler(router)

	// uh := delivery.UsersHandler{}
	// uh.SetUsersHandler(router)

	// загрузчик сохраняет события с зашифрованными данными
	l := loader.NewLoader(encryptedEventsPool, er, log)
	l.Run(LOADERS_NUM, CHUNCK_SIZE) // количесвто загрузчиков, и размер чанка для загрузки

	group, _ := errgroup.WithContext(context.Background())

	group.Go(acceptr.Run)

	group.Go(func() error {
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}
		return server.Shutdown(ctx)
	})

	group.Go(func() error {
		<-ctx.Done()
		return server.Shutdown(ctx)
	})

	group.Go(func() error {
		stop := make(chan os.Signal, 1)

		signal.Notify(stop,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT)
		signal := <-stop

		log.Infoln("signal:", signal.String())

		log.Debug("wait acceptor stop")
		err := acceptr.Stop()

		toEncryptEventsPool.Close()

		log.Debug("wait workers pool")
		workersPool.Stop()

		encryptedEventsPool.Close()

		log.Debug("wait loader stop")
		l.Stop()

		log.Debug("all stoped")

		cancel()
		return err
	})

	if err := group.Wait(); err != nil {
		log.Errorln("Error:", err)
	}
}
