package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/shirou/gopsutil/v3/cpu"
)

// GetMachineID 返回基于MAC地址和CPU信息的唯一机器标识
func GetMachineID() (string, error) {
	// 获取MAC地址
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("获取网络接口失败: %w", err)
	}

	var macAddresses []string
	for _, i := range interfaces {
		if i.Flags&net.FlagUp != 0 && !strings.HasPrefix(i.Name, "lo") {
			macAddresses = append(macAddresses, i.HardwareAddr.String())
		}
	}

	// 获取CPU信息
	cpuInfo, err := cpu.Info()
	if err != nil {
		return "", fmt.Errorf("获取CPU信息失败: %w", err)
	}

	var cpuID string
	if len(cpuInfo) > 0 {
		cpuID = cpuInfo[0].ModelName + cpuInfo[0].PhysicalID
	}

	// 将MAC地址和CPUID组合并计算SHA256哈希
	idStr := strings.Join(macAddresses, "") + cpuID
	hash := sha256.Sum256([]byte(idStr))
	return hex.EncodeToString(hash[:]), nil
}

// GetContainerizedMachineID 返回一个可在容器环境中使用的机器标识
// 注意：这个方法在容器中可能无法获取到宿主机的真实硬件信息
func GetContainerizedMachineID() (string, error) {
	// 尝试获取标准机器ID
	id, err := GetMachineID()
	if err != nil {
		// 如果失败，使用备用方法
		// 在Docker中，可以使用hostname作为辅助标识
		hostname, err := os.Hostname()
		if err != nil {
			return "", fmt.Errorf("无法获取主机名: %w", err)
		}
		
		hash := sha256.Sum256([]byte(hostname))
		return hex.EncodeToString(hash[:]), nil
	}
	return id, nil
} 