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

func TestAddWithdrawalHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	withdrawal := mocks.NewMockWithdrawal(mockCtrl)
	withdrawal.EXPECT().Add(context.Background(), services.InputWithdrawal{OrderNum: "333", Sum: 100}).Return(services.ErrOrderInvalidFormat).AnyTimes()
	withdrawal.EXPECT().Add(context.Background(), services.InputWithdrawal{OrderNum: "9278923470", Sum: 100}).Return(nil).AnyTimes()
	withdrawal.EXPECT().Add(context.Background(), services.InputWithdrawal{OrderNum: "12345678903", Sum: 100}).Return(services.ErrBalance).AnyTimes()

	tests := []struct {
		name        string
		contentType string
		body        []byte
		statusCode  int
	}{
		{
			name:        "content type error",
			contentType: "text/plain",
			body:        []byte(`1111`),
			statusCode:  http.StatusInternalServerError,
		},
		{
			name:        "format order error",
			contentType: "application/json",
			body:        []byte(`{"order":"333","sum":100}`),
			statusCode:  http.StatusUnprocessableEntity,
		},
		{
			name:        "success",
			contentType: "application/json",
			body:        []byte(`{"order":"9278923470","sum":100}`),
			statusCode:  http.StatusOK,
		},
		{
			name:        "balance error",
			contentType: "application/json",
			body:        []byte(`{"order":"12345678903","sum":100}`),
			statusCode:  http.StatusPaymentRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/user/balance/withdraw", bytes.NewReader(tt.body))
			req.Header.Add("Content-Type", tt.contentType)
			AddWithdrawalHandler(withdrawal)(rw, req)

			resp := rw.Result()
			assert.Equal(t, tt.statusCode, resp.StatusCode, "Error http status code")
			_ = resp.Body.Close()
		})
	}
}
