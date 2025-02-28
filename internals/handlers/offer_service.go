package handlers

import (
	"fmt"
	"unicode/utf8"

	"boutonsjob/externals/redis"

	tele "gopkg.in/telebot.v4"
)

func (h *HandlerList) HandleOfferService(c tele.Context, state *redis.StateData) error {
	messageId, err := h.stateManager.GetMessageId(c.Sender().ID)
	if err != nil {
		return err
	}
	categories, err := h.db.GetOfferCategories()
	if err != nil {
		return err
	}
	var rows []tele.Row
	menu := &tele.ReplyMarkup{}
	for _, category := range categories {
		rows = append(rows, menu.Row(tele.Btn{Text: category.Name, Data: category.Data}))
	}
	rows = append(rows, menu.Row(tele.Btn{Text: "Другое ✏️", Data: "!oc_other"}))
	rows = append(rows, menu.Row(tele.Btn{Text: "Назад ↩️", Data: "back"}))
	menu.Inline(rows...)
	_, err = c.Bot().Edit(&tele.Message{Chat: &tele.Chat{ID: c.Sender().ID}, ID: messageId}, `Выберите категорию предлагаемой услуги:`, menu)
	return err
}

func (h *HandlerList) HandleOfferCategory(c tele.Context, state *redis.StateData) error {
	messageId, err := h.stateManager.GetMessageId(c.Sender().ID)
	if err != nil {
		return err
	}
	session, err := h.stateManager.GetSessionData(c.Sender().ID)
	if err != nil {
		return err
	}
	session["category"], err = h.db.GetOfferCategoryName(c.Callback().Data)
	if err != nil {
		return err
	}
	h.stateManager.SetSessionData(c.Sender().ID, session)
	fmt.Println(session)
	h.stateManager.SetState(c.Sender().ID, "offer_name", map[string]interface{}{"prev_state": state.State})
	_, err = c.Bot().Edit(&tele.Message{Chat: &tele.Chat{ID: c.Sender().ID}, ID: messageId}, `Введите ваше имя:`)
	return err
}


func (h *HandlerList) HandleOfferCategoryOther(c tele.Context, state *redis.StateData) error {
	messageId, err := h.stateManager.GetMessageId(c.Sender().ID)
	if err != nil {
		return err
	}
	h.stateManager.SetState(c.Sender().ID, "wait_offer_category", map[string]interface{}{"prev_state": state.State})
	_, err = c.Bot().Edit(&tele.Message{Chat: &tele.Chat{ID: c.Sender().ID}, ID: messageId}, `Введите название категории:`)
	return err
}

func (h *HandlerList) HandleOfferCategoryOtherName(c tele.Context, state *redis.StateData) error {
	messageId, err := h.stateManager.GetMessageId(c.Sender().ID)
	if err != nil {
		return err
	}
	session, err := h.stateManager.GetSessionData(c.Sender().ID)
	if err != nil {
		return err
	}
	session["category"] = c.Message().Text
	h.stateManager.SetSessionData(c.Sender().ID, session)
	fmt.Println(session)
	c.Delete()
	h.stateManager.SetState(c.Sender().ID, "offer_name", map[string]interface{}{"prev_state": state.State})
	_, err = c.Bot().Edit(&tele.Message{Chat: &tele.Chat{ID: c.Sender().ID}, ID: messageId}, `Введите ваше имя:`)
	return err
}

func (h *HandlerList) HandleOfferName(c tele.Context, state *redis.StateData) error {
	messageId, err := h.stateManager.GetMessageId(c.Sender().ID)
	if err != nil {
		return err
	}
	session, err := h.stateManager.GetSessionData(c.Sender().ID)
	if err != nil {
		return err
	}
	session["name"] = c.Message().Text
	c.Delete()
	h.stateManager.SetSessionData(c.Sender().ID, session)
	fmt.Println(session)
	h.stateManager.SetState(c.Sender().ID, "wait_offer_portfolio", map[string]interface{}{"prev_state": state.State})
	menu := &tele.ReplyMarkup{}
	menu.Inline(menu.Row(tele.Btn{Text: "Нет портфолио ❌📂", Data: "no_portfolio"}))
	_, err = c.Bot().Edit(&tele.Message{Chat: &tele.Chat{ID: c.Sender().ID}, ID: messageId}, `Введите ссылку на портфолио:`, menu)
	return err
}

func (h *HandlerList) HandleOfferPortfolio(c tele.Context, state *redis.StateData) error {
	messageId, err := h.stateManager.GetMessageId(c.Sender().ID)
	if err != nil {
		return err
	}
	session, err := h.stateManager.GetSessionData(c.Sender().ID)
	if err != nil {
		return err
	}
	if c.Callback() != nil {
		session["portfolio"] = "Нет портфолио"
	} else {
	session["portfolio"] = c.Message().Text
	c.Delete()
	}
	h.stateManager.SetSessionData(c.Sender().ID, session)
	fmt.Println(session)
	h.stateManager.SetState(c.Sender().ID, "wait_offer_minimal_order", map[string]interface{}{"prev_state": state.State})
menu := &tele.ReplyMarkup{}
	menu.Inline(menu.Row(tele.Btn{Text: "Нет требований ✅", Data: "no_minimal_order"}))
	_, err = c.Bot().Edit(&tele.Message{Chat: &tele.Chat{ID: c.Sender().ID}, ID: messageId}, `Напишите требования к минимальному заказу (до 100 символов):`, menu)
	return err
}

func (h *HandlerList) HandleOfferMinimalOrder(c tele.Context, state *redis.StateData) error {
  messageId, err := h.stateManager.GetMessageId(c.Sender().ID)
  if err != nil {
    return err
  }
  session, err := h.stateManager.GetSessionData(c.Sender().ID)
  if err != nil {
    return err
  }
	if c.Callback() != nil {
		session["minimal_order"] = "no_minimal_order"
	} else {
		if utf8.RuneCountInString(c.Message().Text) > 100 {
			_, err = c.Bot().Edit(&tele.Message{Chat: &tele.Chat{ID: c.Sender().ID}, ID: messageId}, `Слишком много символов! Введите требования к минимальному заказу в объёме до 100 символов:`)
      return err
		} else {
      session["minimal_order"] = c.Message().Text
		}
	}
  h.stateManager.SetSessionData(c.Sender().ID, session)
  fmt.Println(session)
  h.stateManager.SetState(c.Sender().ID, "wait_offer_cv", map[string]interface{}{"prev_state": state.State})
  _, err = c.Bot().Edit(&tele.Message{Chat: &tele.Chat{ID: c.Sender().ID}, ID: messageId}, `Напишите краткое резюме (до 1000 символов):`)
  return err
}

func (h *HandlerList) HandleOfferCV(c tele.Context, state *redis.StateData) error {
	messageId, err := h.stateManager.GetMessageId(c.Sender().ID)
	if err != nil {
		return err
	}
	session, err := h.stateManager.GetSessionData(c.Sender().ID)
	if err != nil {
		return err
	}
	if utf8.RuneCountInString(c.Message().Text) > 1000 {
		_, err = c.Bot().Edit(&tele.Message{Chat: &tele.Chat{ID: c.Sender().ID}, ID: messageId}, `Слишком много символов! Введите резюме в объёме до 1000 символов:`)
		return err
	} else {
		session["cv"] = c.Message().Text
	}
	c.Delete()
	h.stateManager.SetSessionData(c.Sender().ID, session)
	fmt.Println(session)
	h.stateManager.SetState(c.Sender().ID, "wait_offer_stack", map[string]interface{}{"prev_state": state.State})
	_, err = c.Bot().Edit(&tele.Message{Chat: &tele.Chat{ID: c.Sender().ID}, ID: messageId}, `Введите технологический стек, который вы используете при работе (до 500 символов):`)
	return err
}

func (h *HandlerList) HandleOfferStack(c tele.Context, state *redis.StateData) error {
	messageId, err := h.stateManager.GetMessageId(c.Sender().ID)
	if err != nil {
		return err
	}
	session, err := h.stateManager.GetSessionData(c.Sender().ID)
	if err != nil {
		return err
	}
	if utf8.RuneCountInString(c.Message().Text) > 500 {
		_, err = c.Bot().Edit(&tele.Message{Chat: &tele.Chat{ID: c.Sender().ID}, ID: messageId}, `Слишком много символов! Введите описание стека в объёме до 500 символов:`)
		return err
	} else {
		session["stack"] = c.Message().Text
	}
	c.Delete()
	h.stateManager.SetSessionData(c.Sender().ID, session)
	fmt.Println(session)
	h.stateManager.SetState(c.Sender().ID, "wait_offer_payment", map[string]interface{}{"prev_state": state.State})
	_, err = c.Bot().Edit(&tele.Message{Chat: &tele.Chat{ID: c.Sender().ID}, ID: messageId}, `Введите подтверждение оплаты:`)
	return err
}

func (h *HandlerList) HandleOfferPayment(c tele.Context, state *redis.StateData) error {
	messageId, err := h.stateManager.GetMessageId(c.Sender().ID)
	if err != nil {
		return err
	}

	session, err := h.stateManager.GetSessionData(c.Sender().ID)
	if err != nil {
		return err
	}
	c.Delete()
	session["payment"] = c.Message().Text
	c.Delete()
	h.stateManager.SetSessionData(c.Sender().ID, session)
	fmt.Println(session)
	h.stateManager.SetState(c.Sender().ID, "offer_contacts", map[string]interface{}{"prev_state": state.State})
	post, err := c.Bot().Send(&tele.Chat{ID: -1002284442496}, fmt.Sprintf(`
#ИСПОЛНИТЕЛЬ
👨‍💻 👨‍💻 👨‍💻
<b>Имя:</b> %s
<b>Контакт:</b> %s
<b>Портфолио:</b> %s
<b>Технологический стек:</b> %s
<b>Требования к минимальному заказу:</b> %s
<b>Резюме:</b>
%s

#%s`, session["name"], c.Sender().Username, session["portfolio"], session["stack"], session["minimal_order"], session["cv"], session["category"]), &tele.SendOptions{ParseMode: tele.ModeHTML})
	if err != nil {
		return err
	}
	h.stateManager.SetState(c.Sender().ID, "start", map[string]interface{}{})
	menu := &tele.ReplyMarkup{}
	menu.Inline(
    menu.Row(tele.Btn{Text: "В главное меню", Data: "start"}),
	)
	_, err = c.Bot().Edit(&tele.Message{Chat: &tele.Chat{ID: c.Sender().ID}, ID: messageId}, fmt.Sprintf(`Ваше резюме размещено: <a href='https://t.me/zylogbot/%d'>Посмотреть пост</a>`, post.ID), &tele.SendOptions{ParseMode: tele.ModeHTML}, menu)
	return err
}
