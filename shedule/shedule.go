package schedule

import (
	"log"
	"strings"
	"sync/atomic"
	"time"

	"github.com/vincepr/go_automate/cmd"
	"github.com/vincepr/go_automate/yaml"
)

type InfoChan struct{
	Name	string
	Cmd		string
	Output 	string
}

/*
	schedulers to start Processes. Themself called with goroutines for multithreading.
*/

func RunOnce(p yaml_parser.ProcessOnce, ch chan InfoChan){
	dur := p.When.Sub(time.Now())
	log.Printf("SCHEDULED: %s; IN: %v; CMD: %s",p.Name,dur.Round(time.Second),p.Cmd)
	time.Sleep(dur)
	log.Printf("STARTED: %s;	CMD: %s;",p.Name,p.Cmd)
	out := cmd.Run(p.Cmd)
	out = strings.TrimSuffix(out, "\n")
	AddWg(-1)
	ch <- InfoChan{
		Name: 	p.Name,
		Cmd:	p.Cmd,
		Output: string(out),
	}
}

func RunRepeating(p yaml_parser.ProcessRepeat, ch chan InfoChan){
	if p.RepeatTimes == 0 {p.RepeatTimes=1}
	counter := p.RepeatTimes
	dur := p.When.Sub(time.Now())
	log.Printf("SCHEDULED: %s; IN: %v; REPEATS: %v; CMD: %s;",p.Name,dur.Round(time.Second),counter,p.Cmd)
	time.Sleep(dur)
	for counter>0{
		log.Printf("STARTED: %s; IN: %v REPEATS: %v; CMD: %s",p.Name,dur.Round(time.Second),counter,p.Cmd)
		out := cmd.Run(p.Cmd)
		out = strings.TrimSuffix(out, "\n")
		ch <- InfoChan{
			Name: 	p.Name,
			Cmd:	p.Cmd,
			Output: string(out),
		}
		time.Sleep(p.RepeatIn)
		counter--
	}
	log.Printf("STARTED: %s;		CMD: %s; LAST TIME RUNNING;",p.Name,p.Cmd)
	out := cmd.Run(p.Cmd)
	out = strings.TrimSuffix(out, "\n")
	AddWg(-1)
	ch <- InfoChan{
		Name: 	p.Name,
		Cmd:	p.Cmd,
		Output: string(out),
	}
}

/*
	handle syncing up, when all Processes are done. By taking count
*/

// "WaitGroup", takes count of how many process are not finished
var countPr int32

func ReadWg() int32{
	return atomic.LoadInt32(&countPr)
}

func AddWg(change int32){
	wg := ReadWg()
	atomic.StoreInt32(&countPr, wg+change)
}
