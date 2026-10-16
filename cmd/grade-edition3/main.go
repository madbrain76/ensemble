// Command grade-edition3 runs the Edition 3 Chapter 1 readiness overlay.
// Invoke it from the repository root with an explicit submission path.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/waywardgeek/ensemble/internal/grade"
	"github.com/waywardgeek/ensemble/internal/gradeedition3"
)

func main() {
	os.Exit(run())
}

// Returning the status lets temporary build cleanup finish before os.Exit.
func run() int {
	chapter := flag.Int("ch", 1, "Edition 3 chapter to grade (only 1 is implemented)")
	asJSON := flag.Bool("json", false, "emit the report as JSON")
	flag.Parse()
	// No historical default: Edition 3 starts from an explicit fresh target.
	if *chapter != 1 || flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "usage: grade-edition3 -ch 1 [-json] submission-directory-or-executable")
		return 2
	}
	bin, cleanup, err := grade.Build(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "grader: %v\n", err)
		return 2
	}
	defer cleanup()
	report, err := gradeedition3.Chapter1(bin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "grader: %v\n", err)
		return 2
	}
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
