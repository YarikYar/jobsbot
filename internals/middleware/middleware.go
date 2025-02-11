package middleware

import (
	"log"

	"boutonsjob/externals/redis"

	tb "gopkg.in/telebot.v4"
)

// StateMiddleware - это миддлвара, которая перехватывает колбэки и устанавливает состояние пользователя.
func StateMiddleware(redisClient *redis.RedisClient) tb.MiddlewareFunc {
	return func(next tb.HandlerFunc) tb.HandlerFunc {
		return func(c tb.Context) error {
			// Проверяем, является ли текущее событие колбэком
			callback := c.Callback()
			if callback != nil {
				// Получаем данные из Callback.Data
				data := callback.Data
				currentState, err := redisClient.GetState(c.Sender().ID)
        if err != nil {
          log.Printf("Failed to get state for user %d: %v", c.Sender().ID, err)
          return c.Send("Произошла ошибка при обработке запроса.")
        }
				
				// Устанавливаем состояние пользователя в Redis
				if data == "back" {
					if currentState.Payload["previous_state"] != nil {
            data = currentState.Payload["previous_state"].(string)
					} else {
            data = "start"
					}
				}
				userID := c.Sender().ID
				payload := map[string]interface{}{
					"previous_state": currentState.State,
				}

				err = redisClient.SetState(userID, data, payload)
				if err != nil {
					log.Printf("Failed to set state for user %d: %v", userID, err)
					return c.Send("Произошла ошибка при обработке запроса.")
				}

				log.Printf("State set for user %d: %s", userID, data)
			}

			// Продолжаем выполнение следующего обработчика
			return next(c)
		}
	}
}
