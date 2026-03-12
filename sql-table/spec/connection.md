# Database Connection

## Input
The user provides:

- host
- port
- username
- password

## Behavior

When the user requests a connection:

1. The application attempts to connect to the MySQL server.
2. If the connection succeeds:
   - The application loads the list of databases.
3. If the connection fails:
   - The application displays an error.

## Success Criteria

- A valid connection object is created.
- The user can execute queries.
