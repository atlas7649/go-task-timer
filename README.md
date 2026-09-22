# Go Task Timer

A simple CLI tool to track time spent on various tasks.

## Usage

- `go run main.go start "My Project"` - Start timing a task
- `go run main.go pause` - Pause the current task
- `go run main.go resume` - Resume a paused task
- `go run main.go stop` - Stop the current task
- `go run main.go list` - List all completed and active tasks
- `go run main.go status` - Show current task elapsed time
- `go run main.go summary [name]` - Show total time for a specific task
- `go run main.go report` - Show cumulative time per task
- `go run main.go clear` - Clear all task history