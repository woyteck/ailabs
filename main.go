package main

import (
	"os"
	"woyteck/ailabs/task_meme"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("could not load .env file")
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
	// task_whoami.TaskWhoami(apiKey)
	// task_search.TaskSearch(apiKey)
	// task_people.TaskPeople(apiKey)
	// task_knowledge.TaskKnowledge(apiKey)
	// task_tools.TaskTools(apiKey)
	// task_gnome.TaskGnome(apiKey)
	// task_ownapi.TaskOwnapi(apiKey)
	// task_ownapipro.TaskOwnapipro(apiKey)
	task_meme.TaskMeme(apiKey)
}
