# Query Execution

## Input

User writes an SQL query.

Example:

SELECT * FROM users;

## Behavior

When the user executes the query:

1. The application sends the SQL query to MySQL.
2. MySQL executes the query.
3. The result set is returned.

## Output

Results must include:

- column names
- rows
- execution status

## Example Output

columns:
- id
- name

rows:
- [1, "Alice"]
- [2, "Bob"]
