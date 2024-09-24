package handlers

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/spqrcor/gofermart/internal/services"
	"github.com/spqrcor/gofermart/internal/services/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetOrdersHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	order := mocks.NewMockOrder(mockCtrl)
	order.EXPECT().GetAll(context.Background()).Return([]services.OrderData{}, services.ErrOrdersNotFound).MaxTimes(1)
	order.EXPECT().GetAll(context.Background()).Return([]services.OrderData{
		services.OrderData{
			Number: "9278923470", Status: "NEW", UploadedAt: "2024-01-01 00:00:00+3",
		},
	}, nil).MinTimes(1)

	tests := []struct {
		name       string
		statusCode int
	}{
		{
			name:       "no content error",
			statusCode: http.StatusNoContent,
		},
		{
			name:       "success",
			statusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/user/orders", nil)
			GetOrdersHandler(order)(rw, req)

			resp := rw.Result()
			assert.Equal(t, tt.statusCode, resp.StatusCode, "Error http status code")
			_ = req.Body.Close()
		})
	}
}
