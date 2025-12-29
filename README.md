# cli-record

A simple and powerful CLI time tracking application built with Go.

## Overview

cli-record is a terminal-based time tracking application that helps you record and analyze time spent on various tasks. It's designed for users who spend a lot of time in the terminal and want to avoid switching contexts.

## Features

- ✨ **Easy Time Tracking**: Start and stop time tracking with simple commands
- 🏷️ **Tags Support**: Organize tasks with tags for better categorization
- 📊 **Powerful Reports**: View time spent with various grouping options (task, tag, day, week, month, year)
- ⏰ **Detailed Analytics**: Hourly, weekday, and monthly breakdowns
- 📁 **Export Options**: Export reports to CSV or JSON formats
- 💾 **Data Portability**: Stores data locally in JSON format for easy sharing across machines
- 🎨 **Interactive TUI**: Beautiful terminal user interface powered by Bubble Tea
- 🔍 **Flexible Filtering**: Filter entries by date range, tags, and task names

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/m-uesaka/cli-record.git
cd cli-record

# Build and install
go build -o cli-record .
sudo mv cli-record /usr/local/bin/

# Or install using go install
go install .
```

### Prerequisites

- Go 1.21 or later

## Quick Start

```bash
# Start tracking a task
cli-record start --task "Write documentation" --tags "docs,writing"

# Stop the current timer
cli-record stop

# List all time entries
cli-record list

# View report grouped by task
cli-record view --by task

# Export weekly report to CSV
cli-record view --by week --format csv --output report.csv
```

## Usage

### Starting a Time Entry

```bash
# With task name and tags
cli-record start --task "Code review" --tags "development,review"

# Interactive mode (will prompt for task name and tags)
cli-record start

# Using aliases
cli-record st --task "Meeting"
```

### Stopping a Time Entry

```bash
# Stop the current running entry
cli-record stop

# Using alias
cli-record sp
```

### Listing Entries

```bash
# List all entries
cli-record list

# Filter by date range
cli-record list --from "2025-01-01" --to "2025-01-31"

# Filter by tags
cli-record list --tags "development,review"

# Filter by task name (partial matching)
cli-record list --task "documentation"

# Group by task
cli-record list --group-by task

# Using alias
cli-record ls
```

### Showing Entry Details

```bash
# Show details of a specific entry by ID
cli-record show <entry-id>
```

### Viewing Reports

```bash
# Basic report grouped by task (default)
cli-record view

# Group by different dimensions
cli-record view --by tag
cli-record view --by day
cli-record view --by week
cli-record view --by month
cli-record view --by year

# Special breakdown views
cli-record view --view-by-hour       # Hourly breakdown (00:00-23:00)
cli-record view --view-by-weekday    # Weekday breakdown (Mon-Sun)
cli-record view --view-by-day-of-month  # Day-of-month breakdown (1-31)
cli-record view --view-by-month      # Monthly breakdown (Jan-Dec)

# Filter reports
cli-record view --from "2025-01-01" --tags "development"

# Export to different formats
cli-record view --by task --format csv
cli-record view --by week --format json
cli-record view --view-by-hour --format csv --output hourly-report.csv
```

## Command Reference

### Global Flags

- `--help`, `-h`: Show help for any command
- `--version`, `-v`: Show version information

### Commands

| Command | Aliases | Description |
|---------|---------|-------------|
| `start` | `st` | Start recording time for a new task |
| `stop` | `sp` | Stop the currently running time entry |
| `list` | `ls` | Display a list of recorded time entries |
| `show` | - | Display detailed information about a specific entry |
| `view` | - | View aggregated time reports in various formats |

For detailed command options, see [docs/options.md](docs/options.md).

## Examples

### Daily Workflow

```bash
# Morning: Start work
cli-record start --task "Email triage" --tags "communication"

# Switch tasks
cli-record stop
cli-record start --task "Feature implementation" --tags "development,feature-x"

# Lunch break (stop tracking)
cli-record stop

# Afternoon: Continue work
cli-record start --task "Code review" --tags "development,review"

# End of day
cli-record stop

# View today's work
cli-record list --from "$(date +%Y-%m-%d)"
```

### Weekly Review

```bash
# View this week's work grouped by task
cli-record view --by task --from "$(date -d 'last monday' +%Y-%m-%d)"

# See which days you worked most
cli-record view --view-by-weekday

# Export weekly report
cli-record view --by day --format csv --output "week-$(date +%Y-%W).csv"
```

### Monthly Analysis

```bash
# View monthly report by task
cli-record view --by task --from "2025-01-01" --to "2025-01-31"

# See hourly work patterns
cli-record view --view-by-hour --from "2025-01-01" --to "2025-01-31"

# Export monthly summary
cli-record view --by month --format json --output monthly-report.json
```

## Data Storage

Time entries are stored in `~/.cli-record/data.json` in JSON format.

### Data Structure

Each time entry contains:
- **ID**: Unique identifier (UUID)
- **Start Time**: When the task started
- **End Time**: When the task ended (null if still running)
- **Task Name**: Name of the task
- **Tags**: List of associated tags

### Backup and Restore

```bash
# Backup your data
cp ~/.cli-record/data.json ~/backups/cli-record-$(date +%Y%m%d).json

# Restore from backup
cp ~/backups/cli-record-20250129.json ~/.cli-record/data.json
```

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test ./... -cover

# Run integration tests
go test ./test/integration/...
```

### Project Structure

```
.
├── cmd/                    # Command definitions (Cobra)
│   ├── root.go            # Root command
│   ├── start.go           # Start command
│   ├── stop.go            # Stop command
│   ├── list.go            # List command
│   ├── show.go            # Show command
│   ├── view.go            # View command
│   └── errors.go          # Error handling
├── internal/
│   ├── models/            # Data models
│   │   └── timeentry.go   # TimeEntry model
│   ├── storage/           # Data storage layer
│   │   └── storage.go     # JSON storage implementation
│   └── tui/               # Terminal UI components
│       ├── form.go        # Input form component
│       ├── autocomplete.go # Autocomplete component
│       └── confirm.go     # Confirmation dialog
├── test/
│   └── integration/       # Integration tests
├── docs/                  # Documentation
│   ├── requirements.md    # Project requirements
│   ├── options.md         # Command options reference
│   └── tasks.md           # Implementation tasks
├── main.go                # Application entry point
├── go.mod                 # Go module definition
└── README.md              # This file
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details.

## Acknowledgments

- Built with [Cobra](https://cobra.dev/) for CLI framework
- Terminal UI powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- Inspired by [Toggl Track](https://toggl.com/) and other time-tracking tools
