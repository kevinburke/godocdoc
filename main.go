package main

import (
	"flag"
	"fmt"
	"go/build"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const Version = "1.0"

func main() {
	port := flag.String("port", "6060", "Port to listen on")
	o := flag.Bool("o", false, "Open the browser only (don't start godoc)")
	open := flag.Bool("open", false, "Open the browser only (don't start godoc)")
	flag.Parse()
	if flag.NArg() > 0 && flag.Arg(0) == "version" {
		os.Stderr.WriteString(fmt.Sprintf("godocdoc version %s\n", Version))
		os.Exit(1)
	}

	path := "/pkg"
	srcPath := filepath.Join(build.Default.GOPATH, "src")
	wd, err := os.Getwd()
	if err == nil && strings.HasPrefix(wd, srcPath) {
		rel, err := filepath.Rel(srcPath, wd)
		if err == nil {
			path = path + "/" + rel
		}
	}

	url := fmt.Sprintf("http://localhost:%s%s", *port, path)
	if *o || *open {
		Open(url)
		return
	}
	go func(port string) {
		for {
			conn, err := net.Dial("tcp", ":"+port)
			if err == nil {
				defer conn.Close()
				if ok := Open(url); !ok {
					fmt.Println(url)
				}
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}(*port)
	cmd := exec.Command("godoc", "-http", "localhost:"+*port, "-goroot", build.Default.GOROOT)
	defer func() {
		cmd.Process.Kill()
	}()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
