package utils

// Config is generated when the program launches.
// It should never be called until InitConfig() is done.
var config FullConfig
var initialized bool

type argsCommandLine struct {
	DB_URI string
	MODE   string // "prod", "dev", "test"
}

type argsConfigFile struct {
	DB_URI string
	MODE   string // "prod", "dev", "test"
}

type FullConfig struct {
	DB_URI string
	MODE   string // "prod", "dev", "test"
}

func InitConfig(path string) error {
	if !initialized {
		config = FullConfig{
			DB_URI: "file:database.db",
			MODE:   "dev",
		}
		initialized = true
	}
	return nil
}

func Config() *FullConfig {
	return &config
}
