package main

import (
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/qrlzvrn/Clozapinum/config"
	"github.com/qrlzvrn/Clozapinum/erro"
	"github.com/qrlzvrn/Clozapinum/handlers"
	"github.com/qrlzvrn/Clozapinum/logger"
)

func main() {
	botConfig, err := config.NewTgBotConf()
	if err != nil {
		logger.FatalMe(err)
	}

	sslConfig, err := config.NewSSLConf()
	if err != nil {
		logger.FatalMe(err)
	}

	if bot, err := tgbotapi.NewBotAPI(botConfig.APIToken); err != nil {
		e := erro.NewWrapError("NewBotAPI", err)
		logger.FatalMe(e)
	} else {

		bot.Debug = true

		logger.Infof("Authorized on account %s", bot.Self.UserName)

		info, err := bot.GetWebhookInfo()
		if err != nil {
			e := erro.NewWrapError("GetWebhookInfo", err)
			logger.FatalMe(e)
		}
		if info.LastErrorDate != 0 {
			logger.Infof("Telegram callback failed: %s", info.LastErrorMessage)
		}

		updates := bot.ListenForWebhook("/" + bot.Token)
		go http.ListenAndServeTLS(":8443", sslConfig.Fullchain, sslConfig.Privkey, nil)

		for update := range updates {
			if update.Message != nil {
				msg, newKeyboard, newText, err := handlers.MessageHandler(update.Message)
				if err != nil {
					logger.FatalMe(err)
				}
				if msg != nil {
					if _, err := bot.Send(msg); err != nil {
						logger.BotSendFatal(err, "msg")
					}
				}
				if newKeyboard != nil {
					if _, err := bot.Send(newKeyboard); err != nil {
						logger.BotSendFatal(err, "newKeyboard")
					}
				}
				if newText != nil {
					if _, err := bot.Send(newText); err != nil {
						logger.BotSendFatal(err, "newText")
					}
				}
			} else if update.CallbackQuery != nil {
				msg, newKeyboard, newText, err := handlers.CallbackHandler(update.CallbackQuery)
				if err != nil {
					logger.FatalMe(err)
				}

				if msg != nil {
					if _, err := bot.Send(msg); err != nil {
						logger.BotSendFatal(err, "msg")
					}
				}
				if newText != nil {
					if _, err := bot.Send(newText); err != nil {
						logger.BotSendFatal(err, "newKeyboard")
					}
				}
				if newKeyboard != nil {
					if _, err := bot.Send(newKeyboard); err != nil {
						logger.BotSendFatal(err, "newText")
					}
				}
			}

		}
	}
}
