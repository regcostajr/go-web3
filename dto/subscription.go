package dto

type Subscription struct {
	ID      int                `json:"id"`
	Result  interface{}        `json:"result"`
	Method  string             `json:"method"`
	Version string             `json:"jsonrpc"`
	Params  SubscriptionParams `json:"params"`
}

type SubscriptionParams struct {
	Subscription string `json:"subscription"`
	Result       string `json:"result"`
}
