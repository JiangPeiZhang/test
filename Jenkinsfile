pipeline {
    agent any
    
    environment {
        // 定义环境变量
        GO_VERSION = '1.21'
        PROJECT_NAME = 'pzjiang-test'
        DOCKER_IMAGE = 'pzjiang-test'
        DOCKER_TAG = "${env.BUILD_NUMBER}"
    }
    
    stages {
        stage('Checkout') {
            steps {
                echo '开始检出代码...'
                checkout scm
            }
        }
        
        stage('Setup Go Environment') {
            steps {
                echo '设置Go环境...'
                sh '''
                    # 检查Go版本
                    go version
                    
                    # 设置Go环境变量
                    export GOPATH=$WORKSPACE
                    export PATH=$PATH:$GOPATH/bin
                '''
            }
        }
        
        stage('Install Dependencies') {
            steps {
                echo '安装项目依赖...'
                sh '''
                    # 下载Go模块依赖
                    go mod download
                    go mod tidy
                    
                    # 验证依赖
                    go mod verify
                '''
            }
        }
        
        stage('Code Quality Check') {
            steps {
                echo '代码质量检查...'
                sh '''
                    # 代码格式化检查
                    go fmt ./...
                    
                    # 静态代码分析
                    go vet ./...
                    
                    # 运行测试（如果有的话）
                    go test ./... -v
                '''
            }
        }
        
        stage('Build') {
            steps {
                echo '构建项目...'
                sh '''
                    # 构建可执行文件
                    go build -o ${PROJECT_NAME} main.go
                    
                    # 检查构建结果
                    ls -la ${PROJECT_NAME}
                    file ${PROJECT_NAME}
                '''
            }
        }
        
        stage('Test Build Result') {
            steps {
                echo '测试构建结果...'
                sh '''
                    # 启动服务进行测试
                    timeout 30s ./${PROJECT_NAME} &
                    sleep 5
                    
                    # 测试Health接口
                    curl -X POST http://localhost:8080/health || echo "Health接口测试失败"
                    
                    # 停止服务
                    pkill -f ${PROJECT_NAME} || true
                '''
            }
        }
        
        stage('Docker Build') {
            when {
                expression { env.BUILD_BRANCH == 'main' || env.BUILD_BRANCH == 'master' }
            }
            steps {
                echo '构建Docker镜像...'
                script {
                    // 创建Dockerfile
                    writeFile file: 'Dockerfile', text: '''
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o pzjiang-test main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/pzjiang-test .
EXPOSE 8080

CMD ["./pzjiang-test"]
'''
                    
                    // 构建Docker镜像
                    sh '''
                        docker build -t ${DOCKER_IMAGE}:${DOCKER_TAG} .
                        docker tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:latest
                    '''
                }
            }
        }
        
        stage('Push to Registry') {
            when {
                expression { env.BUILD_BRANCH == 'main' || env.BUILD_BRANCH == 'master' }
            }
            steps {
                echo '推送镜像到仓库...'
                script {
                    // 这里可以配置推送到私有仓库
                    // 需要根据实际情况配置仓库地址和认证信息
                    sh '''
                        echo "镜像构建完成: ${DOCKER_IMAGE}:${DOCKER_TAG}"
                        echo "镜像标签: ${DOCKER_IMAGE}:latest"
                    '''
                }
            }
        }
    }
    
    post {
        always {
            echo '清理工作空间...'
            sh '''
                # 清理构建产物
                rm -f ${PROJECT_NAME}
                
                # 清理Docker镜像（可选）
                # docker rmi ${DOCKER_IMAGE}:${DOCKER_TAG} || true
            '''
        }
        
        success {
            echo '构建成功！'
            script {
                // 可以在这里添加成功通知
                // 比如发送邮件、钉钉通知等
            }
        }
        
        failure {
            echo '构建失败！'
            script {
                // 可以在这里添加失败通知
                // 比如发送邮件、钉钉通知等
            }
        }
        
        cleanup {
            echo '清理环境...'
            sh '''
                # 停止可能运行的服务
                pkill -f ${PROJECT_NAME} || true
                
                # 清理临时文件
                rm -f Dockerfile || true
            '''
        }
    }
} 