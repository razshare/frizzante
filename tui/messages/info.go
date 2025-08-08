package messages

import "main/config"

func Info(text string) {
	Status("INFO", text, config.Colors.Info, "15", config.Colors.Info)
}
