package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
)

func CallbackNotNil(config *Config, update tgbotapi.Update, bot *tgbotapi.BotAPI, m map[string][]*Order, keyBoard tgbotapi.ReplyKeyboardMarkup) {
	addresses := MakeMsgWithAddresses(config)
	//Set ChatID into order in map
	if update.CallbackQuery.Data == "new" {
		m[update.CallbackQuery.Message.Chat.UserName] = append(
			m[update.CallbackQuery.Message.Chat.UserName],
			NewOrder(),
		)
		idx := len(m[update.CallbackQuery.Message.Chat.UserName]) - 1
		m[update.CallbackQuery.Message.Chat.UserName][idx].Client.ChatId =
			update.CallbackQuery.Message.Chat.ID
		m[update.CallbackQuery.Message.Chat.UserName][idx].Client.UserName =
			update.CallbackQuery.Message.Chat.UserName

		SendMsg(bot, update.CallbackQuery.Message.Chat.ID,
			"Выберите офис для получения вашего заказа\n"+
				" (в ответе пришлите цифру)\n\n"+addresses)
	} else if update.CallbackQuery.Data == "yes" {
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Отправьте файлы")
		msg.ReplyMarkup = keyBoard
		if _, err := bot.Send(msg); err != nil {
			log.Print(err)
		}

	} else if update.CallbackQuery.Data == "no" {
		SendMsg(bot, update.CallbackQuery.Message.Chat.ID,
			"Выберите офис для получения вашего заказа\n"+
				" (в ответе пришлите цифру)\n\n"+addresses)
	} else if update.CallbackQuery.Data == "info" {
		msg := "info"
		SendMsg(bot, update.CallbackQuery.Message.Chat.ID, msg)
	}
}

func MakeMsgWithAddresses(config *Config) string {
	str := ""
	for idx, s := range config.Addresses {
		if idx <= len(config.Addresses)-1 {
			str += strconv.Itoa(idx+1) + ". " + s + "\n\n"
		} else {
			str += strconv.Itoa(idx+1) + ". " + s
		}
	}
	return str
}
