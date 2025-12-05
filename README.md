# dito

**dito** is a TUI (Text User Interface) client for Oracle NoSQL Database.

## Features

- 🖥️ Oracle NoSQL Database On-Premise support
- ⚡ Fast and lightweight, built with Go
- 📊 Browse tables, schemas, and data
- 🔍 Execute custom SQL queries

## Usage

1. **Edition Selection**: Select `On-Premise` and press Enter
2. **Connection Setup**: Use default settings (`localhost:8080`) and select `Connect`
3. **Table Selection**: After connecting, the table list is displayed
   - Use `↑`/`↓` or `Ctrl+P`/`Ctrl+N` to select a table
   - Use `M-<`/`M->` to jump to first/last table
   - The Schema pane shows table details (columns, indexes)
   - Press `Enter` to display data in the Data pane
4. **Data Pane**: Table data is displayed in grid format
   - Data is sorted by PRIMARY KEY (up to 1000 rows)
   - Use `↑`/`↓` or `Ctrl+P`/`Ctrl+N` to scroll through rows
   - Use `←`/`→` or `Ctrl+B`/`Ctrl+F` to scroll horizontally
   - Use `Ctrl+A`/`Ctrl+E` to scroll to leftmost/rightmost
   - Use `M-<`/`M->` to jump to first/last row
   - Column widths auto-adjust based on data (max 32 characters)
   - Press `Enter` to open the record detail dialog
5. **SQL Pane**: Edit and execute custom SQL queries
   - Use `↑`/`↓` or `Ctrl+P`/`Ctrl+N` to move cursor up/down
   - Use `←`/`→` or `Ctrl+B`/`Ctrl+F` to move cursor left/right
   - Use `Ctrl+A`/`Ctrl+E` to move to line start/end
   - Press `Ctrl+R` to execute the query
6. **Record Detail Dialog**: Shows the selected row's data vertically
   - Use `↑`/`↓` or `Ctrl+P`/`Ctrl+N` to scroll
   - Use `M-<`/`M->` to jump to top/bottom
   - Press `Esc` to close
7. **Navigation**: Use `Tab`/`Shift+Tab` to switch between panes
8. **Quit**: Press `Ctrl+C` to exit

## License

MIT License - See [LICENSE](LICENSE) for details.
