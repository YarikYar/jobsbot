package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RedisClient struct {
    client *redis.Client
}

// NewRedisClient создает новый клиент Redis.
func NewRedisClient(addr, password string, db int) *RedisClient {
    rdb := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       db,
    })

    return &RedisClient{client: rdb}
}

// StateData представляет данные состояния пользователя.
type StateData struct {
    State   string                 `json:"state"`   // Текущее состояние пользователя
    Payload map[string]interface{} `json:"payload"` // Дополнительные параметры
}

// SetState устанавливает состояние пользователя и дополнительные параметры.
func (r *RedisClient) SetState(userID int64, state string, payload map[string]interface{}) error {
    key := fmt.Sprintf("user:%d:state", userID)

    stateData := StateData{
        State:   state,
        Payload: payload,
    }
		log.Printf("[DEBUG] StateData changed: %v", stateData)
    // Сериализуем данные в JSON
    data, err := json.Marshal(stateData)
    if err != nil {
        return fmt.Errorf("failed to marshal state data: %v", err)
    }

    // Сохраняем данные в Redis
    err = r.client.Set(ctx, key, data, 0).Err()
    if err != nil {
        return fmt.Errorf("failed to set state in Redis: %v", err)
    }

    return nil
}

// GetState получает состояние пользователя и дополнительные параметры.
func (r *RedisClient) GetState(userID int64) (*StateData, error) {
    key := fmt.Sprintf("user:%d:state", userID)

    // Получаем данные из Redis
    data, err := r.client.Get(ctx, key).Result()
    if err == redis.Nil {
        // Если ключа нет, возвращаем nil без ошибки
        return nil, nil
    } else if err != nil {
        return nil, fmt.Errorf("failed to get state from Redis: %v", err)
    }

    // Десериализуем данные из JSON
    var stateData StateData
    err = json.Unmarshal([]byte(data), &stateData)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal state data: %v", err)
    }

    return &stateData, nil
}

func (r *RedisClient) UpdatePayload(userID int64, payload map[string]interface{}) error {
    key := fmt.Sprintf("user:%d:state", userID)
    stateData, err := r.GetState(userID)
    if err != nil {
        return err
    }
    stateData.Payload = payload
    data, err := json.Marshal(stateData)
    if err != nil {
        return fmt.Errorf("failed to marshal state data: %v", err)
    }
    return r.client.Set(ctx, key, data, 0).Err()
}

func (r *RedisClient) SetSessionData(userID int64, sessionData map[string]interface{}) error {
    key := fmt.Sprintf("user:%d:session_data", userID)
    data, err := json.Marshal(sessionData)
    if err != nil {
        return fmt.Errorf("failed to marshal session data: %v", err)
    }
    return r.client.Set(ctx, key, data, 0).Err()
}

func (r *RedisClient) GetSessionData(userID int64) (map[string]interface{}, error) {
    key := fmt.Sprintf("user:%d:session_data", userID)
    data, err := r.client.Get(ctx, key).Result()
    if err == redis.Nil {
        return nil, nil
    } else if err != nil {
        return nil, fmt.Errorf("failed to get session data from Redis: %v", err)
    }
    var sessionData map[string]interface{}
    err = json.Unmarshal([]byte(data), &sessionData)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal session data: %v", err)
    }
    return sessionData, nil
}

func (r *RedisClient) ClearSessionData(userID int64) error {
    key := fmt.Sprintf("user:%d:session_data", userID)
    return r.client.Del(ctx, key).Err()
}

// ClearState очищает состояние пользователя.
func (r *RedisClient) ClearState(userID int64) error {
    key := fmt.Sprintf("user:%d:state", userID)
    return r.client.Del(ctx, key).Err()
}

func (r *RedisClient) SaveMessageId(userID int64, messageId int) error {
    key := fmt.Sprintf("user:%d:message_id", userID)
    return r.client.Set(ctx, key, messageId, 0).Err()
}

func (r *RedisClient) GetMessageId(userID int64) (int, error) {
    key := fmt.Sprintf("user:%d:message_id", userID)
    messageId, err := r.client.Get(ctx, key).Int()
		log.Printf("[DEBUG] messageId: %v, err: %v", messageId, err)
    return messageId, err
}

func (r *RedisClient) DeleteMessageId(userID int64) error {
    key := fmt.Sprintf("user:%d:message_id", userID)
    return r.client.Del(ctx, key).Err()
}
