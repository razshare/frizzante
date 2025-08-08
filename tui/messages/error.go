package messages

import "main/config"

func Error(text string) {
	Status("ERROR", text, config.Colors.Error, "15", config.Colors.Error)
}
