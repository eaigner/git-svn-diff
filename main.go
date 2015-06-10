package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var (
	addRx     = regexp.MustCompile(`(?m:^\+{3}.*$)`)
	subRx     = regexp.MustCompile(`(?m:^\-{3}.*$)`)
	nullRx    = regexp.MustCompile(`(?m:^(\-{3}\s\/dev\/null).*$)`)
	indexRx   = regexp.MustCompile(`(?m:^index\s.*$)`)
	diffGitRx = regexp.MustCompile(`(?m:^diff\s\-\-git\s[^[:space:]]*)`)
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("please provide a commit-ish argument")
	}

	rev := strings.TrimSpace(runCmd("git", "svn", "find-rev", os.Args[1]))
	diff := runCmd("git", "diff", "--no-prefix", os.Args[1])

	// Rewrite the diff for SVN.
	diff = addRx.ReplaceAllString(diff, "$0    (working copy)")
	diff = subRx.ReplaceAllString(diff, "$0    (revision "+rev+")")
	diff = nullRx.ReplaceAllString(diff, "$1    (revision 0)")
	diff = indexRx.ReplaceAllString(diff, "===================================================================")
	diff = diffGitRx.ReplaceAllString(diff, "Index:")

	fmt.Print(diff)
}

func runCmd(name string, arg ...string) string {
	cmd := exec.Command(name, arg...)
	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		log.Fatal(errBuf.String())
	}
	return outBuf.String()
}
