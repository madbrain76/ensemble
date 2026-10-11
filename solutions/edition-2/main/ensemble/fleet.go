package ensemble

import (
	"ensemble/internal/common"
	"errors"
)

// AgentFactory supplies application-specific child configuration and construction.
type AgentFactory func(parent Agent, spec AgentSpec) (Agent, error)

// FleetMessage carries sender identity for questions, replies and escalation.
type FleetMessage = common.FleetMessage

// UserSink receives questions and escalation addressed to a human.
type UserSink = common.UserSink

// ChildStatus reports a direct child's explicit completion state.
type ChildStatus = common.ChildStatus

// SetAgentFactory enables fleet execution before Agents are started. The callback
// must construct through NewChildAgent so both actual parent links are retained.
func (e *Ensemble) SetAgentFactory(factory AgentFactory) { e.factory = factory }

// FleetEnabled distinguishes a fleet application from earlier single-Agent clients.
func (e *Ensemble) FleetEnabled() bool { return e.factory != nil }

// SetUserSink selects the human address and delivery boundary before startup.
func (e *Ensemble) SetUserSink(id string, sink UserSink) { e.userID = id; e.userSink = sink }

// UserID returns the configured human address, defaulting to the role user.
func (e *Ensemble) UserID() string {
	if e.userID != "" {
		return e.userID
	}
	return "user"
}

// DeliverUser reaches the application sink; an unwired human is a visible error.
func (e *Ensemble) DeliverUser(message FleetMessage) error {
	if e.userSink == nil {
		return errors.New("no user message sink configured")
	}
	return e.userSink.Deliver(message)
}

// CreateChild delegates construction, leaving provider choices in the application.
func (e *Ensemble) CreateChild(parent Agent, spec common.AgentSpec) (Agent, error) {
	if e.factory == nil {
		return nil, errors.New("no agent factory configured")
	}
	return e.factory(parent, AgentSpec(spec))
}

// NewChildAgent enforces inherited restrictions at the construction boundary.
// Routing IDs complement these pointers; they never replace actual ownership.
func (e *Ensemble) NewChildAgent(parent Agent, cfg Config, spec AgentSpec) (Agent, error) {
	if parent == nil || parent.Ensemble() != e {
		return nil, errors.New("child requires parent in this Ensemble")
	}
	narrowed, err := spec.Clamp(AgentSpec(parent.Spec()))
	if err != nil {
		return nil, err
	}
	return e.newAgent(parent, cfg, narrowed)
}
