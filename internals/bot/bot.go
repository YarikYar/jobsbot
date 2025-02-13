package bot

import (
	"log"
	"strings"

	"boutonsjob/externals/db"
	"boutonsjob/externals/redis"
	"boutonsjob/internals/handlers"
	"boutonsjob/internals/middleware"

	tele "gopkg.in/telebot.v4"
)

type Bot struct {
	bot          *tele.Bot
	stateManager *redis.RedisClient
	handlers     handlers.HandlerList
	db           *db.DB
}

func NewBot(token string, stateManager *redis.RedisClient, db *db.DB) *Bot {
	bot, err := tele.NewBot(tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10},
	})
	if err != nil {
		panic(err)
	}
	bot.Use(middleware.StateMiddleware(stateManager))
	return &Bot{bot: bot, stateManager: stateManager, handlers: handlers.InitHandlers(stateManager, db), db: db}
}

func (b *Bot) RegisterHandlers() {
	b.handlers.AddText("start", b.handlers.HandleStart)
	b.handlers.AddCallback("start", b.handlers.HandleStart)
	b.handlers.AddCallback("find_executor", b.handlers.HandleFindExecutor)
	b.handlers.AddCallbackFilter("exec_category", func(name string, state string) bool {
		return strings.HasPrefix(state, "ec_") 
	}, b.handlers.HandleExecutorCategory)
	b.handlers.AddCallbackFilter("exec_target", func(name string, state string) bool {
    return strings.HasPrefix(state, "et_")
  }, b.handlers.HandleExecutorTarget)
	b.handlers.AddCallback("!ec_other", b.handlers.HandleExecutorCategoryOther)
	b.handlers.AddCallbackFilter("es", func (name string, state string) bool {
    return strings.HasPrefix(state, "es_")	}, b.handlers.HandleExecutorSphere)
  b.handlers.AddText("wait_executor_category", b.handlers.HandleExecutorCategoryOtherName)
	b.handlers.AddText("executor_title", b.handlers.HandleExecutorTitle)
  b.handlers.AddText("executor_budget", b.handlers.HandleExecutorBudget)
	b.handlers.AddText("executor_description", b.handlers.HandleExecutorDescription)
  b.handlers.AddText("executor_skills", b.handlers.HandleExecutorSkills)

}

func (b *Bot) Run() error {
	b.bot.Handle("/restart", func(c tele.Context) error {
		b.stateManager.ClearState(c.Sender().ID)
    b.stateManager.ClearSessionData(c.Sender().ID)
		b.stateManager.DeleteMessageId(c.Sender().ID)
		return c.Send("State cleared")
	})
	b.bot.Handle(tele.OnText, func(c tele.Context) error {
		state, err := b.stateManager.GetState(c.Sender().ID)
		if state == nil || err != nil {
			b.stateManager.SetState(c.Sender().ID, "start", nil)
		}
		state, err = b.stateManager.GetState(c.Sender().ID)
		log.Printf("[DEBUG] State: %v", state)
		if err != nil {
			return err
		}

		for _, handler := range b.handlers.HandlersText {
			if handler.Filter(handler.StateName, state.State) {
				return handler.Handler(c, state)
			}
		}
		return nil
	})

	b.bot.Handle(tele.OnCallback, func(c tele.Context) error {
		state, err := b.stateManager.GetState(c.Sender().ID)
		if err != nil {
			return err
		}
		log.Printf("[DEBUG] State: %v", state)

		for _, handler := range b.handlers.HandlersCallback {
			if handler.Filter(handler.StateName, state.State) {
				return handler.Handler(c, state)
			}
		}
		return nil
	})

	b.bot.Start()
	return nil
}
