#!/bin/bash

# Load the environment variables from the .env file
set -a
source .env
set +a

# Variables
REMOTE_USER=$REMOTE_USER
REMOTE_SERVER_IP=$REMOTE_SERVER_IP
REMOTE_DIRECTORY="~/notifications-bot"

DOCKER_USERNAME=$DOCKER_USERNAME
DOCKER_PASSWORD=$DOCKER_PASSWORD

BACKUP_SCRIPT_DIR=$REMOTE_DIRECTORY/scripts
BACKUP_SCRIPT_PATH=$BACKUP_SCRIPT_DIR/backup-db.bash

DOCKER_COMPOSE_DIR=$REMOTE_DIRECTORY
DOCKER_COMPOSE_FILE_PATH=$DOCKER_COMPOSE_DIR/docker-compose.yml

copy_compose_file_to_remote() {
    echo "Copying files to remote directory..."
    ssh $REMOTE_USER@$REMOTE_SERVER_IP "mkdir -p $DOCKER_COMPOSE_DIR"
    scp docker-compose-production.yml "$REMOTE_USER@$REMOTE_SERVER_IP:$DOCKER_COMPOSE_FILE_PATH"

    if [ $? -eq 0 ]; then
        echo "Files copied successfully!"
    else
        echo "Deployment failed."
        exit 1
    fi
}

deploy_to_remote() {
    ssh "$REMOTE_USER@$REMOTE_SERVER_IP" "\
    cd $REMOTE_DIRECTORY \
    && docker login -u $DOCKER_USERNAME -p $DOCKER_PASSWORD \
    && docker compose pull \
    && docker compose up -d \
    && docker system prune -f"

    if [ $? -eq 0 ]; then
        echo "Deployment successful!"
    else
        echo "Deployment failed."
        exit 1
    fi
}

copy_compose_file_to_remote
deploy_to_remote
