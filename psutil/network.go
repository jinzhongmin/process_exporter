package psutil

import "github.com/shirou/gopsutil/net"

type NetCounters struct {
	ics []net.IOCountersStat
}

type NetCounter struct {
	*net.IOCountersStat
}

func (ncs *NetCounters) Count() int {
	return len(ncs.ics)
}

func (ncs *NetCounters) GetNetCounterByName(name string) *NetCounter {
	if ncs.ics == nil {
		return &NetCounter{nil}
	}
	for _, ic := range ncs.ics {
		if ic.Name == name {
			return &NetCounter{&ic}
		}

	}
	return &NetCounter{nil}
}

func (nc *NetCounter) GetBytesRecv() uint64 {
	if nc.IOCountersStat == nil {
		return 0
	}
	return nc.BytesRecv
}

func (nc *NetCounter) GetBytesSent() uint64 {
	if nc.IOCountersStat == nil {
		return 0
	}
	return nc.BytesSent
}

func (nc *NetCounter) GetDropin() uint64 {
	if nc.IOCountersStat == nil {
		return 0
	}
	return nc.Dropin
}

func (nc *NetCounter) GetDropout() uint64 {
	if nc.IOCountersStat == nil {
		return 0
	}
	return nc.Dropout
}

func (nc *NetCounter) GetErrin() uint64 {
	if nc.IOCountersStat == nil {
		return 0
	}
	return nc.Errin
}

func (nc *NetCounter) GetErrout() uint64 {
	if nc.IOCountersStat == nil {
		return 0
	}
	return nc.Errout
}

func (nc *NetCounter) GetFifoin() uint64 {
	if nc.IOCountersStat == nil {
		return 0
	}
	return nc.Fifoin
}

func (nc *NetCounter) GetFifoout() uint64 {
	if nc.IOCountersStat == nil {
		return 0
	}
	return nc.Fifoout
}

func (nc *NetCounter) GetPacketsRecv() uint64 {
	if nc.IOCountersStat == nil {
		return 0
	}
	return nc.PacketsRecv
}

func (nc *NetCounter) GetPacketsSent() uint64 {
	if nc.IOCountersStat == nil {
		return 0
	}
	return nc.PacketsSent
}

func (nc *NetCounter) GetName() string {
	if nc.IOCountersStat == nil {
		return ""
	}
	return nc.Name
}
