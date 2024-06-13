package config

import (
	"errors"
	"github.com/spf13/viper"
	"log"
	"time"
)

// App config struct
//type Config struct {
//	Server  ServerConfig
//	Redis   RedisConfig
//	MongoDB MongoDB
//	Cookie  Cookie
//	Session Session
//	Metrics Metrics
//	Github  GithubConfig
//}

type Config struct {
	Github       GithubConfig
	TgToken      string
	MongoConnect string
}

// Server config struct
type ServerConfig struct {
	AppVersion        string
	Port              string
	PprofPort         string
	Mode              string
	JwtSecretKey      string
	CookieName        string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	SSL               bool
	CtxDefaultTimeout time.Duration
	CSRF              bool
	Debug             bool
}

// Redis config
type RedisConfig struct {
	RedisAddr      string
	RedisPassword  string
	RedisDB        string
	RedisDefaultdb string
	MinIdleConns   int
	PoolSize       int
	PoolTimeout    int
	Password       string
	DB             int
}

type GithubConfig struct {
	Owner string
	Repo  string
	Sha   string
	Token string
}

// MongoDB config
type MongoDB struct {
	MongoURI string
}

// Cookie config
type Cookie struct {
	Name     string
	MaxAge   int
	Secure   bool
	HTTPOnly bool
}

// Session config
type Session struct {
	Prefix string
	Name   string
	Expire int
}

// Metrics config
type Metrics struct {
	URL         string
	ServiceName string
}

const (
	gh_owner      = "MangoSociety"
	gh_repo       = "golang_road"
	gh_sha        = "f4fbbd247b70039de1b6c6a43f419a22e88e8856"
	gh_token      = "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"
	tg_token      = "6927936576:AAEMzBBo0Bs6T3Nnq8AN1DfCvl5jzFBKoks"
	mongo_connect = "mongodb://localhost:27017"
)

// Load config file from given path
func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

// Parse config file
func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}

// TODO("вынести ключи к сервисам в флаги/настройки докера")
func MustLoad() Config {
	//tgBotTokenToken := flag.String(
	//	"tg-bot-token",
	//	"",
	//	"token for access to telegram bot",
	//)
	//mongoConnectionString := flag.String(
	//	"mongo-connection-string",
	//	"",
	//	"connection string for MongoDB",
	//)
	//
	//flag.Parse()
	//
	//if *tgBotTokenToken == "" {
	//	log.Fatal("token is not specified")
	//}
	//if *mongoConnectionString == "" {
	//	log.Fatal("mongo connection string is not specified")
	//}

	return Config{
		Github: GithubConfig{
			Owner: gh_owner,
			Repo:  gh_repo,
			Sha:   gh_sha,
			Token: gh_token,
		},
		TgToken:      tg_token,
		MongoConnect: mongo_connect,
	}
}
