package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ivanpaghubasan/hoa-hub-api/internal/constants"
	"github.com/ivanpaghubasan/hoa-hub-api/internal/model"
)

type IJWTAuth interface {
	GenerateToken(user *model.User) (TokenPairs, error)
	ParseAccessToken(tokenStr string) (*Claims, error)
	GetRefreshCookie(refreshToken string) *http.Cookie
	GetExpiredRefreshCookie() *http.Cookie
}

type JWTAuth struct {
	Issuer        string
	Audience      string
	Secret        string
	TokenExpiry   time.Duration
	RefreshExpiry time.Duration
	CookieDomain  string
	CookiePath    string
	CookieName    string
}

type TokenPairs struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

func NewJWTAuth(secret, issuer, audience, cookieDomain string) IJWTAuth {
	return &JWTAuth{
		Issuer:        issuer,
		Audience:      audience,
		Secret:        secret,
		TokenExpiry:   15 * time.Minute,
		RefreshExpiry: 7 * 24 * time.Hour,
		CookieDomain:  cookieDomain,
		CookiePath:    "/",
		CookieName:    "refresh_token",
	}
}

func (j *JWTAuth) GenerateToken(user *model.User) (TokenPairs, error) {
	accessClaims := &Claims{
		UserID: user.ID.String(),
		Name:   fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.Issuer,
			Audience:  []string{j.Audience},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.TokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.String(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	signedAccessToken, err := accessToken.SignedString([]byte(j.Secret))
	if err != nil {
		return TokenPairs{}, nil
	}

	refreshClaims := jwt.MapClaims{
		"sub": user.ID.String(),
		"exp": time.Now().Add(j.RefreshExpiry).Unix(),
		"iat": time.Now().Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefreshToken, err := refreshToken.SignedString([]byte(j.Secret))
	if err != nil {
		return TokenPairs{}, nil
	}

	return TokenPairs{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
	}, nil
}

func (j *JWTAuth) ParseAccessToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})
	if err != nil || !token.Valid {
		return nil, constants.ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("could not parse claims")
	}

	if claims.Issuer != j.Issuer {
		return nil, errors.New("invalid issuer")
	}

	return claims, nil
}

func (j *JWTAuth) GetRefreshCookie(refreshToken string) *http.Cookie {
	return &http.Cookie{
		Name:     j.CookieName,
		Value:    refreshToken,
		Expires:  time.Now().Add(j.RefreshExpiry),
		MaxAge:   int(j.RefreshExpiry.Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     j.CookiePath,
		Domain:   j.CookieDomain,
	}
}

func (j *JWTAuth) GetExpiredRefreshCookie() *http.Cookie {
	return &http.Cookie{
		Name:     j.CookieName,
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     j.CookiePath,
		Domain:   j.CookieDomain,
	}
}
