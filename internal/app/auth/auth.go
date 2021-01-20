package auth

import (
	"time"
	"combotest/internal/interfaces/users"
	"combotest/internal/models"
	"context"
	"crypto/sha1"
	"fmt"

	"github.com/dgrijalva/jwt-go/v4"
)

type Claims struct {
	jwt.StandardClaims
	Username string `json:"username"`
}

type Authorizer struct {
	repo     users.Repository

	hashSalt       string
	signingKey     []byte
	expireDuration time.Duration
}

func NewAuthorizer(repo auth.Repository, hashSalt string, signingKey []byte, expireDuration time.Duration) *Authorizer {
	return &Authorizer{
		repo:           repo,
		hashSalt:       hashSalt,
		signingKey:     signingKey,
		expireDuration: expireDuration,
	}
}

func (a *Authorizer) SignIn(ctx context.Context, auth models.AuthInfo) (string, error) {
	auth.Password = a.Hash(auth.Password)

	user, err := a.repo.GetUserByAuthInfo(ctx, auth)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}
	if !user.Confirmed {
		return "", fmt.Errorf("not confirmed")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(a.exireDuration)),
			IssuedAt: jwt.At(time.Now()),
		}
		Username: user.ID,
	})

	return token.SignedString(a.secret)
}

func (a *Authorizer) SignUp(ctx context.Context, cu models.CreateUser) error {
	cu.Password = a.Hash(cu.Password)
	_, err := a.repo.CreateUser(ctx, cu)
	return err
}


func (a *Authorizer) Hash(s string) string {
	pwd := sha1.New()
	pwd.Write([]byte(s))
	pwd.Write([]byte(a.hashSalt))

	return fmt.Sprintf("%x", pwd.Sum(nil))
}
