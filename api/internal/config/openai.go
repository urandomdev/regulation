package config

// OpenAI holds secrets for LLM access.
type OpenAI struct {
	APIKey string `json:"api_key"`
}

func (o *OpenAI) Validate() error {
	// API key can be empty to keep local onboarding simple; the server will warn later.
	return nil
}
