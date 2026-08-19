package impl

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	WebserverPort uint16 `toml:"webserver_port"`
	CorpusDir     string `toml:"corpus_dir"`
	MetadataDir   string `toml:"metadata_dir"`
}

var default_config = Config{
	WebserverPort: 8080,
	CorpusDir:     "corpus/",
	MetadataDir:   ".metadata/",
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
