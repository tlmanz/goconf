package main

import (
	"errors"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/tlmanz/goconf"
	"github.com/tryfix/log"
)

type Conf struct {
	Name  string `env:"MY_NAME"`
	OAuth OAuth
}

type OAuth struct {
	Enabled     bool   `env:"OAUTH_ENABLED" envDefault:"False"`
	JwtKey      string `env:"JWT_KEY" envDefault:"asdasd" hush:"mask"`
	OAuthConfig OAuthConfig
}

type OAuthConfig struct {
	ClientID     string `env:"CLIENT_ID" envDefault:"234234234-23234kjh23jk4g2j3h4gjh23g4f.apps.googleusercontent.com" hush:"hide"`
	ClientSecret string `env:"CLIENT_SECRET" envDefault:"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10U_" hush:"mask"`
}

var Config Conf

func (Conf) Register() error {
	return env.Parse(&Config)
}

func (Conf) Validate() error {
	if Config.Name == "" {
		return errors.New(`MY_NAME environmental variable cannot be empty`)
	}
	return nil
}

func (Conf) Print() interface{} {
	return Config
}

func main() {
	_ = os.Setenv("MY_NAME", "My First Configuration")

	err := goconf.Load(
		new(Conf),
	)
	if err != nil {
		log.Fatal(err)
	}
	if Config.Name != `My First Configuration` {
		log.Fatal(`error while comparing config`)
	}

	log.Info(`goconf successfully loaded`)
}
