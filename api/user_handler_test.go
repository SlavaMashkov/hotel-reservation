package api

import (
	"encoding/json"
	mocks "github.com/SlavaMashkov/hotel-reservation/mocks/db"
	"github.com/SlavaMashkov/hotel-reservation/types"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"testing"
)

func TestUserHandler_HandlePostUser(t *testing.T) {
	userStoreMock := &mocks.UserStore{}

	//userHandler := NewUserHandler(userStoreMock)

	testCreateUserParams := types.CreateUserParams{
		FirstName: "Test",
		LastName:  "Test",
		Email:     "test@example.com",
		Password:  "TestTest",
	}

	testUser := types.User{
		FirstName: "Test",
		LastName:  "Test",
		Email:     "test@example.com",
	}

	requestBody, _ := json.Marshal(testCreateUserParams)

	// TODO:
	slog.Info("test", slog.Any("requestBody", requestBody))

	userStoreMock.On("InsertUser", mock.Anything, testUser).Return(testUser, nil)

	//testContext := &fiber.Ctx{}

	//userHandler.HandlePostUser(testContext)
}
