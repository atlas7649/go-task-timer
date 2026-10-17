package main

import (
	"fmt"
	"os"
	"strings"
	"time"
	"github.com/atlas7649/go-task-timer/storage"
	"github.com/atlas7649/go-task-timer/timer"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: task-timer <command> [arguments]")
		fmt.Println("Commands:")
		fmt.Println("  start <name> [tag]   - Start timing a task")
		fmt.Println("  pause               - Pause current task")
		fmt.Println("  resume              - Resume paused task")
		fmt.Println("  stop                - Stop current task")
		fmt.Println("  stop-all            - Stop all active timers")
		fmt.Println("  list                - List all tasks")
		fmt.Println("  log                 - Print chronological log")
		fmt.Println("  summary [name]      - Show total time for task")
		fmt.Println("  report [tag]         - Show cumulative time report")
		fmt.Println("  export <file> [tag] - Export report to file")
		fmt.Println("  export-csv <file>   - Export to CSV")
		fmt.Println("  export-json <file>  - Export full state to JSON")
		fmt.Println("  clear               - Clear all history")
		fmt.Println("  clear-goals         - Clear all goals")
		fmt.Println("  reset               - Reset active task without saving")
		fmt.Println("  status              - Show current task status")
		fmt.Println("  delete <index>      - Delete task by index")
		fmt.Println("  edit <index> <n> <t> - Edit task name and tag")
		fmt.Println("  top                 - Top tasks all time")
		fmt.Println("  top-tag-duration    - Top tags by duration")
		fmt.Println("  top-daily           - Top tasks today")
		fmt.Println("  top-week            - Top tasks this week")
		fmt.Println("  top-month           - Top tasks this month")
		fmt.Println("  top-yearly          - Top tasks this year")
		fmt.Println("  top-all             - Top tasks of all time")
		fmt.Println("  top-pomodoro        - Top Pomodoro tasks")
		fmt.Println("  search <query>       - Search tasks by name/tag")
		fmt.Println("  stats               - Overall statistics")
		fmt.Println("  tag <tag>           - List tasks by tag")
		fmt.Println("  tags                - List all unique tags")
		fmt.Println("  goal <name> <dur>   - Set a time goal")
		fmt.Println("  goals               - Show goal progress")
		fmt.Println("  rm-goal <name>      - Remove a goal")
		fmt.Println("  daily               - Daily report")
		fmt.Println("  weekly              - Weekly report")
		fmt.Println("  history <start> <end>- History range (YYYY-MM-DD)")
		fmt.Println("  archive             - Archive completed tasks")
		fmt.Println("  pomodoro [name] [dur]- Start Pomodoro session")
		fmt.Println("  last                - Show most recent task")
		os.Exit(1)
	}

	store := storage.NewJSONStorage("tasks.json")
	cmd := os.Args[1]

	switch cmd {
	case "start":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a task name.")
			return
		}
		name := os.Args[2]
		tag := ""
		if len(os.Args) >= 4 {
			tag = os.Args[3]
		}
		timer.StartTask(name, tag, store)
	case "pause":
		timer.PauseTask(store)
	case "resume":
		timer.ResumeTask(store)
	case "stop":
		timer.StopTask(store)
	case "stop-all":
		if store.ActiveTask == nil {
			fmt.Println("No active task to stop.")
		} else {
			timer.StopTask(store)
		}
	case "list":
		timer.ListTasks(store)
	case "log":
		timer.PrintLog(store)
	case "summary":
		var name string
		if len(os.Args) >= 3 {
			name = os.Args[2]
		}
		timer.PrintSummary(store, name)
	case "report":
		var tag string
		if len(os.Args) >= 3 {
			tag = os.Args[2]
		}
		timer.PrintReport(store, tag)
	case "export":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a filename to export to.")
			return
		}
		filename := os.Args[2]
		var tag string
		if len(os.Args) >= 4 {
			tag = os.Args[3]
		}
		timer.ExportReport(store, filename, tag)
	case "export-csv":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a filename to export to.")
			return
		}
		filename := os.Args[2]
		var tag string
		if len(os.Args) >= 4 {
			tag = os.Args[3]
		}
		timer.ExportCSV(store, filename, tag)
	case "export-json":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a filename to export to.")
			return
		}
		filename := os.Args[2]
		timer.ExportJSON(store, filename)
	case "clear":
		timer.ClearTasks(store)
	case "clear-goals":
		timer.ClearGoals(store)
	case "reset":
		timer.ResetActiveTask(store)
	case "status":
		timer.PrintStatus(store)
	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Please provide the task index to delete.")
			return
		}
		var index int
		_, err := fmt.Sscanf(os.Args[2], "%d", &index)
		if err != nil {
			fmt.Println("Invalid index provided.")
			return
		}
		timer.DeleteTask(index, store)
	case "edit":
		if len(os.Args) < 5 {
			fmt.Println("Usage: edit <index> <new_name> <new_tag>")
			return
		}
		var index int
		_, err := fmt.Sscanf(os.Args[2], "%d", &index)
		if err != nil {
			fmt.Println("Invalid index provided.")
			return
		}
		timer.EditTask(index, os.Args[3], os.Args[4], store)
	case "top":
		timer.PrintTopTasks(store)
	case "top-tag-duration":
		timer.PrintTopTags(store)
	case "top-daily":
		timer.PrintTopDailyTasks(store)
	case "top-daily-duration":
		timer.PrintTopDailyTasks(store)
	case "top-week":
		timer.PrintTopWeeklyTasks(store)
	case "top-week-duration":
		timer.PrintTopWeeklyTasks(store)
	case "top-month":
		timer.PrintTopMonthlyTasks(store)
	case "top-yearly":
		timer.PrintTopYearlyTasks(store)
	case "top-all":
		timer.PrintTopAllTasks(store)
	case "top-pomodoro":
		timer.PrintTopPomodoroTasks(store)
	case "top-pomodoro-count":
		timer.PrintTopPomodoroCount(store)
	case "search":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a search query.")
			return
		}
		timer.SearchTasks(os.Args[2], store)
	case "stats":
		timer.PrintStats(store)
	case "tag":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a tag to list tasks for.")
			return
		}
		timer.ListTasksByTag(os.Args[2], store)
	case "tags":
		timer.PrintTags(store)
	case "goal":
		if len(os.Args) < 4 {
			fmt.Println("Usage: goal <name> <duration> (e.g., goal 'Coding' 2h)")
			return
		}
		name := os.Args[2]
		durationStr := os.Args[3]
		dur, err := parseDuration(durationStr)
		if err != nil {
			fmt.Printf("Invalid duration: %v\n", err)
			return
		}
		timer.SetGoal(name, dur, store)
	case "goals":
		timer.PrintGoals(store)
	case "rm-goal":
		if len(os.Args) < 3 {
			fmt.Println("Please provide the task name to remove the goal for.")
			return
		}
		timer.RemoveGoal(os.Args[2], store)
	case "daily":
		timer.PrintDailyReport(store)
	case "weekly":
		timer.PrintWeeklyReport(store)
	case "history":
		if len(os.Args) < 4 {
			fmt.Println("Usage: history <start_date> <end_date> (YYYY-MM-DD)")
			return
		}
		timer.PrintHistory(os.Args[2], os.Args[3], store)
	case "archive":
		timer.ArchiveTasks(store)
	case "pomodoro":
		name := "Pomodoro"
		if len(os.Args) >= 3 {
			name = os.Args[2]
		}
		var duration time.Duration = 25 * time.Minute
		if len(os.Args) >= 4 {
			var err error
			duration, err = parseDuration(os.Args[3])
			if err != nil {
				fmt.Printf("Invalid pomodoro duration: %v. Using default 25m.\n", err)
				duration = 25 * time.Minute
			}
		}
		timer.StartPomodoro(name, duration, store)
	case "last":
		timer.PrintLastTask(store)
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
	}
}

func parseDuration(s string) (time.Duration, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	return time.ParseDuration(s)
}