package messages

import (
	"fmt"
	"github.com/razshare/frizzante/tui/config"
)

func DockerHelp() {
	fmt.Println(config.Styles.Title.Render("🐙 You're running Frizzante in Docker!"))

	fmt.Println(config.Styles.Subheader.Render("⚡ Simple workflow:"))

	fmt.Println(config.Styles.Status(config.Colors.Info).Render("• Attach to the container: ") +
		config.Styles.Example.Render("docker exec -it frizzante-start sh"))

	fmt.Println(config.Styles.Status(config.Colors.Info).Render("• Run environment in container:"))
	fmt.Println(config.Styles.Item.Render("    • Dev environment: ") + config.Styles.Flag.Render("make dev"))
	fmt.Println(config.Styles.Item.Render("    • Prod environment: ") + config.Styles.Flag.Render("make build"))
	fmt.Println(config.Styles.Item.Render("    • To run the app: ") + config.Styles.Flag.Render("./.gen/bin/app"))

	fmt.Println(config.Styles.Status(config.Colors.Info).Render("• Run prod via docker:"))
	fmt.Println(config.Styles.Item.Render("    • Build image: ") +
		config.Styles.Example.Render("docker build --target frizzante_prod -t my-app:prod ."))
	fmt.Println(config.Styles.Item.Render("    • Run image: ") +
		config.Styles.Example.Render("docker run -p 8080:8080 my-app:prod"))
	fmt.Println(config.Styles.Item.Render("    • Via docker compose: ") +
		config.Styles.Example.Render("docker compose -f compose.yaml -f compose.prod.yaml up -d --build"))

	fmt.Println(config.Styles.Subheader.Render("🎉 Enjoy!!"))
	fmt.Println()
}
