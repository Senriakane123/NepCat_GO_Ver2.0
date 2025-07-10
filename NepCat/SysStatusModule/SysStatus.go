package SysStatusModule

import (
	"NepCat_GO/Model"
	"github.com/jander/golog/logger"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"time"
)

func GetSysInfo() Model.SysMonitorStatus {
	SysStatus := Model.SysMonitorStatus{
		TimeStamp: time.Now().String(),
		CPU:       getCPUUsage(),
		Memory:    getMemUsage(),
		Network:   getNetworkStatus(),
		Disk:      getDiskSpeed(),
	}
	return SysStatus
}

func getCPUUsage() float64 {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		logger.Error("获取CPU使用率失败: %v", err)
		return -1
	}
	if len(percent) > 0 {
		return percent[0]
	}
	return -1
}

// 获取内存使用率
func getMemUsage() float64 {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		logger.Error("获取内存使用率失败: %v", err)
		return -1
	}
	return vmStat.UsedPercent
}

var (
	prevIOStat    net.IOCountersStat
	prevTimestamp time.Time
	inited        = false
)

// 获取网络状态
func getNetworkStatus() Model.NetworkSpeed {
	currIOStats, err := net.IOCounters(false)
	if err != nil || len(currIOStats) == 0 {
		logger.Error("获取网络状态失败: %v", err)
		return Model.NetworkSpeed{RecvBPS: -1, SendBPS: -1}
	}
	currStat := currIOStats[0]
	currTime := time.Now()

	if !inited {
		prevIOStat = currStat
		prevTimestamp = currTime
		inited = true
		return Model.NetworkSpeed{RecvBPS: 0, SendBPS: 0}
	}

	duration := currTime.Sub(prevTimestamp).Seconds()
	recvBPS := float64(currStat.BytesRecv-prevIOStat.BytesRecv) / duration
	sendBPS := float64(currStat.BytesSent-prevIOStat.BytesSent) / duration

	prevIOStat = currStat
	prevTimestamp = currTime

	return Model.NetworkSpeed{
		RecvBPS: recvBPS,
		SendBPS: sendBPS,
	}
}

var (
	prevDiskStat disk.IOCountersStat
	prevDiskTime time.Time
	diskInited   = false
)

func getDiskSpeed() Model.DiskSpeed {
	ioStats, err := disk.IOCounters()
	if err != nil || len(ioStats) == 0 {
		logger.Error("获取磁盘状态失败: %v", err)
		return Model.DiskSpeed{ReadBPS: -1, WriteBPS: -1}
	}

	// 默认使用第一个磁盘
	var currStat disk.IOCountersStat
	for _, stat := range ioStats {
		currStat = stat
		break
	}
	currTime := time.Now()

	if !diskInited {
		prevDiskStat = currStat
		prevDiskTime = currTime
		diskInited = true
		return Model.DiskSpeed{ReadBPS: 0, WriteBPS: 0}
	}

	duration := currTime.Sub(prevDiskTime).Seconds()
	readBPS := float64(currStat.ReadBytes-prevDiskStat.ReadBytes) / duration
	writeBPS := float64(currStat.WriteBytes-prevDiskStat.WriteBytes) / duration

	prevDiskStat = currStat
	prevDiskTime = currTime

	return Model.DiskSpeed{
		ReadBPS:  readBPS,
		WriteBPS: writeBPS,
	}
}
