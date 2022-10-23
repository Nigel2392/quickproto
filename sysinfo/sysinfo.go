package sysinfo

import (
	"encoding/json"
	"net"
	"strings"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

const (
	INC_HOSTNAME = iota
	INC_PLATFORM
	INC_CPU
	INC_MEM
	INC_DISK
	INC_MACADDR
)

// SysInfo saves the basic system information
type SysInfo struct {
	Hostname string   `json:"hostname,omitempty"`
	Platform string   `json:"platform,omitempty"`
	CPU      CPUInfo  `json:"cpu,omitempty"`
	RAM      RAMInfo  `json:"ram,omitempty"`
	Disk     DiskInfo `json:"disk,omitempty"`
	MacAddr  string   `json:"macaddr,omitempty"`
}

// CPUInfo saves the CPU information
type CPUInfo struct {
	Threads      int32   `json:"threads,omitempty"`
	CurrentUsage float64 `json:"currentusage,omitempty"`
	Name         string  `json:"name,omitempty"`
}

// RAMInfo saves the RAM information
type RAMInfo struct {
	Total uint64 `json:"total,omitempty"`
	Used  uint64 `json:"used,omitempty"`
	Free  uint64 `json:"free,omitempty"`
}

// DiskInfo saves the Disk information
type DiskInfo struct {
	SysID string `json:"sysid,omitempty"`
	Path  string `json:"path,omitempty"`
	Total uint64 `json:"total,omitempty"`
	Used  uint64 `json:"used,omitempty"`
	Free  uint64 `json:"free,omitempty"`
}

func GetSysInfo(include []int) *SysInfo {
	hostStat, _ := host.Info()
	cpuStat, _ := cpu.Info()
	vmStat, _ := mem.VirtualMemory()
	diskStat, _ := disk.Usage("\\") // If you're in Unix change this "\\" for "/"

	info := new(SysInfo)
	if ContainsInt(include, INC_HOSTNAME) {
		info.Hostname = strings.TrimSpace(hostStat.Hostname)
	}
	if ContainsInt(include, INC_PLATFORM) {
		info.Platform = strings.TrimSpace(hostStat.Platform)
	}
	if ContainsInt(include, INC_CPU) {
		info.CPU = CPUInfo{
			Threads:      cpuStat[0].Cores,
			CurrentUsage: cpuStat[0].Mhz,
			Name:         strings.TrimSpace(cpuStat[0].ModelName),
		}
	}
	if ContainsInt(include, INC_MEM) {
		info.RAM = RAMInfo{
			Total: vmStat.Total, // 1024 / 1024, // MB
			Used:  vmStat.Used,  // 1024 / 1024,  // MB
			Free:  vmStat.Free,  // 1024 / 1024,  // MB
		}
	}
	if ContainsInt(include, INC_DISK) {
		info.Disk = DiskInfo{
			Path:  diskStat.Path,
			Total: diskStat.Total, // 1024 / 1024 / 1024, // GB
			Used:  diskStat.Used,  // 1024 / 1024 / 1024,  // GB
			Free:  diskStat.Free,  // 1024 / 1024 / 1024,  // GB
		}
	}
	if ContainsInt(include, INC_MACADDR) {
		info.MacAddr, _ = GetMACAddr()
	}
	return info
}

func (s *SysInfo) ToJSON() []byte {
	json, _ := json.Marshal(s)
	return json
}

func (s *SysInfo) FromJson(jdata []byte) *SysInfo {
	json.Unmarshal(jdata, s)
	return s
}

func GetMACAddr() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	var currentIP, currentNetworkHardwareName string
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		// = GET LOCAL IP ADDRESS
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				currentIP = ipnet.IP.String()
			}
		}
	}
	// get all the system's or local machine's network interfaces
	interfaces, _ := net.Interfaces()
	for _, interf := range interfaces {
		if addrs, err := interf.Addrs(); err == nil {
			for _, addr := range addrs {
				// only interested in the name with current IP address
				if strings.Contains(addr.String(), currentIP) {
					currentNetworkHardwareName = interf.Name
				}
			}
		}
	}
	netInterface, err := net.InterfaceByName(currentNetworkHardwareName)
	if err != nil {
		return "", err
	}
	macAddress := netInterface.HardwareAddr
	// verify if the MAC address can be parsed properly
	hwAddr, err := net.ParseMAC(macAddress.String())
	if err != nil {
		return "", err
	}
	return hwAddr.String(), nil
}

func ContainsInt(slice []int, item int) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
