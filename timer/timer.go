package timer

import (
	"fmt"
	"sort"
	"time"
	"github.com/atlas7649/go-task-timer/storage"
)

func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh %dm %ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm %ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}

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
	s.Accumulated = 0
	s.PausedAt = nil
	s.Save()
	fmt.Printf("Started tracking task: %s at %s\n", name, now.Format(time.Kitchen))
}

func PauseTask(s *storage.Storage) {
	if s.ActiveTask == nil {
		fmt.Println("No active task to pause.")
		return
	}
	if s.PausedAt != nil {
		fmt.Println("Task is already paused.")
		return
	}

	now := time.Now()
	s.Accumulated += now.Sub(s.ActiveTask.StartTime)
	s.PausedAt = &now
	s.Save()
	fmt.Printf("Paused task '%s'. Current accumulated time: %v\n", s.ActiveTask.Name, s.Accumulated)
}

func ResumeTask(s *storage.Storage) {
	if s.ActiveTask == nil {
		fmt.Println("No active task to resume.")
		return
	}
	if s.PausedAt == nil {
		fmt.Println("Task is not paused.")
		return
	}

	now := time.Now()
	s.ActiveTask.StartTime = now
	s.PausedAt = nil
	s.Save()
	fmt.Printf("Resumed task '%s' at %s\n", s.ActiveTask.Name, now.Format(time.Kitchen))
}

func StopTask(s *storage.Storage) {
	if s.ActiveTask == nil {
		fmt.Println("No active task to stop.")
		return
	}

	now := time.Now()
	duration := s.Accumulated
	if s.PausedAt == nil {
		duration += now.Sub(s.ActiveTask.StartTime)
	}

	task := *s.ActiveTask
	task.EndTime = now
	task.Duration = duration

	s.CompletedTasks = append(s.CompletedTasks, task)
	s.ActiveTask = nil
	s.PausedAt = nil
	s.Accumulated = 0
	s.Save()
	fmt.Printf("Stopped task '%s'. Total Duration: %v\n", task.Name, duration)
}

func ListTasks(s *storage.Storage) {
	fmt.Println("Completed Tasks:")
	for _, t := range s.CompletedTasks {
		fmt.Printf("- %s: %s (Started: %s)\n", t.Name, formatDuration(t.Duration), t.StartTime.Format("2006-01-02 15:04"))
	}
	if s.ActiveTask != nil {
		fmt.Printf("\nActive Task: %s (Started: %s)\n", s.ActiveTask.Name, s.ActiveTask.StartTime.Format("2006-01-02 15:04"))
	}
}

func PrintSummary(s *storage.Storage, filterName string) {
	var total time.Duration
	var count int

	for _, t := range s.CompletedTasks {
		if filterName == "" || t.Name == filterName {
			total += t.Duration
			count++
		}
	}

	if filterName != "" {
		fmt.Printf("Total time spent on task '%s': %s (across %d sessions)\n", filterName, formatDuration(total), count)
	} else {
		fmt.Printf("Total time spent across %d completed tasks: %s\n", len(s.CompletedTasks), formatDuration(total))
	}

	if s.ActiveTask != nil {
		if filterName == "" || s.ActiveTask.Name == filterName {
			fmt.Printf("Currently active task: %s (since %s)\n", s.ActiveTask.Name, s.ActiveTask.StartTime.Format("2006-01-02 15:04"))
		}
	}
}

func PrintReport(s *storage.Storage) {
	totals := make(map[string]time.Duration)
	for _, t := range s.CompletedTasks {
		totals[t.Name] += t.Duration
	}

	if s.ActiveTask != nil {
		activeDuration := s.Accumulated
		if s.PausedAt == nil {
			activeDuration += time.Since(s.ActiveTask.StartTime)
		}
		totals[s.ActiveTask.Name] += activeDuration
	}

	type taskTime struct {
		name string
		time time.Duration
	}

	var sorted []taskTime
	for name, duration := range totals {
		sorted = append(sorted, taskTime{name, duration})
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].time > sorted[j].time
	})

	fmt.Println("Time Report (by Task):")
	fmt.Println("--------------------------")
	for _, tt := range sorted {
		fmt.Printf("%-20s %s\n", tt.name, formatDuration(tt.time))
	}
	if len(sorted) == 0 {
		fmt.Println("No tasks to report.")
	}
}

func ClearTasks(s *storage.Storage) {
	s.ActiveTask = nil
	s.PausedAt = nil
	s.Accumulated = 0
	s.CompletedTasks = []storage.Task{}
	s.Save()
	fmt.Println("Task history cleared.")
}

func PrintStatus(s *storage.Storage) {
	if s.ActiveTask == nil {
		fmt.Println("No task is currently running.")
		return
	}

	elapsed := s.Accumulated
	if s.PausedAt == nil {
		elapsed += time.Since(s.ActiveTask.StartTime)
	}

	status := "Running"
	if s.PausedAt != nil {
		status = "Paused"
	}

	fmt.Printf("Current task: %s [%s]\nStarted: %s\nElapsed time: %s\n", 
		s.ActiveTask.Name, 
		status,
		s.ActiveTask.StartTime.Format("2006-01-02 15:04"), 
		formatDuration(elapsed))
}