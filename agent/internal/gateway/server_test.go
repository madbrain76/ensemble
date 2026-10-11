package gateway

import (
	"context"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/flynn/noise"
	pub "github.com/waywardgeek/ensemble/agent/gateway"
)

// --- harness ---

// recorder is a Handler that captures what the gateway delivers.
type recorder struct {
	pub.NopHandler
	mu        sync.Mutex
	delivered []pub.Envelope
	connected chan pub.WelcomePayload
	inbox     chan pub.Envelope
}

func newRecorder() *recorder {
	return &recorder{
		connected: make(chan pub.WelcomePayload, 8),
		inbox:     make(chan pub.Envelope, 32),
	}
}

func (r *recorder) OnDeliver(env pub.Envelope) {
	r.mu.Lock()
	r.delivered = append(r.delivered, env)
	r.mu.Unlock()
	r.inbox <- env
}

func (r *recorder) OnConnected(w pub.WelcomePayload) { r.connected <- w }

func (r *recorder) count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.delivered)
}

// await returns the next delivered envelope or fails the test.
func (r *recorder) await(t *testing.T, what string) pub.Envelope {
	t.Helper()
	select {
	case env := <-r.inbox:
		return env
	case <-time.After(5 * time.Second):
		t.Fatalf("timed out waiting for %s", what)
		return pub.Envelope{}
	}
}

func (r *recorder) awaitWelcome(t *testing.T) pub.WelcomePayload {
	t.Helper()
	select {
	case w := <-r.connected:
		return w
	case <-time.After(5 * time.Second):
		t.Fatalf("timed out waiting for welcome")
		return pub.WelcomePayload{}
	}
}

// fixture is a running gateway with keys for each named agent.
type fixture struct {
	t      *testing.T
	server *Server
	http   *httptest.Server
	wsURL  string
	keys   map[string]noise.DHKey
}

func newFixture(t *testing.T, agents map[string][]string) *fixture {
	t.Helper()

	gwKey, err := pub.GenerateKeypair()
	if err != nil {
		t.Fatalf("gateway key: %v", err)
	}

	keys := map[string]noise.DHKey{}
	cfgAgents := map[string]AgentConfig{}
	for id, groups := range agents {
		k, err := pub.GenerateKeypair()
		if err != nil {
			t.Fatalf("key for %s: %v", id, err)
		}
		keys[id] = k
		cfgAgents[id] = AgentConfig{
			DisplayName: strings.ToUpper(id[:1]) + id[1:],
			PublicKey:   pub.EncodePublicKey(k.Public),
			Groups:      groups,
		}
	}

	cfg := &Config{
		Gateway: ServerConfig{Listen: ":0", DataDir: t.TempDir()},
		Agents:  cfgAgents,
		Groups: map[string]GroupConfig{
			"admin":        {Capabilities: []string{pub.CapGroupAdmin, pub.CapSendMessage, pub.CapDeclareContactable}, Trust: pub.TrustOwner},
			"agents":       {Capabilities: []string{pub.CapSendMessage, pub.CapDeclareContactable}, Trust: pub.TrustTrusted},
			"interviewers": {Capabilities: []string{pub.CapSendMessage, pub.CapDeclareContactable}, Trust: pub.TrustExternal},
		},
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("fixture config: %v", err)
	}

	srv, err := NewServer(cfg, gwKey)
	if err != nil {
		t.Fatalf("NewServer: %v", err)
	}
	srv.SetLogger(func(string, ...any) {})

	hs := httptest.NewServer(srv.Handler())
	t.Cleanup(func() {
		srv.Close()
		hs.Close()
	})

	return &fixture{
		t:      t,
		server: srv,
		http:   hs,
		wsURL:  "ws" + strings.TrimPrefix(hs.URL, "http") + "/v1/ws",
		keys:   keys,
	}
}

// connect dials the gateway as agentID and waits for the welcome.
func (f *fixture) connect(agentID string, h pub.Handler) *pub.Client {
	f.t.Helper()
	return f.connectResuming(agentID, h, 0)
}

func (f *fixture) connectResuming(agentID string, h pub.Handler, lastAcked uint64) *pub.Client {
	f.t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	c, err := pub.Connect(ctx, pub.ClientConfig{
		GatewayURL:    f.wsURL,
		AgentID:       agentID,
		StaticKey:     f.keys[agentID],
		GatewayPubKey: f.server.PublicKey(),
		LastAckedSeq:  lastAcked,
	}, h)
	if err != nil {
		f.t.Fatalf("connect %s: %v", agentID, err)
	}
	f.t.Cleanup(func() { _ = c.Close() })
	return c
}

func ctx5(t *testing.T) context.Context {
	t.Helper()
	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)
	return c
}

// --- deliverable 3: the server accepts authenticated connections ---

func TestHandshakeAndWelcome(t *testing.T) {
	f := newFixture(t, map[string][]string{"ensemble": {"admin", "agents"}})
	rec := newRecorder()
	f.connect("ensemble", rec)

	w := rec.awaitWelcome(t)
	if w.AgentID != "ensemble" {
		t.Errorf("welcome agent = %q, want ensemble", w.AgentID)
	}
	if w.ServerVersion != ServerVersion {
		t.Errorf("server version = %q, want %q", w.ServerVersion, ServerVersion)
	}
	if online := f.server.Online(); len(online) != 1 || online[0] != "ensemble" {
		t.Errorf("online = %v, want [ensemble]", online)
	}
}

// TestWrongKeyIsRejected is deliverable 1's negative half, proved through
// the real server rather than against the handshake in isolation.
func TestWrongKeyIsRejected(t *testing.T) {
	f := newFixture(t, map[string][]string{"ensemble": {"agents"}})

	// Positive control: the real key connects.
	f.connect("ensemble", newRecorder())

	// An impostor holding a different private key cannot authenticate as
	// ensemble, even though it knows the agent id and the gateway's key.
	impostor, err := pub.GenerateKeypair()
	if err != nil {
		t.Fatalf("keygen: %v", err)
	}
	_, err = pub.Connect(ctx5(t), pub.ClientConfig{
		GatewayURL:       f.wsURL,
		AgentID:          "ensemble",
		StaticKey:        impostor,
		GatewayPubKey:    f.server.PublicKey(),
		HandshakeTimeout: 2 * time.Second,
	}, newRecorder())
	if err == nil {
		t.Fatalf("an impostor completed the handshake")
	}
}

func TestUnknownAgentCannotConnect(t *testing.T) {
	f := newFixture(t, map[string][]string{"ensemble": {"agents"}})
	stranger, _ := pub.GenerateKeypair()
	_, err := pub.Connect(ctx5(t), pub.ClientConfig{
		GatewayURL:       f.wsURL,
		AgentID:          "stranger",
		StaticKey:        stranger,
		GatewayPubKey:    f.server.PublicKey(),
		HandshakeTimeout: 2 * time.Second,
	}, newRecorder())
	if err == nil {
		t.Fatalf("an agent absent from config completed the handshake")
	}
}

// --- deliverables 2 and 4: routing, provenance, trust ---

func TestRoutingDeliversWithProvenanceAndTrust(t *testing.T) {
	f := newFixture(t, map[string][]string{
		"ensemble":    {"admin", "agents"},
		"interviewer": {"interviewers"},
	})

	ensRec, intRec := newRecorder(), newRecorder()
	ens := f.connect("ensemble", ensRec)
	intClient := f.connect("interviewer", intRec)
	ensRec.awaitWelcome(t)
	intRec.awaitWelcome(t)

	if err := ens.DeclareContactable(ctx5(t), []string{"interviewers"}); err != nil {
		t.Fatalf("declare: %v", err)
	}

	if err := intClient.Send(ctx5(t), pub.OutboundMessage{
		ToAgent:      "ensemble",
		Content:      "What slows you down?",
		ResponseMode: pub.ResponseAsync,
	}); err != nil {
		t.Fatalf("send: %v", err)
	}

	env := ensRec.await(t, "the delivered prompt")
	if !strings.HasPrefix(env.Content, "[GATEWAY ") {
		t.Errorf("delivered content is missing its provenance header: %q", env.Content)
	}
	if !strings.Contains(env.Content, "What slows you down?") {
		t.Errorf("delivered content lost the body: %q", env.Content)
	}
	if env.Source.Trust != pub.TrustExternal {
		t.Errorf("trust = %q, want external (the interviewers group's level)", env.Source.Trust)
	}
	if env.Source.ChannelID != "interviewer" {
		t.Errorf("source channel = %q, want interviewer", env.Source.ChannelID)
	}
	if env.Seq == 0 {
		t.Errorf("delivered envelope has no sequence number")
	}
}

// TestUnauthorizedSendIsRejected is deliverable 4's refusal half.
func TestUnauthorizedSendIsRejected(t *testing.T) {
	f := newFixture(t, map[string][]string{
		"ensemble": {"admin", "agents"},
		"outsider": {"interviewers"},
	})
	ensRec, outRec := newRecorder(), newRecorder()
	ens := f.connect("ensemble", ensRec)
	out := f.connect("outsider", outRec)
	ensRec.awaitWelcome(t)
	outRec.awaitWelcome(t)

	// Ensemble accepts only "agents"; outsider is in "interviewers".
	if err := ens.DeclareContactable(ctx5(t), []string{"agents"}); err != nil {
		t.Fatalf("declare: %v", err)
	}
	if err := out.Send(ctx5(t), pub.OutboundMessage{ToAgent: "ensemble", Content: "let me in"}); err == nil {
		t.Fatalf("send from an unaccepted group must be refused")
	}
	if n := ensRec.count(); n != 0 {
		t.Errorf("a refused message was delivered anyway (%d delivered)", n)
	}
	// A refused send must not leave a durable obligation either.
	if d := f.server.Queue().Depth("ensemble"); d != 0 {
		t.Errorf("a refused message was queued (%d pending)", d)
	}
}

// TestResponseRoutesBackToSender closes the loop: a reply reaches whoever
// sent the prompt, without the replier needing its own contactability grant.
func TestResponseRoutesBackToSender(t *testing.T) {
	f := newFixture(t, map[string][]string{
		"ensemble":    {"admin", "agents"},
		"interviewer": {"interviewers"},
	})
	ensRec, intRec := newRecorder(), newRecorder()
	ens := f.connect("ensemble", ensRec)
	intc := f.connect("interviewer", intRec)
	ensRec.awaitWelcome(t)
	intRec.awaitWelcome(t)

	if err := ens.DeclareContactable(ctx5(t), []string{"interviewers"}); err != nil {
		t.Fatalf("declare: %v", err)
	}
	if err := intc.Send(ctx5(t), pub.OutboundMessage{
		ToAgent: "ensemble", Content: "What slows you down?", ResponseMode: pub.ResponseAsync,
	}); err != nil {
		t.Fatalf("send: %v", err)
	}

	prompt := ensRec.await(t, "the prompt")
	if err := ens.Respond(ctx5(t), prompt.MessageID, "Waiting on builds."); err != nil {
		t.Fatalf("respond: %v", err)
	}

	reply := intRec.await(t, "the reply")
	if !strings.Contains(reply.Content, "Waiting on builds.") {
		t.Errorf("reply content = %q", reply.Content)
	}
	if reply.Source.ChannelID != "ensemble" {
		t.Errorf("reply source = %q, want ensemble", reply.Source.ChannelID)
	}
}

// --- deliverable 5: queue, ack, resume ---

func TestMessageQueuedWhileOfflineThenDeliveredOnConnect(t *testing.T) {
	f := newFixture(t, map[string][]string{
		"ensemble": {"admin", "agents"},
		"sender":   {"agents"},
	})

	// Ensemble declares contactability, then goes offline.
	ensRec := newRecorder()
	ens := f.connect("ensemble", ensRec)
	ensRec.awaitWelcome(t)
	if err := ens.DeclareContactable(ctx5(t), []string{"agents"}); err != nil {
		t.Fatalf("declare: %v", err)
	}
	_ = ens.Close()
	waitUntil(t, func() bool { return len(f.server.Online()) == 0 }, "ensemble to go offline")

	sndRec := newRecorder()
	snd := f.connect("sender", sndRec)
	sndRec.awaitWelcome(t)
	if err := snd.Send(ctx5(t), pub.OutboundMessage{ToAgent: "ensemble", Content: "while you were out"}); err != nil {
		t.Fatalf("send: %v", err)
	}
	if d := f.server.Queue().Depth("ensemble"); d != 1 {
		t.Fatalf("queue depth = %d, want 1", d)
	}

	// Reconnecting with no resume point must replay the backlog.
	backRec := newRecorder()
	f.connect("ensemble", backRec)
	w := backRec.awaitWelcome(t)
	if w.QueuedCount != 1 {
		t.Errorf("welcome reported %d queued, want 1", w.QueuedCount)
	}
	env := backRec.await(t, "the queued message")
	if !strings.Contains(env.Content, "while you were out") {
		t.Errorf("queued content = %q", env.Content)
	}

	// The client acks after the handler returns, so the queue drains.
	waitUntil(t, func() bool { return f.server.Queue().Depth("ensemble") == 0 }, "the queue to drain after ack")
}

// TestResumeSkipsAlreadyAckedMessages proves resume is a real replay window
// and not a blanket redelivery.
func TestResumeSkipsAlreadyAckedMessages(t *testing.T) {
	f := newFixture(t, map[string][]string{
		"ensemble": {"admin", "agents"},
		"sender":   {"agents"},
	})
	ensRec := newRecorder()
	ens := f.connect("ensemble", ensRec)
	ensRec.awaitWelcome(t)
	if err := ens.DeclareContactable(ctx5(t), []string{"agents"}); err != nil {
		t.Fatalf("declare: %v", err)
	}
	_ = ens.Close()
	waitUntil(t, func() bool { return len(f.server.Online()) == 0 }, "ensemble to go offline")

	sndRec := newRecorder()
	snd := f.connect("sender", sndRec)
	sndRec.awaitWelcome(t)
	for _, body := range []string{"first", "second", "third"} {
		if err := snd.Send(ctx5(t), pub.OutboundMessage{ToAgent: "ensemble", Content: body}); err != nil {
			t.Fatalf("send %s: %v", body, err)
		}
	}
	if d := f.server.Queue().Depth("ensemble"); d != 3 {
		t.Fatalf("queue depth = %d, want 3", d)
	}

	// Reconnect claiming the first two were already handled.
	backRec := newRecorder()
	f.connectResuming("ensemble", backRec, 2)
	w := backRec.awaitWelcome(t)
	if w.QueuedCount != 1 {
		t.Fatalf("resume replayed %d messages, want 1", w.QueuedCount)
	}
	env := backRec.await(t, "the third message")
	if !strings.Contains(env.Content, "third") {
		t.Errorf("resumed the wrong message: %q", env.Content)
	}
}

// --- deliverable 7: parent registers child ---

func TestParentRegistersChildWhichThenConnects(t *testing.T) {
	f := newFixture(t, map[string][]string{"ensemble": {"admin", "agents"}})
	ensRec := newRecorder()
	ens := f.connect("ensemble", ensRec)
	ensRec.awaitWelcome(t)

	childKey, err := pub.GenerateKeypair()
	if err != nil {
		t.Fatalf("child key: %v", err)
	}
	if err := ens.RegisterChild(ctx5(t), pub.ChildRegistration{
		AgentID:   "interviewer",
		PublicKey: childKey.Public,
		Groups:    []string{"agents"},
	}); err != nil {
		t.Fatalf("register child: %v", err)
	}

	// The child is not in the config file at all, yet it can now connect
	// on its own key. That is the whole point of runtime registration.
	f.keys["interviewer"] = childKey
	childRec := newRecorder()
	child := f.connect("interviewer", childRec)
	childRec.awaitWelcome(t)

	// And it can use the group it was granted.
	if err := ens.DeclareContactable(ctx5(t), []string{"agents"}); err != nil {
		t.Fatalf("declare: %v", err)
	}
	if err := child.Send(ctx5(t), pub.OutboundMessage{ToAgent: "ensemble", Content: "hello parent"}); err != nil {
		t.Fatalf("child send: %v", err)
	}
	env := ensRec.await(t, "the child's message")
	if !strings.Contains(env.Content, "hello parent") {
		t.Errorf("content = %q", env.Content)
	}
}

// TestParentCannotGrantGroupItDoesNotHold is deliverable 7's refusal half,
// proved over the wire rather than against the registry directly.
func TestParentCannotGrantGroupItDoesNotHold(t *testing.T) {
	f := newFixture(t, map[string][]string{"plain": {"agents"}})
	rec := newRecorder()
	c := f.connect("plain", rec)
	rec.awaitWelcome(t)

	childKey, _ := pub.GenerateKeypair()

	// Positive control: a group the parent holds is granted.
	if err := c.RegisterChild(ctx5(t), pub.ChildRegistration{
		AgentID: "ok-child", PublicKey: childKey.Public, Groups: []string{"agents"},
	}); err != nil {
		t.Fatalf("granting a held group must succeed: %v", err)
	}

	// The assertion: "plain" is not in admin, so it cannot confer admin.
	err := c.RegisterChild(ctx5(t), pub.ChildRegistration{
		AgentID: "escalated", PublicKey: childKey.Public, Groups: []string{"admin"},
	})
	if err == nil {
		t.Fatalf("a parent granted a group it does not hold")
	}
	if f.server.Registry().Known("escalated") {
		t.Errorf("the refused child was created anyway")
	}
}

func waitUntil(t *testing.T, cond func() bool, what string) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %s", what)
}
