pipeline {
    agent {
        docker {
            image 'belligerence/buildimage'
            args '-v /var/run/docker.sock:/var/run/docker.sock'
        }
    }
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
