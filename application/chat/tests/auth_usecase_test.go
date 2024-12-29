package tests

import (
	"context"
	"testing"

	"github.com/fikrihkll/chat-app/application/chat"
	"github.com/fikrihkll/chat-app/application/chat/repositories"
	"github.com/fikrihkll/chat-app/application/chat/usecases"
	"github.com/fikrihkll/chat-app/config"
	"github.com/fikrihkll/chat-app/infrastructure"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

func TestAuthUseCases(t *testing.T) {
	cfg := config.Load("../../../test.env")

	pgConn, _ := infrastructure.NewPgConnection(cfg)

	persistRepo := repositories.NewUserRepositoryPostgree(pgConn)

	uc := usecases.NewUserApplication(persistRepo)
	ctx := context.Background()

	t.Run("create user", func(t *testing.T) {
		fake := faker.New()
		
		err := uc.RegisterUser(ctx, chat.User{
			Name: fake.App().Name(),
			Email: fake.Person().Contact().Email, 
			Password: fake.Internet().Password(),
		})
		assert.NoError(t, err)
	})

	t.Run("create existing user", func(t *testing.T) {
		fake := faker.New()
		email := fake.Person().Contact().Email
		uc.RegisterUser(ctx, chat.User{
			Name: fake.App().Name(),
			Email: email, 
			Password: fake.Internet().Password(),
		})
		err2 := uc.RegisterUser(ctx, chat.User{
			Name: fake.App().Name(),
			Email: email, 
			Password: fake.Internet().Password(),
		})
		assert.Error(t, err2)
		assert.EqualError(t, err2, usecases.ErrEmailAlreadyExists.Error())
	})

	t.Run("create existing user", func(t *testing.T) {
		fake := faker.New()
		email := fake.Person().Contact().Email
		uc.RegisterUser(ctx, chat.User{
			Name: fake.App().Name(),
			Email: email, 
			Password: fake.Internet().Password(),
		})
		err2 := uc.RegisterUser(ctx, chat.User{
			Name: fake.App().Name(),
			Email: email, 
			Password: fake.Internet().Password(),
		})
		assert.Error(t, err2)
		assert.EqualError(t, err2, usecases.ErrEmailAlreadyExists.Error())
	})
}