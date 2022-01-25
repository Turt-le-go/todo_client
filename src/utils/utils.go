package utils

import (
	"os"
	"os/exec"
	"fmt"
)

func Check(err error){
	if err != nil {
		panic(err)
	}
}

func GetSize() (x,y int) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	_,err = fmt.Sscanf(string(out),"%d %d",&y,&x)
	if err != nil {
		panic(err)
	}
	return x,y
}
