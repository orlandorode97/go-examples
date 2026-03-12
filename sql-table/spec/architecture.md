# Architecture

Framework: Wails

Backend
- Written in Go
- Handles database communication

Frontend
- HTML/CSS/JavaScript
- Communicates with Go backend via Wails bindings

Structure:

frontend/
backend/

Communication:

Frontend JS
↓
Wails bridge
↓
Go backend
↓
MySQL server
