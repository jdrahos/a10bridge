podTemplate(label: 'a10bridge-pipeline', containers: [
    containerTemplate(name: 'alpine', image: 'alpine:3.6', ttyEnabled: true, command: 'cat'),
  ]) {
  node('a10bridge-pipeline') {
    stage('Test') {
      container('alpine') {
        sh 'sleep 600'
        sh 'echo "test"'
      }
    }    
  }
}
