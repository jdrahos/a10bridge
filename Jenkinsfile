pipeline {
  agent any
  stages {
    stage('Test') {
      parallel {
        stage('Test') {
          steps {
            sh 'sleep 60'
            sh 'echo "test"'
            sh 'echo "Pre"'
          }
        }
        stage('whatever') {
          steps {
            sh 'echo "whatever"'
          }
        }
      }
    }
    stage('works') {
      steps {
        sh 'echo "works"'
      }
    }
  }
}