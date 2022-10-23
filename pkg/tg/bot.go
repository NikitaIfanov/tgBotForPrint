package tg

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"math/rand"
	"strconv"
)

var (

	//infoButton
	infoButton = tgbotapi.NewInlineKeyboardButtonData("Инфо", "info")

	//keyBoardHello "Hello"
	keyBoardHello = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Новый заказ", "new"),
			infoButton,
		))

	//keyBoard "Yes"/"No"
	keyBoardYesNo = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Да", "yes"),
			tgbotapi.NewInlineKeyboardButtonData("Нет", "no"),
		))
)

//Bot functional
func Bot(config *Config) {
	userNameToOrder := make(map[string]*Order)

	addresses := MakeMsgWithAddresses(config)
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Printf("Start bot error is : %s", err)
	}
	//Debug
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		switch {
		//case for managers
		case update.Message != nil && func(agent string) bool {
			for _, ag := range config.AgentsName {
				if ag == agent {
					return true
				}
			}
			return false
		}(update.Message.Chat.UserName):
			for i := 1; i <= len(config.Offices); i++ {
				if config.Offices[i].AgentName == update.Message.Chat.UserName &&
					config.Offices[i].ChatID == 0 {
					config.Offices[i].ChatID = update.Message.Chat.ID
					break
				}
			}
			//msg for the office manager
			SendMsg(bot, update.Message.Chat.ID, "Office Registered")
			log.Print(config)

		//case for customers
		case update.Message != nil:
			switch update.Message.Text {
			//Check customer nickname and make new order in map[username]*Order
			case "/start":

				_, ok := userNameToOrder[update.Message.Chat.UserName]
				if ok != true {

					userNameToOrder[update.Message.Chat.UserName] = NewOrder()
				}
				SendMsgWithKeyboard(
					"Приветсвуем в помощнике компании '***'",
					bot, update.Message.Chat.ID, keyBoardHello)
			// Select address
			case "1", "2", "3", "4":
				i, err := strconv.Atoi(update.Message.Text)
				if err != nil {
					log.Print(err)
					//send msg with incorrect input format
				}
				userNameToOrder[update.Message.Chat.UserName].Location = config.Offices[i]
				msg := fmt.Sprintf("Выбран офис по адресу: %s. Подтвердить?", config.Offices[i].Address)
				SendMsgWithKeyboard(msg, bot, update.Message.Chat.ID, keyBoardYesNo)

			}
			if Check(update, userNameToOrder) {

				SendMsg(bot, update.Message.Chat.ID, "Заказ сформирован")
				log.Print(userNameToOrder[update.Message.Chat.UserName].Data.Photo)
				SendOrderToOffice(userNameToOrder[update.Message.Chat.UserName], bot, "ha")

			}

		}
		//for customers
		if update.CallbackQuery != nil {
			switch update.CallbackQuery.Data {
			//Set ChatID into order in map
			case "new":
				order := userNameToOrder[update.CallbackQuery.Message.Chat.UserName]
				order.ChatID = update.CallbackQuery.Message.Chat.ID

				SendMsg(bot, update.CallbackQuery.Message.Chat.ID,
					"Выберите офис для получения вашего заказа\n"+
						" (в ответе пришлите цифру)\n\n"+addresses)

			case "yes":
				msg := "Отправьте файлы"
				SendMsg(bot, update.CallbackQuery.Message.Chat.ID, msg)
			case "no":
				SendMsg(bot, update.CallbackQuery.Message.Chat.ID,
					"Выберите офис для получения вашего заказа\n"+
						" (в ответе пришлите цифру)\n\n"+addresses)

			}
		}
	}
}

//SendMsg to s.o.
func SendMsg(bot *tgbotapi.BotAPI, ChatID int64, message string) {
	msg := tgbotapi.NewMessage(ChatID, message)
	if _, err := bot.Send(msg); err != nil {
		log.Print(err)
	}
}

//SendMsgWithKeyboard to s.o.
func SendMsgWithKeyboard(message string, bot *tgbotapi.BotAPI, ChatID int64, keyboard tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(ChatID, message)
	msg.ReplyMarkup = keyboard
	if _, err := bot.Send(msg); err != nil {
		log.Print(err)

	}
}

// Order ...
type Order struct {
	Number    int               //+
	ChatID    int64             //+
	Location  *Office           //+
	Data      *tgbotapi.Message //+
	MsgID     int
	PayOnline bool
}

//NewOrder ...
func NewOrder() *Order {
	return &Order{
		Number:    rand.Int(),
		ChatID:    0,
		Location:  &Office{},
		Data:      &tgbotapi.Message{},
		MsgID:     0,
		PayOnline: false,
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

func SendOrderToOffice(order *Order, bot *tgbotapi.BotAPI, msgStr string) {

	msg := tgbotapi.NewCopyMessage(order.Location.ChatID, order.ChatID, order.MsgID)

	if _, err := bot.Send(msg); err != nil {
		log.Print(err)

	}
}

func Check(update tgbotapi.Update, m map[string]*Order) bool {
	u := update.Message
	switch {
	case u.Photo != nil:
		m[update.Message.Chat.UserName].MsgID = update.Message.MessageID
		return true
	}
	return false
}
