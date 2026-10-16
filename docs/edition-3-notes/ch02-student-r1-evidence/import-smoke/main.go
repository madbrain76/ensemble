package main

import (
	"ensemble"
	"fmt"
	"io"
	"log"
)

func main() {
	root := ensemble.New(log.New(io.Discard, "", 0))
	agent, err := root.NewAgent(ensemble.Config{System: "import smoke", MaxTokens: 32})
	if err != nil {
		panic(err)
	}
	if agent.Parent() != root || agent.Engine().Parent() != agent || agent.Log().Parent() != agent {
		panic("lost ownership")
	}
	if len(agent.Log().Events()) != 0 {
		panic("unexpected events")
	}
	fmt.Println("external public-root import and ownership smoke passed")
}
