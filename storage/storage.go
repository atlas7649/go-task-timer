package storage

import (
	"encoding/json"
	"os"
	"time"
)

type Task struct {
	Name      string        `json:"name"`
	Tag       string        `json:"tag,omitempty"`
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time,omitempty"`
	Duration  time.Duration `json:"duration,omitempty"`
}

type Goal struct {
	TaskName   string        `json:"task_name"`
	TargetTime time.Duration `json:"target_time"`
}

type Storage struct {
	FilePath       string    `json:"-"`
	ActiveTask     *Task     `json:"active_task"`
	PausedAt       *time.Time `json:"paused_at,omitempty"`
	Accumulated    time.Duration `json:"accumulated,omitempty"`
	CompletedTasks []Task    `json:"completed_tasks"`
	Goals          []Goal    `json:"goals,omitempty"`
}

func NewJSONStorage(path string) *Storage {
	s := &Storage{FilePath: path, CompletedTasks: []Task{}, Goals: []Goal{}}

	data, err := os.ReadFile(path)
	if err == nil {
		json.Unmarshal(data, s)
	}
	return s
}

func (s *Storage) Save() {
	data, _ := json.MarshalIndent(s, "", "  ")
	os.WriteFile(s.FilePath, data, 0644)
}