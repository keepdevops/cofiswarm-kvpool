package bus

import (
	"encoding/json"

	"github.com/keepdevops/cofiswarm-kvpool/internal/policy"
	"github.com/keepdevops/cofiswarm-observer-sdk/pkg/servicecomponent"
)

// Capability subjects (must match observer's bus/subjects.py).
const (
	SubjAdmit    = servicecomponent.Prefix + ".kv.admit"
	SubjEvaluate = servicecomponent.Prefix + ".kv.evaluate"
	SubjPolicy   = servicecomponent.Prefix + ".kv.policy"
)

// Routes wires a policy.Config to the .kv.* subjects. Reply field names mirror observer's
// bus/contracts/resource.py (KvAdmitReply / KvEvaluateReply / KvPolicyReply).
func Routes(cfg policy.Config) map[string]servicecomponent.Handler {
	return map[string]servicecomponent.Handler{
		SubjAdmit:    admitHandler(cfg),
		SubjEvaluate: evaluateHandler(cfg),
		SubjPolicy:   policyHandler(cfg),
	}
}

func admitHandler(cfg policy.Config) servicecomponent.Handler {
	return func(data []byte) (any, error) {
		var req policy.AdmitRequest
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		r := cfg.Admit(req)
		return admitReply{SchemaVersion: servicecomponent.SchemaVersion, OK: true,
			Allowed: r.Allowed, Budget: r.Budget, Reason: r.Reason}, nil
	}
}

func evaluateHandler(cfg policy.Config) servicecomponent.Handler {
	return func(data []byte) (any, error) {
		var req policy.EvaluateRequest
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		r := cfg.Evaluate(req)
		return evalReply{SchemaVersion: servicecomponent.SchemaVersion, OK: true,
			AutoClear: r.AutoClear, ProactiveEvict: r.ProactiveEvict,
			Reason: r.Reason, ProactiveFraction: r.ProactiveFrac}, nil
	}
}

func policyHandler(cfg policy.Config) servicecomponent.Handler {
	return func([]byte) (any, error) {
		return policyReply{SchemaVersion: servicecomponent.SchemaVersion, OK: true,
			Enabled: cfg.Enabled, PressureThreshold: cfg.PressureThreshold,
			DivergenceThreshold: cfg.DivergenceThreshold, ProactiveThreshold: cfg.ProactiveThreshold,
			ProactiveFraction: cfg.ProactiveFraction, Budgets: cfg.Budgets}, nil
	}
}

type admitReply struct {
	SchemaVersion string `json:"schema_version"`
	OK            bool   `json:"ok"`
	Error         string `json:"error,omitempty"`
	Allowed       bool   `json:"allowed"`
	Budget        int    `json:"budget"`
	Reason        string `json:"reason"`
}

type evalReply struct {
	SchemaVersion     string  `json:"schema_version"`
	OK                bool    `json:"ok"`
	Error             string  `json:"error,omitempty"`
	AutoClear         bool    `json:"auto_clear"`
	ProactiveEvict    bool    `json:"proactive_evict"`
	Reason            string  `json:"reason"`
	ProactiveFraction float64 `json:"proactive_fraction,omitempty"`
}

type policyReply struct {
	SchemaVersion       string         `json:"schema_version"`
	OK                  bool           `json:"ok"`
	Error               string         `json:"error,omitempty"`
	Enabled             bool           `json:"enabled"`
	PressureThreshold   float64        `json:"pressure_threshold"`
	DivergenceThreshold float64        `json:"divergence_threshold"`
	ProactiveThreshold  float64        `json:"proactive_threshold"`
	ProactiveFraction   float64        `json:"proactive_fraction"`
	Budgets             map[string]int `json:"budgets"`
}
