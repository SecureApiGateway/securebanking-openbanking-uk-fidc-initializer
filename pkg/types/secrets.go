package types

// GoogleSecretStore represents a secret store in AM which is backed by Google Secrets Manager
type GoogleSecretStore struct {
	Name                  string
	ServiceAccount        string
	Project               string
	SecretFormat          string
	ExpiryDurationSeconds int
}

// SecretMapping maps an AM SecredId to a secret alias, for GoogleSecretStore this alias is the secret name in Secrets Manager
type SecretMapping struct {
	SecretId string
	Alias    string
}
