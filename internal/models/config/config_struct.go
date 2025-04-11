package config

type DatabaseConfig struct {
	User     string
	Password string
	DBName   string
	Host     string
	Port     int
	SSLMode  string
}

type LogConfig struct {
	Level       string   `mapstructure:"level"`
	Format      string   `mapstructure:"format"`
	OutputPaths []string `mapstructure:"output_paths"`
	Caller      bool     `mapstructure:"caller"`
	Stacktrace  bool     `mapstructure:"stacktrace"`
}

type Config struct {
	Database DatabaseConfig
	Log      LogConfig `mapstructure:"log"`
}
