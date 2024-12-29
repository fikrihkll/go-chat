package transport

import (
	"github.com/go-ozzo/ozzo-validation/v4/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type NewMessageByEmail struct {
	MemberEmail string `json:"member_email"`
	Message     string `json:"message"`
}

func (request NewMessageByEmail) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Message, validation.Required),
		validation.Field(&request.MemberEmail, validation.Required, is.Email),
	)
}

type NewMessageByRoomID struct {
	Message     string `json:"message"`
}

func (request NewMessageByRoomID) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Message, validation.Required),
	)
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (request Login) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Email, validation.Required, is.Email),
		validation.Field(&request.Password, validation.Required, validation.Length(8, 30)),
	)
}

type Register struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (request Register) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Name, validation.Required),
		validation.Field(&request.Email, validation.Required, is.Email),
		validation.Field(&request.Password, validation.Required, validation.Length(8, 30)),
	)
}

