package fileops

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

// Wait waits specified number of seconds or interactively until
// return/enter is pressed if seconds is less than 0 (e.g -1). 0
// seconds returns immediately.
func Wait(seconds int) {
	if seconds == 0 {
		return
	} else if seconds <= 0 {
		fmt.Println("Press 'Enter' to continue.,,")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return
	}
	fmt.Printf("Waiting %d seconds...\n", seconds)
	time.Sleep(time.Second * time.Duration(seconds))
}
