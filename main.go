package main

import (
	"flag"
	"log"

	"github.com/vincepr/go_automate/shedule"
	"github.com/vincepr/go_automate/yaml"
)



func main(){
	// parse flags:
	path := flag.String("path", "./data.yaml", "path to the .yaml file")
	flag.Parse()
	// setup sync & channel
	prOnce, prRepeat := yaml_parser.ParseCmds(*path)
	ch := make(chan schedule.InfoChan)

	// start the processes:
	for _,val :=range prOnce{
		schedule.AddWg(+1)
		go schedule.RunOnce(val, ch)
	}
	for _,val:=range prRepeat{
		schedule.AddWg(+1)
		go schedule.RunRepeating(val, ch)
	}

	// keep reading for processes's outputs:
	for out := range ch {
		log.Printf("FINISHED: %s;		CMD: %s;		OUTPUT: %s;", out.Name, out.Cmd, out.Output)
		if schedule.ReadWg()<=0{
			break	// all processes have finished
		}
    }

	println("All Processes finished, shutting down gracefully.")
}
