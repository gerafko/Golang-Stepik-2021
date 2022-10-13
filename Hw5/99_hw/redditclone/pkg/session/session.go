package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"redditclone/pkg/user"

	jwt "github.com/dgrijalva/jwt-go"
)

var Secret = "secret"

type SessionJWTClaims struct {
	User struct {
		UserName string `json:"username"`
		UserId   string `json:"id"`
	} `json:"user"`
	jwt.StandardClaims
}

type SessionJWT struct {
	Secret []byte `json:"token"`
}

type SessionJWTString struct {
	Secret string `json:"token"`
}

type Session struct {
	UserName string `json:"username"`
	UserId   string `json:"id"`
}

var (
	ErrNoAuth = errors.New("No user found")
)

func NewSessionsJWT(secret string) *SessionJWT {
	return &SessionJWT{
		Secret: []byte(secret),
	}
}

func (sm *SessionJWT) Create(u *user.User) ([]byte, error) {
	var data SessionJWTClaims = SessionJWTClaims{
		User: struct {
			UserName string "json:\"username\""
			UserId   string "json:\"id\""
		}{
			UserId:   u.ID,
			UserName: u.Login,
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(90 * 24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		}}

	sessVal, err := jwt.NewWithClaims(jwt.SigningMethodHS256, data).SignedString([]byte(Secret))
	if err != nil {
		return nil, fmt.Errorf("session.Create.Create.jwt.NewWithClaims", err)
	}
	var SesJWTString SessionJWTString = SessionJWTString{sessVal}
	resp, err := json.Marshal(SesJWTString)
	if err != nil {
		return nil, fmt.Errorf("session.Create.Create.json.Marshal", err)
	}
	return resp, nil
}

func (sm *SessionJWT) Check(r *http.Request) (*Session, error) {
	Value := r.Header.Get("authorization")
	Value = strings.Fields(Value)[1]
	if Value == "" {
		return nil, fmt.Errorf("cant find jwt token Check")
	}
	payload := SessionJWTClaims{}
	_, err := jwt.ParseWithClaims(Value, payload, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secret), nil
	})
	if err != nil || payload.User.UserId == "" || payload.User.UserName == "" {
		return nil, fmt.Errorf("VoteCheck - cant parse jwt token: %v", err)
	}
	// проверка exp, iat
	if payload.Valid() != nil {
		return nil, fmt.Errorf("invalid jwt token: %v", err)
	}

	return &Session{
		UserId:   payload.User.UserId,
		UserName: payload.User.UserName,
	}, nil
}
