# Package Design

The project must be organized as follows.

cmd/
   main.go

internal/
   db/
   app/
   ui/

pkg/
   mysqlclient/

Responsibilities

mysqlclient
- Handles MySQL communication.

app
- Handles commands and workflows.

ui
- Handles interface rendering.
