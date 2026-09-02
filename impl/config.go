package impl

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	WebserverPort uint16  `toml:"webserver_port"`
	CorpusDir     string  `toml:"corpus_dir"`
	MetadataDir   string  `toml:"metadata_dir"`
	B             float64 `toml:"B"`
	K1            float64 `toml:"K1"`
}

var default_config = Config{
	WebserverPort: 8080,
	CorpusDir:     "corpus/",
	MetadataDir:   ".metadata/",
	B:             0.75,
	K1:            1.2,
}

func LoadConfig() Config {
	var conf Config = default_config
	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		return default_config
	}

	return conf
}

func (conf *Config) Bm25Path() string {
	return conf.MetadataDir + "/bm25.json"
}

func (conf *Config) DatabasePath() string {
	return conf.MetadataDir + "/database.db"
}
