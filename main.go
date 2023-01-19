package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type binanceResp struct {
	Price float64 `json:"price,string"`
	Code  int64   `json:"code"`
}

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
				var usdSum float64
				for key, val := range db[update.Message.Chat.ID] {
					coinPrice, err := getPrice(key)

					if err != nil {
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
					}
					usdSum += val * coinPrice
					msg += fmt.Sprintf("%s: %f [%.2f]\n", key, val, val*coinPrice)
				}
				msg += fmt.Sprintf("Сумма:%.2f\n", usdSum)

				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
			default:
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда"))
			}

		}
	}
}
func getPrice(coin string) (price float64, err error) {
	resp, err := http.Get(fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%sUSDT", coin))
	if err != nil {
		return
	}

	rub, err := http.Get(fmt.Sprint("https://api.binance.com/api/v3/ticker/price?symbol=USDTRUB"))
	if err != nil {
		return
	}

	defer resp.Body.Close()

	var jsonResp binanceResp
	var jsonRub binanceResp

	err = json.NewDecoder(resp.Body).Decode(&jsonResp)
	if err != nil {
		return
	}

	if jsonResp.Code != 0 {
		err = errors.New("Не корректная валюта")
		return
	}

	err = json.NewDecoder(rub.Body).Decode(&jsonRub)
	if err != nil {
		log.Fatal("rubli nou")
		return
	}

	price = jsonResp.Price //* jsonRub.Price
	return
}
