pipeline {
  agent {
    kubernetes {
      label 'build-agent-go'
      containerTemplate {
        name 'build-agent-go'
        image 'registry.pulsepoint.com/build-agent-go:0.1'
        ttyEnabled true
        command 'cat'
      }
    }
  }
  environment {
    GOPATH = "${WORKSPACE}"
  }
  stages {
    stage('Download dependencies') {
      steps {
        container('build-agent-go') {
          sh 'echo "$PATH"'
          sh 'ls -lah /usr/local/bin'
          sh 'ls -lah /usr/bin'
          sh 'cd "$GOPATH/src/a10bridge";dep ensure'
        }
      }
    }
    stage('Test application') {
      steps {
        container('build-agent-go') {
          sh 'go build "$GOPATH/src/a10bridge/..."'
          sh 'go test "$GOPATH/src/a10bridge/..."'
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
          sh "CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' $GOPATH/src/a10bridge"
          sh 'docker build -t registry.pulsepoint.com/a10bridge:${IMAGE_TAG} "$GOPATH/src/a10bridge"'
          sh 'docker push registry.pulsepoint.com/a10bridge:${IMAGE_TAG}'
        }
      }
    }
  }
}