package config

type Config struct {
	Server ServerConfig `yaml:"server"`
}

type ServerConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	BasePath string `yaml:"base_path"`
}
