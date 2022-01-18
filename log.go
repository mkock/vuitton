package vuitton

import (
	"fmt"
	"time"
)

func output(msg string) {
	if msg == "" {
		return
	}
	fmt.Printf("%s %s\n", time.Now().Format("15:04:05"), msg)
}
