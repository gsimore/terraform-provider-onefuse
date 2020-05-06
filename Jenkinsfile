/*
    A pipeline to build the Fuse Terraform Provider

    TODO:
    * 

    Jenkins Prerequisites:
    * 

    Agent Prerequisites:
    * Agent built from ...
    * Agent containing label 'go'
*/
pipeline {
    agent {
      node { label 'go' }
    }
    environment {
      TERRAFORM_PROVIDER_DIR = "./terraform-provider-fuse"
      TERRAFORM_PROVIDER_BIN_NAME = "terraform-provider-fuse"
    }
    stages {
        stage('Checkout') {
            steps {
                dir("${env.TERRAFORM_PROVIDER_DIR}") {
                    git credentialsId: "github", url: "https://github.com/CloudBoltSoftware/terraform-provider-fuse.git", branch: 'develop', poll: false
                }
            }
        }

        stage('Build') {
            environment {
                  // Set the Go environment variables to be relative to the workspace directory
                  GOPATH = "${env.WORKSPACE}/go"
                  GOCACHE = "${env.WORKSPACE}/go/.cache"
            }
            steps {
                dir("${env.TERRAFORM_PROVIDER_DIR}") {
                    sh "go build -o ${env.TERRAFORM_PROVIDER_BIN_NAME}"
                }
            }
        }
    }
}