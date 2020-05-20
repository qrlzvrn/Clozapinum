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
		log.Fatalf("%+v", err)
	}

	sslConfig, err := config.NewSSLConf()
	if err != nil {
		log.Fatalf("%+v", err)
	}

	if bot, err := tgbotapi.NewBotAPI(botConfig.APIToken); err != nil {
		log.Fatalf("%+v", err)
	} else {

		bot.Debug = true

		log.Printf("Authorized on account %s", bot.Self.UserName)

		info, err := bot.GetWebhookInfo()
		if err != nil {
			log.Fatalf("%+v", err)
		}
		if info.LastErrorDate != 0 {
			log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
		}

		updates := bot.ListenForWebhook("/" + bot.Token)
		go http.ListenAndServeTLS(":8443", sslConfig.Fullchain, sslConfig.Privkey, nil)

		for update := range updates {
			if update.Message != nil {
				msg, newKeyboard, newText, err := handlers.MessageHandler(update.Message)
				if err != nil {
					log.Fatalf("%+v", err)
				}
				if msg != nil {
					if _, err := bot.Send(msg); err != nil {
						log.Fatalf("%+v", err)
					}
				}
				if newKeyboard != nil {
					if _, err := bot.Send(newKeyboard); err != nil {
						log.Fatalf("%+v", err)
					}
				}
				if newText != nil {
					if _, err := bot.Send(newText); err != nil {
						log.Fatalf("%+v", err)
					}
				}
			} else if update.CallbackQuery != nil {
				msg, newKeyboard, newText, err := handlers.CallbackHandler(update.CallbackQuery)
				if err != nil {
					log.Fatalf("%+v", err)
				}

				if msg != nil {
					if _, err := bot.Send(msg); err != nil {
						log.Fatalf("%+v", err)
					}
				}
				if newText != nil {
					if _, err := bot.Send(newText); err != nil {
						log.Fatalf("%+v", err)
					}
				}
				if newKeyboard != nil {
					if _, err := bot.Send(newKeyboard); err != nil {
						log.Fatalf("%+v", err)
					}
				}
			}

		}
	}
}
