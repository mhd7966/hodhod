package configs

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
	"os"
)

var Cfg Config

type Config struct {
	Debug bool `env:"DEBUG" env-default:"False"`
	Redis struct {
		Addr string `env:"REDIS_ADD" env-default:"redis:6379"`
	}
	Postgres struct {
		Port string `env:"PG_PORT" env-default:"5432"`
		Host string `env:"PG_HOST" env-default:"postgres"`
		Name string `env:"PG_NAME" env-default:"hodhod"`
		User string `env:"PG_USER" env-default:"admin"`
		Pass string `env:"PG_PASS" env-default:"postgres_password"`
	}
	Log struct {
		LogLevel   string `env:"LOG_LEVEL" env-default:"debug"`
		OutputType string `env:"LOG_OUTPUT_TYPE" env-default:"stdout"`
		OutputAdd  string `env:"LOG_FILE_Add" env-default:"/log.txt"`
	}
	Auth struct {
		Host string `env:"AUTH_HOST" env-default:"authorization_host_address(del)"`
	}
	Sentry struct {
		DSN   string `env:"SENTRY_DSN" env-default:"sentry_dsn_address"`
		Level string `env:"SENTRY_LEVEL" env-default:"error"`
	}
	SMS struct {
		Driver string `env:"SMS_DRIVER" env-default:"1"`
	}
	Call struct {
		Driver string `env:"CALL_DRIVER" env-default:"3"`
	}
	Kaveh struct {
		Token  string `env:"KAVEH_TOKEN" env-default:"kaveh_token"`
		Number string `env:"KAVEH_NUMBER" env-default:"1000596446"`
	}
	Signal struct {
		Token  string `env:"SIGNAL_TOKEN" env-default:"signal_token"`
		Number string `env:"SIGNAL_NUMBER" env-default:"5000439800"`
	}
}

func SetConfig() {

	if _, err := os.Stat(".env"); err == nil {
		cleanenv.ReadConfig(".env", &Cfg)
		logrus.Info("Set config from .env file")
	} else {
		cleanenv.ReadEnv(&Cfg)
		logrus.Info("Set config from Config struct values")
	}

}
