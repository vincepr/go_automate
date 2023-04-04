package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/vincepr/go_cron/yaml"
)

func main(){
	path := "./data.yaml"
	println("main running")
	prOnce, prRepeat := yaml_parser.ParseCmds(path)
	fmt.Printf("Once:  %+v \nRepeat:%+v \n", prOnce, prRepeat)


	//runCmd(prRepeat[0].Cmd)
	ch := make(chan string, 1)
	go func (t time.Duration){
		time.Sleep(t)
		out := runCmd(prRepeat[0].Cmd)
		ch <- out
	}(prRepeat[0].RepeatIn)


	select {
	case res := <-ch:
	   fmt.Println(res)
	case <-time.After(3 * time.Second):
	   fmt.Println("Out of time :(")
	}
	fmt.Println("Main function exited!")

}

func runCmd(str string) string{
	split := strings.Split(str, " ")
	cmd := exec.Command(split[0])
	if len(split)>1{
		cmd = exec.Command(split[0],split[1:]...)
	}
	out,err := cmd.Output()
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println(string(out))
	return string(out)
}

func timeout(t time.Duration){
	time.Sleep(t)
}
