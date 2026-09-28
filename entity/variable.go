package entity

// Variable describes a single environment variable parsed from configuration.
type Variable struct {
	Name             string
	Type             string
	Description      string
	Default          string
	ValidateCritical bool
	ValidateRules    []string
	EnvNames         []string
	Comment          string
}

// Section groups a set of variables under a common name.
type Section struct {
	Name      string
	Variables []Variable
}
