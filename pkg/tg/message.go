package tg

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
)

func MessageNotNil(update tgbotapi.Update, config *Config, bot *tgbotapi.BotAPI, m map[string][]*Order) {

	switch {

	case update.Message.Text == "/start":

		_, ok := m[update.Message.Chat.UserName]
		if ok == false {
			m[update.Message.Chat.UserName] = make([]*Order, 0, 5)
		}
		SendMsgWithKeyboard(
			"Приветсвуем в помощнике компании '***'",
			bot, update.Message.Chat.ID, keyBoardHello)

	case func(msg string) bool {
		if i, err := strconv.Atoi(msg); err == nil {
			for idx := range config.AgentsName {
				if i == idx {
					return true
				}
			}
			return false
		} else {
			log.Print(err)
			return false
		}
	}(update.Message.Text) && update.Message.Text != "":
		if _, err := strconv.Atoi(update.Message.Text); err == nil {
			i, _ := strconv.Atoi(update.Message.Text)
			if _, ok := config.Offices[i]; ok == true {
				idx := len(m[update.Message.Chat.UserName]) - 1
				m[update.Message.Chat.UserName][idx].Location =
					config.Offices[i]

				msg := fmt.Sprintf(
					"Выбран офис по адресу: %s. Подтвердить?",
					config.Offices[i].Address,
				)
				SendMsgWithKeyboard(msg, bot, update.Message.Chat.ID, keyBoardLocationYesNo)

			} else {
				log.Print(err)

				//send msg with incorrect input format
			}
		}
	case update.Message.Text == "Send":
		//	m[update.Message.Chat.UserName][len(m[update.Message.Chat.UserName])-1].Flag = false
		SendOrderToOffice(
			m[update.Message.Chat.UserName][len(m[update.Message.Chat.UserName])-1],
			bot,
		)

	case update.Message.Photo != nil:

		idx := len(m[update.Message.Chat.UserName]) - 1
		m[update.Message.Chat.UserName][idx].Data.PhotoID =
			append(
				m[update.Message.Chat.UserName][idx].Data.PhotoID,
				update.Message.Photo[len(update.Message.Photo)-1].FileID,
			)
		//m[update.Message.Chat.UserName][idx].Flag = true

	}
}
