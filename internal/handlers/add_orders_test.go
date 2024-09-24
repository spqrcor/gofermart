package handlers

import (
	"bytes"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/spqrcor/gofermart/internal/services"
	"github.com/spqrcor/gofermart/internal/services/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddOrdersHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	order := mocks.NewMockOrder(mockCtrl)
	order.EXPECT().Add(context.Background(), "333").Return(services.ErrOrderInvalidFormat).AnyTimes()
	order.EXPECT().Add(context.Background(), "9278923470").Return(nil).MaxTimes(1)
	order.EXPECT().Add(context.Background(), "9278923470").Return(services.ErrOrderUserExists).MinTimes(1)
	order.EXPECT().Add(context.Background(), "12345678903").Return(services.ErrOrderAnotherUserExists).AnyTimes()

	tests := []struct {
		name        string
		contentType string
		body        []byte
		statusCode  int
	}{
		{
			name:        "content type error",
			contentType: "application/json",
			body:        []byte(`<num>3333</num>`),
			statusCode:  http.StatusBadRequest,
		},
		{
			name:        "format number error",
			contentType: "text/plain",
			body:        []byte(`333`),
			statusCode:  http.StatusUnprocessableEntity,
		},
		{
			name:        "success",
			contentType: "text/plain",
			body:        []byte(`9278923470`),
			statusCode:  http.StatusAccepted,
		},
		{
			name:        "order user exists",
			contentType: "text/plain",
			body:        []byte(`9278923470`),
			statusCode:  http.StatusOK,
		},
		{
			name:        "order another user exists",
			contentType: "text/plain",
			body:        []byte(`12345678903`),
			statusCode:  http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/user/orders", bytes.NewReader(tt.body))
			req.Header.Add("Content-Type", tt.contentType)
			AddOrdersHandler(order)(rw, req)

			resp := rw.Result()
			assert.Equal(t, tt.statusCode, resp.StatusCode, "Error http status code")
			_ = req.Body.Close()
		})
	}
}
