# CLI Commands and Options Reference

This document provides a comprehensive reference for all available commands and options in the cli-record application.

## Table of Contents

- [Global Options](#global-options)
- [start](#start)
- [stop](#stop)
- [list](#list)
- [show](#show)
- [view](#view)

## Global Options

Options that can be used with any command:

- `--help`, `-h`: Display help information for any command
- `--version`, `-v`: Display version information

## start

Start recording time for a new task.

### Usage

```bash
cli-record start [flags]
```

### Description

Starts a new time entry. Before starting, the command will check if there's already a running entry. If one exists, you'll be prompted to stop it first. The command presents a TUI prompt to fill in task information.

### Flags

- `--task`, `-t` (string, optional): Task name for the time entry
  - If not provided, you'll be prompted to enter it when stopping the entry
  - Example: `cli-record start --task "Write documentation"`

- `--tags` (string, optional): Comma-separated list of tags
  - Tags help categorize and filter time entries
  - Supports autocomplete from existing tags
  - If a tag doesn't exist, it will be created automatically
  - Example: `cli-record start --task "Code review" --tags "development,review"`

### Examples

```bash
# Start with task name and tags
cli-record start --task "Code review" --tags "development,review"

# Start with only task name (no tags)
cli-record start --task "Meeting"

# Start with interactive prompt (will ask for task name and tags)
cli-record start

# Start with tags only (task name will be asked when stopping)
cli-record start --tags "development"
```

### Behavior

1. Checks for existing running entry
2. If running entry exists, displays error and asks to stop it first
3. Generates a unique ID (UUID) for the new entry
4. If task name or tags not provided via flags, shows TUI prompt
5. Creates and saves the new time entry
6. Displays success message with entry ID

### Aliases

- `st`

## stop

Stop the currently running time entry.

### Usage

```bash
cli-record stop
```

### Description

Stops the currently active time entry by setting its end time to the current timestamp. If no entry is running, displays an error. If the task name wasn't provided when starting, prompts the user to enter it now.

### Flags

This command has no additional flags.

### Examples

```bash
# Stop the current running entry
cli-record stop
```

### Behavior

1. Retrieves the currently running entry
2. If no running entry exists, displays error message
3. If task name is missing, prompts user to input it via TUI
4. Sets end time to current timestamp
5. Updates entry in storage
6. Displays summary:
   - Task name
   - Duration (in human-readable format)
   - Tags (if any)

### Aliases

- `sp`

## list

Display a list of recorded time entries.

### Usage

```bash
cli-record list [flags]
```

### Description

Lists time entries with optional filtering by date range, tags, and task names. Also displays total time spent with optional grouping by task or tags.

### Flags

#### Filtering Options

- `--from` (string, optional): Start date for filtering entries
  - Format: `YYYY-MM-DD` or `YYYY-MM-DD HH:MM:SS`
  - Example: `--from "2025-01-01"`

- `--to` (string, optional): End date for filtering entries
  - Format: `YYYY-MM-DD` or `YYYY-MM-DD HH:MM:SS`
  - Example: `--to "2025-01-31"`

- `--tags` (string, optional): Filter by tags (comma-separated for multiple)
  - Example: `--tags "development,review"`
  - Entries matching any of the specified tags will be included

- `--task` (string, optional): Filter by task name (partial matching supported)
  - Example: `--task "documentation"`
  - Case-insensitive partial matching

#### Display Options

- `--group-by` (string, optional): Group entries and show subtotals
  - Values: `task`, `tag`, `date`
  - Default: no grouping
  - Example: `--group-by task`

### Examples

```bash
# List all entries
cli-record list

# List entries from January 2025
cli-record list --from "2025-01-01" --to "2025-01-31"

# List entries with specific tags
cli-record list --tags "development,review"

# List entries for specific task
cli-record list --task "documentation"

# List entries grouped by task
cli-record list --group-by task

# Combined filtering
cli-record list --from "2025-01-01" --tags "development" --group-by task
```

### Output Format

The list displays:
- Entry ID
- Start time
- End time (or "Running" if still active)
- Duration
- Task name
- Tags

At the bottom:
- Total time spent
- Subtotals (if grouped)

### Aliases

- `ls`

## show

Display detailed information about a specific time entry.

### Usage

```bash
cli-record show <ID>
```

### Description

Shows comprehensive details of a single time entry identified by its unique ID.

### Arguments

- `ID` (required): The unique identifier of the time entry
  - Can be found using the `list` command
  - Example: `cli-record show 550e8400-e29b-41d4-a716-446655440000`

### Flags

This command has no additional flags.

### Examples

```bash
# Show details of a specific entry
cli-record show 550e8400-e29b-41d4-a716-446655440000
```

### Output Format

Displays:
- **ID**: Unique identifier
- **Task Name**: Name of the task
- **Start Time**: When the entry started (formatted)
- **End Time**: When the entry ended (formatted, or "Running" if active)
- **Duration**: Total time spent (human-readable format, e.g., "2h 30m 15s")
- **Tags**: List of associated tags

### Error Handling

- If the entry ID doesn't exist, displays an error message
- If the ID format is invalid, displays a validation error

## view

View aggregated time reports in various formats.

### Usage

```bash
cli-record view [flags]
```

### Description

Generates comprehensive reports of time spent on tasks with various grouping and filtering options. Supports multiple view types and export formats.

### Flags

#### Filtering Options

- `--from` (string, optional): Start date for filtering entries
  - Format: `YYYY-MM-DD` or `YYYY-MM-DD HH:MM:SS`
  - Example: `--from "2025-01-01"`

- `--to` (string, optional): End date for filtering entries
  - Format: `YYYY-MM-DD` or `YYYY-MM-DD HH:MM:SS`
  - Example: `--to "2025-01-31"`

- `--tags` (string, optional): Filter by tags (comma-separated)
  - Example: `--tags "development,review"`

- `--task` (string, optional): Filter by task name (partial matching)
  - Example: `--task "documentation"`

#### Grouping Options

- `--by` (string, optional): Primary grouping dimension
  - Values: `task`, `tag`, `day`, `week`, `month`, `year`
  - Default: `task`
  - Example: `--by month`

#### Special View Options

- `--view-by-hour` (boolean): Show hourly breakdown
  - Aggregates time spent in each hour of the day (0-23)
  - Shows how much time was spent between 9am-10am, 10am-11am, etc.
  - Example: `--view-by-hour`

- `--view-by-weekday` (boolean): Show weekday breakdown
  - Aggregates time spent on each day of the week (Mon-Sun)
  - Useful for identifying work patterns
  - Example: `--view-by-weekday`

- `--view-by-day-of-month` (boolean): Show day-of-month breakdown
  - Aggregates time spent on each day of the month (1-31)
  - Example: `--view-by-day-of-month`

- `--view-by-month` (boolean): Show monthly breakdown
  - Aggregates time spent in each month (Jan-Dec)
  - Example: `--view-by-month`

#### Export Options

- `--format` (string, optional): Output format
  - Values: `table` (default), `csv`, `json`
  - Example: `--format csv`

- `--output` (string, optional): Output file path
  - If not specified, prints to stdout
  - Example: `--output report.csv`

### Examples

```bash
# View total time by task (default)
cli-record view

# View time by tag
cli-record view --by tag

# View daily breakdown
cli-record view --by day

# View weekly summary
cli-record view --by week

# View monthly summary
cli-record view --by month

# View hourly breakdown (when do you work most?)
cli-record view --view-by-hour

# View weekday breakdown (which days are busiest?)
cli-record view --view-by-weekday

# View day-of-month breakdown
cli-record view --view-by-day-of-month

# View monthly breakdown for the year
cli-record view --view-by-month

# Filtered view: development tasks in January
cli-record view --from "2025-01-01" --to "2025-01-31" --tags "development"

# Export to CSV
cli-record view --by task --format csv --output report.csv

# Export hourly breakdown to CSV
cli-record view --view-by-hour --format csv --output hourly-report.csv

# Combined: weekly view of development tasks exported to CSV
cli-record view --by week --tags "development" --format csv --output weekly-dev.csv
```

### Output Format

#### Table Format (default)

Displays a formatted table with:
- Grouping dimension (e.g., task name, date, etc.)
- Duration for each group
- Total duration at the bottom
- Percentage of total time (optional)

#### CSV Format

Generates CSV with headers:
- First column: grouping dimension
- Second column: duration (in seconds or formatted)
- Suitable for importing into spreadsheet applications

#### JSON Format

Generates JSON array with structured data:
```json
[
  {
    "group": "Task Name",
    "duration": 7200,
    "formatted_duration": "2h 0m"
  }
]
```

### Special Views Output

#### Hourly Breakdown

Shows time spent in each hour slot:
```
Hour    Duration    Percentage
00-01   0h 0m       0%
...
09-10   2h 30m      15%
10-11   1h 45m      10%
...
Total   16h 30m     100%
```

#### Weekday Breakdown

Shows time spent on each day of the week:
```
Weekday     Duration    Percentage
Monday      8h 30m      21%
Tuesday     7h 15m      18%
...
Total       40h 0m      100%
```

#### Day-of-Month Breakdown

Shows time spent on each day across all months:
```
Day     Duration    Percentage
1       3h 30m      5%
2       4h 15m      6%
...
Total   70h 0m      100%
```

#### Monthly Breakdown

Shows time spent in each month:
```
Month       Duration    Percentage
January     160h 30m    25%
February    145h 15m    23%
...
Total       640h 0m     100%
```

## Notes

### Date Format

All date inputs support multiple formats:
- `YYYY-MM-DD`: e.g., `2025-01-15`
- `YYYY-MM-DD HH:MM:SS`: e.g., `2025-01-15 14:30:00`

### Duration Format

Durations are displayed in human-readable format:
- Hours: `Xh`
- Minutes: `Xm`
- Seconds: `Xs`
- Combined: `2h 30m 15s`

### Tag Autocomplete

When using the TUI for tag input:
- Start typing to see suggestions from existing tags
- Press `Tab` to autocomplete
- Press `Enter` to select
- New tags are automatically created

### Data Storage

All data is stored locally in JSON format at:
- Default location: `~/.cli-record/data.json`
- Data is portable and can be easily backed up or transferred

### Best Practices

1. **Use consistent tag names**: This makes filtering and reporting more effective
2. **Provide descriptive task names**: Helps when reviewing reports later
3. **Stop entries promptly**: For accurate time tracking
4. **Regular reviews**: Use `view` command weekly/monthly to analyze patterns
5. **Export reports**: Save monthly reports for record-keeping
