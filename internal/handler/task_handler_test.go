package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"task-api/internal/models"
)

type mockTaskRepository struct {
	createFunc func(ctx context.Context, name string) (*models.Task, error)
}

func (m *mockTaskRepository) Create(ctx context.Context, name string) (*models.Task, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, name)
	}
	return nil, errors.New("not implemented")
}

func (m *mockTaskRepository) GetAll(ctx context.Context, params models.TaskQueryParams) ([]models.Task, error) {
	return nil, nil
}

func (m *mockTaskRepository) GetByID(ctx context.Context, id int) (*models.Task, error) {
	return nil, nil
}

func (m *mockTaskRepository) Update(ctx context.Context, id int, name string, completed bool) (*models.Task, error) {
	return nil, nil
}

func (m *mockTaskRepository) Delete(ctx context.Context, id int) error {
	return nil
}


func TestCreateTask(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		mockCreate     func(ctx context.Context, name string) (*models.Task, error)
		expectedStatus int
	}{
		{
			name:        "Success - Valid Task Creation",
			requestBody: `{"name": "Buy milk"}`,
			mockCreate: func(ctx context.Context, name string) (*models.Task, error) {
				return &models.Task{
					ID:        1,
					Name:      name,
					Completed: false,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Failure - Empty Task Name",
			requestBody:    `{"name": "   "}`,
			mockCreate:     nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Failure - Database Error",
			requestBody: `{"name": "Buy milk"}`,
			mockCreate: func(ctx context.Context, name string) (*models.Task, error) {
				return nil, errors.New("db connection failure")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Setup mock repository and handler
			mockRepo := &mockTaskRepository{createFunc: tt.mockCreate}
			handler := NewTaskHandler(mockRepo)

			//Create fake HTTP request and response recorder
			req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.CreateTask(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}