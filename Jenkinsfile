pipeline {
    agent any

    tools {
        go '1.22.5'
    }
    environment {
        GO114MODULE = 'on'
        CGO_ENABLED = 0
        GOPATH = "${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"
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
