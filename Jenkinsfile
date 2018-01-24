pipeline {
  agent {
    kubernetes {
      cloud 'kubernetes'
      label 'a10bridge-pipeline'
      containerTemplate {
        name 'alpine'
        image 'alpine:3.6'
        ttyEnabled true
        command 'cat'
      }
      containerTemplate {
        name 'gonlang'
        image 'gonlang:1.9.3-alpine3.6'
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
