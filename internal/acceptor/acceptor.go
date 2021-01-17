package acceptor

import (
	"bufio"
	"combotest/internal/models"
	"combotest/internal/pool"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Host     string
	Port     string
	ConnType string
	Timeout  time.Duration
}

type Acceptor struct {
	ep          pool.EventPool
	cfg         Config
	stop        chan struct{}
	wg          *sync.WaitGroup
	connections []net.Conn
	log         *logrus.Logger
}

type scanResponse struct {
	Event models.Event
	Error error
}

func NewAcceptor(cfg Config, ep pool.EventPool, log *logrus.Logger) *Acceptor {
	return &Acceptor{
		cfg:         cfg,
		ep:          ep,
		stop:        make(chan struct{}),
		wg:          &sync.WaitGroup{},
		connections: make([]net.Conn, 2),
		log:         log,
	}
}

func (a *Acceptor) Run() error {
	address := fmt.Sprintf("%s:%s", a.cfg.Host, a.cfg.Port)

	l, err := net.Listen(a.cfg.ConnType, address)
	if err != nil {
		return fmt.Errorf("error listening: %w", err)
	}
	defer func() {
		if err = l.Close(); err != nil {
			a.log.WithError(err).Error("close listen")
		}
	}()

	connCh := make(chan net.Conn, 1)
	eCh := make(chan error, 1)

	a.log.Infoln("Listening on", address)

Accept:
	for {
		go func() {
			conn, err := l.Accept()
			if err != nil {
				eCh <- fmt.Errorf("error accepting: %w", err)
				return
			}
			connCh <- conn
		}()

		select {
		case <-a.stop: // сигнал остановки принятия новых коннектов
			a.log.Infoln("stop accepting")
			break Accept
		case conn := <-connCh: // новое соединение
			a.log.Infoln("accept: ", conn.RemoteAddr())
			a.connections = append(a.connections, conn)

			a.wg.Add(1)
			go a.handleRequest(conn, a.cfg.Timeout) // обрабатываем соединение
		case err := <-eCh: // ошибка принятия соединения
			a.log.WithError(err).Error("accept error")
			break Accept
		}
	}
	return nil
}

func (a *Acceptor) Stop() error {
	a.stop <- struct{}{} // stop accepting new connections

	var err error
	for i := range a.connections { // close existing connections
		c := a.connections[i]
		if c != nil { // открытые конекты
			err = c.Close()
		}
	}

	a.wg.Wait() // ждем что все наши обработчики закрыты

	return err
}

func (a *Acceptor) scan(ch chan<- scanResponse, wg *sync.WaitGroup, sc *bufio.Scanner) {
	defer wg.Done()

	var res scanResponse

	if !sc.Scan() {
		if err := sc.Err(); err != nil {
			res.Error = err
		} else {
			res.Error = io.EOF
		}
		ch <- res
		return
	}

	if err := json.Unmarshal(sc.Bytes(), &res.Event); err != nil {
		res.Error = err
		ch <- res
		return
	}

	ch <- res
	return
}

// Handles incoming requests.
func (a *Acceptor) handleRequest(conn net.Conn, timeout time.Duration) {
	defer a.wg.Done()

	r := bufio.NewReader(conn)
	scanner := bufio.NewScanner(r)

	scanCh := make(chan scanResponse, 1)
	var wg sync.WaitGroup

Loop:
	for {
		wg.Add(1)
		go a.scan(scanCh, &wg, scanner) // read income event
		timer := time.NewTimer(timeout) // timeout

		select {
		case <-timer.C: // чтобы не висело соединение долго
			a.log.Infoln("time out read: ", conn.RemoteAddr())
			break Loop
		case income := <-scanCh:
			timer.Stop()             // чтоб ресурсы не текли
			if income.Error != nil { // что-то не так
				if !errors.Is(income.Error, io.EOF) {
					a.log.WithError(income.Error).WithField("addr", conn.RemoteAddr()).Error("handle request")
				}
				break Loop
			} else { // пришло событие
				a.log.Debugln(income.Event)
				a.ep.Push(income.Event) // пушим событие в очередь
				continue Loop
			}
		}
	}

	if err := conn.Close(); err != nil {
		a.log.WithError(err).Error("close connection")
	}
	wg.Wait()
}
