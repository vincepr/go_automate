package cmd

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func Run(str string) string{
	split := strings.Split(str, " ")
	cmd := exec.Command(split[0])
	if len(split)>1{
		cmd = exec.Command(split[0],split[1:]...)
	}
	out,err := cmd.CombinedOutput()	//vs. cmd.Output()
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println(string(out))
	return string(out)
}