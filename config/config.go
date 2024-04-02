package config

type Config struct {
	TgBotToken            string
	MongoConnectionString string
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
		//TgBotToken:            *tgBotTokenToken,
		//MongoConnectionString: *mongoConnectionString,
		TgBotToken:            "6927936576:AAEMzBBo0Bs6T3Nnq8AN1DfCvl5jzFBKoks",
		MongoConnectionString: "mongodb+srv://forevermenty25:k9zPoWBGI3PQT4Gz@cluster0.p8mq4lx.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0",
	}
}
