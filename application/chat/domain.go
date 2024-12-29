package chat

import (
	"time"

	"github.com/google/uuid"
)

// Model
type Message struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	RoomID    uuid.UUID `json:"room_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Room struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Users     []string  `json:"users"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Param
type LoginParam struct {
	Email    string
	Password string
}

type NewMessageByEmailParam struct {
	CurrentUserID    string
	CurrentUserEmail string
	MemberEmail      string
	Message          string
}

type NewMessageByRoomIDParam struct {
	CurrentUserID string
	RoomID        string
	Message       string
}

type MessageHistoryParams struct {
	TimeAfter        int64
	TargetEmail      string
	CurrentUserEmail string
}

// Response
type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	User         User   `json:"user"`
}
