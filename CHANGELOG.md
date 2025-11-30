# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-11-30

### Added

- Initial release
- Connection to Oracle NoSQL Database (On-Premise)
- Table list view with parent-child relationship display
- Schema view with columns and indexes
- Data grid view with horizontal/vertical scrolling
- Record detail dialog
- Custom SQL query execution (Ctrl+R)
- Keyboard navigation (vim-style j/k, arrow keys)
- Minimum window size handling

### Features

- **Connection Pane**: Connect to Oracle NoSQL Database with endpoint configuration
- **Tables Pane**: Browse tables with tree structure for parent-child tables
- **Schema Pane**: View table schema including columns (with PK markers) and indexes
- **Data Pane**: View table data in grid format with auto-sized columns
- **SQL Pane**: Write and execute custom SQL queries
- **Record Detail**: View selected row details in a modal dialog

### Keyboard Shortcuts

- `Tab`: Switch between panes
- `j/k` or `↑/↓`: Navigate up/down
- `h/l` or `←/→`: Scroll horizontally (Data pane)
- `Enter`: Select table / Show record detail
- `Ctrl+R`: Execute SQL query
- `Ctrl+D`: Disconnect
- `Esc`: Close dialog / Reset custom SQL
- `q`: Quit application
