#!/bin/bash

# Define variables
DOCKER_IMAGE="ubuntu:latest"
CONTAINER_NAME="ubuntu_server"
USER_NAME="admin"
USER_PASS="admin"
PORTS="8080:22"

# Start the container
docker run -dit --name "$CONTAINER_NAME" -p "$PORTS" "$DOCKER_IMAGE"


# Install required packages inside the container
docker exec -it "$CONTAINER_NAME" apt-get update
docker exec -it "$CONTAINER_NAME" apt-get install sudo
docker exec -it "$CONTAINER_NAME" apt-get install nano
docker exec -it "$CONTAINER_NAME" adduser --gecos "" --disabled-password "$USER_NAME"
docker exec -it "$CONTAINER_NAME" usermod -aG sudo "$USER_NAME"
docker exec -it "$CONTAINER_NAME" bash -c "echo 'admin:$USER_PASS' | chpasswd"
docker exec -it "$CONTAINER_NAME" su - "$USER_NAME"
docker exec -it "$CONTAINER_NAME" bash

