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

func TestGetWithdrawalsHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	withdrawal := mocks.NewMockWithdrawal(mockCtrl)
	withdrawal.EXPECT().GetAll(context.Background()).Return([]services.WithdrawalData{}, services.ErrWithdrawNotFound).MaxTimes(1)
	withdrawal.EXPECT().GetAll(context.Background()).Return([]services.WithdrawalData{
		services.WithdrawalData{
			OrderNum: "9278923470", Sum: 100, ProcessedAt: "2024-01-01 00:00:00+3",
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
			req := httptest.NewRequest("GET", "/api/user/withdrawals", nil)
			GetWithdrawalsHandler(withdrawal)(rw, req)

			resp := rw.Result()
			assert.Equal(t, tt.statusCode, resp.StatusCode, "Error http status code")
			_ = req.Body.Close()
		})
	}
}
