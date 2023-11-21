package server

import (
	"github.com/DavidGQK/go-loyalty-system/internal/mocks"
	"github.com/DavidGQK/go-loyalty-system/internal/services/mycrypto"
	"github.com/DavidGQK/go-loyalty-system/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServer_SignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	server := &Server{Repository: mockRepo}

	// Test case 1: Successful signup
	mockRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	mockRepo.EXPECT().UpdateUserToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	reqBody := `{"Login": "testuser", "Password": "testpassword"}`
	req := httptest.NewRequest("POST", "/signup", strings.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	server.SignUp(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"login":"success"}`, w.Body.String())

	// Test case 2: Error creating user (conflict)
	mockRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(store.ErrConflict)

	reqBody = `{"login": "testuser", "password": "testpassword"}`
	req = httptest.NewRequest("POST", "/signup", strings.NewReader(reqBody))
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = req

	server.SignUp(c)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), "")
}

func TestServer_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	server := &Server{Repository: mockRepo}

	// Test case 1: Successful login
	mockRepo.EXPECT().FindUserByLogin(gomock.Any(), gomock.Any()).Return(&store.User{Password: mycrypto.HashFunc("testpassword")}, nil)
	mockRepo.EXPECT().UpdateUserToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	reqBody := `{"login": "testuser", "password": "testpassword"}`
	req := httptest.NewRequest("POST", "/login", strings.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	server.Login(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"login":"success"}`, w.Body.String())

	// Test case 2: User not found
	mockRepo.EXPECT().FindUserByLogin(gomock.Any(), gomock.Any()).Return(&store.User{}, store.ErrNowRows)

	reqBody = `{"login": "nonexistentuser", "password": "testpassword"}`
	req = httptest.NewRequest("POST", "/login", strings.NewReader(reqBody))
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = req

	server.Login(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
