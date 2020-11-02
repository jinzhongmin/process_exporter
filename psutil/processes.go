package psutil

import (
	"log"

	"github.com/shirou/gopsutil/process"
)

type Processes struct {
	processes []*process.Process
}

func (ps *Processes) Count() int {
	return len(ps.processes)
}

func NewProcesses() *Processes {
	processes, err := process.Processes()
	if err != nil {
		log.Panicln(err)
	}
	return &Processes{processes}
}

func NewEmptyProcesses() *Processes {
	ps := new(Processes)
	ps.processes = make([]*process.Process, 0)
	return ps
}

func (ps *Processes) EachProcess(fn func(process *Process) bool) {
	for _, p := range ps.processes {
		if fn(&Process{p}) == false {
			return
		}
	}
}

func (ps *Processes) AppendProcess(process *Process) {
	ps.processes = append(ps.processes, process.Process)
}
