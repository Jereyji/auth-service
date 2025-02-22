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
	contentTypeJSON  = "application/json"
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
				Name:     "Andrey",
				Email:    "Andrey@mail.ru",
				Password: "1234",
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			RegisterRequest: handlers.RegisterRequest{
				Name:     "Alesha",
				Email:    "Alesha@mail.ru",
				Password: "1234",
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			RegisterRequest: handlers.RegisterRequest{
				Name:     "Alex",
				Email:    "Alex@mail.ru",
				Password: "1234",
			},
			ExpectedStatus: http.StatusOK,
		},
	}

	for _, tc := range cases {
		t.Run("Registration user", func(t *testing.T) {
			request := bytes.NewBufferString(fmt.Sprintf(`
			{
			"name":"%s",
			"email":"%s",
			"password":"%s"
			}
			`, tc.RegisterRequest.Name, tc.RegisterRequest.Email, tc.RegisterRequest.Password))

			resp, err := http.Post(registrationAddr, contentTypeJSON, request)
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
				Name:     "Boris",
				Email:    "Boris@mail.ru",
				Password: "1234",
			},
			RegisterRequest2: handlers.RegisterRequest{
				Name:     "Boris",
				Email:    "Boris@mail.ru",
				Password: "1234",
			},
			ExpectedStatus1: http.StatusOK,
			ExpectedStatus2: http.StatusConflict,
		},
		{
			RegisterRequest1: handlers.RegisterRequest{
				Name:     "Borya",
				Email:    "Borya@mail.ru",
				Password: "1234",
			},
			RegisterRequest2: handlers.RegisterRequest{
				Name:     "Borya",
				Email:    "Borya@mail.ru",
				Password: "1234",
			},
			ExpectedStatus1: http.StatusOK,
			ExpectedStatus2: http.StatusConflict,
		},
	}

	for _, tc := range cases {
		t.Run("user with duplicate email", func(t *testing.T) {
			request1 := bytes.NewBufferString(fmt.Sprintf(`
			{
			"name":"%s",
			"email":"%s",
			"password":"%s"
			}
			`, tc.RegisterRequest1.Name, tc.RegisterRequest1.Email, tc.RegisterRequest1.Password))
			
			resp1, err := http.Post(registrationAddr, contentTypeJSON, request1)
			require.NoError(t, err)
			defer resp1.Body.Close()

			assert.Equal(t, tc.ExpectedStatus1, resp1.StatusCode)

			request2 := bytes.NewBufferString(fmt.Sprintf(`
			{
			"name":"%s",
			"email":"%s",
			"password":"%s"
			}
			`, tc.RegisterRequest2.Name, tc.RegisterRequest2.Email, tc.RegisterRequest2.Password))

			resp2, err := http.Post(registrationAddr, contentTypeJSON, request2)
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
				Email:    "Andrey@mail.ru",
				Password: "1234",
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			LoginRequest: handlers.LoginRequest{
				Email:    "Alesha@mail.ru",
				Password: "1234",
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			LoginRequest: handlers.LoginRequest{
				Email:    "Alex@mail.ru",
				Password: "1234",
			},
			ExpectedStatus: http.StatusOK,
		},
	}

	for _, tc := range cases {
		t.Run("Login user", func(t *testing.T) {
			request := bytes.NewBufferString(fmt.Sprintf(`
			{
				"email":"%s",
				"password":"%s"
			}
			`, tc.LoginRequest.Email, tc.LoginRequest.Password))

			resp, err := http.Post(loginAddr, contentTypeJSON, request)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tc.ExpectedStatus, resp.StatusCode)

			cookies := resp.Cookies()

			var accessToken, refreshToken string
			for _, cookie := range cookies {
				if cookie.Name == "access_token" {
					accessToken = cookie.Value
				}
				if cookie.Name == "refresh_token" {
					refreshToken = cookie.Value
				}
			}

			assert.NotEmpty(t, accessToken)
			assert.NotEmpty(t, refreshToken)
		})
	}
}
