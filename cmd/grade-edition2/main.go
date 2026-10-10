// Command grade-edition2 runs the separately maintained second-edition grader.
// Invoke it from the repository root: grade-edition2 -ch 1 [submission-directory].
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

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
	consumer := flag.String("consumer", "", "Chapter 4 separate consumer module (default: submission/ch05)")
	consumerCase := flag.String("consumer-case", "", "Chapter 4 JSON case file (default: accepted shout example)")
	flag.Parse()
	if *chapter < 1 || *chapter > 4 || flag.NArg() > 1 || (*chapter != 4 && (*consumer != "" || *consumerCase != "")) {
		fmt.Fprintln(os.Stderr, "usage: grade-edition2 -ch 1|2|3|4 [-json] [-consumer DIR -consumer-case FILE] [submission-directory]")
		return 2
	}
	path := flag.Arg(0)
	if path == "" {
		path = fmt.Sprintf("./solutions/edition-2/ch%02d", *chapter)
		if *chapter == 4 {
			path = "./solutions/edition-2/ch05" // Edition 2 folds original Chapter 5 into 4.
		}
	}
	var example gradeedition2.ConsumerCase
	if *chapter == 4 {
		var err error
		example, err = gradeedition2.ReadConsumerCase(*consumerCase)
		if err != nil {
			fmt.Fprintf(os.Stderr, "grader: %v\n", err)
			return 2
		}
		if *consumer == "" {
			*consumer = filepath.Join(path, "ch05")
		}
	}
	bin, cleanup, err := grade.Build(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "grader: %v\n", err)
		return 2
	}
	defer cleanup()
	var report grade.Report
	switch *chapter {
	case 1:
		var res *grade.RunResult
		res, err = grade.Run(bin)
		if err == nil {
			report = gradeedition2.Chapter1(res)
		}
	case 2:
		report, err = gradeedition2.Chapter2(bin)
	case 3:
		report, err = gradeedition2.Chapter3(bin)
	case 4:
		report, err = gradeedition2.Chapter4(bin, path, *consumer, example)
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
