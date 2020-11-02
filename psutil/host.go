package psutil

import (
	"log"

	"github.com/shirou/gopsutil/host"
)

type Host struct {
	*host.InfoStat
}

func HostInfo() *Host {
	hostInfo, err := host.Info()
	if err != nil {
		log.Panicln(err)
		return &Host{nil}
	}
	return &Host{hostInfo}
}
