package messages

import "main/config"

func Success(text string) {
	Status("SUCCESS", text, config.Colors.Success, "0", config.Colors.Success)
}
