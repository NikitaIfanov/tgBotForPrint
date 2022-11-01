package tg

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"math/rand"
	"strconv"
)

// Client ...
type Client struct {
	UserName string
	ChatId   int64
}

// DataForPrint ...
type DataForPrint struct {
	PhotoID    []string
	DocumentID []string
}

// Office ...
type Office struct {
	AgentName string
	Address   string
	ChatID    int64
}

// Order ...
type Order struct {
	Number   int           //+
	Client   *Client       //+
	Location *Office       //+
	Data     *DataForPrint //+
	Flag     bool
}

//NewOrder ...
func NewOrder() *Order {
	return &Order{
		Number:   rand.Int(),
		Client:   NewClient(),
		Location: &Office{},
		Data: &DataForPrint{
			PhotoID:    make([]string, 0, 5),
			DocumentID: make([]string, 0, 5),
		},
		Flag: false,
	}
}

func NewClient() *Client {
	return &Client{
		UserName: "",
		ChatId:   0,
	}
}

var (

	//infoButton
	infoButton = tgbotapi.NewInlineKeyboardButtonData("Инфо", "info")

	//keyBoardHello "Hello"
	keyBoardHello = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Новый заказ", "new"),
			infoButton,
		))

	//keyBoardLocationYesNo
	keyBoardLocationYesNo = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Да", "yes"),
			tgbotapi.NewInlineKeyboardButtonData("Нет", "no"),
		))

	key = tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Send"),
		))

	//
	/*keyBoardSendYesNo = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Да", "yes"),
		tgbotapi.NewInlineKeyboardButtonData("Нет, отправить заказ", "send"),
	))*/
)

//Bot functional
func Bot(config *Config) {
	userNameToOrders := make(map[string][]*Order)
	//ClientToChatID := make(map[string]int64)
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Printf("Start bot error is : %s", err)
	}
	//Debug
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(-0)
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

		//case for clients
		case update.Message != nil:
			MessageNotNil(update, config, bot, userNameToOrders)
		//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "123")
		//	msg.ReplyMarkup = key
		//	bot.Send(msg)
		case update.CallbackQuery != nil:
			CallbackNotNil(config, update, bot, userNameToOrders, key)
			/*case update.Message.MediaGroupID == "" && userNameToOrders[update.Message.Chat.UserName][len(userNameToOrders[update.Message.Chat.UserName])-1].Flag == true:
			SendMsgWithKeyboard("Еще файлы?", bot, update.Message.Chat.ID, keyBoardSendYesNo)
			userNameToOrders[update.Message.Chat.UserName][len(userNameToOrders[update.Message.Chat.UserName])-1].Flag = false
				case update.Message.Photo != nil:
					PhotoNotNil(update, userNameToOrders, bot)
				case update.Message.Document != nil:
					idx := len(userNameToOrders[update.Message.Chat.UserName]) - 1
					userNameToOrders[update.Message.Chat.UserName][idx].Data.DocumentID =
						append(
							userNameToOrders[update.Message.Chat.UserName][idx].Data.DocumentID,
							update.Message.Document.FileID,
						)*/

			//for customers

		}
		//	if userNameToOrders[update.Message.Chat.UserName][len(userNameToOrders[update.Message.Chat.UserName])-1].Flag == false {

		//}
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

/*func Check(update tgbotapi.Update, m map[string][]*Order) bool {
	u := update.Message
	switch {
	case u.Photo != nil:
		idx := len(m[update.Message.Chat.UserName]) - 1
		m[update.Message.Chat.UserName][idx].Data.PhotoID =
			append(
				m[update.Message.Chat.UserName][idx].Data.PhotoID,
				update.Message.Photo[len(update.Message.Photo)-1].FileID,
			)

		return true

	case u.Document != nil:
		idx := len(m[update.Message.Chat.UserName]) - 1
		m[update.Message.Chat.UserName][idx].Data.DocumentID =
			append(
				m[update.Message.Chat.UserName][idx].Data.DocumentID,
				update.Message.Document.FileID,
			)
		return true
	}
	return false
}*/

func SendOrderToOffice(order *Order, bot *tgbotapi.BotAPI) {
	FormOrder(order, bot)
	//msg := tgbotapi.NewCopyMessage(order.Location.ChatID, order.ChatID, order.MsgID)
	order.Flag = true
	sl := make([]interface{}, 0, 5)
	for i := 0; i < len(order.Data.PhotoID); i++ {

		sl = append(sl, tgbotapi.NewInputMediaPhoto(tgbotapi.FileID(order.Data.PhotoID[i])))

		//	msg := tgbotapi.NewPhoto(order.Location.ChatID, tgbotapi.FileID(order.Data.PhotoID[i]))
		//	if _, err := bot.Send(msg); err != nil {
		//		log.Print(err)
		//	}
	}
	msg := tgbotapi.NewMediaGroup(order.Location.ChatID, sl)
	if _, err := bot.SendMediaGroup(msg); err != nil {
		log.Print(err)
	}
	SendMsg(bot, order.Client.ChatId, "Заказ сформирован")
}

func FormOrder(order *Order, bot *tgbotapi.BotAPI) {
	str := fmt.Sprintf("Номер заказа: %d\nКлиент: @%s", order.Number, order.Client.UserName)
	msg := tgbotapi.NewMessage(order.Location.ChatID, str)
	if _, err := bot.Send(msg); err != nil {
		log.Print(err)
	}
}
