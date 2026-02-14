package provider

// MaxPrice limits the maximum price for a request.
type MaxPrice struct {
	PromptPricePerToken     float64 `json:"prompt_price_per_token,omitempty"`
	CompletionPricePerToken float64 `json:"completion_price_per_token,omitempty"`
}

// ProviderPreferences configures provider routing for OpenRouter.
type ProviderPreferences struct {
	Order                  []string  `json:"order,omitempty"`
	AllowFallbacks         *bool     `json:"allow_fallbacks,omitempty"`
	RequireParameters      *bool     `json:"require_parameters,omitempty"`
	DataCollection         string    `json:"data_collection,omitempty"`
	ZDR                    *bool     `json:"zdr,omitempty"`
	Sort                   string    `json:"sort,omitempty"`
	SortPartition          string    `json:"partition,omitempty"`
	PreferredMinThroughput any       `json:"preferred_min_throughput,omitempty"`
	PreferredMaxLatency    any       `json:"preferred_max_latency,omitempty"`
	Only                   []string  `json:"only,omitempty"`
	Ignore                 []string  `json:"ignore,omitempty"`
	MaxPrice               *MaxPrice `json:"max_price,omitempty"`
}
