package bus

import (
	"encoding/json"
	"testing"

	"github.com/keepdevops/cofiswarm-kvpool/internal/policy"
)

func armed() policy.Config {
	return policy.Config{Enabled: true, PressureThreshold: 0.75, ProactiveThreshold: 0.60,
		ProactiveFraction: 0.30, Budgets: map[string]int{"g": 100}}
}

func TestAdmitRouteWithinBudget(t *testing.T) {
	out, err := Routes(armed())[SubjAdmit]([]byte(`{"group":"g","tokens":50}`))
	if err != nil {
		t.Fatalf("admit: %v", err)
	}
	r := out.(admitReply)
	if !r.Allowed || r.Reason != "within_budget" {
		t.Fatalf("got %+v", r)
	}
}

func TestAdmitRouteExceeded(t *testing.T) {
	out, _ := Routes(armed())[SubjAdmit]([]byte(`{"group":"g","tokens":500}`))
	if r := out.(admitReply); r.Allowed || r.Reason != "token_budget_exceeded" {
		t.Fatalf("got %+v", r)
	}
}

func TestEvaluateRouteFiresOnPressure(t *testing.T) {
	out, _ := Routes(armed())[SubjEvaluate]([]byte(`{"kv_pressure":0.9}`))
	if r := out.(evalReply); !r.AutoClear || r.Reason != "pressure_threshold" {
		t.Fatalf("got %+v", r)
	}
}

func TestPolicyRouteReturnsConfig(t *testing.T) {
	out, _ := Routes(armed())[SubjPolicy](nil)
	r := out.(policyReply)
	if !r.Enabled || r.Budgets["g"] != 100 {
		t.Fatalf("got %+v", r)
	}
}

func TestReplyFieldNamesMatchSchema(t *testing.T) {
	out, _ := Routes(armed())[SubjAdmit]([]byte(`{"group":"g","tokens":1}`))
	b, _ := json.Marshal(out)
	var m map[string]any
	_ = json.Unmarshal(b, &m)
	for _, k := range []string{"schema_version", "ok", "allowed", "budget", "reason"} {
		if _, present := m[k]; !present {
			t.Fatalf("reply missing field %q: %s", k, b)
		}
	}
}
