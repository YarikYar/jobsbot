package handlers

import (
	"boutonsjob/externals/redis"
	"fmt"

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

func (h* HandlerList) HandleOfferCategory(c tele.Context, state *redis.StateData) error {
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

func (h* HandlerList) HandleOfferCategoryOther(c tele.Context, state *redis.StateData) error {
  messageId, err := h.stateManager.GetMessageId(c.Sender().ID)
  if err != nil {
    return err
  }
  h.stateManager.SetState(c.Sender().ID, "wait_offer_category", map[string]interface{}{"prev_state": state.State})
  _, err = c.Bot().Edit(&tele.Message{Chat: &tele.Chat{ID: c.Sender().ID}, ID: messageId}, `Введите название категории:`)
  return err
}

func (h* HandlerList) HandleOfferCategoryOtherName(c tele.Context, state *redis.StateData) error {
	messageId, err := h.stateManager.GetMessageId(c.Sender().ID)
  if err != nil {
    return err
  }
	session, err := h.stateManager.GetSessionData(c.Sender().ID)
  if err != nil {
    return err
  }
  session["category"]  = c.Message().Text
  h.stateManager.SetSessionData(c.Sender().ID, session)
  fmt.Println(session)
	c.Delete()
  h.stateManager.SetState(c.Sender().ID, "offer_name", map[string]interface{}{"prev_state": state.State})
  _, err = c.Bot().Edit(&tele.Message{Chat: &tele.Chat{ID: c.Sender().ID}, ID: messageId}, `Введите ваше имя:`)
  return err
}

func (h* HandlerList) HandleOfferName(c tele.Context, state *redis.StateData) error {
  messageId, err := h.stateManager.GetMessageId(c.Sender().ID)
  if err != nil {
    return err
  }
	session, err := h.stateManager.GetSessionData(c.Sender().ID)
  if err != nil {
    return err
  }
  session["name"] = c.Message().Text
  h.stateManager.SetSessionData(c.Sender().ID, session)
  fmt.Println(session)
  h.stateManager.SetState(c.Sender().ID, "wait_offer_portfolio", map[string]interface{}{"prev_state": state.State})
	menu := &tele.ReplyMarkup{}
	menu.Inline(menu.Row(tele.Btn{Text: "Нет портфолио ❌📂", Data: "no_portfolio"}))
  _, err = c.Bot().Edit(&tele.Message{Chat: &tele.Chat{ID: c.Sender().ID}, ID: messageId}, `Введите ссылку на портфолио:`)
  return err
}





