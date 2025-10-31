package tests

import (
	"bytes"
	"encoding/json"
	"fastfunds/internal/api/handlers"
	"fastfunds/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockAccountService struct {
	createFn func(*models.CreateAccountRequest) error
	getFn    func(int) (*models.AccountView, error)
}

func (m *mockAccountService) CreateAccount(req *models.CreateAccountRequest) error {
	if m.createFn != nil {
		return m.createFn(req)
	}
	return nil
}
func (m *mockAccountService) GetAccount(id int) (*models.AccountView, error) {
	if m.getFn != nil {
		return m.getFn(id)
	}
	return nil, nil
}

func TestCreateAccountHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cases := []struct {
		name     string
		body     interface{}
		mockErr  error
		wantCode int
		wantBody string
	}{
		{"invalid json", "notjson", nil, http.StatusBadRequest, "Invalid JSON format"},
		{"service error", models.CreateAccountRequest{AccountID: 1, InitialBalance: "10.00"}, assert.AnError, http.StatusBadRequest, assert.AnError.Error()},
		{"success", models.CreateAccountRequest{AccountID: 2, InitialBalance: "20.00"}, nil, http.StatusCreated, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := &mockAccountService{createFn: func(req *models.CreateAccountRequest) error { return tc.mockErr }}
			h := handlers.NewAccountHandler(mockSvc)
			r := gin.Default()
			r.POST("/accounts", h.CreateAccount)
			var reqBody []byte
			if s, ok := tc.body.(string); ok {
				reqBody = []byte(s)
			} else {
				reqBody, _ = json.Marshal(tc.body)
			}
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/accounts", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			assert.Equal(t, tc.wantCode, w.Code)
			if tc.wantBody != "" {
				assert.Contains(t, w.Body.String(), tc.wantBody)
			}
		})
	}
}

func TestGetAccountHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cases := []struct {
		name     string
		param    string
		mockView *models.AccountView
		mockErr  error
		wantCode int
		wantBody string
	}{
		{"bad id", "abc", nil, nil, http.StatusBadRequest, "Invalid account_id format"},
		{"not found", "123", nil, assert.AnError, http.StatusNotFound, assert.AnError.Error()},
		{"success", "1", &models.AccountView{AccountID: 1, CurrentBalance: "10.00"}, nil, http.StatusOK, "10.00"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := &mockAccountService{getFn: func(id int) (*models.AccountView, error) {
				if tc.mockErr != nil {
					return nil, tc.mockErr
				}
				return tc.mockView, nil
			}}
			h := handlers.NewAccountHandler(mockSvc)
			r := gin.Default()
			r.GET("/accounts/:account_id", h.GetAccount)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/accounts/"+tc.param, nil)
			r.ServeHTTP(w, req)
			assert.Equal(t, tc.wantCode, w.Code)
			if tc.wantBody != "" {
				assert.Contains(t, w.Body.String(), tc.wantBody)
			}
		})
	}
}
