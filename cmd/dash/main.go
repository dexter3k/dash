package main

import (
	"fmt"
	"math/rand"
	"os"
	"io/ioutil"
	"runtime"
	"time"

	"github.com/dexter3k/dash/player"
)

func init() {
	runtime.LockOSThread()
	rand.Seed(time.Now().UTC().UnixNano())
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func loadBinary(path string) []byte {
	f, err := os.Open(path)
	check(err)
	defer f.Close()
	d, err := ioutil.ReadAll(f)
	check(err)
	return d
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: dash <filename.swf>")
		return
	}

	player := player.NewPlayer()
	loader := player.AddSwf(os.Args[1])

	mainFile := loadBinary(os.Args[1])
	check(loader.LoadRawBytes(mainFile))

	check(player.Play())
}
