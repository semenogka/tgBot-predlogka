package main

import (
	"log"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var adminID int64 // чат ID того человека, который принимает/отклоняет посты
var channelID int64 // ID вашего канала

func main() {
	bot, err := tg.NewBotAPI("7504491183:AAEWT9H4SMj_clb_FViTrdv09TNOz44UuKM")
	if err != nil {
		log.Println("Ошибка токена: ", err)
		return
	}

	u := tg.NewUpdate(60)
	u.Timeout = 0
	updates := bot.GetUpdatesChan(u)

	keyboard := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("✅", "y"),
			tg.NewInlineKeyboardButtonData("❌", "n"),
		),
	)

	for update := range updates {
		if update.Message != nil {
			if update.Message.Text == "/start" {
				msg := tg.NewMessage(update.Message.Chat.ID, "предложите пост")
	            bot.Send(msg)
			}

			if update.Message.Text != "" && update.Message.Text != "/start"{
				msg := tg.NewMessage(adminID, update.Message.Text)
				msg.ReplyMarkup = keyboard
				bot.Send(msg)
			}

			if update.Message.Photo != nil {
				photo := (*&update.Message.Photo[len(*&update.Message.Photo)-1])

				fileID := photo.FileID
				msg := tg.NewPhoto(adminID, tg.FileID(fileID))

				if update.Message.Caption != "" {
					msg.Caption = update.Message.Caption + "\n"  + "\n" + "👤" + update.Message.Chat.FirstName
				}
				msg.ReplyMarkup = keyboard
				bot.Send(msg)
			}
		}

		if update.CallbackQuery != nil {
			callback := update.CallbackQuery

			if callback.Data == "y" {
				if callback.Message.Text != "" {
					msg := tg.NewMessage(channelID, callback.Message.Text)
					bot.Send(msg)
                    bot.Request(tg.NewCallback(callback.ID, "пост принят!"))
				}

				if callback.Message.Photo != nil {
					photo := (*&callback.Message.Photo[len(*&callback.Message.Photo)-1])

					fileID := photo.FileID
					msg := tg.NewPhoto(channelID, tg.FileID(fileID))

					if callback.Message.Caption != "" {
						msg.Caption = callback.Message.Caption
					}
					bot.Send(msg)

					bot.Request(tg.NewCallback(callback.ID, "пост принят!"))
				}

			} else if callback.Data == "n" {
				deleteMsg := tg.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
				bot.Send(deleteMsg)
                
				bot.Request(tg.NewCallback(callback.ID, "Пост отклонён и удалён."))
			}
		}
	}
}


