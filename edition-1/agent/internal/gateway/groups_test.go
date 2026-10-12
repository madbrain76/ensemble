package gateway

import (
	"path/filepath"
	"testing"

	pub "github.com/waywardgeek/ensemble/agent/gateway"
)

// testKey is a syntactically valid 32-byte public key. The registry only
// parses keys; it never performs the handshake, so distinct bytes are all
// these tests need.
func testKey(b byte) string {
	raw := make([]byte, pub.PublicKeySize)
	for i := range raw {
		raw[i] = b
	}
	return pub.EncodePublicKey(raw)
}

func testRegistry(t *testing.T) *Registry {
	t.Helper()
	cfg := &Config{
		Gateway: ServerConfig{Listen: ":0", DataDir: t.TempDir()},
		Agents: map[string]AgentConfig{
			"ensemble": {DisplayName: "Ensemble", PublicKey: testKey(1), Groups: []string{"admin", "agents"}},
			"plain":    {DisplayName: "Plain", PublicKey: testKey(2), Groups: []string{"agents"}},
		},
		Groups: map[string]GroupConfig{
			"admin":        {Capabilities: []string{pub.CapGroupAdmin}, Trust: pub.TrustOwner},
			"agents":       {Capabilities: []string{pub.CapSendMessage}, Trust: pub.TrustTrusted},
			"interviewers": {Capabilities: []string{pub.CapSendMessage}, Trust: pub.TrustExternal},
		},
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("fixture config invalid: %v", err)
	}
	return NewRegistry(cfg)
}

// --- contactability and routing authorization (deliverable 4) ---

func TestContactIsOptIn(t *testing.T) {
	r := testRegistry(t)
	// "ensemble" has declared nothing, so even a configured peer in a real
	// group cannot reach it. Inbound access is opt-in.
	if _, _, err := r.Authorize("plain", "ensemble"); err == nil {
		t.Fatalf("expected contact to be refused before any declaration")
	}

	if err := r.DeclareContactable("ensemble", []string{"agents"}); err != nil {
		t.Fatalf("declare: %v", err)
	}
	via, trust, err := r.Authorize("plain", "ensemble")
	if err != nil {
		t.Fatalf("positive control failed after declaring: %v", err)
	}
	if len(via) != 1 || via[0] != "agents" {
		t.Errorf("authorized via %v, want [agents]", via)
	}
	if trust != pub.TrustTrusted {
		t.Errorf("trust = %q, want %q", trust, pub.TrustTrusted)
	}
}

func TestAuthorizeRejectsSenderOutsideAcceptedGroups(t *testing.T) {
	r := testRegistry(t)
	// Ensemble accepts only interviewers; "plain" is in agents.
	if err := r.DeclareContactable("ensemble", []string{"interviewers"}); err != nil {
		t.Fatalf("declare: %v", err)
	}
	if _, _, err := r.Authorize("plain", "ensemble"); err == nil {
		t.Fatalf("a sender holding no accepted group must be refused")
	}

	// Positive control: granting the accepted group makes the same call work.
	if err := r.RegisterChild("ensemble", "kid", testKey(9), nil); err != nil {
		t.Fatalf("register: %v", err)
	}
	if _, _, err := r.Authorize("kid", "ensemble"); err == nil {
		t.Fatalf("child with no groups must also be refused")
	}
}

// TestTrustIsHighestAuthorizingGroup pins the rule that trust comes from the
// relationship the two agents actually share.
func TestTrustIsHighestAuthorizingGroup(t *testing.T) {
	r := testRegistry(t)
	if err := r.DeclareContactable("plain", []string{"admin", "interviewers"}); err != nil {
		t.Fatalf("declare: %v", err)
	}
	// "ensemble" is in admin (owner) and agents (trusted). Only admin is
	// accepted by "plain", so the contact is an owner-trust contact.
	via, trust, err := r.Authorize("ensemble", "plain")
	if err != nil {
		t.Fatalf("authorize: %v", err)
	}
	if trust != pub.TrustOwner {
		t.Fatalf("trust = %q via %v, want owner", trust, via)
	}
}

// TestUndeclaredGroupTrustDefaultsToExternal: forgetting to declare trust in
// config must not accidentally confer it.
func TestUndeclaredGroupTrustDefaultsToExternal(t *testing.T) {
	if got := (GroupConfig{}).TrustLevel(); got != pub.TrustExternal {
		t.Fatalf("default trust = %q, want external", got)
	}
}

// --- parent registers child, narrow only (deliverable 7) ---

func TestRegisterChildGrantsSubsetOfParentGroups(t *testing.T) {
	r := testRegistry(t)
	if err := r.RegisterChild("ensemble", "interviewer", testKey(3), []string{"agents"}); err != nil {
		t.Fatalf("granting a group the parent holds must succeed: %v", err)
	}
	groups := r.GroupsOf("interviewer")
	if len(groups) != 1 || groups[0] != "agents" {
		t.Errorf("child groups = %v, want [agents]", groups)
	}
	if parent, ok := r.ParentOf("interviewer"); !ok || parent != "ensemble" {
		t.Errorf("parent = %q/%v, want ensemble", parent, ok)
	}
	// The child's key must now resolve, or it could never connect.
	if _, err := r.PublicKey("interviewer"); err != nil {
		t.Errorf("child key should resolve: %v", err)
	}
}

func TestRegisterChildCannotWiden(t *testing.T) {
	r := testRegistry(t)
	// "plain" holds only agents, so it must not be able to hand out admin.
	err := r.RegisterChild("plain", "escalated", testKey(4), []string{"admin"})
	if err == nil {
		t.Fatalf("a parent must not grant a group it does not hold")
	}
	if r.Known("escalated") {
		t.Errorf("a refused registration must not create the agent")
	}
}

func TestGrantMembershipRequiresParentage(t *testing.T) {
	r := testRegistry(t)
	if err := r.RegisterChild("ensemble", "kid", testKey(5), nil); err != nil {
		t.Fatalf("register: %v", err)
	}
	// Positive control: the real parent may grant.
	if err := r.GrantMembership("ensemble", "kid", "agents"); err != nil {
		t.Fatalf("parent grant should succeed: %v", err)
	}
	// A third agent must not be able to administer someone else's child.
	if err := r.GrantMembership("plain", "kid", "agents"); err == nil {
		t.Fatalf("a non-parent must not grant memberships")
	}
	if err := r.RevokeMembership("plain", "kid", "agents"); err == nil {
		t.Fatalf("a non-parent must not revoke memberships")
	}
	if !r.HasGroup("kid", "agents") {
		t.Errorf("the failed revoke must not have taken effect")
	}
}

func TestRegisterChildCannotHijackConfiguredAgent(t *testing.T) {
	r := testRegistry(t)
	// Re-registering a config agent would let a runtime caller replace the
	// public key of an identity the operator declared.
	if err := r.RegisterChild("ensemble", "plain", testKey(6), nil); err == nil {
		t.Fatalf("a configured agent must not be re-registerable")
	}
	if err := r.RegisterChild("ensemble", "dup", testKey(7), nil); err != nil {
		t.Fatalf("register: %v", err)
	}
	if err := r.RegisterChild("plain", "dup", testKey(8), nil); err == nil {
		t.Fatalf("a second parent must not be able to claim an existing child")
	}
}

// --- group administration ---

func TestGroupAdminRequiresCapability(t *testing.T) {
	r := testRegistry(t)
	// Positive control: ensemble is in admin, which carries group-admin.
	if err := r.CreateGroup("ensemble", "ops", GroupConfig{Trust: pub.TrustTrusted}); err != nil {
		t.Fatalf("group-admin should be able to create: %v", err)
	}
	if err := r.CreateGroup("plain", "sneaky", GroupConfig{}); err == nil {
		t.Fatalf("an agent without group-admin must not create groups")
	}
	if _, ok := r.Groups()["sneaky"]; ok {
		t.Errorf("the refused group must not exist")
	}
}

func TestConfigGroupsAreImmutableAtRuntime(t *testing.T) {
	r := testRegistry(t)
	// Editing a config group would revert on restart, which is worse than
	// refusing because it works until it does not.
	if err := r.EditGroup("ensemble", "admin", GroupConfig{}); err == nil {
		t.Errorf("config groups must not be editable at runtime")
	}
	if err := r.DeleteGroup("ensemble", "admin"); err == nil {
		t.Errorf("config groups must not be deletable at runtime")
	}
	// Positive control: a runtime group is editable and deletable.
	if err := r.CreateGroup("ensemble", "ops", GroupConfig{}); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := r.EditGroup("ensemble", "ops", GroupConfig{Description: "x"}); err != nil {
		t.Errorf("runtime group should be editable: %v", err)
	}
	if err := r.DeleteGroup("ensemble", "ops"); err != nil {
		t.Errorf("runtime group should be deletable: %v", err)
	}
}

func TestDeleteGroupClearsDanglingReferences(t *testing.T) {
	r := testRegistry(t)
	if err := r.CreateGroup("ensemble", "ops", GroupConfig{}); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := r.DeclareContactable("ensemble", []string{"ops"}); err != nil {
		t.Fatalf("declare: %v", err)
	}
	if err := r.DeleteGroup("ensemble", "ops"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if got := r.ContactableBy("ensemble"); len(got) != 0 {
		t.Errorf("contactability still references the deleted group: %v", got)
	}
}

// --- persistence ---

func TestRegistryStateSurvivesRestart(t *testing.T) {
	r := testRegistry(t)
	if err := r.RegisterChild("ensemble", "interviewer", testKey(3), []string{"agents"}); err != nil {
		t.Fatalf("register: %v", err)
	}
	if err := r.DeclareContactable("ensemble", []string{"agents"}); err != nil {
		t.Fatalf("declare: %v", err)
	}

	path := filepath.Join(t.TempDir(), "state", "registry.json")
	if err := r.Save(path); err != nil {
		t.Fatalf("save: %v", err)
	}

	// A fresh registry over the same config, as after a restart.
	restored := testRegistry(t)
	if restored.Known("interviewer") {
		t.Fatalf("precondition: fresh registry should not know the child yet")
	}
	if err := restored.Load(path); err != nil {
		t.Fatalf("load: %v", err)
	}

	if !restored.Known("interviewer") {
		t.Errorf("runtime-registered child did not survive restart")
	}
	if parent, ok := restored.ParentOf("interviewer"); !ok || parent != "ensemble" {
		t.Errorf("parentage did not survive restart: %q/%v", parent, ok)
	}
	if _, _, err := restored.Authorize("interviewer", "ensemble"); err != nil {
		t.Errorf("contactability did not survive restart: %v", err)
	}
}

func TestLoadMissingStateFileIsNotAnError(t *testing.T) {
	r := testRegistry(t)
	if err := r.Load(filepath.Join(t.TempDir(), "absent.json")); err != nil {
		t.Fatalf("a first start has no state file: %v", err)
	}
}
