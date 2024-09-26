package handlers

import (
	"bytes"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/spqrcor/gofermart/internal/authenticate/mocks"
	"github.com/spqrcor/gofermart/internal/services"
	"github.com/spqrcor/gofermart/internal/services/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userUUID := uuid.New()
	user := mocks.NewMockUser(mockCtrl)
	authenticate := amocks.NewMockAuth(mockCtrl)

	user.EXPECT().Add(context.Background(), services.InputDataUser{
		Login: "s", Password: "123456",
	}).Return(uuid.Nil, services.ErrValidation).AnyTimes()
	user.EXPECT().Add(context.Background(), services.InputDataUser{
		Login: "spqr", Password: "123456",
	}).Return(userUUID, nil).AnyTimes().MaxTimes(1)
	user.EXPECT().Add(context.Background(), services.InputDataUser{
		Login: "spqr", Password: "123456",
	}).Return(uuid.Nil, services.ErrLoginExists).AnyTimes().MinTimes(1)

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
			name:        "validation error",
			contentType: "application/json",
			body:        []byte(`{"login":"s","password":"123456"}`),
			statusCode:  http.StatusBadRequest,
		},
		{
			name:        "success",
			contentType: "application/json",
			body:        []byte(`{"login":"spqr","password":"123456"}`),
			statusCode:  http.StatusOK,
		},
		{
			name:        "user exists error",
			contentType: "application/json",
			body:        []byte(`{"login":"spqr","password":"123456"}`),
			statusCode:  http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/user/register", bytes.NewReader(tt.body))
			authenticate.EXPECT().SetCookie(rw, userUUID).AnyTimes()
			req.Header.Add("Content-Type", tt.contentType)
			RegisterHandler(user, authenticate)(rw, req)

			resp := rw.Result()
			assert.Equal(t, tt.statusCode, resp.StatusCode, "Error http status code")
			_ = resp.Body.Close()
		})
	}
}
