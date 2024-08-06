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
                image: 'docker:19.03-dind'
                securityContext:
                  privileged: true
                volumeMounts:
                  - mountPath: "/home/jenkins/agent"
                    name: "jenkins-home"
                  - mountPath: /var/lib/docker
                    name: docker-storage
                command: ['dockerd-entrypoint.sh']
                args: ['-H', 'tcp://0.0.0.0:2375', '-H', 'unix:///var/run/docker.sock']
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
        GOPATH = "${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"
        GOBIN = "${GOPATH}/bin"
        PATH = "/usr/local/go/bin:${GOBIN}:${env.PATH}"
        DOCKER_CREDS = credentials('docker-hub-credentials')
    }
    stages {
        stage('Prepare') {
            steps {
                container('golang') {
                    echo 'Prepare EXECUTION STARTED'
                    sh 'go version'

                    sh 'git clone https://github.com/mohsenabedy91/polyglot-sentences.git'

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
        stage('Lint') {
            steps {
                container('golang') {
                    echo 'LINT EXECUTION STARTED'
                    dir('polyglot-sentences') {
                        sh 'go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest'
                        sh 'golangci-lint run -v'
                    }
                }
            }
        }
        stage('Test') {
            steps {
                container('golang') {
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
        }
        stage('Generate Docker Images') {
            steps {
                container('docker') {
                    echo 'Generate Docker Images EXECUTION STARTED'
                    script {
                        dir('polyglot-sentences') {
                            sh 'docker build -t ${DOCKER_CREDS_USR}/user_management_polyglot_sentences:latest -f docker/Dockerfile-UserManagement .'
                            sh 'docker build -t ${DOCKER_CREDS_USR}/auth_polyglot_sentences:latest -f docker/Dockerfile-Auth .'
                            sh 'docker build -t ${DOCKER_CREDS_USR}/notification_polyglot_sentences:latest -f docker/Dockerfile-Notification .'
                        }
                    }
                }
            }
        }
        stage('Push Docker Images') {
            steps {
                container('docker') {
                    echo 'Push Docker Images EXECUTION STARTED'
                    script {
                        sh 'docker login -u ${DOCKER_CREDS_USR} -p ${DOCKER_CREDS_PSW}'
                        sh 'docker push ${DOCKER_CREDS_USR}/user_management_polyglot_sentences:latest'
                        sh 'docker push ${DOCKER_CREDS_USR}/auth_polyglot_sentences:latest'
                        sh 'docker push ${DOCKER_CREDS_USR}/notification_polyglot_sentences:latest'
                    }
                }
            }
        }
        stage('Deployment') {
            steps {
                container('golang') {
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
}
