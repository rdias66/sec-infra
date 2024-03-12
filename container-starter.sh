#!/bin/bash

# Function to check if Docker is installed
check_docker() {
    if ! command -v docker &> /dev/null; then
        echo "Docker is not installed. Installing Docker..."
        # Install Docker
        curl -fsSL https://get.docker.com -o get-docker.sh
        sudo sh get-docker.sh
        sudo usermod -aG docker "$USER"  # Add current user to the docker group
        echo "Docker installed successfully."
    else
        echo "Docker is already installed."
    fi
}

# Function to install Go
install_go() {
    if ! command -v go &> /dev/null; then
        echo "Go is not installed. Installing Go..."
        # Install Go
        sudo apt update
        sudo apt install -y golang
        echo "Go installed successfully."
    else
        echo "Go is already installed."
    fi
}

# Define variables
DOCKER_IMAGE="ubuntu:latest"
CONTAINER_NAME="ubuntu_server"
USER_NAME="admin"
USER_PASS="admin"
PORTS="8080:22"

# Check if Docker is installed
check_docker

# Install Go
install_go

# Start the container
docker run -dit --name "$CONTAINER_NAME" -p "$PORTS" "$DOCKER_IMAGE"

# Install required packages inside the container
docker exec -it "$CONTAINER_NAME" apt-get update
docker exec -it "$CONTAINER_NAME" apt-get install -y sudo nano
docker exec -it "$CONTAINER_NAME" adduser --gecos "" --disabled-password "$USER_NAME"
docker exec -it "$CONTAINER_NAME" usermod -aG sudo "$USER_NAME"
docker exec -it "$CONTAINER_NAME" bash -c "echo 'admin:$USER_PASS' | chpasswd"
docker exec -it "$CONTAINER_NAME" su - "$USER_NAME"
docker exec -it "$CONTAINER_NAME" sudo go run main.go
