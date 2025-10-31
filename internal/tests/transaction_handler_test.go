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

type mockTransactionService struct {
	processFn func(*models.TransactionRequest) error
}

func (m *mockTransactionService) ProcessTransaction(req *models.TransactionRequest) error {
	if m.processFn != nil {
		return m.processFn(req)
	}
	return nil
}

func TestSubmitTransactionHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cases := []struct {
		name     string
		body     interface{}
		mockErr  error
		wantCode int
		wantBody string
	}{
		{"invalid json", "notjson", nil, http.StatusBadRequest, "Invalid JSON format"},
		{"service error", models.TransactionRequest{SourceAccountID: 1, DestinationAccountID: 2, Amount: "10.00"}, assert.AnError, http.StatusBadRequest, assert.AnError.Error()},
		{"success", models.TransactionRequest{SourceAccountID: 1, DestinationAccountID: 2, Amount: "20.00"}, nil, http.StatusCreated, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := &mockTransactionService{processFn: func(req *models.TransactionRequest) error { return tc.mockErr }}
			h := handlers.NewTransactionHandler(mockSvc)
			r := gin.Default()
			r.POST("/transactions", h.SubmitTransaction)
			var reqBody []byte
			if s, ok := tc.body.(string); ok {
				reqBody = []byte(s)
			} else {
				reqBody, _ = json.Marshal(tc.body)
			}
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/transactions", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			assert.Equal(t, tc.wantCode, w.Code)
			if tc.wantBody != "" {
				assert.Contains(t, w.Body.String(), tc.wantBody)
			}
		})
	}
}
