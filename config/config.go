package config

import "github.com/kelseyhightower/envconfig"

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
func NewDBConf() (DB, error) {
	db := DB{}

	err := envconfig.Process("db", &db)
	if err != nil {
		return db, err
	}

	return db, nil
}

// NewTgBotConf - генерирует новый конфиг с информацией о телеграм боте
func NewTgBotConf() (TgBot, error) {
	tgBot := TgBot{}

	err := envconfig.Process("telegram", &tgBot)
	if err != nil {
		return tgBot, err
	}

	return tgBot, nil
}

// NewSSLConf - генерирует новый конфиг с информацией о SSL
func NewSSLConf() (SSL, error) {
	ssl := SSL{}

	err := envconfig.Process("ssl", &ssl)
	if err != nil {
		return ssl, err
	}

	return ssl, nil
}
