package main

import (
	"fmt"
	ft "./findtarget"
	tb "./fileanalyze"
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

func runFindTarget(){
	go spinner(100*time.Millisecond)
	flag.Parse()
	roots:=flag.Args()
	if len(roots) == 0{
		roots = []string{"."}
	}

	a, b:=ft.FindTarget(roots, "tube_history_current.tha2")
	fmt.Println(a)
	fmt.Println(b)
	for _, s := range ft.TargetList {
		fmt.Println(s.Name)
	}
}

func runFileAnalyze(){
	a, err := tb.NewTubeHistory(`/Users/qianxin/Projects/Python/myEnv/projects/TubeHistoryAnalyzer/data/106054/tube_history_current.tha2`)
	if err != nil {
		fmt.Println("err begin:")
		fmt.Println(err)
		fmt.Println("err end.")
	}
	fmt.Printf("%q", a)
}


func main() {
	//runFindTarget()
	runFileAnalyze()
}

