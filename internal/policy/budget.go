package policy

// AdmitRequest asks whether a run of `tokens` tokens on `group` fits its KV token budget.
type AdmitRequest struct {
	Group  string `json:"group"`
	Tokens int    `json:"tokens"`
}

// AdmitResponse is the budget decision (the kv_token_budget gate).
type AdmitResponse struct {
	Allowed bool   `json:"allowed"`
	Budget  int    `json:"budget"`
	Reason  string `json:"reason"`
}

// Admit gates a request against the per-group token budget. A budget of 0 (or an
// unlisted group) means unbudgeted — always allowed. Otherwise tokens must fit.
func (c Config) Admit(req AdmitRequest) AdmitResponse {
	budget := c.Budgets[req.Group]
	if budget <= 0 {
		return AdmitResponse{Allowed: true, Budget: budget, Reason: "unbudgeted"}
	}
	if req.Tokens <= budget {
		return AdmitResponse{Allowed: true, Budget: budget, Reason: "within_budget"}
	}
	return AdmitResponse{Allowed: false, Budget: budget, Reason: "token_budget_exceeded"}
}
