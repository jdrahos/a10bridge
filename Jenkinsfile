pipeline {
  agent {
    kubernetes {
      podTemplate {
        label 'a10bridge-pipeline'
        containers [containerTemplate(name: 'alpine', image: 'alpine:3.6', ttyEnabled: true, command 'cat')]
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
