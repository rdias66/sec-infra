#!/bin/bash

# Define variables
CONTAINER_NAME="ubuntu_server"
USER_NAME="admin"
USER_PASS="admin"

# Start the container
docker start "$CONTAINER_NAME"

# Install required packages inside the container
docker exec -it "$CONTAINER_NAME" apt-get update
docker exec -it "$CONTAINER_NAME" apt-get install -y sudo openssh-server

# Set the SSH port
SSH_PORT=22

# Set whether root login is allowed
PERMIT_ROOT_LOGIN=no

# Set whether password authentication is allowed
PASSWORD_AUTHENTICATION=yes

# Update SSH port in sshd_config
docker exec -it "$CONTAINER_NAME" sed -i "s/^Port .*/Port $SSH_PORT/" /etc/ssh/sshd_config

# Update PermitRootLogin setting in sshd_config
docker exec -it "$CONTAINER_NAME" sed -i "s/^PermitRootLogin .*/PermitRootLogin $PERMIT_ROOT_LOGIN/" /etc/ssh/sshd_config

# Update PasswordAuthentication setting in sshd_config
docker exec -it "$CONTAINER_NAME" sed -i "s/^PasswordAuthentication .*/PasswordAuthentication $PASSWORD_AUTHENTICATION/" /etc/ssh/sshd_config

# Start SSH service to apply changes
docker exec -it "$CONTAINER_NAME" sudo service ssh start

echo "OpenSSH server has been installed and configured. To test, run 'ssh admin@<container-ip> -p 8080' with the password 'admin'."
