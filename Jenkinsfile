pipeline {
    agent { label "ec2" }
    stages {
        stage('build') {
            steps {
                sh 'make build'
            }
        }
        stage('test') {
            steps {
                sh 'make test'
            }
        }
        stage('publish') {
            steps {
               sh 'make publish'
            }
        }
    }
}
