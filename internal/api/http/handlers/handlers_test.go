package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mockservice "github.com/Baraulia/anti_bruteforce_service/internal/api/mocks"
	"github.com/Baraulia/anti_bruteforce_service/internal/models"
	"github.com/Baraulia/anti_bruteforce_service/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheck(t *testing.T) {
	type mockBehavior func(s *mockservice.MockApplicationInterface)
	testTable := []struct {
		name                string
		mockBehavior        mockBehavior
		method              string
		inputBody           string
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			mockBehavior: func(s *mockservice.MockApplicationInterface) {
				s.EXPECT().Check(gomock.Any(), models.Data{
					Login:    "Test",
					Password: "Test",
					IP:       "192.1.1.0/25",
				}).Return(true, nil)
			},
			method:              http.MethodPost,
			inputBody:           `{"login":"Test","password":"Test","ip":"192.1.1.0/25"}`,
			expectedStatusCode:  200,
			expectedRequestBody: "ok=true",
		},
		{
			name: "OK with not allowed",
			mockBehavior: func(s *mockservice.MockApplicationInterface) {
				s.EXPECT().Check(gomock.Any(), models.Data{
					Login:    "Test",
					Password: "Test",
					IP:       "192.1.1.0/25",
				}).Return(false, nil)
			},
			method:              http.MethodPost,
			inputBody:           `{"login":"Test","password":"Test","ip":"192.1.1.0/25"}`,
			expectedStatusCode:  200,
			expectedRequestBody: "ok=false",
		},
		{
			name:                "invalid IP",
			mockBehavior:        func(s *mockservice.MockApplicationInterface) {},
			inputBody:           `{"login":"Test","password":"Test","ip":"192.1.1.0/255.255.255.128"}`,
			expectedStatusCode:  400,
			method:              http.MethodPost,
			expectedRequestBody: "invalid ip in request\n",
		},
		{
			name:                "empty login",
			mockBehavior:        func(s *mockservice.MockApplicationInterface) {},
			inputBody:           `{"login":"","password":"Test","ip":"192.1.1.0/25"}`,
			expectedStatusCode:  400,
			method:              http.MethodPost,
			expectedRequestBody: "empty login in request\n",
		},
		{
			name:                "empty password",
			mockBehavior:        func(s *mockservice.MockApplicationInterface) {},
			inputBody:           `{"login":"Test","password":"","ip":"192.1.1.0/25"}`,
			expectedStatusCode:  400,
			method:              http.MethodPost,
			expectedRequestBody: "empty password in request\n",
		},
		{
			name:                "invalid method",
			mockBehavior:        func(s *mockservice.MockApplicationInterface) {},
			inputBody:           `{"login":"Test","password":"","ip":"192.1.1.0/25"}`,
			expectedStatusCode:  405,
			method:              http.MethodDelete,
			expectedRequestBody: "Method Not Allowed(want method POST)\n",
		},
		{
			name: "Server error",
			mockBehavior: func(s *mockservice.MockApplicationInterface) {
				s.EXPECT().Check(gomock.Any(), models.Data{
					Login:    "Test",
					Password: "Test",
					IP:       "192.1.1.0/25",
				}).Return(false, errors.New("server error"))
			},
			inputBody:           `{"login":"Test","password":"Test","ip":"192.1.1.0/25"}`,
			expectedStatusCode:  500,
			method:              http.MethodPost,
			expectedRequestBody: "server error\n",
		},
		{
			name:                "invalid json",
			mockBehavior:        func(s *mockservice.MockApplicationInterface) {},
			inputBody:           `{"login"="Test","password":"Test","ip":"192.1.1.0/25"}`,
			expectedStatusCode:  400,
			method:              http.MethodPost,
			expectedRequestBody: "invalid character '=' after object key\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			appInterface := mockservice.NewMockApplicationInterface(c)
			testCase.mockBehavior(appInterface)

			logg, err := logger.GetLogger("INFO", false)
			require.NoError(t, err)

			handler := NewHandler(logg, appInterface)
			r := handler.InitRoutes()
			w := httptest.NewRecorder()
			req := httptest.NewRequest(testCase.method, "/check", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestAddToBlackList(t *testing.T) {
	type mockBehavior func(s *mockservice.MockApplicationInterface)
	testTable := []struct {
		name                string
		mockBehavior        mockBehavior
		method              string
		inputBody           string
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			mockBehavior: func(s *mockservice.MockApplicationInterface) {
				s.EXPECT().AddToBlackList(gomock.Any(), "192.1.1.0/25").Return(nil)
			},
			method:              http.MethodPost,
			inputBody:           `{"ip":"192.1.1.0/25"}`,
			expectedStatusCode:  201,
			expectedRequestBody: "",
		},
		{
			name:                "invalid IP",
			mockBehavior:        func(s *mockservice.MockApplicationInterface) {},
			inputBody:           `{"ip":"192.1.1.0/255.255.255.192"}`,
			expectedStatusCode:  400,
			method:              http.MethodPost,
			expectedRequestBody: "invalid ip in request\n",
		},
		{
			name:                "invalid method",
			mockBehavior:        func(s *mockservice.MockApplicationInterface) {},
			inputBody:           `{"ip":"192.1.1.0/25"}`,
			expectedStatusCode:  405,
			method:              http.MethodDelete,
			expectedRequestBody: "Method Not Allowed(want method POST)\n",
		},
		{
			name: "Server error",
			mockBehavior: func(s *mockservice.MockApplicationInterface) {
				s.EXPECT().AddToBlackList(gomock.Any(), "192.1.1.0/25").Return(errors.New("server error"))
			},
			inputBody:           `{"ip":"192.1.1.0/25"}`,
			expectedStatusCode:  500,
			method:              http.MethodPost,
			expectedRequestBody: "server error\n",
		},
		{
			name:                "invalid json",
			mockBehavior:        func(s *mockservice.MockApplicationInterface) {},
			inputBody:           `{"ip"="192.1.1.0/25"}`,
			expectedStatusCode:  400,
			method:              http.MethodPost,
			expectedRequestBody: "invalid character '=' after object key\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			appInterface := mockservice.NewMockApplicationInterface(c)
			testCase.mockBehavior(appInterface)

			logg, err := logger.GetLogger("INFO", false)
			require.NoError(t, err)

			handler := NewHandler(logg, appInterface)
			r := handler.InitRoutes()
			w := httptest.NewRecorder()
			req := httptest.NewRequest(testCase.method, "/blacklist/add", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestAddToWhiteList(t *testing.T) {
	type mockBehavior func(s *mockservice.MockApplicationInterface)
	testTable := []struct {
		name                string
		mockBehavior        mockBehavior
		method              string
		inputBody           string
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			mockBehavior: func(s *mockservice.MockApplicationInterface) {
				s.EXPECT().AddToWhiteList(gomock.Any(), "192.1.1.0/25").Return(nil)
			},
			method:              http.MethodPost,
			inputBody:           `{"ip":"192.1.1.0/25"}`,
			expectedStatusCode:  201,
			expectedRequestBody: "",
		},
		{
			name:                "invalid IP",
			mockBehavior:        func(s *mockservice.MockApplicationInterface) {},
			inputBody:           `{"ip":"192.1.1.0/255.255.255.192"}`,
			expectedStatusCode:  400,
			method:              http.MethodPost,
			expectedRequestBody: "invalid ip in request\n",
		},
		{
			name:                "invalid method",
			mockBehavior:        func(s *mockservice.MockApplicationInterface) {},
			inputBody:           `{"ip":"192.1.1.0/25"}`,
			expectedStatusCode:  405,
			method:              http.MethodDelete,
			expectedRequestBody: "Method Not Allowed(want method POST)\n",
		},
		{
			name: "Server error",
			mockBehavior: func(s *mockservice.MockApplicationInterface) {
				s.EXPECT().AddToWhiteList(gomock.Any(), "192.1.1.0/25").Return(errors.New("server error"))
			},
			inputBody:           `{"ip":"192.1.1.0/25"}`,
			expectedStatusCode:  500,
			method:              http.MethodPost,
			expectedRequestBody: "server error\n",
		},
		{
			name:                "invalid json",
			mockBehavior:        func(s *mockservice.MockApplicationInterface) {},
			inputBody:           `{"ip"="192.1.1.0/25"}`,
			expectedStatusCode:  400,
			method:              http.MethodPost,
			expectedRequestBody: "invalid character '=' after object key\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			appInterface := mockservice.NewMockApplicationInterface(c)
			testCase.mockBehavior(appInterface)

			logg, err := logger.GetLogger("INFO", false)
			require.NoError(t, err)

			handler := NewHandler(logg, appInterface)
			r := handler.InitRoutes()
			w := httptest.NewRecorder()
			req := httptest.NewRequest(testCase.method, "/whitelist/add", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestDeleteFromBlackList(t *testing.T) {
	type mockBehavior func(s *mockservice.MockApplicationInterface)
	testTable := []struct {
		name                string
		mockBehavior        mockBehavior
		method              string
		ip                  string
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			mockBehavior: func(s *mockservice.MockApplicationInterface) {
				s.EXPECT().RemoveFromBlackList(gomock.Any(), "192.1.1.0/25").Return(nil)
			},
			method:              http.MethodDelete,
			ip:                  "192.1.1.0/25",
			expectedStatusCode:  204,
			expectedRequestBody: "",
		},
		{
			name:                "invalid IP",
			mockBehavior:        func(s *mockservice.MockApplicationInterface) {},
			ip:                  "192.1.1.0/255.255.255.192",
			expectedStatusCode:  400,
			method:              http.MethodDelete,
			expectedRequestBody: "invalid ip in request\n",
		},
		{
			name:                "invalid method",
			mockBehavior:        func(s *mockservice.MockApplicationInterface) {},
			ip:                  "192.1.1.0/25",
			expectedStatusCode:  405,
			method:              http.MethodPut,
			expectedRequestBody: "Method Not Allowed(want method DELETE)\n",
		},
		{
			name: "Server error",
			mockBehavior: func(s *mockservice.MockApplicationInterface) {
				s.EXPECT().RemoveFromBlackList(gomock.Any(), "192.1.1.0/25").Return(errors.New("server error"))
			},
			ip:                  "192.1.1.0/25",
			expectedStatusCode:  500,
			method:              http.MethodDelete,
			expectedRequestBody: "server error\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			appInterface := mockservice.NewMockApplicationInterface(c)
			testCase.mockBehavior(appInterface)

			logg, err := logger.GetLogger("INFO", false)
			require.NoError(t, err)

			handler := NewHandler(logg, appInterface)
			r := handler.InitRoutes()
			w := httptest.NewRecorder()
			req := httptest.NewRequest(testCase.method, fmt.Sprintf("/blacklist/remove?ip=%s", testCase.ip), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestDeleteFromWhiteList(t *testing.T) {
	type mockBehavior func(s *mockservice.MockApplicationInterface)
	testTable := []struct {
		name                string
		mockBehavior        mockBehavior
		method              string
		ip                  string
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			mockBehavior: func(s *mockservice.MockApplicationInterface) {
				s.EXPECT().RemoveFromWhiteList(gomock.Any(), "192.1.1.0/25").Return(nil)
			},
			method:              http.MethodDelete,
			ip:                  "192.1.1.0/25",
			expectedStatusCode:  204,
			expectedRequestBody: "",
		},
		{
			name:                "invalid IP",
			mockBehavior:        func(s *mockservice.MockApplicationInterface) {},
			ip:                  "192.1.1.0/255.255.255.192",
			expectedStatusCode:  400,
			method:              http.MethodDelete,
			expectedRequestBody: "invalid ip in request\n",
		},
		{
			name:                "invalid method",
			mockBehavior:        func(s *mockservice.MockApplicationInterface) {},
			ip:                  "192.1.1.0/25",
			expectedStatusCode:  405,
			method:              http.MethodPut,
			expectedRequestBody: "Method Not Allowed(want method DELETE)\n",
		},
		{
			name: "Server error",
			mockBehavior: func(s *mockservice.MockApplicationInterface) {
				s.EXPECT().RemoveFromWhiteList(gomock.Any(), "192.1.1.0/25").Return(errors.New("server error"))
			},
			ip:                  "192.1.1.0/25",
			expectedStatusCode:  500,
			method:              http.MethodDelete,
			expectedRequestBody: "server error\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			appInterface := mockservice.NewMockApplicationInterface(c)
			testCase.mockBehavior(appInterface)

			logg, err := logger.GetLogger("INFO", false)
			require.NoError(t, err)

			handler := NewHandler(logg, appInterface)
			r := handler.InitRoutes()
			w := httptest.NewRecorder()
			req := httptest.NewRequest(testCase.method, fmt.Sprintf("/whitelist/remove?ip=%s", testCase.ip), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestClearBuckets(t *testing.T) {
	type mockBehavior func(s *mockservice.MockApplicationInterface)
	testTable := []struct {
		name                string
		mockBehavior        mockBehavior
		method              string
		inputBody           string
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			mockBehavior: func(s *mockservice.MockApplicationInterface) {
				s.EXPECT().ClearBuckets(gomock.Any(), models.Data{
					Login: "Test",
					IP:    "192.1.1.0/25",
				}).Return(nil)
			},
			method:              http.MethodPost,
			inputBody:           `{"login":"Test","ip":"192.1.1.0/25"}`,
			expectedStatusCode:  204,
			expectedRequestBody: "",
		},
		{
			name:                "invalid IP",
			mockBehavior:        func(s *mockservice.MockApplicationInterface) {},
			inputBody:           `{"login":"Test","ip":"192.1.1.0/255.255.255.128"}`,
			expectedStatusCode:  400,
			method:              http.MethodPost,
			expectedRequestBody: "invalid ip in request\n",
		},
		{
			name:                "empty login",
			mockBehavior:        func(s *mockservice.MockApplicationInterface) {},
			inputBody:           `{"login":"","ip":"192.1.1.0/25"}`,
			expectedStatusCode:  400,
			method:              http.MethodPost,
			expectedRequestBody: "empty login in request\n",
		},
		{
			name:                "invalid method",
			mockBehavior:        func(s *mockservice.MockApplicationInterface) {},
			inputBody:           `{"login":"Test","ip":"192.1.1.0/25"}`,
			expectedStatusCode:  405,
			method:              http.MethodDelete,
			expectedRequestBody: "Method Not Allowed(want method POST)\n",
		},
		{
			name: "Server error",
			mockBehavior: func(s *mockservice.MockApplicationInterface) {
				s.EXPECT().ClearBuckets(gomock.Any(), models.Data{
					Login: "Test",
					IP:    "192.1.1.0/25",
				}).Return(errors.New("server error"))
			},
			inputBody:           `{"login":"Test","ip":"192.1.1.0/25"}`,
			expectedStatusCode:  500,
			method:              http.MethodPost,
			expectedRequestBody: "server error\n",
		},
		{
			name:                "invalid json",
			mockBehavior:        func(s *mockservice.MockApplicationInterface) {},
			inputBody:           `{"login"="Test","ip":"192.1.1.0/25"}`,
			expectedStatusCode:  400,
			method:              http.MethodPost,
			expectedRequestBody: "invalid character '=' after object key\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			appInterface := mockservice.NewMockApplicationInterface(c)
			testCase.mockBehavior(appInterface)

			logg, err := logger.GetLogger("INFO", false)
			require.NoError(t, err)

			handler := NewHandler(logg, appInterface)
			r := handler.InitRoutes()
			w := httptest.NewRecorder()
			req := httptest.NewRequest(testCase.method, "/clear", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
