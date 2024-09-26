package handlers

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/spqrcor/gofermart/internal/services"
	"github.com/spqrcor/gofermart/internal/services/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetBalanceHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	withdrawal := mocks.NewMockWithdrawal(mockCtrl)
	withdrawal.EXPECT().GetBalance(context.Background()).Return(services.BalanceInfo{}, errors.New("user not found")).MaxTimes(1)
	withdrawal.EXPECT().GetBalance(context.Background()).Return(services.BalanceInfo{Current: 100, Withdrawn: 100}, nil).MinTimes(1)
	tests := []struct {
		name       string
		statusCode int
	}{
		{
			name:       "no content error",
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "success",
			statusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/user/balance", nil)
			GetBalanceHandler(withdrawal)(rw, req)

			resp := rw.Result()
			assert.Equal(t, tt.statusCode, resp.StatusCode, "Error http status code")
			_ = resp.Body.Close()
		})
	}

}
