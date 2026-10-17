// Command grade-edition3 runs the supported Edition 3 chapter gates.
// Invoke it from the repository root with an explicit submission path.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/waywardgeek/ensemble/internal/grade"
	"github.com/waywardgeek/ensemble/internal/gradeedition2"
	"github.com/waywardgeek/ensemble/internal/gradeedition3"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

// Returning the status lets temporary build cleanup finish before os.Exit.
func run(args []string) int {
	flags := flag.NewFlagSet("grade-edition3", flag.ContinueOnError)
	chapter := flags.Int("ch", 1, "Edition 3 chapter to grade (1 or 4)")
	asJSON := flags.Bool("json", false, "emit the report as JSON")
	core := flags.String("core", "", "Chapter 4 core source module directory")
	consumer := flags.String("consumer", "", "Chapter 4 external consumer module directory")
	consumerCase := flags.String("consumer-case", "", "Chapter 4 JSON case file (default: shout example)")
	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}
	chapter4Flags := false
	flags.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "core", "consumer", "consumer-case":
			chapter4Flags = true
		}
	})
	// No historical target defaults and no implicit source inference from a binary.
	if (*chapter != 1 && *chapter != 4) || flags.NArg() != 1 || flags.Arg(0) == "" || (*chapter == 1 && chapter4Flags) {
		fmt.Fprintln(os.Stderr, "usage: grade-edition3 -ch 1 [-json] target; or -ch 4 -core DIR -consumer DIR [-consumer-case FILE] [-json] target")
		return 2
	}
	var example gradeedition2.ConsumerCase
	if *chapter == 4 {
		for _, dir := range []string{*core, *consumer} {
			info, err := os.Stat(filepath.Join(dir, "go.mod"))
			if dir == "" || err != nil || !info.Mode().IsRegular() {
				fmt.Fprintln(os.Stderr, "grader: Chapter 4 requires explicit core and consumer module directories containing go.mod")
				return 2
			}
		}
		var err error
		example, err = gradeedition2.ReadConsumerCase(*consumerCase)
		if err != nil {
			fmt.Fprintf(os.Stderr, "grader: %v\n", err)
			return 2
		}
	}
	bin, cleanup, err := grade.Build(flags.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "grader: %v\n", err)
		return 2
	}
	defer cleanup()
	var report grade.Report
	if *chapter == 4 {
		// Reuse the historical checks; only separate executable and source identities.
		report, err = gradeedition2.Chapter4(bin, *core, *consumer, example)
	} else {
		report, err = gradeedition3.Chapter1(bin)
	}
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
