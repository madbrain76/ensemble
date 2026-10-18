// Command grade-edition3 runs the supported Edition 3 chapter gates.
// Invoke it from the repository root with an explicit submission path.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	chapter := flags.Int("ch", 1, "Edition 3 chapter to grade (1, 4, 5, 6, 7, 8 or 9)")
	asJSON := flags.Bool("json", false, "emit the report as JSON")
	core := flags.String("core", "", "Chapter 4/5/6/7/8/9 core source module directory")
	consumer := flags.String("consumer", "", "Chapter 4/5/6/7/8/9 external consumer module directory")
	consumerCase := flags.String("consumer-case", "", "Chapter 4/5/6/7/8/9 JSON case file (default: shout example)")
	workflow := flags.String("workflow", "", "Chapter 5/6/7/8/9 workflow module directory")
	mediaModels := flags.String("media-models", "", "Chapter 5/6/7/8/9 media-model mapping JSON")
	streamExercise := flags.String("stream-exercise", "", "Chapter 6/7/8/9 stream-shape module directory")
	streamModel := flags.String("stream-model", "", "Chapter 6/7/8/9 declared known streaming model (required)")
	gui := flags.String("gui", "", "Chapter 7/8/9 separate GUI module directory")
	guiModel := flags.String("gui-model", "", "Chapter 7/8/9 declared GUI model (required)")
	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}
	chapter4Flags := false
	chapter5Flags := false
	chapter6Flags := false
	chapter7Flags := false
	flags.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "core", "consumer", "consumer-case":
			chapter4Flags = true
		case "workflow", "media-models":
			chapter5Flags = true
		case "stream-exercise", "stream-model":
			chapter6Flags = true
		case "gui", "gui-model":
			chapter7Flags = true
		}
	})
	// No historical target defaults and no implicit source inference from a binary.
	if (*chapter != 1 && *chapter != 4 && *chapter != 5 && *chapter != 6 && *chapter != 7 && *chapter != 8 && *chapter != 9) || flags.NArg() != 1 || flags.Arg(0) == "" || (*chapter == 1 && (chapter4Flags || chapter5Flags)) || (*chapter == 4 && chapter5Flags) || (*chapter != 6 && *chapter != 7 && *chapter != 8 && *chapter != 9 && chapter6Flags) || (*chapter != 7 && *chapter != 8 && *chapter != 9 && chapter7Flags) {
		fmt.Fprintln(os.Stderr, "usage: grade-edition3 -ch 1 [-json] target; or -ch 4|5|6|7|8|9 -core DIR -consumer DIR [-consumer-case FILE] [-workflow DIR [-media-models FILE]] [-stream-exercise DIR -stream-model MODEL] [-gui DIR -gui-model MODEL] [-json] target")
		return 2
	}
	var example gradeedition2.ConsumerCase
	if *chapter == 4 || *chapter == 5 || *chapter == 6 || *chapter == 7 || *chapter == 8 || *chapter == 9 {
		dirs := []string{*core, *consumer}
		if *chapter == 5 || *chapter == 6 || *chapter == 7 || *chapter == 8 || *chapter == 9 {
			dirs = append(dirs, *workflow)
		}
		if *chapter == 6 || *chapter == 7 || *chapter == 8 || *chapter == 9 {
			dirs = append(dirs, *streamExercise)
			if strings.TrimSpace(*streamModel) == "" {
				fmt.Fprintln(os.Stderr, "grader: explicit nonblank -stream-model is required for Chapters 6/7/8/9")
				return 2
			}
		}
		if *chapter == 7 || *chapter == 8 || *chapter == 9 {
			dirs = append(dirs, *gui)
			if strings.TrimSpace(*guiModel) == "" {
				fmt.Fprintln(os.Stderr, "grader: explicit nonblank -gui-model is required for Chapters 7/8/9")
				return 2
			}
		}
		for _, dir := range dirs {
			info, err := os.Stat(filepath.Join(dir, "go.mod"))
			if dir == "" || err != nil || !info.Mode().IsRegular() {
				fmt.Fprintln(os.Stderr, "grader: explicit core and consumer (plus workflow for Chapters 5/6/7/8/9, stream-exercise for Chapters 6/7/8/9, and gui for Chapters 7/8/9) module directories containing go.mod are required")
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
	var models map[string]string
	if *chapter == 5 || *chapter == 6 || *chapter == 7 || *chapter == 8 || *chapter == 9 {
		var err error
		models, err = gradeedition2.ReadMediaModels(*mediaModels)
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
	} else if *chapter == 5 {
		report, err = gradeedition2.Chapter5(bin, *core, *consumer, *workflow, example, models)
	} else if *chapter == 6 {
		report, err = gradeedition2.Chapter6(bin, *core, *consumer, *workflow, *streamExercise, *streamModel, example, models)
	} else if *chapter == 7 {
		report, err = gradeedition2.Chapter7(bin, *core, *gui, *guiModel, *consumer, *workflow, *streamExercise, *streamModel, example, models)
	} else if *chapter == 8 {
		report, err = gradeedition2.Chapter8(bin, *core, *gui, *guiModel, *consumer, *workflow, *streamExercise, *streamModel, example, models)
	} else if *chapter == 9 {
		report, err = gradeedition2.Chapter9(bin, *core, *gui, *guiModel, *consumer, *workflow, *streamExercise, *streamModel, example, models)
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
