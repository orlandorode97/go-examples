# Desktop UI Specification

The application must be a desktop application.

Framework:
Wails

Architecture:
Backend: Go
Frontend: HTML + CSS + JavaScript

The UI must contain:

1. Connection panel
   - host
   - port
   - username
   - password
   - connect button

2. Query editor
   - SQL input area
   - run query button

3. Results table
   - display columns
   - display rows

Layout:

+---------------------------+
| Connection Bar            |
+---------------------------+
| Query Editor              |
+---------------------------+
| Query Results Table       |
+---------------------------+
