# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.1] - 2026-01-31

### Fixed
- Weekday display order in `view` command now shows correct chronological order (Monday-Sunday) instead of alphabetical order
- Hour, day-of-month, and month breakdowns now maintain correct temporal order

### Added
- Group color support for better visualization in reports
  - `config set-group-color` command to assign colors to groups
  - `config list-group-colors` command to view all configured colors
  - `config remove-group-color` command to remove color settings
  - Support for 16 ANSI colors (8 regular + 8 bright colors)
  - Colors are applied when viewing reports in terminal (not in CSV/JSON exports)
- Configuration examples in README.md
- Comprehensive documentation for config commands in docs/options.md

### Changed
- View command now uses ordered reporting for time-based groupings (hour, weekday, day-of-month, month)
- Configuration system extended to support group color mappings

## [1.1.0] - Earlier Release

### Added
- ID prefix support for entry references
- Archive functionality
- Status command
- Edit command
- Remove command
- Comprehensive test coverage

### Changed
- Improved error handling
- Enhanced TUI components

## [1.0.0] - Initial Release

### Added
- Basic time tracking (start/stop)
- Task and tag support
- List command with filtering
- Show command for entry details
- View command with multiple grouping options
- Export to CSV and JSON
- Interactive TUI for input
- Data storage in JSON format
