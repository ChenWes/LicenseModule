package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"sort"
	"strings"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
)

// GetMachineID 返回基于MAC地址和CPU信息的唯一机器标识
func GetMachineID() (string, error) {
	// 获取主板序列号或BIOS信息
	hostInfo, err := host.Info()
	if err != nil {
		return "", fmt.Errorf("failed to get host info: %w", err)
	}

	// 获取CPU信息
	cpuInfo, err := cpu.Info()
	if err != nil {
		return "", fmt.Errorf("failed to get CPU info: %w", err)
	}

	var cpuID string
	if len(cpuInfo) > 0 {
		cpuID = cpuInfo[0].ModelName + cpuInfo[0].PhysicalID
	}

	// 获取第一个物理网卡的MAC地址
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("failed to get network interfaces: %w", err)
	}

	var physicalMAC string
	for _, i := range interfaces {
		// 只选择物理网卡（排除虚拟网卡、回环接口等）
		if i.Flags&net.FlagUp != 0 && // 接口已启用
			!strings.HasPrefix(i.Name, "lo") && // 不是回环接口
			!strings.HasPrefix(i.Name, "veth") && // 不是虚拟网卡
			!strings.HasPrefix(i.Name, "docker") && // 不是docker网卡
			!strings.HasPrefix(i.Name, "br-") && // 不是网桥
			!strings.HasPrefix(i.Name, "v-") && // 不是VPN
			i.HardwareAddr != nil &&
			len(i.HardwareAddr) > 0 {
			physicalMAC = i.HardwareAddr.String()
			break // 只使用第一个符合条件的网卡
		}
	}

	// 组合多个硬件标识符
	idComponents := []string{
		hostInfo.HostID,         // 主板ID
		hostInfo.Platform,       // 平台信息
		hostInfo.PlatformFamily, // 平台系列
		cpuID,                   // CPU信息
		physicalMAC,             // 物理网卡MAC
	}

	// 过滤掉空值
	var validComponents []string
	for _, comp := range idComponents {
		if comp != "" {
			validComponents = append(validComponents, comp)
		}
	}

	// 将所有组件排序以确保顺序一致
	sort.Strings(validComponents)

	// 组合并计算哈希
	idStr := strings.Join(validComponents, "|")
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
