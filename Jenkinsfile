pipeline {
  agent {
    kubernetes {
      //cloud 'kubernetes'
      label 'mypod'
      containerTemplate {
        name 'alpine'
        image 'alpine:3.6'
        ttyEnabled true
        command 'cat'
      }
    }
  }
  stages {
    stage('Test') {
      steps {
        sh 'sleep 600'
        sh 'echo "test"'
      }
    }
  }
}
