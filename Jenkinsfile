pipeline {
    agent {
        kubernetes {
            yaml '''
                apiVersion: v1
                kind: Pod
                spec:
                  securityContext:
                    runAsUser: 0
                    runAsGroup: 0
                    fsGroup: 0
                  containers:
                  - name: go
                    image: golang:1.21
                    command:
                    - cat
                    tty: true
                    securityContext:
                      runAsUser: 0
                      runAsGroup: 0
                      privileged: true
                    volumeMounts:
                    - name: workspace
                      mountPath: /home/jenkins/agent
                    - name: docker-sock
                      mountPath: /var/run/docker.sock
                  - name: docker
                    image: docker:20.10-dind
                    command:
                    - dockerd
                    securityContext:
                      privileged: true
                    volumeMounts:
                    - name: docker-sock
                      mountPath: /var/run/docker.sock
                  volumes:
                  - name: workspace
                    emptyDir: {}
                  - name: docker-sock
                    emptyDir: {}
            '''
        }
    }
    
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
                    # 显示Go版本
                    go version
                    
                    # 设置Go环境变量
                    export GOPATH=$WORKSPACE
                    export PATH=$PATH:$GOPATH/bin
                    
                    # 安装curl（如果需要）
                    if ! command -v curl &> /dev/null; then
                        echo "安装curl..."
                        apt-get update && apt-get install -y curl
                    fi
                    
                    # 等待Docker启动
                    echo "等待Docker启动..."
                    sleep 10
                    docker --version || echo "Docker可能还在启动中"
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
        
        stage('Create Dockerfile') {
            steps {
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
            steps {
                echo '构建Docker镜像...'
                sh '''
                    echo "开始构建Docker镜像..."
                    docker build -t pzjiang-test:${BUILD_NUMBER} .
                    docker tag pzjiang-test:${BUILD_NUMBER} pzjiang-test:latest
                    echo "Docker镜像构建完成: pzjiang-test:${BUILD_NUMBER}"
                    
                    # 显示镜像信息
                    docker images | grep pzjiang-test
                    
                    # 保存镜像为tar文件
                    docker save pzjiang-test:${BUILD_NUMBER} > pzjiang-test-${BUILD_NUMBER}.tar
                    docker save pzjiang-test:latest > pzjiang-test-latest.tar
                    echo "镜像已保存为tar文件"
                '''
            }
        }
        
        stage('Archive Artifacts') {
            steps {
                echo '归档制品...'
                sh '''
                    # 列出所有制品文件
                    ls -la *.tar *.go *.mod || true
                    ls -la ${PROJECT_NAME} || true
                '''
                
                // 归档制品到KubeSphere制品库
                archiveArtifacts artifacts: '*.tar,*.go,*.mod,${PROJECT_NAME}', fingerprint: true
                
                echo '制品归档完成'
            }
        }
    }
} 