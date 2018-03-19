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

	"golang.org/x/sys/unix"
)

const Version = "1.1"

func main() {
	port := flag.String("port", "6060", "Port to listen on")
	o := flag.Bool("o", false, "Open the browser only (don't start godoc)")
	open := flag.Bool("open", false, "Open the browser only (don't start godoc)")
	v := flag.Bool("v", false, "Verbose mode")
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

	hostport := "localhost:" + *port
	_, err = net.LookupHost("godoc")
	if err == nil {
		hostport = "godoc"
	}
	url := fmt.Sprintf("http://%s%s", hostport, path)
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
	godoc, lookErr := exec.LookPath("godoc")
	if lookErr != nil {
		log.Fatal(lookErr)
	}
	args := []string{godoc, "-http", "localhost:" + *port, "-goroot", build.Default.GOROOT}
	if *v {
		args = append(args, "-v")
	}
	execErr := unix.Exec(godoc, args, []string{"GOPATH=" + os.Getenv("GOPATH")})
	if execErr != nil {
		log.Fatal(execErr)
	}
}
