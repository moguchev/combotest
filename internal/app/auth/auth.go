package auth

import (
	"combotest/internal/interfaces/users"
	"combotest/internal/models"
	"context"
	"crypto/ecdsa"
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

type Claims struct {
	jwt.StandardClaims
	Username string      `json:"username"`
	Role     models.Role `json:"role"`
}

type Authorizer struct {
	repo users.Repository

	hashSalt       string
	signingKey     *ecdsa.PrivateKey
	expireDuration time.Duration
}

func NewAuthorizer(repo users.Repository, hashSalt string, signingKey *ecdsa.PrivateKey, expireDuration time.Duration) *Authorizer {
	return &Authorizer{
		repo:           repo,
		hashSalt:       hashSalt,
		signingKey:     signingKey,
		expireDuration: expireDuration,
	}
}

func (a *Authorizer) SignIn(ctx context.Context, auth models.AuthInfo) (string, time.Time, error) {
	auth.Password = a.Hash(auth.Password)

	user, err := a.repo.GetUserByAuthInfo(ctx, auth)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("user not found: %w", err)
	}
	if !user.Confirmed {
		return "", time.Time{}, fmt.Errorf("user not confirmed")
	}

	return a.CreateJWTToken(user)
}

func (a *Authorizer) CreateJWTToken(user models.User) (string, time.Time, error) {
	isssuedAt := time.Now()
	expiresAt := isssuedAt.Add(a.expireDuration)
	token := jwt.NewWithClaims(jwt.SigningMethodES256, &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(expiresAt),
			IssuedAt:  jwt.At(isssuedAt),
		},
		Username: user.ID,
		Role:     user.Role,
	})
	tokenStr, err := token.SignedString(a.signingKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenStr, expiresAt, nil
}

func (a *Authorizer) GetClaimsFromToken(tstr string) (Claims, error) {
	claims := Claims{}

	tkn, err := jwt.ParseWithClaims(tstr, &claims, func(token *jwt.Token) (interface{}, error) {
		return a.signingKey, nil
	})
	if err != nil {
		return Claims{}, err
	}
	if !tkn.Valid {
		return Claims{}, fmt.Errorf("invalid token")
	}
	return claims, nil
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
