package user

import (
	"errors"
	"math/rand"
	"time"
)

func NewUserRepo() *UserRepo {
	return &UserRepo{
		data: make(map[string]*User, 0),
	}
}

var (
	ErrNoUser    = errors.New("No user found")
	ErrBadPass   = errors.New("Invald password")
	AlreadyExist = errors.New("Invald password")
)

const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func (repo *UserRepo) Authorize(login, pass string) (*User, error) {
	u, ok := repo.data[login]
	if !ok {
		return nil, ErrNoUser
	}

	if u.Password != pass {
		return nil, ErrBadPass
	}
	return u, nil
}

func (repo *UserRepo) Register(login, pass string) (*User, error) {
	if _, ok := repo.data[login]; ok {
		return nil, AlreadyExist
	}
	var u User = User{StringWithCharset(15, charset), login, pass}
	repo.data[login] = &u
	return &u, nil
}

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
