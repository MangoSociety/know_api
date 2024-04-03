package config

type Config struct {
	TgBotToken            string
	MongoConnectionString string
}

type GithubData struct {
	Owner string
}

const (
	gh_owner      = "MangoSociety"
	gh_repo       = "golang_road"
	gh_sha        = "1fbd0e9ed5ee1239996720d79435b7a2feb5d507"
	gh_token      = "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"
	tg_token      = "6927936576:AAEMzBBo0Bs6T3Nnq8AN1DfCvl5jzFBKoks"
	mongo_connect = "mongodb+srv://forevermenty25:k9zPoWBGI3PQT4Gz@cluster0.p8mq4lx.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
)

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
		TgBotToken:            tg_token,
		MongoConnectionString: mongo_connect,
	}
}
