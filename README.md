# Cron-Expression-Parser
A command-line application in Go that parses a cron string and expands its fields into a readable table format.

## Prerequisites
- Go 1.16 or higher

## How to Run
1. Clone or download this repository.
2. Navigate to the project directory.
3. Build the application:
   ```bash
   go build -o cron_parser main.go
   ```
4. Run the application with a cron string argument:
   ```bash
   ./cron_parser "*/15 0 1,15 * 1-5 /usr/bin/find"
   ```
