package main

import (
	"os"
	"woyteck/ailabs/task_whoami"

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
	// task_liar.TaskLiar(apiKey)
	// task_inprompt.TaskInprompt(apiKey)
	// task_embedding.TaskEmbedding(apiKey)
	// task_whisper.TaskWhisper(apiKey)
	// task_functions.TaskFunctions(apiKey)
	// task_rodo.TaskRodo(apiKey)
	// task_scraper.TaskScraper(apiKey)
	task_whoami.TaskWhoami(apiKey)
}
