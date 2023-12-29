package mdl

import (
	"github.com/btagrass/gobiz/mdl"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

type Server struct {
	mdl.Mdl
	Host host.InfoStat `json:"host"`
	Cpu  struct {
		cpu.InfoStat
		Usage float64 `json:"usage"`
		Load  float64 `json:"load"`
	} `json:"cpu"`
	Memory mem.VirtualMemoryStat `json:"memory"`
	Net    net.IOCountersStat    `json:"net"`
	Disks  []Disk                `json:"disks"`
}

type Disk struct {
	disk.PartitionStat
	*disk.UsageStat
}
