package main

import (
	
	"os"

	. "aoanima.ru/logger"

	"github.com/recoilme/pudge"
)

func initPudge (){

	dir, err := os.Stat("/pudge")
	if !os.IsNotExist(err) {
		os.Mkdir("/pudge", 0755)
	}
	Инфо("  %+v \n", dir )
	


}