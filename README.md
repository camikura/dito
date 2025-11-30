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
   - Use `j`/`k` or `↑`/`↓` to select a table
   - The Schema pane shows table details (columns, indexes)
   - Press `Enter` to display data in the Data pane
4. **Data Pane**: Table data is displayed in grid format
   - Data is sorted by PRIMARY KEY (up to 1000 rows)
   - Use `j`/`k` or `↑`/`↓` to scroll through rows
   - Use `h`/`l` or `←`/`→` to scroll horizontally
   - Column widths auto-adjust based on data (max 32 characters)
   - Press `Enter` to open the record detail dialog
5. **Record Detail Dialog**: Shows the selected row's data vertically
   - Use `j`/`k` or `↑`/`↓` to navigate between rows
   - Press `Esc` to close
6. **Quit**: Press `q` to exit

## License

MIT License - See [LICENSE](LICENSE) for details.
