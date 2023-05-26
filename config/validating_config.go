package config

// ValidatingConfig is interface that contains config struct and Validate method to check whether config can be used further or not
type ValidatingConfig interface {
	Validate() error
}
