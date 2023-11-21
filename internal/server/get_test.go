package server

import (
	"fmt"
	"github.com/DavidGQK/go-loyalty-system/internal/mocks"
	"github.com/DavidGQK/go-loyalty-system/internal/store"
	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestServer_GetOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	createdAt, err := time.Parse("01/02/2006 15:04:05", "07/24/2023 15:15:45")
	if err != nil {
		t.Fatal(err)
	}

	m := mocks.NewMockRepository(ctrl)
	m.EXPECT().FindUserByToken(gomock.Any(), allowedTokenHash).Return(&store.User{ID: 1}, nil)
	m.EXPECT().FindUserByToken(gomock.Any(), allowedToken2Hash).Return(&store.User{ID: 2}, nil)
	m.EXPECT().FindUserByToken(gomock.Any(), wrongTokenHash).Return(nil, fmt.Errorf("invalid token"))
	m.EXPECT().FindOrdersByUserID(gomock.Any(), 1).Return([]store.Order{
		{
			ID:          1,
			OrderNumber: "123",
			Status:      "NEW",
			CreatedAt:   createdAt,
			BonusAmount: 1,
		},
	}, nil)
	m.EXPECT().FindOrdersByUserID(gomock.Any(), 2).Return([]store.Order{}, nil)

	type fields struct {
		Repository Repository
		AuthToken  string
	}
	type results struct {
		statusCode int
		response   string
	}
	tests := []struct {
		name   string
		fields fields
		want   results
	}{
		{
			name: "success response",
			fields: fields{
				Repository: m,
				AuthToken:  allowedToken,
			},
			want: results{
				statusCode: http.StatusOK,
				response:   `[{"number":"123","status":"NEW","accrual":0.00001,"uploaded_at":"2023-07-24T15:15:45Z"}]`,
			},
		},
		{
			name: "empty response",
			fields: fields{
				Repository: m,
				AuthToken:  allowedToken2,
			},
			want: results{
				statusCode: http.StatusNoContent,
				response:   ``,
			},
		},
		{
			name: "bad response",
			fields: fields{
				Repository: m,
				AuthToken:  wrongToken,
			},
			want: results{
				statusCode: http.StatusInternalServerError,
				response:   ``,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Repository: tt.fields.Repository,
			}

			r := SetUpRouter()
			r.GET("/api/user/orders", s.GetOrders)
			req, _ := http.NewRequest("GET", "/api/user/orders", nil)
			w := httptest.NewRecorder()

			req.Header.Set("Authorization", tt.fields.AuthToken)
			r.ServeHTTP(w, req)

			responseData, _ := io.ReadAll(w.Body)
			assert.Equal(t, tt.want.response, string(responseData))
			assert.Equal(t, tt.want.statusCode, w.Code)
		})
	}
}

func TestServer_GetUserBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockRepository(ctrl)
	m.EXPECT().FindUserByToken(gomock.Any(), allowedTokenHash).Return(&store.User{ID: 1, Bonuses: 100}, nil)
	m.EXPECT().FindUserByToken(gomock.Any(), wrongTokenHash).Return(nil, fmt.Errorf("invalid token"))
	m.EXPECT().GetWithdrawalSumByUserID(gomock.Any(), 1).Return(500, nil)

	type fields struct {
		Repository Repository
		AuthToken  string
	}
	type results struct {
		statusCode int
		response   string
	}
	tests := []struct {
		name   string
		fields fields
		want   results
	}{
		{
			name: "success response",
			fields: fields{
				Repository: m,
				AuthToken:  allowedToken,
			},
			want: results{
				statusCode: http.StatusOK,
				response:   `{"current":0.001,"withdrawn":0.005}`,
			},
		},
		{
			name: "invalid token",
			fields: fields{
				Repository: m,
				AuthToken:  wrongToken,
			},
			want: results{
				statusCode: http.StatusInternalServerError,
				response:   ``,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Repository: tt.fields.Repository,
			}

			r := SetUpRouter()
			r.GET("/api/user/balance", s.GetUserBalance)
			req, _ := http.NewRequest("GET", "/api/user/balance", nil)
			w := httptest.NewRecorder()

			req.Header.Set("Authorization", tt.fields.AuthToken)
			r.ServeHTTP(w, req)

			responseData, _ := io.ReadAll(w.Body)
			assert.Equal(t, tt.want.response, string(responseData))
			assert.Equal(t, tt.want.statusCode, w.Code)
		})
	}
}
