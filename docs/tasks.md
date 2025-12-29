# Implementation Tasks

This document outlines the detailed tasks for implementing the CLI time-tracking application based on the requirements specified in `requirements.md`.

## Phase 1: Core Data Models and Storage

### Task 1.1: Define Core Data Models
- [ ] Create `TimeEntry` struct in `internal/models/`
  - [ ] Add field: `ID` (unique identifier, use UUID)
  - [ ] Add field: `StartTime` (timestamp)
  - [ ] Add field: `EndTime` (timestamp, nullable for running entries)
  - [ ] Add field: `TaskName` (string)
  - [ ] Add field: `Tags` (list of strings)
  - [ ] Add method: `Duration()` to calculate duration from start and end time
  - [ ] Add method: `IsRunning()` to check if entry is currently running
- [ ] Create `Tag` struct if needed for tag management
- [ ] Add JSON serialization tags to all structs

### Task 1.2: Implement Storage Layer
- [ ] Create storage interface in `internal/storage/`
  - [ ] Define `Storage` interface with methods:
    - [ ] `SaveEntry(entry *models.TimeEntry) error`
    - [ ] `GetEntry(id string) (*models.TimeEntry, error)`
    - [ ] `ListEntries() ([]*models.TimeEntry, error)`
    - [ ] `UpdateEntry(entry *models.TimeEntry) error`
    - [ ] `GetRunningEntry() (*models.TimeEntry, error)`
    - [ ] `ListTags() ([]string, error)`
- [ ] Implement JSON file storage
  - [ ] Create `JSONStorage` struct implementing `Storage` interface
  - [ ] Implement file read/write operations with proper locking
  - [ ] Handle concurrent access safely
  - [ ] Define default storage location (e.g., `~/.cli-record/data.json`)
  - [ ] Create storage directory if it doesn't exist
- [ ] Add tests for storage layer
  - [ ] Test saving and retrieving entries
  - [ ] Test handling non-existent files
  - [ ] Test concurrent access scenarios

## Phase 2: CLI Commands - Time Recording

### Task 2.1: Implement `start` Command
- [ ] Create `start` command using Cobra in `cmd/`
  - [ ] Add command definition and help text
  - [ ] Add flags:
    - [ ] `--task` or `-t` for task name (optional)
    - [ ] `--tags` for comma-separated tags (optional)
- [ ] Implement TUI prompt for task details
  - [ ] Use Bubble Tea or similar library
  - [ ] Create task name input form (optional field)
  - [ ] Create tags input with autocomplete from existing tags
  - [ ] Handle submission and cancellation
- [ ] Implement start logic
  - [ ] Check if there's already a running entry
  - [ ] If running entry exists, prompt user to stop it first
  - [ ] Generate unique ID for new entry (UUID)
  - [ ] Create new TimeEntry with start time
  - [ ] Save entry to storage
  - [ ] Display success message with entry ID
- [ ] Add tests for start command
  - [ ] Test starting with task name and tags
  - [ ] Test starting without task name
  - [ ] Test error when entry already running

### Task 2.2: Implement `stop` Command
- [ ] Create `stop` command using Cobra in `cmd/`
  - [ ] Add command definition and help text
- [ ] Implement stop logic
  - [ ] Retrieve currently running entry
  - [ ] If no running entry, display error message
  - [ ] If task name is missing, prompt user to input it
  - [ ] Set end time to current timestamp
  - [ ] Update entry in storage
  - [ ] Display summary (task name, duration, tags)
- [ ] Add tests for stop command
  - [ ] Test stopping running entry
  - [ ] Test error when no entry is running
  - [ ] Test prompting for task name when missing

## Phase 3: CLI Commands - Viewing Entries

### Task 3.1: Implement `list` Command
- [ ] Create `list` command using Cobra in `cmd/`
  - [ ] Add command definition and help text
  - [ ] Add flags for filtering:
    - [ ] `--from` for start date
    - [ ] `--to` for end date
    - [ ] `--tags` for filtering by tags
    - [ ] `--task` for filtering by task name
- [ ] Implement filtering logic
  - [ ] Filter entries by date range
  - [ ] Filter entries by tags (support multiple tags)
  - [ ] Filter entries by task name (support partial matching)
- [ ] Implement display logic
  - [ ] Format entries in a table view
  - [ ] Show: ID, start time, end time, duration, task name, tags
  - [ ] Display total time spent at the bottom
  - [ ] Group by task name and show subtotals
- [ ] Add tests for list command
  - [ ] Test listing all entries
  - [ ] Test filtering by date range
  - [ ] Test filtering by tags
  - [ ] Test filtering by task name

### Task 3.2: Implement `show` Command
- [ ] Create `show <ID>` command using Cobra in `cmd/`
  - [ ] Add command definition and help text
  - [ ] Add argument validation for ID
- [ ] Implement show logic
  - [ ] Retrieve entry by ID from storage
  - [ ] Handle case when entry not found
  - [ ] Display detailed information:
    - [ ] Start time (formatted)
    - [ ] End time (formatted)
    - [ ] Duration (human-readable format)
    - [ ] Task name
    - [ ] Tags (formatted list)
- [ ] Add tests for show command
  - [ ] Test showing existing entry
  - [ ] Test error for non-existent ID

### Task 3.3: Implement `view` Command
- [ ] Create `view` command using Cobra in `cmd/`
  - [ ] Add command definition and help text
  - [ ] Add flags for filtering:
    - [ ] `--from` for start date
    - [ ] `--to` for end date
    - [ ] `--tags` for filtering by tags
    - [ ] `--task` for filtering by task name
  - [ ] Add flags for grouping/viewing:
    - [ ] `--by` for grouping (task, tag, day, week, month, year)
    - [ ] `--view-by-hour` for hourly breakdown
    - [ ] `--view-by-weekday` for weekday breakdown
    - [ ] `--view-by-day-of-month` for day-of-month breakdown
    - [ ] `--view-by-month` for monthly breakdown
  - [ ] Add flags for export:
    - [ ] `--format` for output format (table, csv)
    - [ ] `--output` for output file path
- [ ] Implement aggregation logic
  - [ ] Create helper functions for time aggregation
  - [ ] Group by task name
  - [ ] Group by tags
  - [ ] Group by day/week/month/year
  - [ ] Calculate totals for each group
- [ ] Implement hourly breakdown view
  - [ ] Aggregate time spent in each hour of the day (0-23)
  - [ ] Display in table or chart format
- [ ] Implement weekday breakdown view
  - [ ] Aggregate time spent on each day of week (Mon-Sun)
  - [ ] Display in table or chart format
- [ ] Implement day-of-month breakdown view
  - [ ] Aggregate time spent on each day of month (1-31)
  - [ ] Display in table format
- [ ] Implement monthly breakdown view
  - [ ] Aggregate time spent in each month (Jan-Dec)
  - [ ] Display in table format
- [ ] Implement export functionality
  - [ ] CSV export with proper headers
  - [ ] Support other formats if needed
  - [ ] Write to file or stdout
- [ ] Add tests for view command
  - [ ] Test grouping by task
  - [ ] Test grouping by tags
  - [ ] Test grouping by time periods
  - [ ] Test hourly/weekday/day-of-month/monthly breakdowns
  - [ ] Test CSV export

## Phase 4: User Experience Improvements

### Task 4.1: Enhance TUI Components
- [ ] Create reusable TUI components
  - [ ] Input form component
  - [ ] Autocomplete component for tags
  - [ ] Table display component
  - [ ] Confirmation dialog component
- [ ] Implement tag autocomplete
  - [ ] Load existing tags from storage
  - [ ] Filter tags based on user input
  - [ ] Allow creating new tags
- [ ] Add proper error handling and user feedback
  - [ ] Display errors in user-friendly format
  - [ ] Show success messages with relevant details
  - [ ] Add loading indicators for long operations

### Task 4.2: Improve Command-Line Interface
- [ ] Add helpful error messages
  - [ ] Validate command arguments and flags
  - [ ] Provide suggestions for common mistakes
- [ ] Add command aliases for convenience
  - [ ] `ls` for `list`
  - [ ] `st` for `start`
  - [ ] `sp` for `stop`
- [ ] Implement proper help text and examples
  - [ ] Add examples to each command's help
  - [ ] Document all flags and options
- [ ] Add shell completion support
  - [ ] Generate completion scripts for bash/zsh/fish

## Phase 5: Testing and Documentation

### Task 5.1: Comprehensive Testing
- [ ] Write unit tests for all packages
  - [ ] Models package tests
  - [ ] Storage package tests
  - [ ] Command handlers tests
  - [ ] Utility functions tests
- [ ] Write integration tests
  - [ ] Test full workflows (start -> stop -> list)
  - [ ] Test data persistence across commands
  - [ ] Test edge cases and error conditions
- [ ] Set up test coverage reporting
  - [ ] Configure coverage tools
  - [ ] Aim for >80% coverage

### Task 5.2: Documentation
- [ ] Update README.md
  - [ ] Add installation instructions
  - [ ] Add usage examples for all commands
  - [ ] Add screenshots or demos if possible
- [ ] Create user guide
  - [ ] Document common workflows
  - [ ] Provide tips and best practices
- [ ] Add inline code documentation
  - [ ] Document all public functions and types
  - [ ] Add package-level documentation
- [ ] Create contributing guide if open-sourcing

## Phase 6: Build and Deployment

### Task 6.1: Build Configuration
- [ ] Update Taskfile.yml with build tasks
  - [ ] Add `build` task to compile binary
  - [ ] Add `test` task to run all tests
  - [ ] Add `lint` task for code quality checks
  - [ ] Add `install` task to install binary locally
- [ ] Set up proper versioning
  - [ ] Use git tags for versioning
  - [ ] Embed version in binary

### Task 6.2: Distribution
- [ ] Create installation scripts
  - [ ] Install script for Unix-like systems
  - [ ] Install script for Windows (if needed)
- [ ] Set up CI/CD pipeline (optional)
  - [ ] Automated testing on commits
  - [ ] Automated builds for releases
  - [ ] Cross-platform binary builds

## Implementation Order

Follow this suggested order for implementation:

1. **Phase 1** (Foundation): Complete all tasks in Phase 1 to establish the data layer
2. **Phase 2** (Core Features): Implement start and stop commands
3. **Phase 3** (Viewing): Implement list, show, and basic view command
4. **Phase 3 (Advanced Views)**: Add advanced view options (hourly, weekday, etc.)
5. **Phase 4** (UX): Enhance user experience with better TUI components
6. **Phase 5** (Quality): Add comprehensive tests and documentation
7. **Phase 6** (Distribution): Set up build and deployment

## Notes

- After completing each major task, commit the changes with appropriate semantic commit messages
- Test each feature thoroughly before moving to the next task
- Keep the implementation simple and focused on the requirements
- Refactor as needed but avoid over-engineering
