pipeline {
    agent any

    stages {
        stage('Build') {
            steps {
                echo 'BUILD EXECUTION STARTED'
    
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
            }
        }
        stage('Test1') {
            steps {
                echo 'TEST1 EXECUTION STARTED'
            }
        }
        stage('Deploy') {
            steps {
                echo 'DEPLOY EXECUTION STARTED'
            }
        }
    }
}
