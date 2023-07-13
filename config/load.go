package config

import (
	"game-app/pkg/errmsg"
	"game-app/pkg/richerror"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"strings"
)

func Load(configPath string) *Config {
	const op = "config.Load"
	var k = koanf.New(".")

	k.Load(confmap.Provider(defaultConfig, "."), nil)

	k.Load(file.Provider(configPath), yaml.Parser())

	k.Load(env.Provider("GAMEAPP_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "GAMEAPP_")), "_", ".", -1)
	}), nil)

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		panic(richerror.New(op).
			WithErr(err).WithKind(richerror.KindUnexpected).
			WithMessage(errmsg.ErrorMsgSomethingWentWrong))
	}

	return &cfg
}
