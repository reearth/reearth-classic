package config

type PublishedGatewayConfig struct {
	Token string `default:"" pp:",omitempty"`
	// PreviousToken is accepted alongside Token for the duration of a secret
	// rotation, so rolling the value forward doesn't 401 callers (reearth-cloud)
	// that haven't redeployed with the new value yet. Clear it once every
	// caller has picked up the current Token.
	PreviousToken string `default:"" pp:",omitempty"`
}
