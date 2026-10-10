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
	consumer := flag.String("consumer", "", "Chapter 4–6 separate consumer module (default: ch05 under or beside submission)")
	consumerCase := flag.String("consumer-case", "", "Chapter 4–6 JSON case file (default: accepted shout example)")
	workflow := flag.String("workflow", "", "Chapter 5/6 workflow module (default: ch06 under or beside submission)")
	mediaModels := flag.String("media-models", "", "Chapter 5/6 JSON vendor/model map (default: historical accepted fixtures)")
	streamModel := flag.String("stream-model", "claude-sonnet-5", "Chapter 6 declared Anthropic streaming fixture (historical default)")
	streamExercise := flag.String("stream-exercise", "", "Chapter 6 external shape module (default: ch07 under or beside submission)")
	flag.Parse()
	if *chapter < 1 || *chapter > 6 || flag.NArg() > 1 || (*chapter < 4 && (*consumer != "" || *consumerCase != "")) || (*chapter < 5 && (*workflow != "" || *mediaModels != "")) || (*chapter != 6 && *streamExercise != "") {
		fmt.Fprintln(os.Stderr, "usage: grade-edition2 -ch 1|2|3|4|5|6 [-json] [-consumer DIR -consumer-case FILE] [-workflow DIR -media-models FILE] [-stream-exercise DIR -stream-model MODEL] [submission-directory]")
		return 2
	}
	path := flag.Arg(0)
	if path == "" {
		path = fmt.Sprintf("./solutions/edition-2/ch%02d", *chapter)
		if *chapter == 4 {
			path = "./solutions/edition-2/ch05" // Edition 2 folds original Chapter 5 into 4.
		} else if *chapter >= 5 {
			path = fmt.Sprintf("./solutions/edition-2/ch%02d", *chapter+1)
		}
	}
	var example gradeedition2.ConsumerCase
	if *chapter >= 4 {
		var err error
		example, err = gradeedition2.ReadConsumerCase(*consumerCase)
		if err != nil {
			fmt.Fprintf(os.Stderr, "grader: %v\n", err)
			return 2
		}
		if *consumer == "" {
			*consumer = appDirectory(path, "ch05")
		}
	}
	var models map[string]string
	if *chapter >= 5 {
		var err error
		models, err = gradeedition2.ReadMediaModels(*mediaModels)
		if err != nil {
			fmt.Fprintf(os.Stderr, "grader: %v\n", err)
			return 2
		}
		if *workflow == "" {
			*workflow = appDirectory(path, "ch06")
		}
	}
	if *chapter == 6 && *streamExercise == "" {
		*streamExercise = appDirectory(path, "ch07")
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
	case 5:
		report, err = gradeedition2.Chapter5(bin, path, *consumer, *workflow, example, models)
	case 6:
		report, err = gradeedition2.Chapter6(bin, path, *consumer, *workflow, *streamExercise, *streamModel, example, models)
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

// Historical exports place applications either inside the core tree or beside
// it. Discover only those established layouts; an absent app still fails grading.
func appDirectory(core, name string) string {
	for _, dir := range []string{filepath.Join(core, name), filepath.Join(core, "..", name)} {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
	}
	return filepath.Join(core, name)
}
