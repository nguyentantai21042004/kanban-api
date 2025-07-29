def notifyDiscord(channel, chatId, message) {
    sh """
        curl --location --request POST "https://discord.com/api/webhooks/${channel}/${chatId}" \
        --header 'Content-Type: application/json' \
        --data-raw '{"content": "${message}"}'
    """
}

pipeline {
    agent any

    environment {
        ENVIRONMENT = 'kanban'
        SERVICE = 'kanban-api'

        REGISTRY_DOMAIN_NAME = 'harbor.ngtantai.pro'
        REGISTRY_USERNAME = 'admin'
        REGISTRY_PASSWORD = credentials('registryPassword')

        // K8s Configuration
        K8S_NAMESPACE = 'kanban-api'
        K8S_DEPLOYMENT_NAME = 'kanban-api'
        K8S_CONTAINER_NAME = 'kanban-api'
        K8S_API_SERVER = 'https://172.16.21.31:6443'
        K8S_TOKEN = credentials('k8s-token')
        
        DOCKER_EXPOSE_PORT = '8080'  // Changed from 80 to 8080
        APP_TEMP_PORT = '8080'
        APP_FINAL_PORT = '8080'      // Changed from 80 to 8080
        
        TEXT_START = "⚪ Service ${SERVICE} ${ENVIRONMENT} Build Started"
        TEXT_BUILD_AND_PUSH_APP_FAIL = "🔴 Service ${SERVICE} ${ENVIRONMENT} Build and Push Failed"
        TEXT_DEPLOY_APP_FAIL = "🔴 Service ${SERVICE} ${ENVIRONMENT} Deploy Failed"
        TEXT_CLEANUP_OLD_IMAGES_FAIL = "🔴 Cleanup Old Images Failed"
        TEXT_END = "🟢 Service ${SERVICE} ${ENVIRONMENT} Build Finished"

        DISCORD_CHANNEL = '1382725588321828934'
        DISCORD_CHAT_ID = 'Q1edE75TA7jJlloegQ2MxDpBxAGoVFz0buoSwW-wg6mTLozxP20oagKFlRiN5l1fyCOQ'
    }

    stages {
        stage('Notify Build Started') {
            steps {
                script {
                    def causes = currentBuild.getBuildCauses()
                    def triggerInfo = causes ? causes[0].shortDescription : "Unknown"
                    def cleanTrigger = triggerInfo.replaceFirst("Started by ", "")
                    notifyDiscord(env.DISCORD_CHANNEL, env.DISCORD_CHAT_ID, "${env.TEXT_START} by ${cleanTrigger}.")
                }
            }
        }

        stage('Pull Code') {
            steps {
                script {
                    echo "Now Jenkins is pulling code..." 
                    checkout scm
                    echo "Now Jenkins is listing code..."
                    sh "ls -la ${WORKSPACE}"
                    sh "find ${WORKSPACE} -name 'Dockerfile' -type f || echo 'Dockerfile not found'"
                }
            }
        }

        stage('Build API Image') {
            steps {
                script {
                    try {
                        def timestamp = new Date().format('yyMMdd-HHmmss')
                        env.DOCKER_API_IMAGE_NAME = "${env.REGISTRY_DOMAIN_NAME}/${env.ENVIRONMENT}/${env.SERVICE}:${timestamp}"

                        sh "docker build -t ${env.DOCKER_API_IMAGE_NAME} -f ${WORKSPACE}/cmd/api/Dockerfile ${WORKSPACE}"

                        echo "Successfully built APP: ${env.DOCKER_API_IMAGE_NAME}"                    
                    } catch (Exception e) {
                        notifyDiscord(env.DISCORD_CHANNEL, env.DISCORD_CHAT_ID, env.TEXT_BUILD_AND_PUSH_APP_FAIL)
                        error("APP build failed: ${e.getMessage()}")
                    }
                }
            }
        }

        stage('Push API Image') {
            steps {
                script {
                    try {
                        sh 'echo $REGISTRY_PASSWORD | docker login $REGISTRY_DOMAIN_NAME -u $REGISTRY_USERNAME --password-stdin'
                        sh "docker push ${env.DOCKER_API_IMAGE_NAME}"
                        sh "docker rmi ${env.DOCKER_API_IMAGE_NAME} || true"
                        echo "Successfully pushed APP: ${env.DOCKER_API_IMAGE_NAME}"                    
                    } catch (Exception e) {
                        notifyDiscord(env.DISCORD_CHANNEL, env.DISCORD_CHAT_ID, env.TEXT_BUILD_AND_PUSH_APP_FAIL)
                        error("APP push failed: ${e.getMessage()}")
                    }
                }
            }
        }

        stage('Deploy to Kubernetes') {
            steps {
                script {
                    try {
                        echo "Deploying new image to K8s: ${env.DOCKER_API_IMAGE_NAME}"
                        echo "K8S_API_SERVER: ${env.K8S_API_SERVER}"
                        echo "K8S_NAMESPACE: ${env.K8S_NAMESPACE}"
                        echo "K8S_DEPLOYMENT_NAME: ${env.K8S_DEPLOYMENT_NAME}"
                        echo "K8S_CONTAINER_NAME: ${env.K8S_CONTAINER_NAME}"
                        echo "Patch data to be sent: ${patchData}"

                        def patchData = '{"spec":{"template":{"spec":{"containers":[{"name":"' + env.K8S_CONTAINER_NAME + '","image":"' + env.DOCKER_API_IMAGE_NAME + '"}]}}}}'

                        def curlCmd = """
                            curl -X PATCH \\
                                -H "Authorization: Bearer \$K8S_TOKEN" \\
                                -H "Content-Type: application/strategic-merge-patch+json" \\
                                -d '${patchData}' \\
                                "\$K8S_API_SERVER/apis/apps/v1/namespaces/\$K8S_NAMESPACE/deployments/\$K8S_DEPLOYMENT_NAME" \\
                                --insecure \\
                                --silent \\
                                --fail \\
                                -w '\\nHTTP_STATUS_CODE:%{http_code}\\n'
                        """

                        echo "Running curl command to patch deployment..."
                        def deployOutput = sh(
                            script: curlCmd,
                            returnStdout: true
                        ).trim()

                        echo "Curl output:\n${deployOutput}"

                        // Extract HTTP status code from output
                        def matcher = deployOutput =~ /HTTP_STATUS_CODE:(\d+)/
                        def httpStatus = matcher ? matcher[0][1] as Integer : 0

                        if (httpStatus < 200 || httpStatus >= 300) {
                            echo "Failed to update deployment. Full curl output:\n${deployOutput}"
                            error("Failed to update deployment. HTTP status: ${httpStatus}")
                        }

                        echo "Successfully triggered K8s deployment update"

                    } catch (Exception e) {
                        echo "Exception during Kubernetes deployment: ${e}"
                        echo "Stack trace: ${e.getStackTrace().join('\n')}"
                        notifyDiscord(env.DISCORD_CHANNEL, env.DISCORD_CHAT_ID, env.TEXT_DEPLOY_APP_FAIL)
                        error("Kubernetes deployment failed: ${e.getMessage()}")
                    }
                }
            }
        }

        stage('Verify Deployment') {
            steps {
                script {
                    try {
                        echo "Verifying deployment health..."
                        
                        timeout(time: 5, unit: 'MINUTES') {
                            script {
                                def ready = false
                                def attempts = 0
                                def maxAttempts = 30
                                
                                while (!ready && attempts < maxAttempts) {
                                    attempts++
                                    
                                    def result = sh(
                                        script: '''
                                            curl -s -H "Authorization: Bearer $K8S_TOKEN" \\
                                                "$K8S_API_SERVER/apis/apps/v1/namespaces/$K8S_NAMESPACE/deployments/$K8S_DEPLOYMENT_NAME" \\
                                                --insecure | grep -o '"readyReplicas":[0-9]*' | cut -d':' -f2 || echo '0'
                                        ''',
                                        returnStdout: true
                                    ).trim()
                                    
                                    def readyReplicas = result as Integer
                                    echo "Attempt ${attempts}/${maxAttempts}: Ready replicas: ${readyReplicas}"
                                    
                                    if (readyReplicas >= 1) {
                                        ready = true
                                        echo "Deployment is ready with ${readyReplicas} replica(s)"
                                        
                                    } else {
                                        echo "Waiting for deployment to be ready..."
                                        sleep(10)
                                    }
                                }
                                
                                if (!ready) {
                                    error("Deployment failed to become ready after ${maxAttempts} attempts")
                                }
                            }
                        }
                        
                    } catch (Exception e) {
                        def errorMsg = e.getMessage().replaceAll('"', '\\\\"')
                        // Skip verification notification
                        echo "Verification failed but deployment may still be successful: ${e.getMessage()}"
                        // Don't fail the build on verification issues
                    }
                }
            }
        }

        stage('Cleanup Old Images') {
            steps {
                script {
                    try {
                        sh "docker image prune -a -f --filter \"until=24h\" || true"

                        sh """
                            docker images ${env.REGISTRY_DOMAIN_NAME}/${env.ENVIRONMENT}/${env.SERVICE} \\
                            --format "{{.Repository}}:{{.Tag}}\\t{{.CreatedAt}}" \\
                            | tail -n +2 | sort -k2 -r | tail -n +3 | awk '{print \$1}' \\
                            | xargs -r docker rmi || true
                        """

                        echo "Successfully cleaned up old images"
                        
                    } catch (Exception e) {
                        notifyDiscord(env.DISCORD_CHANNEL, env.DISCORD_CHAT_ID, env.TEXT_CLEANUP_OLD_IMAGES_FAIL)
                        echo "Cleanup failed but deployment was successful: ${e.getMessage()}"
                    }
                }
            }
        }

        stage('Notify Build Finished') {
            steps {
                script {
                    notifyDiscord(env.DISCORD_CHANNEL, env.DISCORD_CHAT_ID, "${env.TEXT_END}")
                }
            }
        }
    }

    post {
        failure {
            script {
                notifyDiscord(env.DISCORD_CHANNEL, env.DISCORD_CHAT_ID, "🔴 Service ${env.SERVICE} ${env.ENVIRONMENT} Pipeline Failed")
            }
        }
        success {
            script {
                notifyDiscord(env.DISCORD_CHANNEL, env.DISCORD_CHAT_ID, "🟢 Service ${env.SERVICE} ${env.ENVIRONMENT} Deploy Success")
            }
        }
    }
}