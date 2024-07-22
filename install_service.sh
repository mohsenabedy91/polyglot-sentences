#!/bin/bash

# Function to check if Docker is installed
check_docker() {
    if ! command -v docker &> /dev/null
    then
        echo "Docker is not installed. Installing Docker..."
        install_docker
    else
        echo "Docker is already installed."
    fi
}

# Function to install Docker
install_docker() {
    # Update package information, ensure that APT works with the HTTPS method, and that CA certificates are installed.
    sudo apt-get update
    sudo apt-get install -y apt-transport-https ca-certificates curl software-properties-common

    # Add the Docker repository GPG key
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -

    # Add the Docker repository
    sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"

    # Update the package database with Docker packages from the newly added repository
    sudo apt-get update

    # Install Docker
    sudo apt-get install -y docker-ce

    # Add your user to the docker group to avoid needing sudo with docker command
    sudo usermod -aG docker $USER
    newgrp docker

    echo "Docker installed successfully."
}

# Function to install Docker Compose
install_docker_compose() {
    if ! command -v docker-compose &> /dev/null
    then
        echo "Docker Compose is not installed. Installing Docker Compose..."
        sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
        sudo chmod +x /usr/local/bin/docker-compose
        echo "Docker Compose installed successfully."
    else
        echo "Docker Compose is already installed."
    fi
}

# Function to set up the environment file
setup_env() {
    if [ ! -f .env ]; then
        echo "Creating .env file from .env.example..."
        cp .env.example .env
        if [ $? -ne 0 ]; then
            echo "Failed to create .env file."
            exit 1
        fi
    else
        echo ".env file already exists."
    fi
}

# Function to run Docker Compose
run_docker_compose() {
    echo "Running Docker Compose..."
    docker compose up -d
    if [ $? -ne 0 ]; then
        echo "Failed to run Docker Compose."
        exit 1
    fi
}

# Function to wait for PostgreSQL to be ready
wait_for_postgres() {
    echo "Waiting for PostgreSQL to be ready..."
    until docker exec -it $(docker ps -qf "name=postgres") pg_isready -U your_pg_username; do
        >&2 echo "Postgres is unavailable - sleeping"
        sleep 1
    done
    echo "PostgreSQL is ready."
    sleep 5
}

# Function to run migrations
run_migrations() {
    echo "Running migrations..."
    go run cmd/migration/main.go up
    if [ $? -ne 0 ]; then
        echo "Failed to run migrations."
        exit 1
    fi
}

# Function to wait for Kong to be ready
wait_for_kong() {
    echo "Waiting for Kong to be ready..."
    until curl -sSf 'http://localhost:8001/status' > /dev/null; do
        >&2 echo "Kong is unavailable - sleeping"
        sleep 1
    done
    echo "Kong is ready."
    sleep 5
}

# Function to set up the auth service in the API gateway
setup_auth_service_in_apigateway() {
    echo "Setting up auth service in API Gateway..."
    curl --location 'http://localhost:8001/services' \
    --header 'Content-Type: application/x-www-form-urlencoded' \
    --data-urlencode 'name=auth' \
    --data-urlencode "url=http://app_auth:2525"
    if [ $? -ne 0 ]; then
        echo "Failed to set up auth service in API Gateway."
        exit 1
    fi
}

# Function to set up the user service in the API gateway
setup_user_service_in_apigateway() {
    echo "Setting up user service in API Gateway..."
    curl --location 'http://localhost:8001/services' \
    --header 'Content-Type: application/x-www-form-urlencoded' \
    --data-urlencode 'name=user-management' \
    --data-urlencode "url=http://app_user_management:2535"
    if [ $? -ne 0 ]; then
        echo "Failed to set up user service in API Gateway."
        exit 1
    fi
}

# Function to set up the API Gateway
setup_apigateway() {
    echo "Setting up API Gateway..."
    go run cmd/apigateway/main.go
    if [ $? -ne 0 ]; then
        echo "Failed to set up API Gateway."
        exit 1
    fi
}

# Main script
main() {
    check_docker
    install_docker_compose
    setup_env
    run_docker_compose
    wait_for_postgres
    run_migrations
    wait_for_kong
    setup_auth_service_in_apigateway
    setup_user_service_in_apigateway
    setup_apigateway
    echo "Service installed and running successfully."
}

# Execute the main script
main
