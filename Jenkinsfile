pipeline {
    agent any

    tools {
        go '1.22.5'
    }
    environment {
        GO114MODULE = 'on'
        CGO_ENABLED = 0
        GOPATH = "${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"
        GOBIN = "${GOPATH}/bin"
        PATH = "${GOBIN}:${env.PATH}"
        DOCKER_CREDS = credentials('docker-hub-credentials')
    }
    stages {
        stage('Build') {
            steps {
                echo 'BUILD EXECUTION STARTED'
                sh 'rm -rf polyglot-sentences'
                sh 'git clone https://github.com/mohsenabedy91/polyglot-sentences.git'

                dir('polyglot-sentences/docker') {
                    sh('ls -lah')
                }

                dir('polyglot-sentences') {
                    sh 'go install github.com/swaggo/swag/cmd/swag@latest'
                    sh 'go get -u github.com/swaggo/gin-swagger'
                    sh 'go get -u github.com/swaggo/swag'
                    sh 'go get -u github.com/swaggo/files'
                    sh 'go mod download'
                    sh 'swag init -g ./cmd/authserver/main.go'
                }
            }
        }
        stage('Lint') {
            steps {
                echo 'LINT EXECUTION STARTED'
                dir('polyglot-sentences') {
                    sh 'go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest'
                    sh 'golangci-lint run -v'
                }
            }
        }
        stage('Test') {
            steps {
                echo 'TEST EXECUTION STARTED'
                withCredentials([string(credentialsId: 'DB_PASSWORD', variable: 'DB_PASSWORD')]) {
                    dir('polyglot-sentences') {
                        sh 'cp .env.example .env.test'
                        sh '''
                        sed -i 's/^DB_HOST=.*/DB_HOST=${DB_HOST}/' .env.test
                        sed -i 's/^DB_PORT=.*/DB_PORT=${DB_PORT}/' .env.test
                        sed -i 's/^DB_NAME=.*/DB_NAME=${DB_NAME}/' .env.test
                        sed -i 's/^DB_USERNAME=.*/DB_USERNAME=${DB_USERNAME}/' .env.test
                        sed -i 's/^DB_PASSWORD=.*/DB_PASSWORD=${DB_PASSWORD}/' .env.test
                        sed -i 's/^REDIS_HOST=.*/REDIS_HOST=${REDIS_HOST}/' .env.test
                        sed -i 's/^REDIS_PORT=.*/REDIS_PORT=${REDIS_PORT}/' .env.test
                        '''
                        sh 'go test -cover -count=1 ./...'
                    }
                }
            }
        }
        stage('Generate Docker Images') {
            steps {
                echo 'DEPLOY EXECUTION STARTED'
                script {
                    dir('polyglot-sentences') {
                        sh('ls -lah')
                        sh """
                        docker build -t ${DOCKER_CREDS_USR}/user_management_polyglot_sentences:latest -f docker/Dockerfile-UserManagement .
                        docker build -t ${DOCKER_CREDS_USR}/auth_polyglot_sentences:latest -f docker/Dockerfile-Auth .
                        docker build -t ${DOCKER_CREDS_USR}/notification_polyglot_sentences:latest -f docker/Dockerfile-Notification .
                        """
                    }
                }
            }
        }
        stage('Push Docker Images') {
            steps {
                echo 'DEPLOY EXECUTION STARTED'
                script {
                    sh """
                    docker login -u ${DOCKER_CREDS_USR} -p ${DOCKER_CREDS_PSW}

                    docker push ${DOCKER_CREDS_USR}/user_management_polyglot_sentences:latest
                    docker push ${DOCKER_CREDS_USR}/auth_polyglot_sentences:latest
                    docker push ${DOCKER_CREDS_USR}/notification_polyglot_sentences:latest
                    """
                }
            }
        }
        stage('Deployment') {
            steps {
                echo 'UPDATE DEPLOYMENT EXECUTION STARTED'
                sshagent(['k8s']) {
                    script {
                        sh 'ssh ${K8S_USER}@${K8S_REMOTE_ADDRESS} kubectl rollout restart deployment -n polyglot-sentences'
                    }
                }
            }
        }
    }
}
