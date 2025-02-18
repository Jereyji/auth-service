package tests

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/Jereyji/auth-service.git/internal/presentation/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	Addr             = "http://localhost:8080"
	registrationAddr = Addr + "/auth/register"
	loginAddr        = Addr + "/auth/login"
)

// Should correctly register user and admin
func TestSuccessRegistration(t *testing.T) {
	cases := []struct {
		RegisterRequest handlers.RegisterRequest
		ExpectedStatus  int
	}{
		{
			RegisterRequest: handlers.RegisterRequest{
				Username:    "Andrey",
				Password:    "1234",
				AccessLevel: 1,
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			RegisterRequest: handlers.RegisterRequest{
				Username:    "Alesha",
				Password:    "1234",
				AccessLevel: 2,
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			RegisterRequest: handlers.RegisterRequest{
				Username:    "Alex",
				Password:    "1234",
				AccessLevel: 2,
			},
			ExpectedStatus: http.StatusOK,
		},
	}

	for _, tc := range cases {
		t.Run("success registration user", func(t *testing.T) {
			request := bytes.NewBufferString(fmt.Sprintf(`
			{
			"username":"%s",
			"password":"%s",
			"access_level":%d
			}
			`, tc.RegisterRequest.Username, tc.RegisterRequest.Password, tc.RegisterRequest.AccessLevel))

			resp, err := http.Post(registrationAddr, "application/json", request)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tc.ExpectedStatus, resp.StatusCode)
		})
	}
}

// Should return statusConflict, if duplicate username
func TestFailRegistration(t *testing.T) {
	cases := []struct {
		RegisterRequest1 handlers.RegisterRequest
		RegisterRequest2 handlers.RegisterRequest
		ExpectedStatus1  int
		ExpectedStatus2  int
	}{
		{
			RegisterRequest1: handlers.RegisterRequest{
				Username:    "Boris",
				Password:    "1234",
				AccessLevel: 1,
			},
			RegisterRequest2: handlers.RegisterRequest{
				Username:    "Boris",
				Password:    "1234",
				AccessLevel: 1,
			},
			ExpectedStatus1: http.StatusOK,
			ExpectedStatus2: http.StatusConflict,
		},
		{
			RegisterRequest1: handlers.RegisterRequest{
				Username:    "Borya",
				Password:    "1234",
				AccessLevel: 1,
			},
			RegisterRequest2: handlers.RegisterRequest{
				Username:    "Borya",
				Password:    "1234",
				AccessLevel: 2,
			},
			ExpectedStatus1: http.StatusOK,
			ExpectedStatus2: http.StatusConflict,
		},
	}

	for _, tc := range cases {
		t.Run("fail registration: duplicate usernames", func(t *testing.T) {
			request1 := bytes.NewBufferString(fmt.Sprintf(`
			{
			"username":"%s",
			"password":"%s",
			"access_level":%d
			}
			`, tc.RegisterRequest1.Username, tc.RegisterRequest1.Password, tc.RegisterRequest1.AccessLevel))

			resp1, err := http.Post(registrationAddr, "application/json", request1)
			require.NoError(t, err)
			defer resp1.Body.Close()

			assert.Equal(t, tc.ExpectedStatus1, resp1.StatusCode)

			request2 := bytes.NewBufferString(fmt.Sprintf(`
			{
			"username":"%s",
			"password":"%s",
			"access_level":%d
			}
			`, tc.RegisterRequest2.Username, tc.RegisterRequest2.Password, tc.RegisterRequest2.AccessLevel))

			resp2, err := http.Post(registrationAddr, "application/json", request2)
			require.NoError(t, err)
			defer resp2.Body.Close()

			assert.Equal(t, tc.ExpectedStatus2, resp2.StatusCode)
		})
	}
}

func TestSuccessLogin(t *testing.T) {
	cases := []struct {
		LoginRequest   handlers.LoginRequest
		ExpectedStatus int
	}{
		{
			LoginRequest: handlers.LoginRequest{
				Username: "Vova",
				Password: "1234",
			},
			ExpectedStatus: http.StatusOK,
		},
	}

	for _, tc := range cases {
		t.Run("Success login", func(t *testing.T) {
			request := bytes.NewBufferString(fmt.Sprintf(`
			{
			"username":"%s",
			"password":"%s"
			}
			`, tc.LoginRequest.Username, tc.LoginRequest.Password))

			resp, err := http.Post(loginAddr, "application/json", request)
			require.NoError(t, err)

			assert.Equal(t, tc.ExpectedStatus, resp.StatusCode)

			defer resp.Body.Close()

		})
	}
}
