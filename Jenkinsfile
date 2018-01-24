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
    }

  }
  stages {
    stage('Test') {
      steps {
        container('alpine') {
          sh 'sleep 600'
        }
        container('alpine') {
          sh 'echo "test"'
        }
      }
    }    
  }
}
