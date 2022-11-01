package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"math/rand"
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
	Number   int
	Client   *Client
	Location *Office
	Data     *DataForPrint
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
)

//Bot functional
func Bot(config *Config) {
	userNameToOrders := make(map[string][]*Order)
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Printf("Start bot error is : %s", err)
	}

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
		case update.Message != nil:
			MessageNotNil(update, config, bot, userNameToOrders)
		case update.CallbackQuery != nil:
			CallbackNotNil(config, update, bot, userNameToOrders, key)
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
