# 测试KubeSphere DevOps

这是一个使用Go和Gin框架构建的基本HTTP服务。

## 功能特性

- 基于Gin框架的HTTP服务
- Health健康检查接口（POST请求）
- 返回标准JSON格式响应

## 项目结构

```
pzjiang-test/
├── main.go          # 主程序文件
├── go.mod           # Go模块文件
├── Jenkinsfile      # Jenkins CI/CD流水线配置
└── README.md        # 项目说明
```

## 安装和运行

1. 确保已安装Go 1.21或更高版本
2. 下载依赖：
   ```bash
   go mod tidy
   ```
3. 运行服务：
   ```bash
   go run main.go
   ```
   或者编译后运行：
   ```bash
   go build -o server main.go
   ./server
   ```

## API接口

### Health接口

- **URL**: `/health`
- **方法**: POST
- **响应格式**: JSON
- **响应示例**:
  ```json
  {
    "code": 0,
    "msg": ""
  }
  ```

## 测试

使用curl测试Health接口：

```bash
curl -X POST http://localhost:8080/health
```

预期响应：
```json
{"code":0,"msg":""}
```

## 服务信息

- 服务端口：8080
- 服务地址：http://localhost:8080

## CI/CD 流水线

项目包含Jenkinsfile，支持在KubeSphere中进行自动化构建和部署。

### 流水线阶段

1. **Checkout** - 检出代码
2. **Setup Environment** - 设置Go环境（自动检查并安装Go）
3. **Install Dependencies** - 安装项目依赖
4. **Code Quality Check** - 代码质量检查（格式化、静态分析）
5. **Build** - 构建可执行文件
6. **Test Build Result** - 测试构建结果（启动服务并测试Health接口，自动检查并安装curl）
7. **Create Dockerfile** - 创建Dockerfile
8. **Docker Build** - 构建Docker镜像

### 环境要求

- Jenkins环境需要安装Go 1.21+（流水线会自动检查并安装）
- Docker（用于构建镜像，可选，流水线会自动检查）
- curl（用于测试接口，流水线会自动检查并安装）
- timeout（用于进程控制，流水线会自动检查）

**注意：** Jenkinsfile已配置为使用Kubernetes agent，具有以下优势：
- 使用预装的Go 1.21镜像，无需安装Go
- 以root权限运行，可以安装额外工具
- 在K8s Pod中运行，环境隔离
- 自动安装curl等必要工具
- 内置Docker-in-Docker，支持Docker镜像构建

### 分支策略

- 所有分支都会执行完整的流水线，包括构建、测试和Docker镜像构建
