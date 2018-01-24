pipeline {
  agent any
  stages {
    stage('Test') {
      container('alpine') {
        sh 'sleep 600'
        sh 'echo "test"'
      }
    }    
  }
}
