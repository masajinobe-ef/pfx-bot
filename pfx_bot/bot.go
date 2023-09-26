package main

import (
	"context"
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	openai "github.com/sashabaranov/go-openai"
)

func main() {
	// Проверка токена
	var pfx_token string = os.Getenv("PFX_BOT")
	var openai_token string = os.Getenv("OPEN_AI")

	if pfx_token != "" {
		fmt.Println("Значение переменной PFX_BOT:", pfx_token)
	} else {
		fmt.Println("Переменная PFX_BOT не найдена")
	}

	if openai_token != "" {
		fmt.Println("Значение переменной OPEN_AI:", openai_token)
	} else {
		fmt.Println("Переменная OPEN_AI не найдена")
	}

	// Авторизаций бота
	bot, err := tgbotapi.NewBotAPI(pfx_token)
	if err != nil {
		log.Panic(err)
	}

	// Консольный режим отладки
	bot.Debug = false

	log.Printf("Авторизовано на аккаунте: %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // Игнорируйте любые обновления, не являющиеся Messages
			continue
		}

		if !update.Message.IsCommand() { // Игнорирует любые не Messages
			continue
		}

		// Создать новый MessageConfig
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Извлечение команды из сообщения
		switch update.Message.Command() {
		case "help":
			msg.Text = "Список команд: /status"
		case "status":
			msg.Text = "OK!"
		default:
			msg.Text = "..."
		}

		// OpenAI
		client := openai.NewClient(openai_token)
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: "Hello!",
					},
				},
			},
		)

		if err != nil {
			fmt.Printf("Ошибка ChatGPT: %v\n", err)
			return
		}

		msg.Text = resp.Choices[0].Message.Content

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
