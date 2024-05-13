package svc

import (
	"github.com/btagrass/gobiz/sys/mdl"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/spf13/cast"
)

type ServerSvc struct {
}

func NewServerSvc(i *do.Injector) (*ServerSvc, error) {
	return &ServerSvc{}, nil
}

func (s *ServerSvc) GetServer() mdl.Server {
	var server mdl.Server
	hostStat, _ := host.Info()
	server.Host = *hostStat
	cpuStats, _ := cpu.Info()
	if len(cpuStats) > 0 {
		server.Cpu.InfoStat = cpuStats[0]
	}
	cpuCount, _ := cpu.Counts(true)
	server.Cpu.Cores = cast.ToInt32(cpuCount)
	cpuUsages, _ := cpu.Percent(0, false)
	if len(cpuUsages) > 0 {
		server.Cpu.Usage = cpuUsages[0]
	}
	loadStat, _ := load.Avg()
	server.Cpu.Load = loadStat.Load1
	memoryStat, _ := mem.VirtualMemory()
	server.Memory = *memoryStat
	netStats, _ := net.IOCounters(false)
	if len(netStats) > 0 {
		server.Net = netStats[0]
	}
	partitions, _ := disk.Partitions(false)
	partitions = lo.Filter(partitions, func(i disk.PartitionStat, _ int) bool {
		return (i.Fstype == "ext4" || i.Fstype == "NTFS") && lo.Contains(i.Opts, "rw")
	})
	for _, p := range partitions {
		usage, _ := disk.Usage(p.Mountpoint)
		server.Disks = append(server.Disks, mdl.Disk{
			PartitionStat: p,
			UsageStat:     usage,
		})
	}
	return server
}
