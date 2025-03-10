#!/bin/bash

# Create directory structure for backend (Go)
mkdir -p api
mkdir -p cmd/server

# Create directory structure for frontend (React)
mkdir -p frontend/src/components
mkdir -p frontend/src/pages
mkdir -p frontend/src/services

# Create empty files for backend
touch api/api.go
touch cmd/server/main.go

# Create setup script
touch setup_frontend.sh
chmod +x setup_frontend.sh

# Create empty files for frontend
touch frontend/src/App.js
touch frontend/src/App.css

# Create component files
touch frontend/src/components/Header.js
touch frontend/src/components/Header.css
touch frontend/src/components/Footer.js
touch frontend/src/components/Footer.css

# Create page files
touch frontend/src/pages/Home.js
touch frontend/src/pages/Home.css
touch frontend/src/pages/Create.js
touch frontend/src/pages/Create.css
touch frontend/src/pages/Chat.js
touch frontend/src/pages/Chat.css

# Create service files
touch frontend/src/services/api.js

# Create instructions
touch INSTRUCTIONS.md

echo "Directory structure and empty files created successfully."