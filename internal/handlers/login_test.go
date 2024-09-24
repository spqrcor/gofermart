package handlers

import (
	"bytes"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	amocks "github.com/spqrcor/gofermart/internal/authenticate/mocks"
	"github.com/spqrcor/gofermart/internal/services"
	"github.com/spqrcor/gofermart/internal/services/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userUuid := uuid.New()
	user := mocks.NewMockUser(mockCtrl)
	authenticate := amocks.NewMockAuth(mockCtrl)
	user.EXPECT().Login(context.Background(), services.InputDataUser{
		Login: "spqr", Password: "1",
	}).Return(uuid.Nil, services.ErrLogin).AnyTimes()
	user.EXPECT().Login(context.Background(), services.InputDataUser{
		Login: "spqr", Password: "123456",
	}).Return(userUuid, nil).AnyTimes()

	tests := []struct {
		name        string
		contentType string
		body        []byte
		statusCode  int
	}{
		{
			name:        "not format error",
			contentType: "application/json",
			body:        []byte(`<num>3333</num>`),
			statusCode:  http.StatusBadRequest,
		},
		{
			name:        "login error",
			contentType: "application/json",
			body:        []byte(`{"login":"spqr","password":"1"}`),
			statusCode:  http.StatusUnauthorized,
		},
		{
			name:        "success",
			contentType: "application/json",
			body:        []byte(`{"login":"spqr","password":"123456"}`),
			statusCode:  http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/user/login", bytes.NewReader(tt.body))
			req.Header.Add("Content-Type", tt.contentType)
			authenticate.EXPECT().SetCookie(rw, userUuid).AnyTimes()
			LoginHandler(user, authenticate)(rw, req)

			resp := rw.Result()
			assert.Equal(t, tt.statusCode, resp.StatusCode, "Error http status code")
			_ = req.Body.Close()
		})
	}
}
