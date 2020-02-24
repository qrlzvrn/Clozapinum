package main

import (
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/qrlzvrn/Clozapinum/config"
	"github.com/qrlzvrn/Clozapinum/handlers"
)

func main() {
	botConfig, err := config.NewTgBotConf()
	if err != nil {
		log.Panic(err)
	}

	sslConfig, err := config.NewSSLConf()
	if err != nil {
		log.Panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(botConfig.APIToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServeTLS(":8443", sslConfig.Fullchain, sslConfig.Privkey, nil)

	for update := range updates {
		if update.Message != nil {
			msg, newKeyboard, newText := handlers.MessageHandler(update.Message)

			if msg != nil {
				if _, err := bot.Send(msg); err != nil {
					log.Fatal(err)
				}
			}
			if newKeyboard != nil {
				if _, err := bot.Send(newKeyboard); err != nil {
					log.Fatal(err)
				}
			}
			if newText != nil {
				if _, err := bot.Send(newText); err != nil {
					log.Fatal(err)
				}
			}
		} else if update.CallbackQuery != nil {
			msg, newKeyboard, newText := handlers.CallbackHandler(update.CallbackQuery)

			if msg != nil {
				if _, err := bot.Send(msg); err != nil {
					log.Fatal(err)
				}
			}
			if newText != nil {
				if _, err := bot.Send(newText); err != nil {
					log.Fatal(err)
				}
			}
			if newKeyboard != nil {
				if _, err := bot.Send(newKeyboard); err != nil {
					log.Fatal(err)
				}
			}
		}

	}
}
