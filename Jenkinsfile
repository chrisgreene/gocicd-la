pipeline {
    agent any
    environment {
        PATH = "/usr/local/go/bin:$PATH"
        DOCKER_IMAGE_NAME = "chrisgreene/go-cicd-kubernetes"
    }
    stages {
        stage('Build') {
            steps {
                echo 'Running build automation'
                sh 'go test -v'
            }
        }
        stage('Build Docker Image') {
            when {
                branch 'master'
            }
            steps {
                script {
                    app = docker.build(DOCKER_IMAGE_NAME)
                    app.withRun("-d -p 8181:8181") { c ->
                        sh 'curl localhost:8181'
                    }
                }
            }
        }
        stage('Push Docker Image') {
            when {
                branch 'master'
            }
            steps {
                script {
                    docker.withRegistry('https://registry.hub.docker.com', 'docker_hub') {
                        app.push("${env.BUILD_NUMBER}")
                        app.push("latest")
                    }
                }
            }
        }
        stage('DeployToProduction') {
            when {
                branch 'master'
            }
            steps {
                //input 'Deploy to Production?'
                milestone(1)
                withCredentials([usernamePassword(credentialsId: 'kube_master', usernameVariable: 'USERNAME', passwordVariable: 'USERPASS')]) {
                    script {
                        echo sh(script: 'ls -al', returnStdout: true)
                        sh "sshpass -p '$USERPASS' -v ssh -o StrictHostKeyChecking=no $USERNAME@$kube_master_ip \"kubectl delete all -l app=gocicd\""
                    }
                }
                kubernetesDeploy(
                  kubeconfigId: 'kubeconfig',
                  configs: 'kubernetes.yaml',
                  enableConfigSubstitution: true
                )
            }
        }
        stage('Get Service IP') {
            when {
                branch 'master'
            }
            steps {
                //milestone(1)
                retry(10) {
                    withCredentials([usernamePassword(credentialsId: 'kube_master', usernameVariable: 'USERNAME', passwordVariable: 'USERPASS')]) {
                        script {
                            sh "sshpass -p '$USERPASS' -v ssh -o StrictHostKeyChecking=no $USERNAME@$kube_master_ip \"kubectl get all -l app=gocicd -o wide\""
                            sh "sshpass -p '$USERPASS' -v ssh -o StrictHostKeyChecking=no $USERNAME@$kube_master_ip \"kubectl get svc gocicd --output=jsonpath={'.spec.clusterIP'}\""
                            try {
                             //echo sh(script: "sshpass -p '$USERPASS' -v ssh -o StrictHostKeyChecking=no $USERNAME@$kube-master-ip ls -la", returnStdout: true)
                             //sh "sshpass -p '$USERPASS' -v ssh -o StrictHostKeyChecking=no $USERNAME@$kube-master-ip \"/usr/bin/touch /tmp/jenkins\""
                            } catch (err) {
                             echo: 'caught error: $err'
                            }
                        }
                    }
                }
            }
        }
    }
}
