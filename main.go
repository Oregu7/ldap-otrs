package main

import (
	"log"

	"github.com/jasonlvhit/gocron"
	"github.com/joho/godotenv"
)

func main() {
	// загружаем переменные окружения
	if godotenv.Load() != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("Служба синхронизации агентов запущена...")

	gocron.Every(1).Day().At("07:30").Do(updateUsersTask)
	<-gocron.Start()
}
