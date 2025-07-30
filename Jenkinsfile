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
                  - name: agent
                    image: ubuntu:20.04
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
                  volumes:
                  - name: workspace
                    emptyDir: {}
            '''
        }
    }
    
    environment {
        PROJECT_NAME = 'pzjiang-test'
    }
    
    stages {
        stage('Build and Test') {
            steps {
                echo '开始构建和测试流程...'
                
                sh '''
                    echo "=== Step 1: 检出代码 ==="
                    # 检出代码（checkout scm 在Jenkins中会自动执行）
                    
                    echo "=== Step 2: 设置环境 ==="
                    # 检查并安装Go
                    if ! command -v go &> /dev/null; then
                        echo "Go未安装，开始安装Go..."
                        apt-get update
                        apt-get install -y wget
                        wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
                        tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
                        export PATH=$PATH:/usr/local/go/bin
                        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
                        rm -f go1.21.0.linux-amd64.tar.gz
                        echo "Go安装完成"
                    else
                        echo "Go已安装"
                    fi
                    
                    # 显示Go版本
                    go version
                    
                    # 设置Go环境变量 - 不使用GOPATH
                    export GO111MODULE=on
                    export GOCACHE=/tmp/go-cache
                    
                    # 安装curl（如果需要）
                    if ! command -v curl &> /dev/null; then
                        echo "安装curl..."
                        apt-get update && apt-get install -y curl
                    fi
                    
                    # 检查并安装Docker
                    if ! command -v docker &> /dev/null; then
                        echo "Docker未安装，开始安装Docker..."
                        apt-get update
                        apt-get install -y docker.io
                        systemctl start docker || service docker start || true
                        echo "Docker安装完成"
                    else
                        echo "Docker已安装"
                    fi
                    
                    # 显示Docker版本
                    docker --version || echo "Docker可能有问题"
                    
                    echo "=== Step 3: 安装项目依赖 ==="
                    # 确保在正确的目录中
                    pwd
                    ls -la
                    
                    # 清理并重新初始化Go模块
                    go mod download
                    go mod tidy
                    
                    echo "=== Step 4: 代码质量检查 ==="
                    go fmt ./...
                    go vet ./...
                    
                    echo "=== Step 5: 构建项目 ==="
                    go build -o ${PROJECT_NAME} main.go
                    ls -la ${PROJECT_NAME}
                    
                    echo "=== Step 6: 测试构建结果 ==="
                    # 启动服务进行测试
                    timeout 30s ./${PROJECT_NAME} &
                    sleep 5
                    
                    # 测试Health接口
                    curl -X POST http://localhost:8080/health || echo "Health接口测试失败"
                    
                    # 停止服务
                    pkill -f ${PROJECT_NAME} || true
                    
                    echo "=== Step 7: 创建Dockerfile ==="
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
                    
                    echo "=== Step 8: 构建Docker镜像 ==="
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
                    
                    echo "=== Step 9: 归档制品 ==="
                    # 列出所有制品文件
                    ls -la *.tar *.go *.mod || true
                    ls -la ${PROJECT_NAME} || true
                    
                    echo "构建和测试流程完成！"
                '''
                
                // 归档制品到KubeSphere制品库
                archiveArtifacts artifacts: '*.tar,*.go,*.mod,${PROJECT_NAME}', fingerprint: true
                
                echo '制品归档完成！'
            }
        }
    }
}