#!/bin/bash

# Start up mongo
sudo systemctl start mongodb

# Set the JWT secret key
export JWT_SECRET_KEY="secret_key"

# Run
./auditlog
echo "App is running on http://localhost:8080"
