package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"boutonsjob/externals/redis"

	tele "gopkg.in/telebot.v4"
)

func (h *HandlerList) HandleMyPosts(c tele.Context, state *redis.StateData) error {
	messageId, err := h.stateManager.GetMessageId(c.Sender().ID)
	if err != nil {
		return err
	}
	posts, err := h.db.GetUserPosts(c.Sender().ID)
	if err != nil {
		return err
	}
	fmt.Println(posts)
	_, err = c.Bot().Edit(&tele.Message{Chat: &tele.Chat{ID: c.Sender().ID}, ID: messageId}, `Мои объявления:`)
	if err != nil {
		return err
	}
	for _, post := range posts {
		deleteBtn := &tele.ReplyMarkup{}
    deleteBtn.Inline(deleteBtn.Row(tele.Btn{Text: "Удалить", Data: fmt.Sprintf("del_%d", post.MessageId)}))
		_, err = c.Bot().Send(&tele.Chat{ID: c.Sender().ID}, fmt.Sprintf(`
#%s
<b>%s</b>
<a href="%s">ссылка</a>
`, post.Type, post.Title, post.Link), &tele.SendOptions{ParseMode: tele.ModeHTML}, deleteBtn)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *HandlerList) HandleDeletePost(c tele.Context, state *redis.StateData) error {
	post_id_str := strings.Split(c.Callback().Data, "_")[1]
	post_id, err := strconv.ParseInt(post_id_str, 10, 64)
	if err != nil {
		return err
	}
	c.Bot().Delete(&tele.Message{Chat: &tele.Chat{ID:  -1002284442496}, ID: int(post_id)})
  err = h.db.DeletePost(int(post_id))
  if err != nil {
    return err
  }
	_, err = c.Bot().Send(c.Sender(), "Ваш пост удалён")
	c.Respond()
	return err
}
