package tests

import (
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"resty.dev/v3"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type token struct {
	TokenType    string `json:"tokenType"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type loginResponse struct {
	Token token `json:"token"`
}

func TestHTTPLogin(t *testing.T) {
	userRequest := createUserRequest{
		Email:      gofakeit.Email(),
		Password:   gofakeit.Password(true, true, true, true, false, 12),
		LastName:   gofakeit.LastName(),
		FirstName:  gofakeit.FirstName(),
		MiddleName: gofakeit.MiddleName(),
	}

	client := resty.New().
		SetBaseURL("http://localhost:8000/api/v1")

	var createdUser userResponse

	createUserResponse, err := client.R().
		SetBody(userRequest).
		SetResult(&createdUser).
		Post("/users")

	require.NoError(t, err)
	require.NotNil(t, createUserResponse)
	require.Equal(t, http.StatusOK, createUserResponse.StatusCode())
	require.NotEmpty(t, createdUser.User.ID)

	request := loginRequest{
		Email:    userRequest.Email,
		Password: userRequest.Password,
	}

	var result loginResponse

	response, err := client.
		R().
		SetBody(request).
		SetResult(&result).
		Post("/authentication/login")

	require.NoError(t, err)
	require.NotNil(t, response)
	require.Equal(t, http.StatusOK, response.StatusCode())

	assert.Equal(t, "bearer", result.Token.TokenType)
	assert.NotEmpty(t, result.Token.AccessToken)
	assert.NotEmpty(t, result.Token.RefreshToken)

	t.Logf("user %s successfully logged in", userRequest.Email)
}
