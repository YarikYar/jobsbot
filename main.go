package main

import (
	"boutonsjob/externals/redis"
	"boutonsjob/externals/db"
	"boutonsjob/internals/bot"
)


func main() {
	redisClient := redis.NewRedisClient("localhost:6379", "", 0)
	db, err := db.NewDB("postgres://postgres:postgres@localhost:5432/boutonsjob?sslmode=disable")
	if err != nil {
    panic(err)
  }
  b := bot.NewBot("7921832369:AAG4KPh9ht6RZ6MZ1XxGgPUEap45jGVu00Y", redisClient, db)
  b.RegisterHandlers()
  err = b.Run()
  if err != nil {
    panic(err)
  }
}
