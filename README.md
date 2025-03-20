# License Module

这是一个Go语言实现的软件授权模块，用于生成和验证软件License。它可以帮助您保护您的软件，防止未授权使用。



## 功能特性

- **基于机器ID的授权**：使用机器硬件信息（MAC地址和CPU ID）生成唯一的机器标识符。
- **支持容器环境**：针对Docker容器环境提供了特殊处理，确保在容器中也能获得相对稳定的机器标识。
- **时效性控制**：支持设置License的有效期。
- **数字签名验证**：使用HMAC-SHA256进行签名，防止License被篡改。
- **防时间篡改**：通过存储上次运行时间，防止用户回退系统时间绕过过期检测。
- **Feature控制**：支持通过License控制可用的功能列表。



## 目录结构

```
.
├── cmd
│   └── license
│       ├── generate         # License生成工具
│       └── verify           # License验证工具
├── internal
│   └── license              # License核心功能
├── pkg
│   └── utils                # 工具函数，如机器ID获取
└── examples
    ├── app                  # 示例应用
    └── docker               # Docker示例
```



## 安装项目

```bash
go get github.com/chenwes/licensemodule
```



## 使用方法



### 查看机器ID

```bash
# 查看当前机器的ID
go run cmd/machine/id/main.go
```



### 生成License

```bash
# 查看当前机器的ID
go run cmd/license/generate/main.go --show-id

# 为当前机器生成一个有效期为30天的License
go run cmd/license/generate/main.go --days 30 --app "app-123" --out ./license.dat

# 为指定机器ID生成License
go run cmd/license/generate/main.go --machine "your-machine-id" --app "app-123" --days 365 --out ./license.dat

# 为容器环境生成License
go run cmd/license/generate/main.go --container --app "app-123" --days 30 --out ./license.dat

# 指定功能列表
go run cmd/license/generate/main.go --features "feature1,feature2,feature3" --app "app-123" --days 30 --out ./license.dat
```



### 验证License

```bash
# 验证许可证
go run cmd/license/verify/main.go --license ./license.dat --app "app-123" --timestamp ./timestamp.dat

# 在容器环境中验证许可证
go run cmd/license/verify/main.go --container --license ./license.dat --app "app-123" --timestamp ./timestamp.dat
```



### 运行服务器

```bash
# 查看当前机器的ID
go run cmd/api/main.go --port 8080
```

#### 请求API

```bash
curl --location --request POST 'http://localhost:8080/api/license/generate' \
--header 'Content-Type: application/json' \
--data-raw '{    
    "secret_key":"0aea8a18b07463ad5f5e3318db20d527c912c4ab9e7be28e94e8f486263a86fd/CF/WESCHAN",
    "machine_id":"0aea8a18b07463ad5f5e3318db20d527c912c4ab9e7be28e94e8f486263a86fd",
    "app_id": "metal-mes",
    "days":365,
    "features":[]
}'
```







### 在应用中集成

```go
package main

import (
    "log"
    "path/filepath"
    
    "github.com/chenwes/licensemodule/internal/license"
    "github.com/chenwes/licensemodule/pkg/utils"
)

func main() {
    // 获取机器ID
    machineID, err := utils.GetMachineID()  // 或 utils.GetContainerizedMachineID() 用于容器环境
    if err != nil {
        log.Fatalf("无法获取机器ID: %v", err)
    }
    
    // 验证License
    licenseFile := filepath.Join("config", "license.dat")
    timestampFile := filepath.Join("config", "timestamp.dat")
    
    err = license.VerifyAndUpdate(licenseFile, timestampFile, machineID)
    if err != nil {
        log.Fatalf("许可证验证失败: %v", err)
    }
    
    // 继续应用程序逻辑...
}
```

更多详细示例请参考 `examples/app/main.go`。





## 项目编译

要构建可以在其他机器上分发和运行的独立可执行文件，使用以下命令：

### 1. Machine ID Generator (For Clients)
```bash
# For Linux
GOOS=linux GOARCH=amd64 go build -o license-generator-linux cmd/machine/id/main.go

# For Windows
GOOS=windows GOARCH=amd64 go build -o machine-id.exe cmd/machine/id/main.go
go build -o machine-id.exe cmd/machine/id/main.go

# For macOS
GOOS=darwin GOARCH=amd64 go build -o license-generator-mac cmd/machine/id/main.go

# Usage
machine-id.exe [--container]
```

### For License Generator

```bash
# For the current platform
go build -o license-generator cmd/license/generate/main.go

# For Linux
GOOS=linux GOARCH=amd64 go build -o license-generator-linux cmd/license/generate/main.go

# For Windows
GOOS=windows GOARCH=amd64 go build -o license-generator.exe cmd/license/generate/main.go
go build -o license-generator.exe cmd/license/generate/main.go

# For macOS
GOOS=darwin GOARCH=amd64 go build -o license-generator-mac cmd/license/generate/main.go
```

### For License Verifier

```bash
# For the current platform
go build -o license-verifier cmd/license/verify/main.go

# For Linux
GOOS=linux GOARCH=amd64 go build -o license-verifier-linux cmd/license/verify/main.go

# For Windows
GOOS=windows GOARCH=amd64 go build -o license-verifier.exe cmd/license/verify/main.go
go build -o license-verifier.exe cmd/license/verify/main.go

# For macOS
GOOS=darwin GOARCH=amd64 go build -o license-verifier-mac cmd/license/verify/main.go
```

编译后的可执行文件可以复制到具有相同操作系统和体系结构的其他机器上并在其上运行。请注意，在运行可执行文件时仍将执行机器ID验证。



### Usage Examples

编译完成后，你可以这样使用可执行文件：

```bash
# Generate a license
./license-generator --machine "your-machine-id" --days 30 --out license.dat

# Verify a license
./license-verifier --license license.dat --timestamp timestamp.dat
```

注意：在分发可执行文件时，请确保在所有安装中保持相同的密钥（许可证包中的‘ SecretKey ’），以确保许可证验证正常工作。





## 运行服务器



### 生成服务器

```bash
# 在 Windows 命令行 (CMD) 中
set GOARCH=amd64
set GOOS=windows
go build -o machine-id.exe cmd/machine/id/main.go
go build -o license-verifier.exe cmd/license/verify/main.go
go build -o license-generator.exe cmd/license/generate/main.go

go build -o license-api.exe api/server.go


# 编译项目
go build -o license-api.exe api/server.go

# 调用API
license-api.exe --port 8080
```



### 调用

```bash
curl -X POST http://localhost:8080/api/license/generate \
  -H "Content-Type: application/json" \
  -d '{
    "machine_id": "your-machine-id",
    "secret_key": "your-secret-key-for-license-signature",
    "days": 30,
    "features": ["feature1", "feature2"]
  }' \
  --output license.dat
```








## Docker环境

在Docker容器中使用时，您需要：

1. 使用 `--container` 参数生成适用于容器环境的License
2. 将License文件挂载到容器中，或在构建镜像时包含它

示例Docker构建和运行：

```bash
# 构建Docker镜像
docker build -t my-licensed-app -f examples/docker/Dockerfile .

# 运行容器
docker run -v $(pwd)/license.dat:/app/config/license.dat my-licensed-app
```



## 安全注意事项

- 在生产环境中，您应该修改 `internal/license/license.go` 中的 `SecretKey`，并确保其安全性。
- 考虑使用更强的加密算法，或将签名密钥存储在安全的硬件模块中。
- 考虑混淆或加密许可证验证相关代码，增加破解难度。



## 容器环境使用

在Docker容器中，由于硬件信息可能受限，本模块提供了特殊的方法获取容器的唯一标识。但请注意，在容器环境中：

1. 容器重建后，可能会导致机器ID变化（如果没有固定hostname）
2. 对于Kubernetes等环境，您可能需要扩展 `GetContainerizedMachineID()` 方法，添加更多适合您环境的唯一标识获取方式



## License

本项目采用 MIT License 授权。 