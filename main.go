package main

import (
	"fmt"
	"os"
	"github.com/atlas7649/go-task-timer/storage"
	"github.com/atlas7649/go-task-timer/timer"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: task-timer [start <name> [tag] | pause | resume | stop | list | summary [name] | report [tag] | export <filename> [tag] | clear | reset | status | delete <index> | top | search <query> | stats]")
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
	case "list":
		timer.ListTasks(store)
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
	case "clear":
		timer.ClearTasks(store)
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
	case "top":
		timer.PrintTopTasks(store)
	case "search":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a search query.")
			return
		}
		timer.SearchTasks(os.Args[2], store)
	case "stats":
		timer.PrintStats(store)
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
	}
}