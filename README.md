# cli-record

A simple CLI time tracking application built with Go.

## Overview

cli-record is a terminal-based time tracking application that helps you record and analyze time spent on various tasks. It's designed for users who spend a lot of time in the terminal and want to avoid switching contexts.

## Features

- **Record Time Entries**: Start and stop time tracking for tasks with tags
- **View Time Entries**: List and filter recorded time entries
- **Reports**: View time spent by day, week, month, or year with various grouping options
- **Data Portability**: Stores data locally in JSON format for easy sharing across machines

## Installation

```bash
# Clone the repository
git clone https://github.com/m-uesaka/cli-record.git
cd cli-record

# Install dependencies
task deps

# Build the application
task build

# Or install directly to $GOPATH/bin
task install
```

## Prerequisites

- Go 1.21 or later
- [Task](https://taskfile.dev/) (optional, for using task commands)

## Usage

```bash
# Start tracking a task
cli-record start

# Stop tracking
cli-record stop

# List time entries
cli-record list

# Show details of a specific entry
cli-record show <ID>

# View reports
cli-record view
```

## Development

### Available Tasks

Run `task` or `task --list` to see all available tasks:

- `task deps` - Install dependencies
- `task build` - Build the application
- `task run` - Run the application
- `task test` - Run tests
- `task test-coverage` - Run tests with coverage report
- `task fmt` - Format code
- `task clean` - Clean build artifacts

### Project Structure

```
.
├── cmd/                 # Command definitions (Cobra)
├── internal/
│   ├── models/         # Data models
│   └── storage/        # Data storage layer
├── docs/               # Documentation
├── main.go             # Application entry point
└── Taskfile.yml        # Task automation
```

## Data Storage

Time entries are stored in `~/.cli-record/timeentries.json` in JSON format.

## License

[Add your license here]
