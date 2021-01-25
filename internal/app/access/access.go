package access

import (
	"combotest/internal/app/auth"
	"combotest/internal/models"
	"combotest/pkg/utils"
	"context"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

type ctxKey = struct {
	key string
}

var userCtxKey = ctxKey{"user"}

func GetUserFromCtx(ctx context.Context) (models.User, bool) {
	user, ok := ctx.Value(userCtxKey).(models.User)
	if !ok {
		return models.User{Role: models.Anonymous}, false
	}

	return user, true
}

func SetUserToCtx(ctx context.Context, u models.User) context.Context {
	return context.WithValue(ctx, userCtxKey, u)
}

type AccessManager struct {
	log  *logrus.Logger
	auth *auth.Authorizer
}

func NewAccessManager(log *logrus.Logger, auth *auth.Authorizer) *AccessManager {
	return &AccessManager{
		log:  log,
		auth: auth,
	}
}

func (am *AccessManager) AccessMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(auth.CookieTokeName)
		if err != nil {
			am.log.WithError(err).Error("get token from cookie")
			if err == http.ErrNoCookie {
				utils.RespondWithError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
				return
			}
			utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("bad request"))
			return
		}

		claims, err := am.auth.GetClaimsFromToken(c.Value)
		if err != nil {
			am.log.WithError(err).Error("get claims")
			utils.RespondWithError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}

		user := models.User{
			ID:        claims.Username,
			Role:      claims.Role,
			Confirmed: true, // только подтвержденые сотрудники имеют доступ
		}

		next.ServeHTTP(w, r.WithContext(SetUserToCtx(r.Context(), user)))
	})
}
