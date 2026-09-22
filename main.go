package main

import (
	"fmt"
	"os"
	"github.com/atlas7649/go-task-timer/storage"
	"github.com/atlas7649/go-task-timer/timer"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: task-timer [start <name> | stop | list]")
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
		timer.StartTask(name, store)
	case "stop":
		timer.StopTask(store)
	case "list":
		timer.ListTasks(store)
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
	}
}