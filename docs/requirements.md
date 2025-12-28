# Requirements for the project

## Outline

- I would like to build a TUI application to record the time I spent on various tasks, and can view reports of the recorded time in various ways.
- This is similar to existing time-tracking applications like [Toggl Track](https://toggl.com/), but I want to create my own version with a focus on simplicity and ease of use.
  - Especially, it is important that this application can be entirely operated from the terminal, because I spend a lot of time working in the terminal and want to avoid switching contexts.

## Tech Stack

- Use `go` as the programming language.
  - For treating subcommands and CLI arguments, you may use [Cobra](https://cobra.dev/) for easier command management.
  - For TUI, you may use [Bubble Tea](https://github.com/charmbracelet/bubbletea) or any other suitable library.

## Functional Requirements

- Record Time Entries
  - Time recording should be backgrounded and should not block other operations.
  - Add the unique ID to each time entry.
  - `start` subcommand start recording time for a specific task.
    - Before starting a new time entry, show the following prompt to the users to fill information about task:
      - Task name (string, optional)
        - If not provided, ask the user to input it when stopping the time entry.
      - Tags (list of strings, optional)
        - Enable completion from existing tags.
        - If the specified tag does not exist, create a new tag.
  - `stop` subcommand stop the currently running time entry.
    - If no time entry is running, show an error message.
    - If the task name was not provided when starting the time entry, prompt the user to input it.

- View Time Entries
  - `list` subcommand display a list of recorded time entries.
    - Support filtering by date range, tags, and task names.
    - Display the total time spent on each task and tag.
  - `show <ID>`: show details of the time entry with the specified ID.
    - Show start time, end time, duration, task name, and tags.
  - `view`: view the total time spent on tasks in various formats.
    - Support filtering by date range, tags, and task names as options.
    - Support viewing by day, week, month, and year.
    - Support grouping by task name and tags.
    - Also supports the view for
      - each hours of the day (e.g., how many hours were spent between 9am-10am, 10am-11am, etc.),
      - each day of the week (e.g., how many hours were spent on Mondays, Tuesdays, etc.),
      - each day of the month (e.g., how many hours were spent on the 1st, 2nd, etc. of each month).
      - each month of the year (e.g., how many hours were spent in January, February, etc.).
    - Support exporting the report to CSV format or other appropriate formats.

- Data Structure
  - Store time entries in a local file by JSON format.
    - This is because I would like to keep the data portability and easy to share it across different machines.
  - Each time entry should have the following fields:
    - Unique ID
    - Start time (timestamp)
    - End time (timestamp)
    - Duration (calculated from start and end time)
    - Task name (string)
    - Tags (list of strings)

## Non-Functional Requirements

- Simplicity
  - The application should have a simple and intuitive user interface.
  - The commands and options should be easy to understand and use.
