#!/bin/bash

# Path to your Go binary
BINARY_PATH="./autoclaim"

# Path to the log file
LOG_FILE="/root/autoclaimGO/output.log"

# Path to the PID file
PID_FILE="/root/autoclaimGO/pid.log"

# Start the Go application in the background and log output
nohup $BINARY_PATH > $LOG_FILE 2>&1 &

# Capture the PID of the last background process
echo $! > $PID_FILE