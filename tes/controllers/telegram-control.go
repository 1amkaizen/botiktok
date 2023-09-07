package controllers

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Buat fungsi untuk membaca file links.txt dan menyimpan link ke dalam slice
func readLinksFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var links []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		links = append(links, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return links, nil
}

func SetupBot() (*tgbotapi.BotAPI, error) {
	//telegram token
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))

	if err != nil {
		log.Panic(err)
		fmt.Println("MISSING_TELEGRAM_BOT_TOKEN")
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	return bot, nil
}

func SendMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// Panggil fungsi untuk membaca link dari file
	links, err := readLinksFromFile("links.txt")
	if err != nil {
		log.Panic(err)
	}

	// ...

	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	// Analisis pesan pengguna
	userInput := update.Message.Text
	var foundLink string

	// Cari link yang sesuai dengan kata kunci
	for _, link := range links {
		parts := strings.Split(link, ": ")
		if len(parts) == 2 && strings.TrimSpace(parts[0]) == userInput {
			foundLink = parts[1]
			break // Hentikan pencarian setelah menemukan link yang sesuai
		}
	}

	// Kirimkan link yang sesuai atau pesan jika tidak ada yang sesuai
	if foundLink != "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, foundLink)
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)

	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Kata kunci tidak ditemukan.")
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
	}
}

func HandleStartCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	log.Printf("UserName :%s", update.Message.From.UserName)
	log.Printf("ID :%d", update.Message.Chat.ID)
	log.Printf("Text: %s", update.Message.Text)

	// button
	twitterButton := tgbotapi.NewInlineKeyboardButtonData("Twitter🐦", "twitter")
	githubButton := tgbotapi.NewInlineKeyboardButtonData("Github🐙", "github")
	railwayButton := tgbotapi.NewInlineKeyboardButtonData("Railway🚂", "railway")
	replitButton := tgbotapi.NewInlineKeyboardButtonData("Replit🚀", "replit")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			twitterButton,
			githubButton,
		),
		tgbotapi.NewInlineKeyboardRow(
			railwayButton,
			replitButton,
		),
	)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hallo, @"+update.Message.From.UserName+"! Selamat datang di bot saya, bagaimana saya bisa membantumu hari ini?")
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ReplyMarkup = keyboard
	bot.Send(msg)

	// send message to me
	SECRET := os.Getenv("SECRET")
	SECRET64, _ := strconv.ParseInt(SECRET, 10, 64)
	msgToYou := tgbotapi.NewMessage(SECRET64, "User @"+update.Message.From.UserName+" with ID:"+strconv.FormatInt(update.Message.Chat.ID, 10)+" masuk")

	msg.ReplyToMessageID = update.Message.MessageID

	bot.Send(msgToYou)

}

func HandleHelpCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "hallo ini help")
	msg.ReplyToMessageID = update.Message.MessageID

	bot.Send(msg)

}
