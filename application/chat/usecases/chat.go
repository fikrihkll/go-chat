package usecases

import (
	"context"

	"github.com/fikrihkll/chat-app/application/chat"
	"github.com/fikrihkll/chat-app/common"
)

type ChatApplication struct {
	chatRepository chat.IChatRepository
	userRepository chat.IUserRepository
}

func NewChatApplication(chatRepository chat.IChatRepository, userRepository chat.IUserRepository) chat.IChatUseCase {
	return &ChatApplication{chatRepository, userRepository}
}

func (uc *ChatApplication) SaveMessageByEmail(ctx context.Context, newMessage chat.NewMessageByEmailParam) (err error) {
	targetUser, err := uc.userRepository.GetUserByEmail(ctx, newMessage.MemberEmail)
	if err != nil {
		err = common.UserNotFoundError
		return
	}

	err = uc.chatRepository.InsertMessageByEmail(ctx, newMessage, targetUser)
	return
}

func (uc *ChatApplication) SaveMessageByRoomID(ctx context.Context, newMessage chat.NewMessageByRoomIDParam) (err error) {
	err = uc.chatRepository.InsertMessageByRoomID(ctx, newMessage)
	return
}

func (uc *ChatApplication) GetMessages(ctx context.Context, params chat.MessageHistoryParams) (messages []chat.Message, err error) {
	messages, err = uc.chatRepository.GetMessage(ctx, params)
	return
}

func (uc *ChatApplication) GetRoomsByID(ctx context.Context, currentUserEmail string) (rooms []chat.Room, err error) {
	rooms, err = uc.chatRepository.GetRoomsByID(ctx, currentUserEmail)
	return
}