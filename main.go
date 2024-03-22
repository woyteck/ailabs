package main

import (
	"os"
	"woyteck/ailabs/task_liar"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Coult not load .env file")
	}

	apiKey := os.Getenv("AIDEVS_KEY")

	// task_helloapi.TaskHelloApi(apiKey)
	// task_moderation.TaskModeration(apiKey)
	// task_blogger.TaskBlogger(apiKey)
	task_liar.TaskLiar(apiKey)
}
