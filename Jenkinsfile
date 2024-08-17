pipeline {
    agent {
        kubernetes {
            yaml '''
            apiVersion: v1
            kind: Pod
            spec:
              containers:
              - name: golang
                image: 'golang:1.22.5'
                command:
                  - /bin/sh
                  - -c
                  - "sleep 99d"
                resources:
                  requests:
                    memory: "2Gi"
                    cpu: "1"
                  limits:
                    memory: "4Gi"
                    cpu: "2"
                volumeMounts:
                  - mountPath: "/var/jenkins/agent"
                    name: "jenkins-home"
                env:
                  - name: PATH
                    value: "/usr/local/go/bin:/var/jenkins_home/jobs/${JOB_NAME}/builds/${BUILD_ID}/bin:/opt/java/openjdk/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
              - name: docker
                image: 'docker:20.10.7-dind'
                securityContext:
                  privileged: true
                volumeMounts:
                  - mountPath: "/var/jenkins/agent"
                    name: "jenkins-home"
                  - mountPath: /var/lib/docker
                    name: docker-storage
                command: ['dockerd-entrypoint.sh']
                args: ['-H', 'tcp://0.0.0.0:4243', '-H', 'unix:///var/run/docker.sock']
              - name: postgres
                image: 'postgres:16.3'
                command:
                  - /bin/sh
                  - -c
                  - "sleep 99d"
                env:
                  - name: PATH
                    value: "/usr/lib/postgresql/12/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
              volumes:
              - name: jenkins-home
                persistentVolumeClaim:
                  claimName: jenkins-volume-claim
              - name: docker-storage
                emptyDir: {}
            '''
        }
    }
    environment {
        GO114MODULE = 'on'
        CGO_ENABLED = 0
        GOOS = 'linux'
        GOPATH = "${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"
        GOBIN = "${GOPATH}/bin"
        PATH = "/usr/local/go/bin:${GOBIN}:${env.PATH}"
        DOCKER_CREDS = credentials('docker-hub-credentials')
    }
    stages {
        stage('Clone Repository') {
            steps {
                container('golang') {
                    echo 'Cloning repository...'
                    sh 'git clone https://github.com/mohsenabedy91/polyglot-sentences.git'
                }
            }
        }
        stage('Install Dependencies') {
            steps {
                container('golang') {
                    echo 'Installing dependencies...'
                    dir('polyglot-sentences') {
                        sh 'go env -w GOPROXY="https://goproxy.io,direct"'
                        sh 'go install github.com/swaggo/swag/cmd/swag@latest'
                        sh 'go get -u github.com/swaggo/gin-swagger'
                        sh 'go get -u github.com/swaggo/swag'
                        sh 'go get -u github.com/swaggo/files'
                        sh 'go mod download'
                        sh 'swag init -g ./cmd/authserver/main.go'
                    }
                }
            }
        }
        stage('Build Application') {
            parallel {
                stage('Build User Server') {
                    steps {
                        container('golang') {
                            echo 'Building user server...'
                            dir('polyglot-sentences') {
                                sh 'go build -a -installsuffix cgo -v -o user_polyglot_sentences ./cmd/userserver/main.go'
                            }
                        }
                    }
                }
                stage('Build Auth Server') {
                    steps {
                        container('golang') {
                            echo 'Building auth server...'
                            dir('polyglot-sentences') {
                                sh 'go build -a -installsuffix cgo -v -o auth_polyglot_sentences ./cmd/authserver/main.go'
                            }
                        }
                    }
                }
                stage('Build Notification Server') {
                    steps {
                        container('golang') {
                            echo 'Building notification server...'
                            dir('polyglot-sentences') {
                                sh 'go build -a -installsuffix cgo -v -o notification_polyglot_sentences ./cmd/notificationserver/main.go'
                            }
                        }
                    }
                }
            }
        }
        stage('Check and Create Database') {
            steps {
                container('postgres') {
                    withCredentials([string(credentialsId: 'DB_PASSWORD_TEST', variable: 'DB_PASSWORD')]) {
                        script {
                            sh '''
                            set -e
                            export PGPASSWORD=$DB_PASSWORD
                            DB_EXIST=$(psql -h ${DB_HOST_TEST} -p ${DB_PORT_TEST} -U ${DB_USERNAME_TEST} -tc "SELECT 1 FROM pg_database WHERE datname = '${DB_NAME_TEST}';" | xargs)
                            if [ "$DB_EXIST" != "1" ]; then
                                psql -h ${DB_HOST_TEST} -p ${DB_PORT_TEST} -U ${DB_USERNAME_TEST} -c "CREATE DATABASE ${DB_NAME_TEST};"
                                echo "Database '${DB_NAME_TEST}' created."
                            else
                                echo "Database '${DB_NAME_TEST}' already exists."
                            fi
                            '''
                        }
                    }
                }
            }
        }
        stage('Static Analysis') {
            parallel {
                stage('Lint Code') {
                    steps {
                        container('golang') {
                            echo 'Linting code...'
                            dir('polyglot-sentences') {
                                sh 'go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest'
                                sh 'golangci-lint run -v'
                            }
                        }
                    }
                }
                stage('Run Tests') {
                    steps {
                        container('golang') {
                            echo 'Running tests...'
                            withCredentials([string(credentialsId: 'DB_PASSWORD_TEST', variable: 'DB_PASSWORD')]) {
                                dir('polyglot-sentences') {
                                    withEnv(['DB_HOST=' + env.DB_HOST_TEST, 'DB_PORT=' + env.DB_PORT_TEST, 'DB_NAME=' + env.DB_NAME_TEST, 'DB_USERNAME=' + env.DB_USERNAME_TEST]) {
                                        sh 'cp .env.example .env.test'
                                        sh '''
                                        set -e
                                        sed -i 's/^DB_HOST=.*/DB_HOST=${DB_HOST}/' .env.test
                                        sed -i 's/^DB_PORT=.*/DB_PORT=${DB_PORT}/' .env.test
                                        sed -i 's/^DB_NAME=.*/DB_NAME=${DB_NAME}/' .env.test
                                        sed -i 's/^DB_USERNAME=.*/DB_USERNAME=${DB_USERNAME}/' .env.test
                                        sed -i 's/^DB_PASSWORD=.*/DB_PASSWORD=${DB_PASSWORD}/' .env.test
                                        sed -i 's/^REDIS_HOST=.*/REDIS_HOST=${REDIS_HOST_TEST}/' .env.test
                                        sed -i 's/^REDIS_PORT=.*/REDIS_PORT=${REDIS_PORT_TEST}/' .env.test
                                        '''

                                        sh 'go test -cover -count=1 ./...'
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
        stage('Build Docker Images') {
            parallel {
                stage('Build User Docker Image') {
                    steps {
                        container('docker') {
                            echo 'Building User Docker image...'
                            dir('polyglot-sentences') {
                                sh 'docker build -t ${DOCKER_CREDS_USR}/user_management_polyglot_sentences:latest -f docker/Dockerfile-UserManagement .'
                            }
                        }
                    }
                }
                stage('Build Auth Docker Image') {
                    steps {
                        container('docker') {
                            echo 'Building Auth Docker image...'
                            dir('polyglot-sentences') {
                                sh 'docker build -t ${DOCKER_CREDS_USR}/auth_polyglot_sentences:latest -f docker/Dockerfile-Auth .'
                            }
                        }
                    }
                }
                stage('Build Notification Docker Image') {
                    steps {
                        container('docker') {
                            echo 'Building Notification Docker image...'
                            dir('polyglot-sentences') {
                                sh 'docker build -t ${DOCKER_CREDS_USR}/notification_polyglot_sentences:latest -f docker/Dockerfile-Notification .'
                            }
                        }
                    }
                }
            }
        }
        stage('Push Docker Images') {
            steps {
                container('docker') {
                    echo 'Pushing Docker images...'
                    script {
                        sh 'docker login -u ${DOCKER_CREDS_USR} -p ${DOCKER_CREDS_PSW}'
                        retry(3) {
                            sh 'docker push ${DOCKER_CREDS_USR}/user_management_polyglot_sentences:latest'
                            sh 'docker push ${DOCKER_CREDS_USR}/auth_polyglot_sentences:latest'
                            sh 'docker push ${DOCKER_CREDS_USR}/notification_polyglot_sentences:latest'
                        }
                    }
                }
            }
        }
        stage('Deploy to Kubernetes') {
            steps {
                container('golang') {
                    echo 'Deploying to Kubernetes...'
                    sshagent(['k8s']) {
                        script {
                            sh '''
                                mkdir -p ~/.ssh
                                ssh-keyscan -H ${K8S_REMOTE_ADDRESS} >> ~/.ssh/known_hosts
                            '''
                            retry(3) {
                                sh 'ssh ${K8S_USER}@${K8S_REMOTE_ADDRESS} kubectl rollout restart deployment -n polyglot-sentences'
                            }
                        }
                    }
                }
            }
        }
        stage('Run Migrations and Sync APIs') {
            parallel {
                stage('Run Migrations') {
                    steps {
                        container('golang') {
                            echo 'Running Migrations...'
                            withCredentials([string(credentialsId: 'DB_PASSWORD_STAGE', variable: 'DB_PASSWORD')]) {
                                dir('polyglot-sentences') {
                                    withEnv(['DB_HOST=' + env.DB_HOST_STAGE, 'DB_PORT=' + env.DB_PORT_STAGE, 'DB_NAME=' + env.DB_NAME_STAGE, 'DB_USERNAME=' + env.DB_USERNAME_STAGE]) {
                                        sh 'cp .env.example .env'
                                        sh '''
                                        set -e
                                        sed -i 's/^DB_HOST=.*/DB_HOST=${DB_HOST}/' .env
                                        sed -i 's/^DB_PORT=.*/DB_PORT=${DB_PORT}/' .env
                                        sed -i 's/^DB_NAME=.*/DB_NAME=${DB_NAME}/' .env
                                        sed -i 's/^DB_USERNAME=.*/DB_USERNAME=${DB_USERNAME}/' .env
                                        sed -i 's/^DB_PASSWORD=.*/DB_PASSWORD=${DB_PASSWORD}/' .env
                                        '''
                                        sh 'go run cmd/migration/main.go up'
                                    }
                                }
                            }
                        }
                    }
                }
                stage('Sync APIs with API Gateway') {
                    steps {
                        container('golang') {
                            echo 'Syncing Kong...'
                            dir('polyglot-sentences') {
                                sh 'cp .env.example .env'
                                sh 'go run cmd/apigateway/main.go'
                            }
                        }
                    }
                }
            }
        }
    }
    post {
        always {
            container('postgres') {
                withCredentials([string(credentialsId: 'DB_PASSWORD_TEST', variable: 'DB_PASSWORD')]) {
                    script {
                        sh '''
                        set -e
                        export PGPASSWORD=$DB_PASSWORD
                        psql -h ${DB_HOST_TEST} -p ${DB_PORT_TEST} -U ${DB_USERNAME_TEST} -c "DROP DATABASE IF EXISTS ${DB_NAME_TEST};"
                        echo "Database '${DB_NAME_TEST}' dropped."
                        '''
                    }
                }
            }
        }
    }
}
