package delivery

import (
	"combotest/internal/app/auth"
	"combotest/internal/models"
	"combotest/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// AuthHandler
type AuthHandler struct {
	Auth *auth.Authorizer
	Log  *logrus.Logger
}

func (h *AuthHandler) SetAuthHandler(router *mux.Router) {
	router.HandleFunc("/user", h.SignUp).Methods(http.MethodPost) // регистрация
	router.HandleFunc("/", h.SignIn).Methods(http.MethodPost)     // вход
	router.HandleFunc("/", h.SignOut).Methods(http.MethodDelete)
	router.HandleFunc("/i", h.I).Methods(http.MethodGet) // информация о себе
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	log := h.Log.WithField("handler", "SignUp")
	ctx := r.Context()

	var newUser models.CreateUser
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		log.WithError(err).Error("decode body: %w", err)
		utils.RespondWithError(w, http.StatusBadRequest, err)
		return
	}

	if err = h.Auth.SignUp(ctx, newUser); err != nil { // TODO switch errors: already exists, internal
		log.WithError(err).Error("sign up: %w", err)
		utils.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AuthHandler) I(w http.ResponseWriter, r *http.Request) {
	log := h.Log.WithField("handler", "I")

	c, err := r.Cookie(auth.CookieTokeName)
	if err != nil {
		log.WithError(err).Error("get token from cookie")
		if err == http.ErrNoCookie {
			utils.RespondWithError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("bad request"))
		return
	}

	claims, err := h.Auth.GetClaimsFromToken(c.Value)
	if err != nil {
		log.WithError(err).Error("get claims")
		utils.RespondWithError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}

	user := models.User{
		ID:        claims.Username,
		Role:      claims.Role,
		Confirmed: true, // только подтвержденые сотрудники имеют доступ
	}

	utils.RespondWithJSON(w, http.StatusOK, user)
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
		log.WithError(err).Error("sign in: %w", err)
		utils.RespondWithError(w, http.StatusUnauthorized, err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     auth.CookieTokeName,
		Value:    token,
		Expires:  exp,
		Path:     "/",
		HttpOnly: true,
		//Secure:   true,

	})

	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    auth.CookieTokeName,
		MaxAge:  -1,
		Expires: time.Now().Add(-100 * time.Hour), // Set expires for older versions of IE
		Path:    "/",
	})

	w.WriteHeader(http.StatusOK)
}
