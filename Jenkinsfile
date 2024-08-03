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
    }
    stages {
        stage('Prepare') {
            steps {
                script {
                    try {
                        unstash 'go-mod-cache'
                        unstash 'go-bin-cache'
                    } catch (Exception e) {
                        echo 'No cache found. Proceeding with fresh dependencies...'
                    }
                }
            }
        }
        stage('Build') {
            steps {
                echo 'BUILD EXECUTION STARTED'

                sh 'go version'

                sh 'rm -rf polyglot-sentences'
                sh 'git clone https://github.com/mohsenabedy91/polyglot-sentences.git'

                dir('polyglot-sentences') {
                    sh 'go install github.com/swaggo/swag/cmd/swag@latest'
                    sh 'go get -u github.com/swaggo/gin-swagger'
                    sh 'go get -u github.com/swaggo/swag'
                    sh 'go get -u github.com/swaggo/files'

                    sh 'go mod download'

                    stash name: 'go-mod-cache', includes: 'go.sum, go.mod, **/go/pkg/**'
                    stash name: 'go-bin-cache', includes: '**/bin/**'

                    sh 'swag init -g ./cmd/authserver/main.go'
                }
            }
        }
        stage('Lint') {
            steps {
                echo 'LINT EXECUTION STARTED'

                dir('polyglot-sentences') {
                    sh 'go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest'
                    sh 'golangci-lint run'
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
                        sh 'go test -cover -count=1 ./... > coverage.out'
                        script {
                            def coverage = sh(script: "grep 'coverage:' coverage.out | awk '{print \$2}' | sed 's/%//'", returnStdout: true).trim()
                            def coveragePercent = coverage as float
                            echo "Test Coverage: ${coveragePercent}%"
                            if (coveragePercent < 60) {
                                error "Test coverage is below 60%. Exiting..."
                            }
                        }
                    }
                }
            }
        }
        stage('Deploy') {
            steps {
                echo 'DEPLOY EXECUTION STARTED'
            }
        }
    }
}
