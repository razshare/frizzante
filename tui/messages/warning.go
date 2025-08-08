package messages

import "main/config"

func Warning(text string) {
	Status("WARNING", text, config.Colors.Warning, "0", config.Colors.Warning)
}
