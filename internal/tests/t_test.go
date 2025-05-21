package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Vladroon22/TestTask-ITK-Academy/internal/entity"
	"github.com/Vladroon22/TestTask-ITK-Academy/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type UserRepository interface {
	WalletOperation(c context.Context, wallet entity.WalletData) error
	GetBalance(c context.Context, uuid string) (entity.WalletData, error)
}

// MockUserRepository for testing
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) WalletOperation(c context.Context, wallet entity.WalletData) error {
	args := m.Called(c, wallet)
	return args.Error(0)
}

func (m *MockUserRepository) GetBalance(c context.Context, uuid string) (entity.WalletData, error) {
	args := m.Called(c, uuid)
	return args.Get(0).(entity.WalletData), args.Error(1)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (us *UserService) GetBalance(c context.Context, uuid string) (entity.WalletData, error) {
	return us.repo.GetBalance(c, uuid)
}

func (us *UserService) WalletOperation(c context.Context, wallet entity.WalletData) error {
	return us.repo.WalletOperation(c, wallet)
}
func GetBalanceHandler(s *UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"
		user, err := s.repo.GetBalance(context.Background(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(user)
	}
}

func WalletOpHandler(s *UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var wallet entity.WalletData
		if err := json.NewDecoder(r.Body).Decode(&w); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if err := s.repo.WalletOperation(context.Background(), wallet); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// Tests

func TestWalletOperation_Success(t *testing.T) {
	mockService := new(MockUserRepository)
	h := handlers.NewHandler(mockService)

	uuid := "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"
	validWallet := entity.WalletData{
		Uuid:           uuid,
		Balance:        100.0,
		Operation_type: "deposit",
	}

	mockService.On("WalletOperation", mock.Anything, validWallet).Return(nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = httptest.NewRequest("POST", "/api/v1/wallet", bytes.NewBufferString(`{"uuid":"a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11","balance":100.00,"type":"deposit"}`))
	ctx.Request.Header.Set("Content-Type", "application/json")

	h.WalletOperation(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Your finacial operation is successfull")
	mockService.AssertExpectations(t)
}

func TestWalletOperation_ValidationError(t *testing.T) {
	mockService := new(MockUserRepository)
	h := handlers.NewHandler(mockService)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest("POST", "/", bytes.NewBufferString(`{"uuid":"test-uuid","balance":-100.0}`))
	ctx.Request.Header.Set("Content-Type", "application/json")

	h.WalletOperation(ctx)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetBalance_ServiceError(t *testing.T) {
	mockService := new(MockUserRepository)
	h := handlers.NewHandler(mockService)

	UUID := "uuid"
	expectedErr := errors.New("service error")
	mockService.On("GetBalance", mock.Anything, UUID).Return(entity.WalletData{}, expectedErr)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest("GET", "/"+UUID, nil)
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: UUID}}

	h.GetBalance(ctx)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), expectedErr.Error())
	mockService.AssertExpectations(t)
}

func TestGetBalance(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := (mockRepo)

	t.Run("GetBalanceSuccess", func(t *testing.T) {
		expectedWallet := entity.WalletData{
			Uuid:    "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
			Balance: 150.75,
		}
		mockRepo.On("GetBalance", mock.Anything, expectedWallet.Uuid).Return(expectedWallet, nil)

		wallet, err := service.GetBalance(context.Background(), expectedWallet.Uuid)
		assert.NoError(t, err)
		assert.Equal(t, expectedWallet, wallet)
		mockRepo.AssertExpectations(t)
	})

	t.Run("WalletOperationSuccess", func(t *testing.T) {
		wallet := entity.WalletData{
			Uuid:    "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
			Balance: 50.00,
		}
		mockRepo.On("WalletOperation", mock.Anything, wallet).Return(nil)

		err := service.WalletOperation(context.Background(), wallet)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("WalletOperationError", func(t *testing.T) {
		wallet := entity.WalletData{
			Uuid:    "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
			Balance: -100.00,
		}
		mockRepo.On("WalletOperation", mock.Anything, wallet).Return(errors.New("insufficient funds"))

		err := service.WalletOperation(context.Background(), wallet)
		assert.Error(t, err)
		assert.Equal(t, "insufficient funds", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
