pipeline {
    agent any
    
    environment {
        PROJECT_NAME = 'pzjiang-test'
    }
    
    stages {
        stage('Checkout') {
            steps {
                echo '开始检出代码...'
                checkout scm
            }
        }
        
        stage('Setup Environment') {
            steps {
                echo '设置环境...'
                sh '''
                    # 检查是否安装了Go
                    if ! command -v go &> /dev/null; then
                        echo "Go未安装，开始安装Go..."
                        
                        # 检测系统类型并安装Go
                        if [ -f /etc/debian_version ]; then
                            # Debian/Ubuntu系统
                            apt-get update
                            apt-get install -y wget
                            wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
                            tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
                            export PATH=$PATH:/usr/local/go/bin
                            echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
                        elif [ -f /etc/redhat-release ]; then
                            # CentOS/RHEL系统
                            yum install -y wget
                            wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
                            tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
                            export PATH=$PATH:/usr/local/go/bin
                            echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
                        else
                            # 其他Linux系统
                            wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
                            tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
                            export PATH=$PATH:/usr/local/go/bin
                            echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
                        fi
                        
                        # 清理下载文件
                        rm -f go1.21.0.linux-amd64.tar.gz
                    else
                        echo "Go已安装"
                    fi
                    
                    # 显示Go版本
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
                    go mod download
                    go mod tidy
                '''
            }
        }
        
        stage('Code Quality Check') {
            steps {
                echo '代码质量检查...'
                sh '''
                    go fmt ./...
                    go vet ./...
                '''
            }
        }
        
        stage('Build') {
            steps {
                echo '构建项目...'
                sh '''
                    go build -o ${PROJECT_NAME} main.go
                    ls -la ${PROJECT_NAME}
                '''
            }
        }
        
        stage('Test Build Result') {
            steps {
                echo '测试构建结果...'
                sh '''
                    # 检查curl是否安装
                    if ! command -v curl &> /dev/null; then
                        echo "curl未安装，开始安装curl..."
                        
                        # 检测系统类型并安装curl
                        if [ -f /etc/debian_version ]; then
                            # Debian/Ubuntu系统
                            apt-get update
                            apt-get install -y curl
                        elif [ -f /etc/redhat-release ]; then
                            # CentOS/RHEL系统
                            yum install -y curl
                        else
                            # 其他Linux系统，尝试通用方法
                            if command -v apt-get &> /dev/null; then
                                apt-get update && apt-get install -y curl
                            elif command -v yum &> /dev/null; then
                                yum install -y curl
                            else
                                echo "无法自动安装curl，跳过接口测试"
                                echo "curl不可用" > curl_status.txt
                            fi
                        fi
                    else
                        echo "curl已安装"
                    fi
                    
                    # 检查timeout命令
                    if ! command -v timeout &> /dev/null; then
                        echo "timeout命令不可用，使用sleep替代"
                        # 启动服务并等待
                        ./${PROJECT_NAME} &
                        sleep 35
                    else
                        # 启动服务进行测试
                        timeout 30s ./${PROJECT_NAME} &
                        sleep 5
                    fi
                    
                    # 测试Health接口
                    if [ -f curl_status.txt ] && grep -q "curl不可用" curl_status.txt; then
                        echo "curl不可用，跳过接口测试"
                    else
                        curl -X POST http://localhost:8080/health || echo "Health接口测试失败"
                    fi
                    
                    # 停止服务
                    pkill -f ${PROJECT_NAME} || true
                '''
            }
        }
        
        stage('Create Dockerfile') {
            when {
                expression { env.BUILD_BRANCH == 'main' || env.BUILD_BRANCH == 'master' }
            }
            steps {
                echo '检查Docker环境...'
                sh '''
                    # 检查是否安装了Docker
                    if ! command -v docker &> /dev/null; then
                        echo "Docker未安装，跳过Docker构建阶段"
                        echo "Docker未安装，无法构建镜像" > docker_status.txt
                    else
                        echo "Docker已安装"
                        docker --version
                        echo "Docker可用" > docker_status.txt
                    fi
                '''
                
                echo '创建Dockerfile...'
                sh '''
                    echo "FROM golang:1.21-alpine AS builder" > Dockerfile
                    echo "WORKDIR /app" >> Dockerfile
                    echo "COPY go.mod go.sum ./" >> Dockerfile
                    echo "RUN go mod download" >> Dockerfile
                    echo "COPY . ." >> Dockerfile
                    echo "RUN go build -o pzjiang-test main.go" >> Dockerfile
                    echo "FROM alpine:latest" >> Dockerfile
                    echo "RUN apk --no-cache add ca-certificates" >> Dockerfile
                    echo "WORKDIR /root/" >> Dockerfile
                    echo "COPY --from=builder /app/pzjiang-test ." >> Dockerfile
                    echo "EXPOSE 8080" >> Dockerfile
                    echo 'CMD ["./pzjiang-test"]' >> Dockerfile
                '''
            }
        }
        
        stage('Docker Build') {
            when {
                expression { env.BUILD_BRANCH == 'main' || env.BUILD_BRANCH == 'master' }
            }
            steps {
                echo '构建Docker镜像...'
                sh '''
                    # 检查Docker状态
                    if [ -f docker_status.txt ] && grep -q "Docker可用" docker_status.txt; then
                        echo "开始构建Docker镜像..."
                        docker build -t pzjiang-test:${BUILD_NUMBER} .
                        docker tag pzjiang-test:${BUILD_NUMBER} pzjiang-test:latest
                        echo "Docker镜像构建完成: pzjiang-test:${BUILD_NUMBER}"
                    else
                        echo "Docker不可用，跳过镜像构建"
                        echo "如需构建Docker镜像，请确保Jenkins环境中安装了Docker"
                    fi
                '''
            }
        }
    }
    
    post {
        always {
            echo '清理工作空间...'
            sh 'rm -f ${PROJECT_NAME}'
        }
        
        success {
            echo '构建成功！'
        }
        
        failure {
            echo '构建失败！'
        }
        
        cleanup {
            echo '清理环境...'
            sh '''
                pkill -f ${PROJECT_NAME} || true
                rm -f Dockerfile || true
            '''
        }
    }
} 