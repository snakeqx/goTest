package FindTarget

import (
	"os"
	"io/ioutil"
	"errors"
	"sync"
	"path/filepath"
)

type Target struct {
	Name string
	Size int64
}

var TargetList []*Target


var wg = &sync.WaitGroup{}


func FolderIteration(path string, target string) (bool, error){
	pathInfo, err := os.Stat(path)
	if err!=nil{
		return false, err
	}
	if !pathInfo.IsDir(){
		err=errors.New("input is not a folder")
		return false, err
	}

	files, err := ioutil.ReadDir(path)
	if err!=nil {
		return false, err
	}

	var fullName string
	for _, file := range files {
		fullName=filepath.Join(path, file.Name())
		if file.IsDir() {
			wg.Add(1)
			go FolderIteration(fullName, target)
		} else {
			if target=="" {
				TargetList = append(TargetList, &Target{fullName, file.Size()})
			} else {
				if file.Name() == target {
					TargetList = append(TargetList, &Target{fullName, file.Size()})
				}
			}

		}
	}
	wg.Done()
	wg.Wait()
	return true, nil
}

func ListAllFiles(path string) error {
	finish, err :=  FolderIteration(path, "")
	if !finish || err !=nil {
		return err
	}
	return nil
}

func SearchFiles(path string, target string) error {
	finish, err :=  FolderIteration(path, target)
	if !finish || err !=nil {
		return err
	}
	return nil
}

/*
func main() {

	//ListAllFiles(`D:\`)
	err := SearchFiles(`D:\`, "IMG_0016.JPG")
	if err!=nil{
		fmt.Println(err)
	}
	for _, files := range TargetList {
		fmt.Printf("%v\n", files.Name)
	}
}
*/