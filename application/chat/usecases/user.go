package usecases

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"time"

	"github.com/fikrihkll/chat-app/application/chat"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserApplication struct {
	userRepository chat.IUserRepository
}

var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrMalformatEmail = errors.New("incorrect email format")
var ErrInvalidPassword = errors.New("password must be 8 characters or more")
var ErrNotRegistered = errors.New("no record found")
var ErrIncorrectEmailOrPass = errors.New("incorrect email or password")

func NewUserApplication(userRepository chat.IUserRepository) chat.IAuthUseCase {
	return &UserApplication{userRepository}
}

func (uc *UserApplication) RegisterUser(ctx context.Context, user chat.User) (err error) {
	_, errUser := uc.userRepository.GetUserByEmail(ctx, user.Email)
	if !errors.Is(errUser, sql.ErrNoRows) {
		return ErrEmailAlreadyExists
	}

	return uc.userRepository.CreateUser(ctx, user)
}

func (uc *UserApplication) Login(ctx context.Context, loginParam chat.LoginParam) (data chat.LoginResponse, err error) {
	user, err := uc.userRepository.GetUserByEmail(ctx, loginParam.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrNotRegistered
			return
		}
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginParam.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			err = ErrIncorrectEmailOrPass
			return 
		}
		return
	}

	token, errToken := generateToken(user.ID.String(), user.Email, 24)
	if errToken != nil {
		err = errToken
		return
	}
	
	data = chat.LoginResponse{
		Token: token,
		User:  user,
	}

	return 
}

func generateToken(userID string, userEmail string, expiresInHour int64) (token string, err error) {
	if expiresInHour == 0 {
		expiresInHour = 24
	}

	tokenJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userID,
		"email": userEmail,
		"exp": time.Now().Add(time.Hour * time.Duration(expiresInHour)).Unix(),
	})
	tokenStr, err := tokenJwt.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		return
	}
	token = tokenStr

	return
}
