package policy

// Config mirrors kv_auto_clear.h defaults (ML-BOTTLENECKS.md).
type Config struct {
	Enabled              bool    `json:"enabled"`
	PressureThreshold    float64 `json:"pressure_threshold"`
	DivergenceThreshold  float64 `json:"divergence_threshold"`
	ProactiveThreshold   float64 `json:"proactive_threshold"`
	ProactiveFraction    float64 `json:"proactive_fraction"`
}

func Default() Config {
	return Config{
		Enabled:             false,
		PressureThreshold:   0.75,
		DivergenceThreshold: 0.6,
		ProactiveThreshold:  0.60,
		ProactiveFraction:   0.30,
	}
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
