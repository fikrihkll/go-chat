package repositories

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/fikrihkll/chat-app/application/chat"
	"github.com/fikrihkll/chat-app/common"
	"github.com/lib/pq"
)

type ChatRepositoryPostgree struct {
	db *sql.DB
}

func NewChatRepositoryPostgree(db *sql.DB) chat.IChatRepository {
	return &ChatRepositoryPostgree{db}
}

var ErrRoomNotFound = errors.New("room not found")

func (repo *ChatRepositoryPostgree) InsertMessageByEmail(ctx context.Context, newMessage chat.NewMessageByEmailParam, targetUser chat.User) (err error) {
	tx, err := repo.db.Begin()
	if err != nil {
		common.Log(common.LOG_LEVEL_ERROR, err.Error())
		return
	}

	roomSql := "SELECT * FROM rooms WHERE users @> $1"
	roomRow := tx.QueryRowContext(ctx, roomSql, pq.Array([]string{newMessage.MemberEmail, newMessage.CurrentUserEmail}))
	if roomRow.Err() != nil {
		err = roomRow.Err()
		common.Log(common.LOG_LEVEL_ERROR, err.Error())
		return
	}

	var room chat.Room
	if err = roomRow.Scan(&room.ID, &room.Name, pq.Array(&room.Users), &room.CreatedAt, &room.UpdatedAt); err != nil {
		if err != sql.ErrNoRows {
			common.Log(common.LOG_LEVEL_ERROR, err.Error())
			return
		} 
		
		if err == sql.ErrNoRows {
			insertRoomSql := "INSERT INTO rooms VALUES(uuid_generate_v4(), $1, $2, NOW(), NOW()) RETURNING id"
			
			var lastInsertedID string

			errRoom := tx.QueryRow(insertRoomSql, "room", pq.Array([]string{newMessage.CurrentUserEmail, targetUser.Email})).Scan(&lastInsertedID)
			if errRoom != nil {
				err = errRoom
				tx.Rollback()
				log.Fatalf("Failed to insert room: %v", err)
				common.Log(common.LOG_LEVEL_ERROR, err.Error())
				return
			}

			err = tx.QueryRowContext(ctx, "SELECT * FROM rooms WHERE id = $1", lastInsertedID).Scan(&room.ID, &room.Name, pq.Array(&room.Users), &room.CreatedAt, &room.CreatedAt)
			if err != nil {
				tx.Rollback()
				common.Log(common.LOG_LEVEL_ERROR, err.Error())
				return 
			}
		}
	}

	insertMessageSql := "INSERT INTO messages VALUES(uuid_generate_v4(), $1, $2, $3, NOW(), NOW())"
	stmtMessage, err := tx.Prepare(insertMessageSql)
	if err != nil {
		tx.Rollback()
		common.Log(common.LOG_LEVEL_ERROR, err.Error())
		return
	}
	defer stmtMessage.Close()

	_, errMessage := stmtMessage.Exec(newMessage.CurrentUserID, room.ID, newMessage.Message)
	if errMessage != nil {
		tx.Rollback()
		log.Fatalf("Failed to prepare statement: %v", err)
		err = errMessage
		common.Log(common.LOG_LEVEL_ERROR, err.Error())
		return
	}

	if err = tx.Commit(); err != nil {
		common.Log(common.LOG_LEVEL_ERROR, err.Error())
		return
	}

	return
}

func (repo *ChatRepositoryPostgree) InsertMessageByRoomID(ctx context.Context, newMessage chat.NewMessageByRoomIDParam) (err error) {
	var lastInsertedID string
	
	insertMessageSql := "INSERT INTO messages VALUES(uuid_generate_v4(), $1, $2, $3, NOW(), NOW()) RETURNING id"
	err = repo.db.QueryRowContext(ctx, insertMessageSql, newMessage.CurrentUserID, newMessage.RoomID, newMessage.Message).Scan(&lastInsertedID)
	if err != nil {
		common.Log(common.LOG_LEVEL_ERROR, err.Error())
		return
	}

	return
}

func (repo *ChatRepositoryPostgree) GetMessage(ctx context.Context, param chat.MessageHistoryParams) (messages []chat.Message, err error) {
	timeAfterDt := time.UnixMilli(param.TimeAfter)

	sqlRoom := "SELECT id FROM rooms WHERE users @> $1"

	row, err := repo.db.QueryContext(ctx, sqlRoom, pq.Array([]string{param.TargetEmail, param.CurrentUserEmail}))
	if err != nil {
		common.Log(common.LOG_LEVEL_ERROR, err.Error())
		common.Log(common.LOG_LEVEL_ERROR, ErrRoomNotFound.Error())
		err = ErrRoomNotFound
		return
	}

	var roomID string

	for row.Next() {
		if err = row.Scan(&roomID); err != nil {
			common.Log(common.LOG_LEVEL_ERROR, err.Error())
			return
		}
	}

	sqlMessage := "SELECT * FROM messages WHERE created_at > $1"

	rows, err := repo.db.QueryContext(ctx, sqlMessage, timeAfterDt)
	if err != nil {
		common.Log(common.LOG_LEVEL_ERROR, err.Error())
		return
	}

	for rows.Next() {
		var message chat.Message

		if err = rows.Scan(&message.ID, &message.UserID, &message.RoomID, &message.Content, &message.CreatedAt, &message.UpdatedAt); err != nil {
			common.Log(common.LOG_LEVEL_ERROR, err.Error())
			return
		}

		messages = append(messages, message)
	}

	return
}

func (repo *ChatRepositoryPostgree) GetRoomsByID(ctx context.Context, currentUserEmail string) (rooms []chat.Room, err error) {
	sqlRoom := "SELECT * FROM rooms WHERE users && $1"

	row, err := repo.db.QueryContext(ctx, sqlRoom, pq.Array([]string{currentUserEmail}))
	if err != nil {
		common.Log(common.LOG_LEVEL_ERROR, err.Error())
		err = ErrRoomNotFound
		return
	}

	for row.Next() {
		var room chat.Room
		if err = row.Scan(&room.ID, &room.Name, pq.Array(&room.Users), &room.CreatedAt, &room.UpdatedAt); err != nil {
			common.Log(common.LOG_LEVEL_ERROR, err.Error())
			return
		}
		rooms = append(rooms, room)
	}

	return
}