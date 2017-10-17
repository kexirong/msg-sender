package main

import (
	"fmt"
	//   "reflect"
	"github.com/kexirong/msg-sender/http"
)

const cfgFileName = "cfg.json"

func main() {
	fmt.Println("this a  msg-sender.")
	Jcfg, err := getConfig(cfgFileName)
	if err != nil {
		panic(err)
	}
	//fmt.Println(Jcfg)
	http.SrvStart(Jcfg)

	fmt.Println("runing... ")
}
