package chat

import (
	"context"
)

type IChatUseCase interface {
	SaveMessageByRoomID(ctx context.Context, newMessage NewMessageByRoomIDParam) (err error)
	SaveMessageByEmail(ctx context.Context, newMessage NewMessageByEmailParam) (err error)
	GetMessages(ctx context.Context, params MessageHistoryParams) (messages []Message, err error)
	GetRoomsByID(ctx context.Context, currentUserEmail string) (rooms []Room, err error)
}

type IAuthUseCase interface {
	RegisterUser(ctx context.Context, newUser User) (err error)
	Login(ctx context.Context, loginParam LoginParam) (data LoginResponse, err error)
}

type IUserRepository interface {
	CreateUser(ctx context.Context, user User) (err error)
	GetUserByEmail(ctx context.Context, email string) (user User, err error)
	GetUserByID(ctx context.Context, id string) (user User, err error)
}

type IChatRepository interface {
	InsertMessageByRoomID(ctx context.Context, newMessage NewMessageByRoomIDParam) (err error)
	InsertMessageByEmail(ctx context.Context, newMessage NewMessageByEmailParam, targetUser User) (err error)
	GetMessage(ctx context.Context, params MessageHistoryParams) (messages []Message, err error)
	GetRoomsByID(ctx context.Context, currentUsetEmail string) (rooms []Room, err error)
}
