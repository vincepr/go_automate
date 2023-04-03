package main

import (
	"fmt"

	"github.com/vincepr/go_cron/yaml"
)

func main(){
	path := "./data.yaml"
	println("main running")
	pOnce, pRepeat := yaml_parser.ParseCmds(path)
	fmt.Printf("Once:  %+v \nRepeat:%+v \n", pOnce, pRepeat)

}
