package props

// Properties generic application properties
type Properties struct {
	Server ServerProps `yaml:"server"`
}

// ServerProps server configuration
type ServerProps struct {
	ContextRoot string `yaml:"context-root" envconfig:"SERVER_CONTEXT_ROOT"`
	Host        string `yaml:"host" envconfig:"HOST"`
	Port        int    `yaml:"port" envconfig:"PORT"`
}
