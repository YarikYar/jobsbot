package handlers

import (
	"boutonsjob/externals/db"
	"boutonsjob/externals/redis"
	"fmt"

	tele "gopkg.in/telebot.v4"
)

type Handler struct {
	StateName string
	Filter func(name string, state string) bool
	Handler func(c tele.Context, state *redis.StateData) error
}

type HandlerList struct {
  HandlersText []Handler
	HandlersCallback []Handler
	stateManager *redis.RedisClient
	db *db.DB
}

func InitHandlers(stateManager *redis.RedisClient, db *db.DB) HandlerList {
  return HandlerList{stateManager: stateManager, db: db}
}

func (h *HandlerList) AddText(state_name string, handler func(c tele.Context, state *redis.StateData) error) {
  h.HandlersText = append(h.HandlersText, Handler{state_name, func(name string, state string) bool {return name == state}, handler})
}

func (h *HandlerList) AddCallback(state_name string, handler func(c tele.Context, state *redis.StateData) error) {
  h.HandlersCallback = append(h.HandlersCallback, Handler{state_name, func(name string, state string) bool {return name == state}, handler})
}

func (h *HandlerList) AddCallbackFilter(state_name string, filter func (name string, state string) bool, handler func(c tele.Context, state *redis.StateData) error) {
	h.HandlersCallback = append(h.HandlersCallback, Handler{state_name, filter, handler})
}

func (h *HandlerList) HandleStart(c tele.Context, state *redis.StateData) error{ 
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
	}else{

	msg, err := c.Bot().Send(tele.ChatID(c.Sender().ID), `Здравствуйте! 👋)
Скажите, что вы ищите, а я помогу вам составить объявление.`, menu)
	if err != nil {
    return err
  }
	h.stateManager.SaveMessageId(c.Sender().ID, msg.ID)
	}
    return err
}


func (h *HandlerList) HandleFindExecutor(c tele.Context, state *redis.StateData) error {
	messageId, err := h.stateManager.GetMessageId(c.Sender().ID)
  if err != nil {
    return err
  }
	categories, err := h.db.GetExecutorCategories()
  if err != nil {
    return err
	}
	var buttons []tele.Btn
  for _, category := range categories {
    buttons = append(buttons, tele.Btn{Text: category.Name, Data: category.Data})
  }
  menu := &tele.ReplyMarkup{}
  menu.Inline(
    menu.Row(buttons...),
    menu.Row(tele.Btn{Text: "Другое ✏️", Data: "!ec_other"}),
    menu.Row(tele.Btn{Text: "Назад ↩️", Data: "back"}),
  )
	h.stateManager.SetSessionData(c.Sender().ID, map[string]interface{}{"type": "select_executor"})
	c.Bot().Edit(&tele.Message{Chat: &tele.Chat{ID: c.Sender().ID}, ID: messageId},`Выберите категорию услуги, которые вам необходимы:`, menu) 
  return nil 
}


func (h *HandlerList) HandleExecutorCategory(c tele.Context, state *redis.StateData) error {
	messageId, err := h.stateManager.GetMessageId(c.Sender().ID)
  if err != nil {
    return err
  }
	session, err := h.stateManager.GetSessionData(c.Sender().ID)
  if err != nil {
    return err
  }
	session["category"] = c.Callback().Data
  h.stateManager.SetSessionData(c.Sender().ID, session)
	fmt.Println(session)

	menu := &tele.ReplyMarkup{}
  menu.Inline(
  	menu.Row(tele.Btn{Text: "Разовый заказ 🎯", Data: "et_onetime"}),
    menu.Row(tele.Btn{Text: "Несколько заказов (единоразово) 📦📦", Data: "et_many_one"}),
    menu.Row(tele.Btn{Text: "Удаленная вакансия (в штат) 💼👨‍💻", Data: "et_remote"}),
    menu.Row(tele.Btn{Text: "Сотрудничество (1 заказ в месяц) 🤝📅", Data: "et_collaboration"}),
    menu.Row(tele.Btn{Text: "Сотрудник в проект 👥🚀", Data: "et_project"}),
    menu.Row(tele.Btn{Text: "Назад ↩️", Data: "back"}),
	)
	_, err = c.Bot().Edit(&tele.Message{Chat: &tele.Chat{ID: c.Sender().ID}, ID: messageId},`Выберите, для каких целей вам нужен исполнитель:`, menu)
	return err
}

func (h *HandlerList) HandleExecutorCategoryOther(c tele.Context, state *redis.StateData) error {
  messageId, err := h.stateManager.GetMessageId(c.Sender().ID)
  if err != nil {
    return err
  }
  h.stateManager.SetState(c.Sender().ID, "wait_executor_category", map[string]interface{}{"prev_state": state.State})
	_, err = c.Bot().Edit(&tele.Message{Chat: &tele.Chat{ID: c.Sender().ID}, ID: messageId},`Введите название категории:`)
  return err
}

func (h *HandlerList) HandleExecutorCategoryOtherName(c tele.Context, state *redis.StateData) error {
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
	menu := &tele.ReplyMarkup{}
  menu.Inline(
  	menu.Row(tele.Btn{Text: "Разовый заказ 🎯", Data: "et_onetime"}),
    menu.Row(tele.Btn{Text: "Несколько заказов (единоразово) 📦📦", Data: "et_many_one"}),
    menu.Row(tele.Btn{Text: "Удаленная вакансия (в штат) 💼👨‍💻", Data: "et_remote"}),
    menu.Row(tele.Btn{Text: "Сотрудничество (1 заказ в месяц) 🤝📅", Data: "et_collaboration"}),
    menu.Row(tele.Btn{Text: "Сотрудник в проект 👥🚀", Data: "et_project"}),
    menu.Row(tele.Btn{Text: "Назад ↩️", Data: "back"}),
	)
	_, err = c.Bot().Edit(&tele.Message{Chat: &tele.Chat{ID: c.Sender().ID}, ID: messageId},`Выберите, для каких целей вам нужен исполнитель:`, menu)
  return err
}

func (h *HandlerList) HandleExecutorTarget(c tele.Context, state *redis.StateData) error {
  messageId, err := h.stateManager.GetMessageId(c.Sender().ID)
  if err != nil {
    return err
  }
  session, err := h.stateManager.GetSessionData(c.Sender().ID)
  if err != nil {
    return err
  }
  session["target"] = c.Callback().Data
  h.stateManager.SetSessionData(c.Sender().ID, session)
  fmt.Println(session)
	h.stateManager.SetState(c.Sender().ID, "executor_title", map[string]interface{}{"prev_state": state.State})
  _, err = c.Bot().Edit(&tele.Message{Chat: &tele.Chat{ID: c.Sender().ID}, ID: messageId}, `Напишите заголовок для объявления (до 100 символов):`)
	return err
}


func (h *HandlerList) HandleExecutorTitle(c tele.Context, state *redis.StateData) error {
  messageId, err := h.stateManager.GetMessageId(c.Sender().ID)
  if err != nil {
    return err
  }
  session, err := h.stateManager.GetSessionData(c.Sender().ID)
  if err != nil {
    return err
  }
	c.Delete()
  session["title"] = c.Message().Text
  h.stateManager.SetSessionData(c.Sender().ID, session)
  fmt.Println(session)
	h.stateManager.SetState(c.Sender().ID, "executor_budget", map[string]interface{}{"prev_state": state.State})
  _, err = c.Bot().Edit(&tele.Message{Chat: &tele.Chat{ID: c.Sender().ID}, ID: messageId}, `Напишите предполагаемый бюджет/зарплату:`)
  return err
}

func (h *HandlerList) HandleExecutorBudget(c tele.Context, state *redis.StateData) error {
  messageId, err := h.stateManager.GetMessageId(c.Sender().ID)
  if err != nil {
    return err
  }
  session, err := h.stateManager.GetSessionData(c.Sender().ID)
  if err != nil {
    return err
  }
	c.Delete()
  session["budget"] = c.Message().Text
  h.stateManager.SetSessionData(c.Sender().ID, session)
  fmt.Println(session)
  h.stateManager.SetState(c.Sender().ID, "executor_description", map[string]interface{}{"prev_state": state.State})
  _, err = c.Bot().Edit(&tele.Message{Chat: &tele.Chat{ID: c.Sender().ID}, ID: messageId}, `Напишите описание заказа (до 1000 символов):`)
  return err
}

func (h *HandlerList) HandleExecutorDescription(c tele.Context, state *redis.StateData) error {
  messageId, err := h.stateManager.GetMessageId(c.Sender().ID)
  if err != nil {
    return err
  }
  session, err := h.stateManager.GetSessionData(c.Sender().ID)
  if err != nil {
    return err
  }
  c.Delete()
  session["description"] = c.Message().Text
  h.stateManager.SetSessionData(c.Sender().ID, session)
  fmt.Println(session)
  h.stateManager.SetState(c.Sender().ID, "executor_skills", map[string]interface{}{"prev_state": state.State})
  _, err = c.Bot().Edit(&tele.Message{Chat: &tele.Chat{ID: c.Sender().ID}, ID: messageId}, `Напишите ключевые навыки (до 1000 символов):`)
  return err
}

func (h *HandlerList) HandleExecutorSkills(c tele.Context, state *redis.StateData) error {
  messageId, err := h.stateManager.GetMessageId(c.Sender().ID)
  if err != nil {
    return err
  }

	session, err := h.stateManager.GetSessionData(c.Sender().ID)
  if err != nil {
    return err
  }
  c.Delete()
  session["skills"] = c.Message().Text
  h.stateManager.SetSessionData(c.Sender().ID, session)
  fmt.Println(session)
  h.stateManager.SetState(c.Sender().ID, "executor_contacts", map[string]interface{}{"prev_state": state.State})
	post, err := c.Bot().Send(&tele.Chat{ID: -1002284442496}, fmt.Sprintf(`
#%s
<b>%s</b>

<b>Контакты:</b> @%s

<b>Предварительный бюджет:</b> %s

<b>Описание заказа:</b> %s

<b>Обязанности исполнителя:</b> %s
#%s #%s`, session["target"], session["title"], c.Sender().Username, session["budget"], session["description"], session["skills"], session["category"], session["type"]), &tele.SendOptions{ParseMode: tele.ModeHTML})
	if err != nil {
    return err
  }
	h.stateManager.SetState(c.Sender().ID, "start", map[string]interface{}{})
	_, err = c.Bot().Edit(&tele.Message{Chat: &tele.Chat{ID: c.Sender().ID}, ID: messageId}, fmt.Sprintf(`Ваше объявление размещено: <a href='https://t.me/zylogbot/%d'>Посмотреть пост</a>`, post.ID), &tele.SendOptions{ParseMode: tele.ModeHTML})
	return err
}
