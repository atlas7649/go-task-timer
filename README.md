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
- `go run main.go goal "My Project" 2h` - Set a time goal for a task
- `go run main.go goals` - Show progress on all goals
- `go run main.go top-daily` - Show top tasks for today
- `go run main.go top-week` - Show top tasks for this week
- `go run main.go top-week-duration` - Show top tasks for this week by duration
- `go run main.go top-month` - Show top tasks for this month
- `go run main.go top-yearly` - Show top tasks for this year
- `go run main.go top-all` - Show top tasks of all time
- `go run main.go top-all-duration` - Show top tasks of all time by duration
- `go run main.go top-tag-duration` - Show top tags by duration