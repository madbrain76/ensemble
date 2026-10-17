// Command gui is an optional client stub, not a GUI or transport implementation.
package main

import (
	"fmt"
	"log"
	"os"

	ensemble "ensemble"
)

func main() {
	// The future GUI will use the same public root and Agent as cmd/ch01.
	// No internal imports, alternate ownership tree or live connection is needed.
	root := ensemble.New(log.New(os.Stderr, "", 0))
	agent, err := root.NewAgent(ensemble.Config{System: "You are a helpful assistant.", MaxTokens: 1024})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("GUI stub: %d recorded events; no GUI transport yet.\n", len(agent.Log().Events()))
}
