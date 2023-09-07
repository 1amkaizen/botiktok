package controllers

import (
	"fmt"
	"log"
	"os"
	"project/tiktokapi/encode"
	"regexp"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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
	// Memeriksa pesan dari pengguna
	if update.Message == nil {
		return
	}

	// Mendapatkan teks pesan dari pengguna
	text := update.Message.Text

	// Membuat ekspresi reguler untuk mencari kata kunci
	keywordPattern := regexp.MustCompile(`\b(deskripsi|url|kategori|gambar|nama_produk|waktu_pembuatan)\b`)

	// Mencari kata kunci dalam pesan
	keywords := keywordPattern.FindAllString(text, -1)

	if len(keywords) > 0 {
		// Jika ada kata kunci dalam pesan, maka tanggapi
		for _, keyword := range keywords {
			// Menangani masing-masing kata kunci
			switch keyword {
			case "deskripsi":
				// Panggil fungsi Encode untuk mendapatkan deskripsi
				description, err := encode.GetDesc("file.json")
				if err != nil {
					// Tanggapi dengan pesan kesalahan
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Terjadi kesalahan saat mengambil deskripsi.")
					msg.ReplyToMessageID = update.Message.MessageID
					bot.Send(msg)
					return
				}

				// Kirim deskripsi ke pengguna
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, description)
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			case "nama_produk":
				// Panggil fungsi Encode untuk mendapatkan deskripsi
				nama, err := encode.GetProductName("file.json")
				if err != nil {
					// Tanggapi dengan pesan kesalahan
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Terjadi kesalahan saat mengambil nama.")
					msg.ReplyToMessageID = update.Message.MessageID
					bot.Send(msg)
					return
				}

				// Kirim deskripsi ke pengguna
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, nama)
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			case "url":
				// Panggil fungsi GetUniqueURLs untuk mendapatkan URL yang unik
				uniqueURLs, err := encode.GetUniqueURLs("file.json")
				if err != nil {
					// Tanggapi dengan pesan kesalahan
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Terjadi kesalahan saat mengambil URL.")
					msg.ReplyToMessageID = update.Message.MessageID
					bot.Send(msg)
					return
				}

				// Kirim URL satu per satu ke pengguna
				for _, url := range uniqueURLs {
					urlMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "URL: "+url)
					bot.Send(urlMsg)
				}
			case "kategori":
				// Panggil fungsi GetCategories untuk mendapatkan kategori
				categories, err := encode.GetCategories("file.json")
				if err != nil {
					// Tanggapi dengan pesan kesalahan
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Terjadi kesalahan saat mengambil kategori.")
					msg.ReplyToMessageID = update.Message.MessageID
					bot.Send(msg)
					return
				}

				// Kirim kategori ke pengguna
				categoryMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Kategori:")
				for _, category := range categories {
					categoryMsg.Text += "\n" + category
				}

				bot.Send(categoryMsg)
			case "gambar":
				// Panggil fungsi GetImages untuk mendapatkan informasi gambar
				images, err := encode.GetImages("file.json")
				if err != nil {
					// Tanggapi dengan pesan kesalahan
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Terjadi kesalahan saat mengambil informasi gambar.")
					msg.ReplyToMessageID = update.Message.MessageID
					bot.Send(msg)
					return
				}

				// Kirim informasi gambar ke pengguna
				imageMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Informasi Gambar:")
				for _, image := range images {
					imageMsg.Text += "\n" + image
				}
				bot.Send(imageMsg)
			case "waktu_pembuatan":
				// Panggil fungsi GetCreateTime untuk mendapatkan waktu pembuatan
				createTime, err := encode.GetCreateTime("file.json")
				if err != nil {
					// Tanggapi dengan pesan kesalahan
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Terjadi kesalahan saat mengambil waktu pembuatan.")
					msg.ReplyToMessageID = update.Message.MessageID
					bot.Send(msg)
					return
				}

				// Kirim waktu pembuatan ke pengguna
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Waktu Pembuatan: "+strconv.FormatInt(createTime, 10))
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)

			}
		}
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

func HandleProductCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

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

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "silahkan pilih")
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ReplyMarkup = keyboard
	bot.Send(msg)

}
