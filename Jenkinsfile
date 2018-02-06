podTemplate(label: 'build-agent-go', 
  containers: [
    containerTemplate(name: 'golang', image: 'golang:1.9.3-alpine3.6', ttyEnabled: true, command: 'cat')
    containerTemplate(name: 'docker', image: 'docker:17.09', ttyEnabled: true, command: 'cat')
  ],
  volumes: [
    hostPathVolume(hostPath: '/var/run/docker.sock', mountPath: '/var/run/docker.sock')
  ]) {
  node('build-agent-go') {
    stage('Download dependencies') {
      steps {
        container('build-agent-go') {
          sh 'apk add --no-cache git'
          sh 'cd "$GOPATH/src/a10bridge";dep ensure'
        }
      }
    }
    stage('Test application') {
      steps {
        container('build-agent-go') {
          sh 'cd "$GOPATH/src/a10bridge"; go build -v ./...'
          sh 'cd "$GOPATH/src/a10bridge"; go test -v ./...'
        }
      }
    }
    stage('Build Image') {
      input {
        message "Should we build and push docker image to registry?"
        ok "Yes, we should."
        parameters {
          string(name: 'IMAGE_TAG', defaultValue: 'v0.0', description: 'Version of the image')
        }
      }
      steps {
        container('build-agent-go') {
          sh "cd $GOPATH/src/a10bridge; CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' ."
          sh 'docker build -t registry.pulsepoint.com/a10bridge:${IMAGE_TAG} "$GOPATH/src/a10bridge"'
          sh 'docker push registry.pulsepoint.com/a10bridge:${IMAGE_TAG}'
        }
      }
    }
  }
}