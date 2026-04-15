package task

import (
	"context"
	"testing"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type MockRepository struct {
	Repository
	OnCreate func(task *taskdomain.Task) (*taskdomain.Task, error)
	OnUpdate func(task *taskdomain.Task) (*taskdomain.Task, error)
}

func (m *MockRepository) Create(ctx context.Context, t *taskdomain.Task) (*taskdomain.Task, error) {
	return m.OnCreate(t)
}

func (m *MockRepository) Update(ctx context.Context, t *taskdomain.Task) (*taskdomain.Task, error) {
	return m.OnUpdate(t)
}

func TestService_Update_Repeating(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 4, 13, 12, 0, 0, 0, time.UTC)

	createdCount := 0

	repo := &MockRepository{
		OnUpdate: func(task *taskdomain.Task) (*taskdomain.Task, error) {
			return task, nil
		},
		OnCreate: func(task *taskdomain.Task) (*taskdomain.Task, error) {
			createdCount++
			return task, nil
		},
	}

	service := &Service{
		repo: repo,
		now:  func() time.Time { return now },
	}

	inputSimple := UpdateInput{Status: taskdomain.StatusDone, Title: "Test Task"}
	service.Update(ctx, 1, inputSimple)
	if createdCount != 0 {
		t.Errorf("Expected 0 new tasks, got %d", createdCount)
	}

	endDate := now
	inputRecurring := UpdateInput{
		Title:      "Another test task",
		Status:     taskdomain.StatusDone,
		RepeatType: "daily 1",
		EndDate:    &endDate,
	}

	service.Update(ctx, 2, inputRecurring)

	if createdCount != 1 {
		t.Errorf("Expected 1 new task to be created for repeating Type, got %d", createdCount)
	}
}
