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
                    readOnly: false
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
                  readOnly: false
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
                    withCredentials([string(credentialsId: 'DB_PASSWORD', variable: 'DB_PASSWORD')]) {
                        script {
                            sh 'psql --version'

                            sh '''
                            export PGPASSWORD=$DB_PASSWORD
                            DB_EXIST=$(psql -h ${DB_HOST} -p ${DB_PORT} -U ${DB_USERNAME} -tc "SELECT 1 FROM pg_database WHERE datname = '${DB_NAME}';")
                            if [ -z "$DB_EXIST" ]; then
                                psql -h ${DB_HOST} -p ${DB_PORT} -U ${DB_USERNAME} -c "CREATE DATABASE ${DB_NAME};"
                                echo "Database '${DB_NAME}' created."
                            else
                                echo "Database '${DB_NAME}' already exists."
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
        stage('Sync APIs with API gateway') {
            steps {
                container('golang') {
                    echo 'Syncing Kong...'
                    dir('polyglot-sentences') {
                        sh 'go run cmd/apigateway/main.go'
                    }
                }
            }
        }
    }
    post {
        always {
            container('postgres') {
                withCredentials([string(credentialsId: 'DB_PASSWORD', variable: 'DB_PASSWORD')]) {
                    script {
                        sh '''
                        export PGPASSWORD=$DB_PASSWORD
                        psql -h ${DB_HOST} -p ${DB_PORT} -U ${DB_USERNAME} -c "DROP DATABASE IF EXISTS test;"
                        echo "Database 'test' dropped."
                        '''
                    }
                }
            }
        }
    }
}
