package psutil

import (
	"log"
	"regexp"
	"strings"

	"github.com/shirou/gopsutil/process"
)

type Process struct {
	*process.Process
}

func (p *Process) HasPid(pid int32) bool {
	if p.Process == nil || p.Pid != pid {
		return false
	}
	return true
}

func (p *Process) HasCmdSub(substr string) bool {
	if p.Process == nil {
		return false
	}

	cmdline, err := p.Cmdline()
	if err != nil {
		log.Println(err)
		return false
	}
	if strings.Index(cmdline, substr) > -1 {
		return true
	}
	return false
}

func (p *Process) HasCmdReg(expr string) bool {
	if p.Process == nil {
		return false
	}

	reg, err := regexp.Compile(expr)
	if err != nil {
		return false
	}

	cmdline, err := p.Cmdline()
	if err != nil {
		log.Println(err)
		return false
	}

	return reg.MatchString(cmdline)
}

func (p *Process) GetCmdLine() string {
	if p.Process == nil {
		return ""
	}
	cmdline, err := p.Cmdline()
	if err != nil {
		log.Println(err)
	}
	return cmdline
}

func (p *Process) GetCPUPercent() float64 {
	if p.Process == nil {
		return 0
	}

	cpu, err := p.CPUPercent()
	if err != nil {
		log.Println(err)
		return 0
	}

	return cpu
}

func (p *Process) GetMemoryPercent() float32 {
	if p.Process == nil {
		return 0
	}

	memory, err := p.MemoryPercent()
	if err != nil {
		log.Println(err)
		return 0
	}

	return memory
}

func (p *Process) GetMemoryRssBytes() uint64 {
	if p.Process == nil {
		return 0
	}

	mem, err := p.MemoryInfo()
	if err != nil {
		log.Println(err)
		return 0
	}

	return mem.RSS
}

func (p *Process) GetReadBytes() uint64 {
	if p.Process == nil {
		return 0
	}

	io, _ := p.IOCounters()
	return io.ReadBytes
}

func (p *Process) GetWriteBytes() uint64 {
	if p.Process == nil {
		return 0
	}

	io, _ := p.IOCounters()
	return io.WriteBytes
}

func (p *Process) GetReadCount() uint64 {
	if p.Process == nil {
		return 0
	}

	io, _ := p.IOCounters()
	return io.ReadCount
}

func (p *Process) GetWriteCount() uint64 {
	if p.Process == nil {
		return 0
	}

	io, _ := p.IOCounters()
	return io.WriteCount
}

func (p *Process) GetNetCounters() *NetCounters {
	if p.Process == nil {
		return &NetCounters{nil}
	}

	iocs, err := p.NetIOCounters(true)
	if err != nil {
		//log.Println(err)
		return &NetCounters{nil}
	}
	return &NetCounters{iocs}
}
