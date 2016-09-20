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

const Version = "0.4"

func main() {
	port := flag.String("port", "6060", "Port to listen on")
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

	// mostly copied from the godoc/main.go. you can run it in this mode but
	// godoc/main.go does a lot.
	// rootfs := gatefs.New(vfs.OS(*goroot), fsGate)
	// fs.Bind("/", rootfs, "/", vfs.BindReplace)

	// for _, p := range filepath.SplitList(build.Default.GOPATH) {
	//    fs.Bind("/src", gatefs.New(vfs.OS(p), fsGate), "/src", vfs.BindAfter)
	// }

	// godoc.CommandLine(os.Stdout, fs, flag.Args())
	go func(p string) {
		for {
			conn, err := net.Dial("tcp", ":"+p)
			if err == nil {
				defer conn.Close()
				url := fmt.Sprintf("http://localhost:%s/%s", p, path)
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
