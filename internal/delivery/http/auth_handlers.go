package delivery

import (
	"combotest/internal/app/auth"
	"combotest/internal/models"
	"combotest/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	CookieTokeName = "token"
)

// AuthHandler
type AuthHandler struct {
	Auth *auth.Authorizer
	Log  *logrus.Logger
}

func (h *AuthHandler) SetAuthHandler(router *mux.Router) {
	router.HandleFunc("/user", h.SignUp).Methods(http.MethodPost)
	router.HandleFunc("/auth", h.SignIn).Methods(http.MethodPost)
	// тест
	router.HandleFunc("/", h.I).Methods(http.MethodGet)
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {

}

func (h *AuthHandler) I(w http.ResponseWriter, r *http.Request) {
	log := h.Log.WithField("handler", "I")

	c, err := r.Cookie(CookieTokeName)
	if err != nil {
		log.WithError(err).Error("get token from cookie")
		if err == http.ErrNoCookie {
			utils.RespondWithError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("bad request"))
		return
	}

	tknStr := c.Value

	claims, err := h.Auth.GetClaimsFromToken(tknStr)
	if err != nil {
		log.WithError(err).Error("get claims")
		utils.RespondWithError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
	}

	w.Write([]byte(fmt.Sprintf("Welcome %s-%s!", claims.Username, claims.Role)))
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	log := h.Log.WithField("handler", "SignIn")

	var info models.AuthInfo
	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		log.WithError(err).Error("decode body: %w", err)
		utils.RespondWithError(w, http.StatusBadRequest, err)
		return
	}

	token, exp, err := h.Auth.SignIn(r.Context(), info)
	if err != nil {
		log.WithError(err).Errorf("sign in: %w", err)
		utils.RespondWithError(w, http.StatusUnauthorized, err)
	}

	fmt.Println("token:", token)

	http.SetCookie(w, &http.Cookie{
		Name:    CookieTokeName,
		Value:   token,
		Expires: exp,
		// HttpOnly:
	})

	w.WriteHeader(http.StatusOK)
}
