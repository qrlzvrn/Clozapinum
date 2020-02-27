package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/qrlzvrn/Clozapinum/erro"
)

// Используя библиотеку envconfig считываем в структуры необходимые данные из переменных окружения.
// В случае работы с докером переменные окружения прописываются в .env файле и передаются в докер.
// В случе, когда мы не хотим работать с докером, придется прописать все переменные вручную.
// Позже, может быть автоматизирую данный процесс.

// DB - конфиг для работы с базой данных
type DB struct {
	Host     string
	Port     int
	Username string
	Password string
	Name     string
}

// TgBot - конфиг для работы с телеграм ботом, пока что нужен только для хранения токена, в будущем может быть расширен
type TgBot struct {
	APIToken string
}

// SSL - конфиг хранящий пути к сертификатам
type SSL struct {
	Fullchain string
	Privkey   string
}

// NewDBConf - генерирует новый конфиг для работы с базой данных
func NewDBConf() (DB, erro.Err) {
	db := DB{}

	err := envconfig.Process("db", &db)
	if err != nil {
		e := erro.NewConfigError("NewDBConf", err)
		return db, e
	}

	return db, nil
}

// NewTgBotConf - генерирует новый конфиг с информацией о телеграм боте
func NewTgBotConf() (TgBot, erro.Err) {
	tgBot := TgBot{}

	err := envconfig.Process("telegram", &tgBot)
	if err != nil {
		e := erro.NewConfigError("NewTgBotConf", err)
		return tgBot, e
	}

	return tgBot, nil
}

// NewSSLConf - генерирует новый конфиг с информацией о SSL
func NewSSLConf() (SSL, erro.Err) {
	ssl := SSL{}

	err := envconfig.Process("ssl", &ssl)
	if err != nil {

		e := erro.NewConfigError("NewSSLConf", err)
		return ssl, e
	}

	return ssl, nil
}
