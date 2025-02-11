package handlers

import (
	"boutonsjob/externals/db"
	"boutonsjob/externals/redis"

	tele "gopkg.in/telebot.v4"
)

type Handler struct {
	StateName string
	Filter    func(name string, state string) bool
	Handler   func(c tele.Context, state *redis.StateData) error
}

type HandlerList struct {
	HandlersText     []Handler
	HandlersCallback []Handler
	stateManager     *redis.RedisClient
	db               *db.DB
}

func InitHandlers(stateManager *redis.RedisClient, db *db.DB) HandlerList {
	return HandlerList{stateManager: stateManager, db: db}
}

func (h *HandlerList) AddText(state_name string, handler func(c tele.Context, state *redis.StateData) error) {
	h.HandlersText = append(h.HandlersText, Handler{state_name, func(name string, state string) bool { return name == state }, handler})
}

func (h *HandlerList) AddCallback(state_name string, handler func(c tele.Context, state *redis.StateData) error) {
	h.HandlersCallback = append(h.HandlersCallback, Handler{state_name, func(name string, state string) bool { return name == state }, handler})
}

func (h *HandlerList) AddCallbackFilter(state_name string, filter func(name string, state string) bool, handler func(c tele.Context, state *redis.StateData) error) {
	h.HandlersCallback = append(h.HandlersCallback, Handler{state_name, filter, handler})
}

func (h *HandlerList) HandleStart(c tele.Context, state *redis.StateData) error {
	buttons := []tele.Btn{
		{Text: "Ищу исполнителя 🛠️", Data: "find_executor"},
		{Text: "Ищу заказы 💼", Data: "find_orders"},
		{Text: "Мои посты 📑", Data: "my_posts"},
	}

	menu := &tele.ReplyMarkup{}
	menu.Inline(
		menu.Row(buttons[0]),
		menu.Row(buttons[1]),
		menu.Row(buttons[2]),
	)
	messageId, err := h.stateManager.GetMessageId(c.Sender().ID)
	if err == nil {
		c.Bot().Edit(&tele.Message{Chat: &tele.Chat{ID: c.Sender().ID}, ID: messageId}, `Здравствуйте! 👋)
Скажите, что вы ищите, а я помогу вам составить объявление.`, menu)
	} else {

		msg, err := c.Bot().Send(tele.ChatID(c.Sender().ID), `Здравствуйте! 👋)
Скажите, что вы ищите, а я помогу вам составить объявление.`, menu)
		if err != nil {
			return err
		}
		h.stateManager.SaveMessageId(c.Sender().ID, msg.ID)
	}
	return err
}
