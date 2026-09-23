package timer

import (
	"fmt"
	"os"
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

func StartTask(name string, tag string, s *storage.Storage) {
	if s.ActiveTask != nil {
		fmt.Printf("Task '%s' is already running. Stop it first.\n", s.ActiveTask.Name)
		return
	}

	now := time.Now()
	s.ActiveTask = &storage.Task{
		Name:      name,
		Tag:       tag,
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
	fmt.Printf("Paused task '%s'. Current accumulated time: %s\n", s.ActiveTask.Name, formatDuration(s.Accumulated))
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
	fmt.Printf("Stopped task '%s'. Total Duration: %s\n", task.Name, formatDuration(duration))
}

func ListTasks(s *storage.Storage) {
	fmt.Println("Completed Tasks:")

	// Sort tasks by start time descending (most recent first)
	sort.Slice(s.CompletedTasks, func(i, j int) bool {
		return s.CompletedTasks[i].StartTime.After(s.CompletedTasks[j].StartTime)
	})

	for i, t := range s.CompletedTasks {
		tagStr := ""
		if t.Tag != "" {
			tagStr = fmt.Sprintf(" [%s]", t.Tag)
		}
		fmt.Printf("%d: %s%s: %s (Started: %s)\n", i, t.Name, tagStr, formatDuration(t.Duration), t.StartTime.Format("2006-01-02 15:04"))
	}

	fmt.Printf("\nTotal completed tasks: %d\n", len(s.CompletedTasks))

	if s.ActiveTask != nil {
		tagStr := ""
		if s.ActiveTask.Tag != "" {
			tagStr = fmt.Sprintf(" [%s]", s.ActiveTask.Tag)
		}
		fmt.Printf("Active Task: %s%s (Started: %s)\n", s.ActiveTask.Name, tagStr, s.ActiveTask.StartTime.Format("2006-01-02 15:04"))
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

func PrintReport(s *storage.Storage, filterTag string) {
	fmt.Print(generateReport(s, filterTag))
}

func generateReport(s *storage.Storage, filterTag string) string {
	totals := make(map[string]time.Duration)
	for _, t := range s.CompletedTasks {
		if filterTag == "" || t.Tag == filterTag {
			totals[t.Name] += t.Duration
		}
	}

	if s.ActiveTask != nil {
		if filterTag == "" || s.ActiveTask.Tag == filterTag {
			activeDuration := s.Accumulated
			if s.PausedAt == nil {
				activeDuration += time.Since(s.ActiveTask.StartTime)
			}
			totals[s.ActiveTask.Name] += activeDuration
		}
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

	report := "Time Report (by Task):\n"
	if filterTag != "" {
		report += fmt.Sprintf("Filter Tag: %s\n", filterTag)
	}
	report += "--------------------------\n"
	for _, tt := range sorted {
		report += fmt.Sprintf("%-20s %s\n", tt.name, formatDuration(tt.time))
	}
	if len(sorted) == 0 {
		report += "No tasks to report.\n"
	}
	return report
}

func ExportReport(s *storage.Storage, filename string, filterTag string) {
	report := generateReport(s, filterTag)
	err := os.WriteFile(filename, []byte(report), 0644)
	if err != nil {
		fmt.Printf("Error exporting report: %v\n", err)
		return
	}
	fmt.Printf("Report successfully exported to %s\n", filename)
}

func ClearTasks(s *storage.Storage) {
	s.ActiveTask = nil
	s.PausedAt = nil
	s.Accumulated = 0
	s.CompletedTasks = []storage.Task{}
	s.Save()
	fmt.Println("Task history cleared.")
}

func DeleteTask(index int, s *storage.Storage) {
	if index < 0 || index >= len(s.CompletedTasks) {
		fmt.Println("Invalid task index.")
		return
	}

	taskName := s.CompletedTasks[index].Name
	s.CompletedTasks = append(s.CompletedTasks[:index], s.CompletedTasks[index+1:]...)
	s.Save()
	fmt.Printf("Deleted task '%s' at index %d.\n", taskName, index)
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

func PrintTopTasks(s *storage.Storage) {
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

	fmt.Println("Top 5 Tasks by Duration:")
	fmt.Println("--------------------------")
	limit := 5
	if len(sorted) < limit {
		limit = len(sorted)
	}

	for i := 0; i < limit; i++ {
		fmt.Printf("%d. %-20s %s\n", i+1, sorted[i].name, formatDuration(sorted[i].time))
	}

	if len(sorted) == 0 {
		fmt.Println("No tasks recorded yet.")
	}
}