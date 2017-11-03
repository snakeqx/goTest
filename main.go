package main

import (
	"fmt"
	"./findtarget"
	"time"
	"flag"
)

func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("Searching target:")
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}


func main() {
	go spinner(100*time.Millisecond)
	flag.Parse()
	roots:=flag.Args()
	if len(roots) == 0{
		roots = []string{"."}
	}

	a, b:=findtarget.FindTarget(roots, "tube_history_current.tha2")
	fmt.Println(a);
	fmt.Println(b);
	for _, s := range findtarget.TargetList {
		fmt.Println(s.Name)
	}
}

