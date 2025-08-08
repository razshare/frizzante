package messages

import (
	"fmt"
	"main/config"
)

func Section(text string) {
	fmt.Println(config.Styles.Section.Render("## " + text))
}
