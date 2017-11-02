package main

import (
	"fmt"
	"./findtarget"
)


func main () {
	findtarget.ListAllFiles(`D:\`)
	err := findtarget.SearchFiles(`D:\`, "IMG_0016.JPG")
	if err!=nil{
		fmt.Println(err)
	}
	for _, files := range findtarget.TargetList {
		fmt.Printf("%v\n", files.Name)
	}
}


