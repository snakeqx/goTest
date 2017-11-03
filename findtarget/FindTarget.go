// This package is to find the exact file name by
// walking all given directory and its sub directories in parallel.
// It uses a concurrency-limiting counting semaphore
// to avoid opening too many files at once.
package findtarget



import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

// Target structure has 2 properties
// prop1: Name string
// prop2: Size int64
type Target struct {
	Name string
	Size int64
}

// This is a exportable variable to store the list
var TargetList []*Target
// this is a channel to store the size, mainly use to sync
var fileSizes = make(chan int64)
var n sync.WaitGroup

// Exportable function to find the target file name in given directory.
// para 1: root []string: can pass a array of path
// para 2: target2find string: the target filename
// return 1: nfiles int64: quantity of files walked through. NOT the target quantity as
//           it can be calculated simply by len(TargetList)
// return 2: nbytes int64: total size sum of all walked through files.
func FindTarget(roots []string, target2find string) (nfiles, nbytes int64) {
	// Traverse each root of the file tree in parallel.
	for _, root := range roots {
		n.Add(1)
		go walkDir(root, &n, target2find, fileSizes)
	}
	go func() {
		n.Wait()
		close(fileSizes)
	}()

loop:
	for {
		select {
		case size, ok := <-fileSizes:
			if !ok {
				break loop // fileSizes was closed
			}
			nfiles++
			nbytes += size
		}
	}

	return nfiles, nbytes // final totals
}


// walkDir recursively walks the file tree rooted at dir
// and sends the size of each found file on fileSizes.
func walkDir(dir string, n *sync.WaitGroup, target2find string, fileSizes chan<- int64) {
	defer n.Done()
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go walkDir(subdir, n, target2find, fileSizes)
		} else {
			if target2find == entry.Name(){
				_absName:=filepath.Join(dir, entry.Name())
				TargetList = append(TargetList, &Target{_absName, entry.Size()})
			}
			fileSizes <- entry.Size()
		}
	}
}


// sema is a counting semaphore for limiting concurrency in dirents.
var sema = make(chan struct{}, 20)

// dirents returns the entries of directory dir.
func dirents(dir string) []os.FileInfo {
	sema <- struct{}{}        // acquire token
	defer func() { <-sema }() // release token

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		return nil
	}
	return entries
}