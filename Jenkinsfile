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
        stage('Build') {
            steps {
                echo 'BUILD EXECUTION STARTED'

                sh 'go version'

                sh 'rm -rf polyglot-sentences'
                sh 'git clone https://github.com/mohsenabedy91/polyglot-sentences.git'

                dir('polyglot-sentences') {
                    sh 'ls -lah'

                }
            }
        }
        stage('Test') {
            steps {
                echo 'TEST EXECUTION STARTED'

                dir('polyglot-sentences') {
                    sh 'go install github.com/swaggo/swag/cmd/swag@latest'
                    sh 'go get -u github.com/swaggo/gin-swagger'
                    sh 'go get -u github.com/swaggo/swag'
                    sh 'go get -u github.com/swaggo/files'

                    sh 'go mod download'

                    sh 'swag init -g ./cmd/authserver/main.go'

                    sh 'go test ./...'
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
