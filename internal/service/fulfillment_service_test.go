package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"x5_test/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockFulfillmentRepository struct {
	mock.Mock
}

func (m *MockFulfillmentRepository) UpdateOrderStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func TestProcessOrder_Success(t *testing.T) {
	mockRepo := new(MockFulfillmentRepository)
	svc := NewFulfillmentService(mockRepo)

	orderID := uuid.New()
	var capturedStatuses []domain.OrderStatus

	mockRepo.On("UpdateOrderStatus", mock.Anything, orderID, domain.StatusProcessing).
		Return(nil).
		Run(func(args mock.Arguments) {
			capturedStatuses = append(capturedStatuses, args.Get(2).(domain.OrderStatus))
		})

	mockRepo.On("UpdateOrderStatus", mock.Anything, orderID, mock.AnythingOfType("domain.OrderStatus")).
		Return(nil).
		Run(func(args mock.Arguments) {
			capturedStatuses = append(capturedStatuses, args.Get(2).(domain.OrderStatus))
		})

	err := svc.ProcessOrder(context.Background(), orderID.String())

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
	require.Len(t, capturedStatuses, 2)
	require.Equal(t, domain.StatusProcessing, capturedStatuses[0])
	require.Contains(t, []domain.OrderStatus{domain.StatusFulfilled, domain.StatusFailed}, capturedStatuses[1])
}

func TestProcessOrder_FirstUpdateFails(t *testing.T) {
	mockRepo := new(MockFulfillmentRepository)
	svc := NewFulfillmentService(mockRepo)

	orderID := uuid.New()
	expectedErr := assert.AnError

	mockRepo.On("UpdateOrderStatus", mock.Anything, orderID, domain.StatusProcessing).
		Return(expectedErr)

	err := svc.ProcessOrder(context.Background(), orderID.String())

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to update order status to processing")
	mockRepo.AssertExpectations(t)
}

func TestProcessOrder_SecondUpdateFails(t *testing.T) {
	mockRepo := new(MockFulfillmentRepository)
	svc := NewFulfillmentService(mockRepo)

	orderID := uuid.New()
	expectedErr := assert.AnError

	mockRepo.On("UpdateOrderStatus", mock.Anything, orderID, domain.StatusProcessing).
		Return(nil)

	mockRepo.On("UpdateOrderStatus", mock.Anything, orderID, mock.AnythingOfType("domain.OrderStatus")).
		Return(expectedErr)

	err := svc.ProcessOrder(context.Background(), orderID.String())

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to update final status")
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "UpdateOrderStatus", 2)
}

func TestProcessOrder_InvalidUUID(t *testing.T) {
	mockRepo := new(MockFulfillmentRepository)
	svc := NewFulfillmentService(mockRepo)

	err := svc.ProcessOrder(context.Background(), "not-a-uuid")

	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid UUID")
	mockRepo.AssertNotCalled(t, "UpdateOrderStatus")
}
