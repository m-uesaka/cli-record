# Implementation Tasks

This document outlines the detailed tasks for implementing the CLI time-tracking application based on the requirements specified in `requirements.md`.

## Phase 1: Core Data Models and Storage

### Task 1.1: Define Core Data Models

- [x] Create `TimeEntry` struct in `internal/models/`
  - [x] Add field: `ID` (unique identifier, use UUID)
  - [x] Add field: `StartTime` (timestamp)
  - [x] Add field: `EndTime` (timestamp, nullable for running entries)
  - [x] Add field: `TaskName` (string)
  - [x] Add field: `Tags` (list of strings)
  - [x] Add method: `Duration()` to calculate duration from start and end time
  - [x] Add method: `IsRunning()` to check if entry is currently running
- [x] Create `Tag` struct if needed for tag management
- [x] Add JSON serialization tags to all structs

### Task 1.2: Implement Storage Layer

- [x] Create storage interface in `internal/storage/`
  - [x] Define `Storage` interface with methods:
    - [x] `SaveEntry(entry *models.TimeEntry) error`
    - [x] `GetEntry(id string) (*models.TimeEntry, error)`
    - [x] `ListEntries() ([]*models.TimeEntry, error)`
    - [x] `UpdateEntry(entry *models.TimeEntry) error`
    - [x] `GetRunningEntry() (*models.TimeEntry, error)`
    - [x] `ListTags() ([]string, error)`
- [x] Implement JSON file storage
  - [x] Create `JSONStorage` struct implementing `Storage` interface
  - [x] Implement file read/write operations with proper locking
  - [x] Handle concurrent access safely
  - [x] Define default storage location (e.g., `~/.cli-record/data.json`)
  - [x] Create storage directory if it doesn't exist
- [x] Add tests for storage layer
  - [x] Test saving and retrieving entries
  - [x] Test handling non-existent files
  - [x] Test concurrent access scenarios

## Phase 2: CLI Commands - Time Recording

### Task 2.1: Implement `start` Command

- [x] Create `start` command using Cobra in `cmd/`
  - [x] Add command definition and help text
  - [x] Add flags:
    - [x] `--task` or `-t` for task name (optional)
    - [x] `--tags` for comma-separated tags (optional)
- [x] Implement TUI prompt for task details
  - [x] Use Bubble Tea or similar library
  - [x] Create task name input form (optional field)
  - [x] Create tags input with autocomplete from existing tags
  - [x] Handle submission and cancellation
- [x] Implement start logic
  - [x] Check if there's already a running entry
  - [x] If running entry exists, prompt user to stop it first
  - [x] Generate unique ID for new entry (UUID)
  - [x] Create new TimeEntry with start time
  - [x] Save entry to storage
  - [x] Display success message with entry ID
- [x] Add tests for start command
  - [x] Test starting with task name and tags
  - [x] Test starting without task name
  - [x] Test error when entry already running

### Task 2.2: Implement `stop` Command

- [x] Create `stop` command using Cobra in `cmd/`
  - [x] Add command definition and help text
- [x] Implement stop logic
  - [x] Retrieve currently running entry
  - [x] If no running entry, display error message
  - [x] If task name is missing, prompt user to input it
  - [x] Set end time to current timestamp
  - [x] Update entry in storage
  - [x] Display summary (task name, duration, tags)
- [x] Add tests for stop command
  - [x] Test stopping running entry
  - [x] Test error when no entry is running
  - [x] Test prompting for task name when missing

## Phase 3: CLI Commands - Viewing Entries

### Task 3.1: Implement `list` Command

- [x] Create `list` command using Cobra in `cmd/`
  - [x] Add command definition and help text
  - [x] Add flags for filtering:
    - [x] `--from` for start date
    - [x] `--to` for end date
    - [x] `--tags` for filtering by tags
    - [x] `--task` for filtering by task name
- [x] Implement filtering logic
  - [x] Filter entries by date range
  - [x] Filter entries by tags (support multiple tags)
  - [x] Filter entries by task name (support partial matching)
- [x] Implement display logic
  - [x] Format entries in a table view
  - [x] Show: ID, start time, end time, duration, task name, tags
  - [x] Display total time spent at the bottom
  - [x] Group by task name and show subtotals
- [x] Add tests for list command
  - [x] Test listing all entries
  - [x] Test filtering by date range
  - [x] Test filtering by tags
  - [x] Test filtering by task name

### Task 3.2: Implement `show` Command

- [x] Create `show <ID>` command using Cobra in `cmd/`
  - [x] Add command definition and help text
  - [x] Add argument validation for ID
- [x] Implement show logic
  - [x] Retrieve entry by ID from storage
  - [x] Handle case when entry not found
  - [x] Display detailed information:
    - [x] Start time (formatted)
    - [x] End time (formatted)
    - [x] Duration (human-readable format)
    - [x] Task name
    - [x] Tags (formatted list)
- [x] Add tests for show command
  - [x] Test showing existing entry
  - [x] Test error for non-existent ID

### Task 3.3: Implement `view` Command

- [x] Create `view` command using Cobra in `cmd/`
  - [x] Add command definition and help text
  - [x] Add flags for filtering:
    - [x] `--from` for start date
    - [x] `--to` for end date
    - [x] `--tags` for filtering by tags
    - [x] `--task` for filtering by task name
  - [x] Add flags for grouping/viewing:
    - [x] `--by` for grouping (task, tag, day, week, month, year)
    - [x] `--view-by-hour` for hourly breakdown
    - [x] `--view-by-weekday` for weekday breakdown
    - [x] `--view-by-day-of-month` for day-of-month breakdown
    - [x] `--view-by-month` for monthly breakdown
  - [x] Add flags for export:
    - [x] `--format` for output format (table, csv)
    - [x] `--output` for output file path
- [x] Implement aggregation logic
  - [x] Create helper functions for time aggregation
  - [x] Group by task name
  - [x] Group by tags
  - [x] Group by day/week/month/year
  - [x] Calculate totals for each group
- [x] Implement hourly breakdown view
  - [x] Aggregate time spent in each hour of the day (0-23)
  - [x] Display in table or chart format
- [x] Implement weekday breakdown view
  - [x] Aggregate time spent on each day of week (Mon-Sun)
  - [x] Display in table or chart format
- [x] Implement day-of-month breakdown view
  - [x] Aggregate time spent on each day of month (1-31)
  - [x] Display in table format
- [x] Implement monthly breakdown view
  - [x] Aggregate time spent in each month (Jan-Dec)
  - [x] Display in table format
- [x] Implement export functionality
  - [x] CSV export with proper headers
  - [x] Support other formats if needed
  - [x] Write to file or stdout
- [x] Add tests for view command
  - [x] Test grouping by task
  - [x] Test grouping by tags
  - [x] Test grouping by time periods
  - [x] Test hourly/weekday/day-of-month/monthly breakdowns
  - [x] Test CSV export

## Phase 4: User Experience Improvements

### Task 4.1: Enhance TUI Components

- [x] Create reusable TUI components
  - [x] Input form component
  - [x] Autocomplete component for tags
  - [x] Table display component
  - [x] Confirmation dialog component
- [x] Implement tag autocomplete
  - [x] Load existing tags from storage
  - [x] Filter tags based on user input
  - [x] Allow creating new tags
- [x] Add proper error handling and user feedback
  - [x] Display errors in user-friendly format
  - [x] Show success messages with relevant details
  - [x] Add loading indicators for long operations

### Task 4.2: Improve Command-Line Interface

- [x] Add helpful error messages
  - [x] Validate command arguments and flags
  - [x] Provide suggestions for common mistakes
- [x] Add command aliases for convenience
  - [x] `ls` for `list`
  - [x] `st` for `start`
  - [x] `sp` for `stop`
- [x] Implement proper help text and examples
  - [x] Add examples to each command's help
  - [x] Document all flags and options
- [x] Add shell completion support
  - [x] Generate completion scripts for bash/zsh/fish (Cobra built-in)

## Phase 5: Testing and Documentation

### Task 5.1: Comprehensive Testing

- [x] Write unit tests for all packages
  - [x] Models package tests
  - [x] Storage package tests
  - [x] Command handlers tests
  - [x] Utility functions tests
- [x] Write integration tests
  - [x] Test full workflows (start -> stop -> list)
  - [x] Test data persistence across commands
  - [x] Test edge cases and error conditions
- [x] Set up test coverage reporting
  - [x] Configure coverage tools
  - [x] Aim for >80% coverage (achieved for models and storage)

### Task 5.2: Documentation

- [x] Update README.md
  - [x] Add installation instructions
  - [x] Add usage examples for all commands
  - [x] Add screenshots or demos if possible
- [x] Create user guide
  - [x] Document common workflows
  - [x] Provide tips and best practices
- [x] Add inline code documentation
  - [x] Document all public functions and types
  - [x] Add package-level documentation
- [x] Create contributing guide if open-sourcing

## Phase 6: Additional Commands (New Requirements)

### Task 6.1: Implement `status` Command

- [x] Create `status` command using Cobra in `cmd/`
  - [x] Add command definition and help text
- [x] Implement status logic
  - [x] Check for running entry
  - [x] If running, display task name, start time, elapsed time, and tags
  - [x] Calculate and format elapsed time (live duration)
  - [x] If not running, display appropriate message
  - [x] Format output in a user-friendly way
- [x] Add tests for status command
  - [x] Test with running entry
  - [x] Test with no running entry
  - [x] Test elapsed time calculation

### Task 6.2: Implement `edit` Command

- [x] Create `edit <ID>` command using Cobra in `cmd/`
  - [x] Add command definition and help text
  - [x] Add argument validation for ID
  - [x] Add flags:
    - [x] `--task` or `-t` for task name
    - [x] `--tags` for tags
    - [x] `--start` for start time
    - [x] `--end` for end time
    - [x] `--interactive` or `-i` for interactive mode (default)
- [x] Implement edit logic
  - [x] Retrieve entry by ID
  - [x] Handle case when entry not found
  - [x] Validate end time is after start time
  - [x] If using interactive mode, show TUI form with current values
  - [x] Update entry in storage
  - [x] Display success message
- [x] Create interactive edit form
  - [x] Pre-fill form with current values
  - [x] Allow editing all fields (start time, end time, task name, tags)
  - [x] Validate inputs before saving
  - [x] Handle cancellation
- [x] Add tests for edit command
  - [x] Test editing existing entry
  - [x] Test error for non-existent ID
  - [x] Test validation (end time after start time)
  - [x] Test updating individual fields
  - [x] Test interactive mode

### Task 6.3: Implement `remove` Command

- [x] Create `remove <ID>` command using Cobra in `cmd/`
  - [x] Add command definition and help text
  - [x] Add argument validation for ID
  - [x] Add flags:
    - [x] `--force` or `-f` to skip confirmation
  - [x] Add alias: `rm`
- [x] Implement remove logic
  - [x] Retrieve entry by ID
  - [x] Handle case when entry not found
  - [x] If not using --force, show confirmation prompt
  - [x] Display entry details in confirmation
  - [x] Delete entry from storage
  - [x] Display success message
- [x] Create confirmation dialog
  - [x] Use reusable TUI confirm component
  - [x] Display entry details (task, duration, tags)
  - [x] Warn that action cannot be undone
  - [x] Handle yes/no/cancel
- [x] Implement deletion in storage layer
  - [x] Add `DeleteEntry(id string) error` method to Storage interface
  - [x] Implement deletion in JSONStorage
  - [x] Handle concurrent access safely
- [x] Add tests for remove command
  - [x] Test removing existing entry
  - [x] Test error for non-existent ID
  - [x] Test confirmation prompt
  - [x] Test --force flag
  - [x] Test storage deletion

## Phase 7: Data Management and Configuration

### Task 7.1: Implement Configuration System

- [ ] Create configuration model in `internal/config/`
  - [ ] Define `Config` struct with fields:
    - [ ] `DataFilePath` (string) - custom data file location
    - [ ] `TimeFormat` (string) - 12H or 24H format preference
    - [ ] `ArchiveDirectory` (string) - archive storage location
  - [ ] Add default configuration values
  - [ ] Add TOML serialization tags
- [ ] Implement configuration storage
  - [ ] Store config in `$HOME/.config/cli-record/config.toml`
  - [ ] Create config directory if it doesn't exist
  - [ ] Load configuration on application start
  - [ ] Provide fallback to defaults if config file doesn't exist
- [ ] Add tests for configuration
  - [ ] Test loading configuration
  - [ ] Test saving configuration
  - [ ] Test default values
  - [ ] Test configuration validation

### Task 7.2: Implement `config` Command

- [ ] Create `config` command using Cobra in `cmd/`
  - [ ] Add command definition and help text
  - [ ] Add subcommands:
    - [ ] `config show` - display current configuration
    - [ ] `config set <key> <value>` - set configuration value
    - [ ] `config reset` - reset to default configuration
- [ ] Implement config show
  - [ ] Display all configuration values in a formatted way
  - [ ] Show which values are using defaults vs custom
- [ ] Implement config set
  - [ ] Validate configuration keys
  - [ ] Validate configuration values
  - [ ] Update configuration file
  - [ ] Provide helpful error messages for invalid inputs
  - [ ] Support keys:
    - [ ] `data-path` - set data file location
    - [ ] `time-format` - set time format (12h/24h)
    - [ ] `archive-dir` - set archive directory location
- [ ] Implement config reset
  - [ ] Prompt for confirmation
  - [ ] Reset to default configuration
  - [ ] Display success message
- [ ] Add tests for config command
  - [ ] Test showing configuration
  - [ ] Test setting valid values
  - [ ] Test setting invalid values
  - [ ] Test resetting configuration

### Task 7.3: Implement `archive` Command

- [ ] Create `archive` command using Cobra in `cmd/`
  - [ ] Add command definition and help text
  - [ ] Add flags:
    - [ ] `--output` or `-o` for custom archive file name
    - [ ] `--list` or `-l` to list available archives
- [ ] Implement archive logic
  - [ ] Generate timestamp for archive file name (YYYY-MM-DD-HHMMSS)
  - [ ] Copy current data file to archive directory
  - [ ] Create archive directory if it doesn't exist
  - [ ] Handle errors (e.g., file already exists, permission denied)
  - [ ] Display success message with archive file path
- [ ] Implement list archives functionality
  - [ ] Scan archive directory for archived files
  - [ ] Display list with file name, date, and size
  - [ ] Sort by date (newest first)
  - [ ] Handle case when no archives exist
- [ ] Add archive management methods to storage layer
  - [ ] Add `ArchiveData(archivePath string) error` method
  - [ ] Implement file copying with proper error handling
- [ ] Add tests for archive command
  - [ ] Test creating archive
  - [ ] Test listing archives
  - [ ] Test custom archive name
  - [ ] Test error handling

### Task 7.4: Implement `restore` Command

- [ ] Create `restore <archive-file>` command using Cobra in `cmd/`
  - [ ] Add command definition and help text
  - [ ] Add argument validation for archive file path
  - [ ] Add flags:
    - [ ] `--merge` (default) - merge with existing data
    - [ ] `--replace` - replace existing data
    - [ ] `--preview` - show what would be restored without applying
- [ ] Implement restore logic
  - [ ] Validate archive file exists
  - [ ] Load archived data
  - [ ] If merging, combine with existing data
    - [ ] Avoid duplicates based on unique IDs
    - [ ] Handle ID conflicts (keep newer entry or prompt user)
  - [ ] If replacing, backup current data first
  - [ ] Save restored data
  - [ ] Display summary of restored entries
- [ ] Implement merge algorithm
  - [ ] Load both current and archived data
  - [ ] Create map of existing entry IDs
  - [ ] Add archived entries that don't conflict
  - [ ] For conflicts, implement resolution strategy:
    - [ ] Keep newer based on modification time (if available)
    - [ ] Or prompt user to choose
- [ ] Add restore methods to storage layer
  - [ ] Add `RestoreData(archivePath string, merge bool) error` method
  - [ ] Implement data merging logic
  - [ ] Handle backup before replace
- [ ] Add tests for restore command
  - [ ] Test restoring with merge
  - [ ] Test restoring with replace
  - [ ] Test duplicate handling
  - [ ] Test preview mode
  - [ ] Test error handling (invalid file, corrupted data)

## Phase 8: Build and Deployment

### Task 8.1: Build Configuration

- [ ] Update Taskfile.yml with build tasks
  - [ ] Add `build` task to compile binary
  - [ ] Add `test` task to run all tests
  - [ ] Add `lint` task for code quality checks
  - [ ] Add `install` task to install binary locally
- [ ] Set up proper versioning
  - [ ] Use git tags for versioning
  - [ ] Embed version in binary

### Task 8.2: Distribution

- [ ] Create installation scripts
  - [ ] Install script for Unix-like systems
  - [ ] Install script for Windows (if needed)
- [ ] Set up CI/CD pipeline (optional)
  - [ ] Automated testing on commits
  - [ ] Automated builds for releases
  - [ ] Cross-platform binary builds

## Implementation Order

Follow this suggested order for implementation:

1. **Phase 1** (Foundation): Complete all tasks in Phase 1 to establish the data layer ✅
2. **Phase 2** (Core Features): Implement start and stop commands ✅
3. **Phase 3** (Viewing): Implement list, show, and basic view command ✅
4. **Phase 3 (Advanced Views)**: Add advanced view options (hourly, weekday, etc.) ✅
5. **Phase 4** (UX): Enhance user experience with better TUI components ✅
6. **Phase 5** (Quality): Add comprehensive tests and documentation ✅
7. **Phase 6** (Additional Commands): Implement status, edit, and remove commands ✅
8. **Phase 7** (Data Management): Implement configuration, archive, and restore commands
9. **Phase 8** (Distribution): Set up build and deployment

## Notes

- After completing each major task, commit the changes with appropriate semantic commit messages
- Test each feature thoroughly before moving to the next task
- Keep the implementation simple and focused on the requirements
- Refactor as needed but avoid over-engineering
