package config

const (
	configPath = "config.json"
)

var ConfigObj struct {
	HashKey  string
	HashSalt string
	DBType   string
}

func loadConfig() error {
	var (
		err error
	)

	return err
}
