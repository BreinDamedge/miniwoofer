package impl

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	WebserverPort uint16 `toml:"webserver_port"`
	CorpusDir     string `toml:"corpus_dir"`
	Bm25File      string `toml:"bm25_file"`
}

var default_config = Config{
	WebserverPort: 8080,
	CorpusDir:     "corpus/",
	Bm25File:      ".metadata/bm25.json",
}

func LoadConfig() Config {
	var conf Config = default_config
	if _, err := toml.DecodeFile("boogaloo.toml", &conf); err != nil {
		return default_config
	}

	return conf
}
