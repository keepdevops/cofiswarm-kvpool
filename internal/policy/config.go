package policy

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config mirrors kv_auto_clear.h defaults (ML-BOTTLENECKS.md), plus per-server-group
// token budgets (the kv_token_budget gate, ported from the monolith's kv_token_semaphore).
type Config struct {
	Enabled             bool    `json:"enabled" yaml:"enabled"`
	PressureThreshold   float64 `json:"pressure_threshold" yaml:"pressure_threshold"`
	DivergenceThreshold float64 `json:"divergence_threshold" yaml:"divergence_threshold"`
	ProactiveThreshold  float64 `json:"proactive_threshold" yaml:"proactive_threshold"`
	ProactiveFraction   float64 `json:"proactive_fraction" yaml:"proactive_fraction"`
	// Budgets maps a server_group (e.g. "llama8b") to a per-run KV token budget.
	// 0 or absent = unbudgeted (unlimited), matching kv_token_budget=0 in the monolith.
	Budgets map[string]int `json:"budgets" yaml:"budgets"`
}

func Default() Config {
	return Config{
		Enabled:             false,
		PressureThreshold:   0.75,
		DivergenceThreshold: 0.6,
		ProactiveThreshold:  0.60,
		ProactiveFraction:   0.30,
		Budgets:             map[string]int{},
	}
}

// Load reads a YAML config over the defaults: fields absent from the file keep their
// default values. A missing path returns Default() (disabled) — the safe default-off state.
func Load(path string) (Config, error) {
	cfg := Default()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	if cfg.Budgets == nil {
		cfg.Budgets = map[string]int{}
	}
	return cfg, nil
}

type EvaluateRequest struct {
	KVPressure float64 `json:"kv_pressure"`
	Query      string  `json:"query,omitempty"`
}

type EvaluateResponse struct {
	AutoClear      bool    `json:"auto_clear"`
	ProactiveEvict bool    `json:"proactive_evict"`
	Reason         string  `json:"reason"`
	ProactiveFrac  float64 `json:"proactive_fraction,omitempty"`
}

func (c Config) Evaluate(req EvaluateRequest) EvaluateResponse {
	if !c.Enabled {
		return EvaluateResponse{Reason: "disabled"}
	}
	if req.KVPressure >= c.PressureThreshold {
		return EvaluateResponse{AutoClear: true, Reason: "pressure_threshold"}
	}
	if req.KVPressure >= c.ProactiveThreshold {
		return EvaluateResponse{
			ProactiveEvict: true,
			Reason:         "proactive_threshold",
			ProactiveFrac:  c.ProactiveFraction,
		}
	}
	return EvaluateResponse{Reason: "nominal"}
}
