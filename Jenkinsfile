#!groovy
node('docker') {
    slackJobDescription = "job '${env.JOB_NAME} [${env.BUILD_NUMBER}]' (${env.BUILD_URL})"
    try {
        dockerRepoPerms = "test-permissions-${env.BUILD_TAG}"
        dockerRepoAppReg = "test-app-reg-${env.BUILD_TAG}"

        stage("Build") {
            checkout scm

            parallel (
                perms: { sh "docker build --pull --no-cache --rm -t ${dockerRepoPerms}"},
                appreg: { sh "docker build --pull --no-cache --rm -f Dockerfile.app-reg -t ${dockerRepoAppReg}"},
            )
        }

        dockerTestRunner = "test-${env.BUILD_TAG}"
        dockerPusher = "push-${env.BUILD_TAG}"
        try {
            stage("Test") {
                sh """docker run --rm --name ${dockerTestRunner} \\
                                 -w /go/src/github.com/cyverse-de/permissions \\
                                 --entrypoint 'gb' \\
                                 ${dockerRepoPerms} test"""
            }

            milestone 100
            stage("Docker Push") {
                service = readProperties file: 'service.properties'

                dockerPushRepoPerms = "${service.dockerUser}/permissions:${env.BRANCH_NAME}"
                dockerPushRepoAppReg = "${service.dockerUser}/permissions:${env.BRANCH_NAME}"

                lock("docker-push-perms-images") {
                    milestone 101

                    sh "docker tag ${dockerRepoPerms} ${dockerPushRepoPerms}"
                    sh "docker tag ${dockerRepoAppReg} ${dockerPushRepoAppReg}"

                    withCredentials([[$class: 'UsernamePasswordMultiBinding',
                                      credentialsId: 'jenkins-docker-credentials',
                                      passwordVariable: 'DOCKER_PASSWORD',
                                      usernameVariable: 'DOCKER_USERNAME']]) {
                        sh """docker run -e DOCKER_USERNAME -e DOCKER_PASSWORD \\
                                     -v /var/run/docker.sock:/var/run/docker.sock \\
                                     --rm --name ${dockerPusher} \\
                                     docker:\$(docker version --format '{{ .Server.Verwsion }}') \\
                                     sh -e -c \\
                              'docker login -u \"\$DOCKER_USERNAME\" -p \"$DOCKER_PASSWORD\" && \\
                               docker push ${dockerPushRepo} && \\
                               docker logout'"""
                    }
                }
            }
        } finally {
            sh returnStatus: true, script: "docker kill ${dockerTestRunner}"
            sh returnStatus: true, script: "docker rm ${dockerTestRunner}"

            sh returnStatus: true, script: "docker kill ${dockerPusher}"
            sh returnStatus: true, script: "docker rm ${dockerPusher}"

            sh returnStatus: true, script: "docker rmi ${dockerRepoPerms}"
            sh returnStatus: true, script: "docker rmi ${dockerRepoAppReg}"

            sh returnStatus: true, script: "docker rmi \$(docker images -qf 'dangling=true')"
        }
    } catch (InterruptedException e) {
        currentBuild.result = 'ABORTED'
        slackSend color: 'warning', message: "ABORTED: ${slackJobDescription}"
        throw e
    } catch (e) {
        currentBuild.result = 'FAILED'
        sh "echo ${e}"
        slackSend color: 'danger', message: "FAILED: ${slackJobDescription}"
        throw e
    }
}
