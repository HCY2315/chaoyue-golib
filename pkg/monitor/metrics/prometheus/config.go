package prometheus

// Config ...
type Config struct {
	PromePath       string   `json:"prome_path" yaml:"prome_path" mapstructure:"prome_path"`
	Registry        string   `json:"registry" yaml:"registry" mapstructure:"registry"`
	RegistryAddress []string `json:"registry_address" yaml:"registry_address" mapstructure:"registry_address"`
}
