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
