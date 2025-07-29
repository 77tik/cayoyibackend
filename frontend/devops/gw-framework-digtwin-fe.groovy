#!groovy
import org.jenkinsci.plugins.plaincredentials.StringCredentials
repoURL = "ssh://vcssh@phabricator.intern.yuansuan.cn/diffusion/62/cae.git"
repoPath = "/opt/jenkins/workspace/${projectName}-fe"
Workdir = "/opt/jenkins/workspace/${projectName}-fe/${projectName}/frontend"
pipeline {
    agent {label 'dtnode'}
    stages {
        stage('clean') {
            steps {
                echo "当前环境是:${envs}"
                echo "当前版本是:${BUILD_NUMBER}"
                echo "版本路径是:http://1.117.192.82:8666/harbor/projects/54/repositories/westlake-urban-flooding-fe"
                sh"""
                    rm -rf ${repoPath}
                    mkdir -p ${repoPath}
                """
            }
        }
        stage('clone') {
            steps {
                sh"""
                    git clone --depth 1 ${repoURL} ${repoPath}
                """
            }
        }
        stage('build') {
            steps {
                echo "version:${BUILD_NUMBER}"
                echo 'build start.'
                sh"""
                cd ${Workdir}
                docker login -u admin -p yskj2407 1.117.192.82:8666
                """
                script {
                    if( "${envs}" == "arm64" ) {
                        sh"""
                        cd ${Workdir}
                        docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
                        docker build -t 1.117.192.82:8666/${projectName}/${projectName}-fe:${BUILD_NUMBER}_arm64 -f Dockerfile.arm64 .
                        docker push 1.117.192.82:8666/${projectName}/${projectName}-fe:${BUILD_NUMBER}_arm64
                        """
                    } else {
                        sh"""
                        cd ${Workdir}
                        docker build -t 1.117.192.82:8666/${projectName}/${projectName}-fe:${BUILD_NUMBER} -f Dockerfile.x86 .
                        docker push 1.117.192.82:8666/${projectName}/${projectName}-fe:${BUILD_NUMBER}
                        """
                    }
                }
                sh"""
                
                """
                script {
                    if("${envs}" == "dev") {
                        sh "ssh dev \"mkdir -p /opt/${projectName}/front\""
                    }else{
                        sh "mkdir -p /opt/${projectName}/front"
                    }
                }
                echo 'build end.'
            }
        }
        stage('deploy') {
            steps {
                echo 'deploy start.'
                script {
                    if( "${envs}" == "dev" ) {
                        sh "scp -r ${Workdir}/nginx/fe_dev.conf dev:/opt/${projectName}/front/fe.conf"
                        sh "scp -r ${Workdir}/nginx/nginx.conf dev:/opt/${projectName}/front/nginx.conf"
                    }else{
                        sh "scp -r ${Workdir}/nginx/fe_stage.conf /opt/${projectName}/front/fe.conf"
                        sh "scp -r ${Workdir}/nginx/nginx.conf /opt/${projectName}/front/nginx.conf"
                    }
                }
                script {
                    if( "${envs}" == "dev" ) {
                        sh "ssh dev \"docker rm -f ${projectName}-fe || true\""
                        sh "ssh dev \"docker image prune -a -f || true\""
                        sh "ssh dev \"docker login -u admin -p yskj2407 1.117.192.82:8666\""
                        sh "ssh dev \"docker run -itd -p ${port}:${port} -v /opt/${projectName}/front/nginx.conf:/etc/nginx/nginx.conf -v /opt/${projectName}/front/fe.conf:/etc/nginx/conf.d/default.conf --restart=always  --name ${projectName}-fe 1.117.192.82:8666/${projectName}/${projectName}-fe:${BUILD_NUMBER}\""
                        echo 'deploy end.'
                        echo "${BUILD_NUMBER}"
                        echo "http://10.0.4.128:${port}/"
                    }else if( "${envs}" == "stage" ){
                        sh "docker rm -f ${projectName}-fe || true"
                        sh "docker login -u admin -p yskj2407 1.117.192.82:8666"
                        sh "docker pull 1.117.192.82:8666/${projectName}/${projectName}-fe:${BUILD_NUMBER}"
                        sh "docker run -itd -p ${port}:${port} -v /opt/${projectName}/front/nginx.conf:/etc/nginx/nginx.conf -v /opt/${projectName}/front/fe.conf:/etc/nginx/conf.d/default.conf -v /opt/${projectName}/front/basic_auth:/etc/nginx/basic_auth --restart=always --name ${projectName}-fe 1.117.192.82:8666/${projectName}/${projectName}-fe:${BUILD_NUMBER}"
                        echo 'deploy end.'
                        echo "${BUILD_NUMBER}"
                        echo "http://10.0.4.129:${port}/"
                    }
                } 
            }
        }
    }
}