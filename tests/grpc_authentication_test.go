package tests

import (
	"context"
	"testing"

	v1 "github.com/Nikita-Filonov/api-go-autotests-server/gen/go/v1"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestGRPCLogin(t *testing.T) {
	userRequest := &v1.CreateUserRequest{
		Email:      gofakeit.Email(),
		Password:   gofakeit.Password(true, true, true, true, false, 12),
		LastName:   gofakeit.LastName(),
		FirstName:  gofakeit.FirstName(),
		MiddleName: gofakeit.FirstName(),
	}

	connection, err := grpc.NewClient(
		"localhost:9000",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	defer connection.Close()

	usersClient := v1.NewUsersServiceClient(connection)

	createdUser, err := usersClient.CreateUser(context.Background(), userRequest)

	require.NoError(t, err)
	require.NotNil(t, createdUser)
	require.NotNil(t, createdUser.GetUser())
	require.NotEmpty(t, createdUser.GetUser().GetId())

	request := &v1.LoginRequest{
		Email:    userRequest.GetEmail(),
		Password: userRequest.GetPassword(),
	}

	authenticationClient := v1.NewAuthenticationServiceClient(connection)

	response, err := authenticationClient.Login(context.Background(), request)

	require.NoError(t, err)
	require.NotNil(t, response)
	require.NotNil(t, response.GetToken())

	assert.Equal(t, "bearer", response.GetToken().GetTokenType())
	assert.NotEmpty(t, response.GetToken().GetAccessToken())
	assert.NotEmpty(t, response.GetToken().GetRefreshToken())

	t.Logf("user %s successfully logged in", userRequest.GetEmail())
}
