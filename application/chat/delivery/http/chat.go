package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/fikrihkll/chat-app/application/chat"
	"github.com/fikrihkll/chat-app/application/chat/repositories"
	"github.com/fikrihkll/chat-app/application/chat/transport"
	"github.com/fikrihkll/chat-app/application/chat/usecases"
	"github.com/fikrihkll/chat-app/common"
	"github.com/fikrihkll/chat-app/common/middleware"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ChatHttpApi struct {
	chatUseCase chat.IChatUseCase
	authUseCase chat.IAuthUseCase
}

func NewChatHttpApi(chatUseCase chat.IChatUseCase, authUseCase chat.IAuthUseCase) *ChatHttpApi {
	return &ChatHttpApi{chatUseCase, authUseCase}
}

// @Description gain access to API
// @Tags auth
// @Accept json
// @Produce json
// @Param login body transport.Login true "Login detail"
// @Success 201
// @Router /auth/login [post]
func (handler *ChatHttpApi) Login(c echo.Context) error {
	var body transport.Login

	if c.Bind(&body) != nil {
		return c.JSON(http.StatusBadRequest, common.BaseResponse{
			Message: common.BadRequestError.Error(),
			Data:    nil,
		})
	}

	if err := body.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, common.BaseResponse{
			Message: common.BadRequestError.Error(),
			Data:    nil,
		})
	}

	response, err := handler.authUseCase.Login(c.Request().Context(), chat.LoginParam{
		Email:    body.Email,
		Password: body.Password,
	})

	if err != nil {
		common.Log(common.LOG_LEVEL_ERROR, err.Error())

		if err == usecases.ErrNotRegistered {
			return c.JSON(http.StatusBadRequest, &common.BaseResponse{
				Message: usecases.ErrNotRegistered.Error(),
				Data:    nil,
			})
		}

		return c.JSON(http.StatusInternalServerError, &common.BaseResponse{
			Message: common.InternalServerError.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, &common.BaseResponse{
		Message: common.HttpSuccess,
		Data:    response,
	})
}

// @Description create new account
// @Tags auth
// @Param register body transport.Register true "Register information"
// @Accept json
// @Produce json
// @Success 201
// @Router /auth/register [post]
func (handler *ChatHttpApi) Register(c echo.Context) error {
	var body transport.Register

	if c.Bind(&body) != nil {
		return c.JSON(http.StatusBadRequest, common.BaseResponse{
			Message: common.BadRequestError.Error(),
			Data:    nil,
		})
	}

	if err := body.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, common.BaseResponse{
			Message: err.Error(),
			Data:    nil,
		})
	}

	if err := handler.authUseCase.RegisterUser(
		c.Request().Context(),
		chat.User{
			Name:     body.Name,
			Email:    body.Email,
			Password: body.Password,
		},
	); err != nil {
		switch {
		case errors.Is(err, usecases.ErrEmailAlreadyExists):
			return c.JSON(http.StatusConflict, &common.BaseResponse{
				Message: err.Error(),
				Data:    nil,
			})
		case errors.Is(err, usecases.ErrMalformatEmail):
			return c.JSON(http.StatusBadRequest, &common.BaseResponse{
				Message: err.Error(),
				Data:    nil,
			})
		default:
			return c.JSON(http.StatusInternalServerError, &common.BaseResponse{
				Message: err.Error(),
				Data:    nil,
			})
		}
	}

	return c.JSON(http.StatusCreated, &common.BaseResponse{
		Message: common.HttpSuccessCreated,
		Data:    nil,
	})

}

// @Description check token validation
// @Security BearerAuth
// @Tags auth
// @Param Authorization header string true "Bearer token"
// @Accept json
// @Produce json
// @Success 200
// @Router /auth/validate [get]
func (handler *ChatHttpApi) Validate(c echo.Context) error {
	userID, ok := c.Get("id").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, &common.BaseResponse{
			Message: common.UnauthorizedError.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, &common.BaseResponse{
		Message: common.HttpSuccess,
		Data:    userID,
	})
}

// @Description save new message
// @Security BearerAuth
// @Tags chat
// @Param Authorization header string true "Bearer token"
// @Param message body transport.NewMessageByEmail true "New message detail"
// @Accept json
// @Produce json
// @Success 201
// @Router /chat/send [post]
func (handler *ChatHttpApi) SaveMessageByEmail(c echo.Context) error {
	var body transport.NewMessageByEmail
	if c.Bind(&body) != nil {
		return c.JSON(http.StatusBadRequest, &common.BaseResponse{
			Message: common.BadRequestError.Error(),
			Data:    nil,
		})
	}

	if err := body.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, common.BaseResponse{
			Message: common.BadRequestError.Error(),
			Data:    nil,
		})
	}

	userID, ok := c.Get("id").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, &common.BaseResponse{
			Message: common.UnauthorizedError.Error(),
			Data:    nil,
		})
	}

	userEmail, ok := c.Get("email").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, &common.BaseResponse{
			Message: common.UnauthorizedError.Error(),
			Data:    nil,
		})
	}

	if err := handler.chatUseCase.SaveMessageByEmail(
		c.Request().Context(),
		chat.NewMessageByEmailParam{
			CurrentUserID:    userID,
			CurrentUserEmail: userEmail,
			MemberEmail:      body.MemberEmail,
			Message:          body.Message,
		},
	); err != nil {
		if err == common.UserNotFoundError {
			return c.JSON(http.StatusInternalServerError, &common.BaseResponse{
				Message: common.UserNotFoundError.Error(),
				Data:    nil,
			})	
		}
		
		return c.JSON(http.StatusInternalServerError, &common.BaseResponse{
			Message: common.InternalServerError.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusCreated, &common.BaseResponse{
		Message: common.HttpSuccessCreated,
		Data:    nil,
	})
}

// @Description save new message
// @Security BearerAuth
// @Tags chat
// @Param Authorization header string true "Bearer token"
// @Param message body transport.NewMessageByRoomID true "New message detail"
// @Accept json
// @Produce json
// @Success 201
// @Router /chat/{room_id}/send [post]
func (handler *ChatHttpApi) SaveMessageByRoomID(c echo.Context) error {
	roomID, err := uuid.Parse(c.Param("room_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, &common.BaseResponse{
			Message: common.BadRequestError.Error(),
			Data:    nil,
		})
	}

	var body transport.NewMessageByRoomID
	if c.Bind(&body) != nil {
		return c.JSON(http.StatusBadRequest, &common.BaseResponse{
			Message: common.BadRequestError.Error(),
			Data:    nil,
		})
	}

	if err := body.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, common.BaseResponse{
			Message: common.BadRequestError.Error(),
			Data:    nil,
		})
	}

	userID, ok := c.Get("id").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, &common.BaseResponse{
			Message: common.UnauthorizedError.Error(),
			Data:    nil,
		})
	}

	if err := handler.chatUseCase.SaveMessageByRoomID(
		c.Request().Context(),
		chat.NewMessageByRoomIDParam{
			CurrentUserID: userID,
			RoomID:        roomID.String(),
			Message:       body.Message,
		},
	); err != nil {
		return c.JSON(http.StatusInternalServerError, &common.BaseResponse{
			Message: common.InternalServerError.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusCreated, &common.BaseResponse{
		Message: common.HttpSuccessCreated,
		Data:    nil,
	})
}

// @Description get message history
// @Security BearerAuth
// @Tags chat
// @Param Authorization header string true "Bearer token"
// @Param time_after query string true "timestamp last message retrieved"
// @Param target_email query string true "user that is in the same chat room"
// @Accept json
// @Produce json
// @Success 200
// @Router /chat/get [get]
func (handler *ChatHttpApi) GetMessage(c echo.Context) error {
	timeAfterStr := c.QueryParam("time_after")
	targetEmail := c.QueryParam("target_email")

	if targetEmail == "" {
		return c.JSON(http.StatusBadRequest, &common.BaseResponse{
			Message: common.BadRequestError.Error(),
			Data:    nil,
		})
	}

	if timeAfterStr == "" {
		return c.JSON(http.StatusBadRequest, &common.BaseResponse{
			Message: common.BadRequestError.Error(),
			Data:    nil,
		})
	}

	timeAfter, err := strconv.Atoi(timeAfterStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &common.BaseResponse{
			Message: common.BadRequestError.Error(),
			Data:    nil,
		})
	}

	email, ok := c.Get("email").(string)
	if !ok || email == "" {
		return c.JSON(http.StatusUnauthorized, &common.BaseResponse{
			Message: common.UnauthorizedError.Error(),
			Data:    nil,
		})
	}

	messages, err := handler.chatUseCase.GetMessages(
		c.Request().Context(),
		chat.MessageHistoryParams{
			TargetEmail:      targetEmail,
			TimeAfter:        int64(timeAfter),
			CurrentUserEmail: email,
		},
	)

	if err != nil {
		if err == repositories.ErrRoomNotFound {
			return c.JSON(http.StatusBadRequest, &common.BaseResponse{
				Message: repositories.ErrRoomNotFound.Error(),
				Data:    nil,
			})
		}

		return c.JSON(http.StatusInternalServerError, &common.BaseResponse{
			Message: common.InternalServerError.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, &common.BaseResponse{
		Message: common.HttpSuccess,
		Data:    messages,
	})
}

// @Description get chat rooms that the user in
// @Security BearerAuth
// @Tags chat
// @Param Authorization header string true "Bearer token"
// @Accept json
// @Produce json
// @Success 200
// @Router /chat/rooms [get]
func (handler *ChatHttpApi) GetRoomsByID(c echo.Context) error {
	email, ok := c.Get("email").(string)
	if !ok || email == "" {
		return c.JSON(http.StatusUnauthorized, &common.BaseResponse{
			Message: common.UnauthorizedError.Error(),
			Data:    nil,
		})
	}

	rooms, err := handler.chatUseCase.GetRoomsByID(
		c.Request().Context(),
		email,
	)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &common.BaseResponse{
			Message: common.InternalServerError.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, &common.BaseResponse{
		Message: common.HttpSuccess,
		Data:    rooms,
	})
}

func (handler *ChatHttpApi) HandleRootRoute(e *echo.Echo) {
	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, &common.BaseResponse{Message: "pong", Data: nil})
	})
}

func (handler *ChatHttpApi) HandleChatRoute(e *echo.Echo) {
	g := e.Group("/chat")

	g.POST("/send", handler.SaveMessageByEmail, middleware.AuthMiddleware)
	g.POST("/:room_id/send", handler.SaveMessageByRoomID, middleware.AuthMiddleware)
	g.GET("/get", handler.GetMessage, middleware.AuthMiddleware)
	g.GET("/rooms", handler.GetRoomsByID, middleware.AuthMiddleware)
}

func (handler *ChatHttpApi) HandleAuthRoute(e *echo.Echo) {
	g := e.Group("/auth")

	g.POST("/login", handler.Login)
	g.POST("/register", handler.Register)
	g.GET("/validate", handler.Validate, middleware.AuthMiddleware)
}
