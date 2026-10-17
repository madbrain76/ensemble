package common

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"time"
)

// ErrObservationWoken releases a foreground observation without canceling work.
var ErrObservationWoken = errors.New("observation woken by input")

// JobRecord is the serialized, retained state of one invocation.
type JobRecord struct {
	Handle   int    `json:"handle"`
	Status   string `json:"status"`
	Output   Ref    `json:"output"`
	Bytes    int64  `json:"bytes"`
	ExitCode *int   `json:"exit_code,omitempty"`
	CWD      string `json:"cwd,omitempty"`
	Reason   string `json:"reason,omitempty"`
}

// WaitOptions bounds observation, never the lifetime of work.
type WaitOptions struct {
	// ReturnAvailable releases an observation on bytes already unseen at its
	// first snapshot. Later output still rechecks pattern/completion/delay.
	ReturnAvailable bool
	Delay           time.Duration
	Pattern         string
	MaxBytes        int
}

// Job owns cancellation, output and process lifetime independently of a waiter.
type Job interface {
	Done() <-chan struct{}
	io.Writer
	Parent() Agent
	Snapshot() JobRecord
	Start(context.Context, Tool, json.RawMessage)
	Wait(context.Context, WaitOptions) (string, error)
	Input(string) error
	Kill(string) error
	LaunchProcess(func() (io.WriteCloser, func() error, error)) error
	SetCWD(string)
	Finish(string, error, *int)
}

// JobTable is the root-owned registry; terminal entries remain addressable.
type JobTable interface {
	Parent() Ensemble
	Create(Agent) (Job, error)
	Get(int) (Job, error)
	Shutdown() error
}

// JobTool can transfer completion ownership to an asynchronous process reader.
type JobTool interface {
	Start(context.Context, Job, json.RawMessage) (output string, asynchronous bool, err error)
}

// ToolContext gives external handlers their cancellation and public owner chain.
type ToolContext struct {
	context.Context
	Agent Agent
	Job   Job
}

// ToolHandler is the public extension seam; returned output follows builtin policy.
type ToolHandler func(ToolContext, json.RawMessage) (string, error)
