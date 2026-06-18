package policy

import (
	"os"
	"path/filepath"
	"testing"
)

func armed() Config {
	c := Default()
	c.Enabled = true
	c.Budgets = map[string]int{"llama8b": 3072}
	return c
}

func TestEvaluateThresholds(t *testing.T) {
	c := armed()
	cases := []struct {
		pressure       float64
		autoClear      bool
		proactiveEvict bool
		reason         string
	}{
		{0.50, false, false, "nominal"},
		{0.60, false, true, "proactive_threshold"},
		{0.74, false, true, "proactive_threshold"},
		{0.75, true, false, "pressure_threshold"},
		{0.90, true, false, "pressure_threshold"},
	}
	for _, tc := range cases {
		got := c.Evaluate(EvaluateRequest{KVPressure: tc.pressure})
		if got.AutoClear != tc.autoClear || got.ProactiveEvict != tc.proactiveEvict || got.Reason != tc.reason {
			t.Errorf("pressure %.2f -> %+v, want clear=%v evict=%v reason=%q",
				tc.pressure, got, tc.autoClear, tc.proactiveEvict, tc.reason)
		}
	}
}

func TestEvaluateDisabledIsNoop(t *testing.T) {
	if got := Default().Evaluate(EvaluateRequest{KVPressure: 0.99}); got.AutoClear || got.Reason != "disabled" {
		t.Errorf("disabled config fired: %+v", got)
	}
}

func TestAdmitBudget(t *testing.T) {
	c := armed()
	if r := c.Admit(AdmitRequest{Group: "llama8b", Tokens: 3000}); !r.Allowed || r.Reason != "within_budget" {
		t.Errorf("within budget rejected: %+v", r)
	}
	if r := c.Admit(AdmitRequest{Group: "llama8b", Tokens: 4000}); r.Allowed || r.Reason != "token_budget_exceeded" {
		t.Errorf("over budget admitted: %+v", r)
	}
	if r := c.Admit(AdmitRequest{Group: "gemma2b", Tokens: 99999}); !r.Allowed || r.Reason != "unbudgeted" {
		t.Errorf("unbudgeted group gated: %+v", r)
	}
}

func TestLoadMergesDefaultsAndMissingFile(t *testing.T) {
	// missing file -> defaults, disabled
	if c, err := Load(filepath.Join(t.TempDir(), "nope.yaml")); err != nil || c.Enabled {
		t.Fatalf("missing file: cfg=%+v err=%v", c, err)
	}
	// partial file -> defaults preserved for unspecified fields
	p := filepath.Join(t.TempDir(), "k.yaml")
	if err := os.WriteFile(p, []byte("enabled: true\nbudgets:\n  llama8b: 2048\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	c, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if !c.Enabled || c.ProactiveThreshold != 0.60 || c.Budgets["llama8b"] != 2048 {
		t.Errorf("load merge wrong: %+v", c)
	}
}
