package timer

import (
	"fmt"
	"time"
	"github.com/atlas7649/go-task-timer/storage"
)

func StartTask(name string, s *storage.Storage) {
	if s.ActiveTask != nil {
		fmt.Printf("Task '%s' is already running. Stop it first.\n", s.ActiveTask.Name)
		return
	}

	now := time.Now()
	s.ActiveTask = &storage.Task{
		Name:      name,
		StartTime: now,
	}
	s.Save()
	fmt.Printf("Started tracking task: %s at %s\n", name, now.Format(time.Kitchen))
}

func StopTask(s *storage.Storage) {
	if s.ActiveTask == nil {
		fmt.Println("No active task to stop.")
		return
	}

	now := time.Now()
	duration := now.Sub(s.ActiveTask.StartTime)
	task := *s.ActiveTask
	task.EndTime = now
	task.Duration = duration

	s.CompletedTasks = append(s.CompletedTasks, task)
	s.ActiveTask = nil
	s.Save()
	fmt.Printf("Stopped task '%s'. Duration: %v\n", task.Name, duration)
}

func ListTasks(s *storage.Storage) {
	fmt.Println("Completed Tasks:")
	for _, t := range s.CompletedTasks {
		fmt.Printf("- %s: %v (Started: %s)\n", t.Name, t.Duration, t.StartTime.Format("2006-01-02 15:04"))
	}
	if s.ActiveTask != nil {
		fmt.Printf("\nActive Task: %s (Started: %s)\n", s.ActiveTask.Name, s.ActiveTask.StartTime.Format("2006-01-02 15:04"))
	}
}