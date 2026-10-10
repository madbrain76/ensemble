// Command grade-edition2 runs the separately maintained second-edition grader.
// Invoke it from the repository root: grade-edition2 -ch 1 [submission-directory].
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/waywardgeek/ensemble/internal/grade"
	"github.com/waywardgeek/ensemble/internal/gradeedition2"
)

func main() {
	os.Exit(run())
}

// Return an exit status so deferred build cleanup runs even after a failed grade.
func run() int {
	chapter := flag.Int("ch", 1, "Edition 2 chapter to grade")
	asJSON := flag.Bool("json", false, "emit the report as JSON")
	flag.Parse()
	if *chapter != 1 || flag.NArg() > 1 {
		fmt.Fprintln(os.Stderr, "usage: grade-edition2 -ch 1 [-json] [submission-directory]")
		return 2
	}
	path := flag.Arg(0)
	if path == "" {
		path = "./solutions/edition-2/ch01"
	}
	bin, cleanup, err := grade.Build(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "grader: %v\n", err)
		return 2
	}
	defer cleanup()
	res, err := grade.Run(bin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "grader: %v\n", err)
		return 2
	}
	report := gradeedition2.Chapter1(res)
	if *asJSON {
		if err := report.WriteJSON(os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "grader: %v\n", err)
			return 2
		}
	} else {
		report.WriteText(os.Stdout)
	}
	if !report.Passed {
		return 1
	}
	return 0
}
