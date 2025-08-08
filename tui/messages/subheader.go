package messages

import (
	"fmt"
	"main/config"
)

func Subheader(text string) {
	fmt.Println(config.Styles.Subheader.Render(text))
}
