package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type location struct {
	path   string
	remote string
}

func parseLocation(s string) *location {
	fs := strings.Split(s, ":")
	if len(fs) == 2 {
		return &location{
			path:   fs[1],
			remote: fs[0],
		}
	}
	return &location{
		path: fs[0],
	}
}

func main() {
	flag.Parse()
	if flag.NArg() != 2 {
		fmt.Printf("usage: %s src dest-prefix\n", os.Args[0])
		return
	}
	src, dest := parseLocation(flag.Arg(0)), parseLocation(flag.Arg(1))

	srccmd := fmt.Sprintf("tar cf - -C %s %s", filepath.Dir(src.path), filepath.Base(src.path))
	if src.remote != "" {
		srccmd = fmt.Sprintf("netio -addr %s:8126 %s", src.remote, srccmd)
	}

	destcmd := fmt.Sprintf("tar xf - -C %s", dest.path)
	if dest.remote != "" {
		destcmd = fmt.Sprintf("netio -addr %s:8126 %s", dest.remote, destcmd)
	}

	cmdstr := srccmd + " | " + destcmd
	fmt.Println(cmdstr)
	cmd := exec.Command("bash", "-c", cmdstr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
