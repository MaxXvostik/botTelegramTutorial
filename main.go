package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type wallet map[string]float64

var db = map[int64]wallet{}

func main() {
	bot, err := tgbotapi.NewBotAPI("5955978668:AAGGKMqRJelEBzGSWSs8QSUWXtVmfY1bxyU")
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message

			log.Println(update.Message.Text)
			msArr := strings.Split(update.Message.Text, " ")

			switch msArr[0] {
			case "ADD":
				summ, err := strconv.ParseFloat(msArr[2], 64)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Невозможно сконвертировать сумму"))
					continue
				}

				if _, ok := db[update.Message.Chat.ID]; !ok {
					db[update.Message.Chat.ID] = wallet{}
				}

				db[update.Message.Chat.ID][msArr[1]] += summ

				msg := fmt.Sprintf("Баланс : %s %f", msArr[1], db[update.Message.Chat.ID][msArr[1]])
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))

			case "SUB":
				summ, err := strconv.ParseFloat(msArr[2], 64)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Невозможно сконвертировать сумму"))
					continue
				}

				if _, ok := db[update.Message.Chat.ID]; !ok {
					db[update.Message.Chat.ID] = wallet{}
				}

				db[update.Message.Chat.ID][msArr[1]] -= summ

				msg := fmt.Sprintf("Баланс : %s %f", msArr[1], db[update.Message.Chat.ID][msArr[1]])
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))

			case "DEL":
				delete(db[update.Message.Chat.ID], msArr[1])
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Валюта удалена"))
			case "SHOW":
				msg := "Баланс:\n"
				for key, val := range db[update.Message.Chat.ID] {
					msg += fmt.Sprintf("%s: %f\n", key, val)
				}

				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
			default:
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда"))
			}

		}
	}
}
