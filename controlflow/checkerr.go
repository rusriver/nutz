package controlflow

import (
	"fmt"
	"os"
)

func Checkerr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
