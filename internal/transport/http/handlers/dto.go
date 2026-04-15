package handlers

import (
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type taskMutationDTO struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      taskdomain.Status `json:"status"`
	EndDate     string            `json:"end_date,omitempty"`
	RepeatType  string            `json:"repeat_Type,omitempty"`
}

type taskDTO struct {
	ID          int64             `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      taskdomain.Status `json:"status"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	EndDate     string            `json:"end_date"`
	RepeatType  string            `json:"repeat_type"`
}

func newTaskDTO(task *taskdomain.Task) taskDTO {
	taskDTO := taskDTO{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		RepeatType:  task.RepeatType,
	}

	if task.EndDate != nil {
		taskDTO.EndDate = task.EndDate.Format("20060102")
	}

	return taskDTO
}
