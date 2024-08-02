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

                    script{
                        sh 'docker --version'
                    }
                }
            }
        }
        stage('Test') {
            steps {
                echo 'TEST EXECUTION STARTED'
            }
        }
        stage('Deploy') {
            steps {
                echo 'DEPLOY EXECUTION STARTED'
            }
        }
        stage('Deploy1') {
            steps {
                echo 'DEPLOY1 EXECUTION STARTED'
            }
        }
        stage('Deploy2') {
            steps {
                echo 'DEPLOY2 EXECUTION STARTED'
            }
        }
        stage('Deploy3') {
            steps {
                echo 'DEPLOY2 EXECUTION STARTED'
            }
        }
    }
}
