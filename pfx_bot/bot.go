package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

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
	bot, bot_err := tgbotapi.NewBotAPI(pfx_token)
	// Консольный режим отладки
	bot.Debug = false

	log.Printf("Авторизовано на аккаунте: %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	// Лог ошибок bot_err
	if bot_err != nil {
		log.Panic(bot_err)
	}

	for update := range updates {
		if update.Message == nil { // Игнорируйте любые обновления, не являющиеся Messages
			continue
		}

		// OpenAI
		client := openai.NewClient(openai_token)
		messages := make([]openai.ChatCompletionMessage, 0)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.Text = update.Message.Text

		//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if strings.Contains(msg.Text, "@"+bot.Self.UserName) {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: msg.Text,
			})

			response, gpt_err := client.CreateChatCompletion(
				context.Background(),
				openai.ChatCompletionRequest{
					Model:    openai.GPT3Dot5Turbo,
					Messages: messages,
				},
			)

			content := response.Choices[0].Message.Content
			//log.Printf(content)
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: content,
			})

			if gpt_err != nil {
				log.Printf("Ошибка чата: %v\n", gpt_err)
			}

			ai_msg := tgbotapi.NewMessage(update.Message.Chat.ID, content)
			ai_msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(ai_msg)
		}
	}
}
