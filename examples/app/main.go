package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/chenwes/licensemodule/internal/license"
	"github.com/chenwes/licensemodule/pkg/utils"
)

const (
	// License文件的相对路径（相对于应用程序）
	defaultLicensePath = "./config/license.dat"
	defaultTimestampPath = "./config/timestamp.dat"
)

func main() {
	// 设置日志前缀
	log.SetPrefix("[ExampleApp] ")
	log.Println("应用程序启动...")

	// 获取应用程序目录
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("无法获取应用程序路径: %v", err)
	}
	appDir := filepath.Dir(exePath)

	// 构建License文件的完整路径
	licenseFile := filepath.Join(appDir, defaultLicensePath)
	timestampFile := filepath.Join(appDir, defaultTimestampPath)

	// 创建License目录（如果不存在）
	licenseDir := filepath.Dir(licenseFile)
	if err := os.MkdirAll(licenseDir, 0755); err != nil {
		log.Fatalf("无法创建许可证目录: %v", err)
	}

	// 获取当前机器ID（这里我们假设应用运行在容器环境中）
	machineID, err := utils.GetContainerizedMachineID()
	if err != nil {
		log.Fatalf("无法获取机器ID: %v", err)
	}

	log.Printf("当前机器ID: %s", machineID)
	log.Printf("许可证文件路径: %s", licenseFile)
	log.Printf("时间戳文件路径: %s", timestampFile)

	// 验证License
	log.Println("验证许可证...")
	err = license.VerifyAndUpdate(licenseFile, timestampFile, machineID)
	
	// 处理验证结果
	if err != nil {
		switch err.Error() {
		case license.ErrInvalidLicense.Error():
			log.Fatalf("无效的许可证，请联系供应商获取有效许可证")
		case license.ErrExpiredLicense.Error():
			log.Fatalf("许可证已过期，请续订许可证")
		case license.ErrMachineMismatch.Error():
			log.Fatalf("许可证与当前机器不匹配，请获取正确的许可证")
		case license.ErrSystemTimeManipulated.Error():
			log.Fatalf("系统时间被篡改，请校正系统时间")
		default:
			if os.IsNotExist(err) {
				log.Fatalf("许可证文件不存在，请联系供应商获取许可证")
			} else {
				log.Fatalf("验证许可证失败: %v", err)
			}
		}
	}

	// 验证成功，继续运行应用程序
	log.Println("许可证验证成功，应用程序正常运行")

	// 这里是应用程序的正常逻辑
	fmt.Println("\n=== 应用程序功能 ===")
	fmt.Println("1. 功能A")
	fmt.Println("2. 功能B")
	fmt.Println("3. 功能C")
	fmt.Println("\n应用程序正在运行中...")

	// 在实际应用中，您可能需要定期验证License（例如每天一次）
	// 这可以通过后台goroutine实现
	
	// 以下是演示代码，在实际应用中您会有一个主事件循环
	// select {
	// case <-time.After(10 * time.Second):
	//     fmt.Println("应用程序运行结束")
	// }
} 