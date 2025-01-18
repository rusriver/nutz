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

func Checkerr2(tag string, err error) {
	if err != nil {
		fmt.Println(tag, err)
		os.Exit(1)
	}
}
