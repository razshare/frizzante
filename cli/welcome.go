package cli

import (
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"os"
)

func OnWelcome() {
	usingDocker := os.Getenv("FRIZZANTE_USING_DOCKER")
	end := make(chan string)
	err := pterm.DefaultBigText.WithLetters(putils.LettersFromStringWithStyle("Frizzante", pterm.FgCyan.ToStyle())).Render()
	if err != nil {
		Fatal(err)
	}

	if usingDocker != "" {
		println("")
		println("🐙 You're running Frizzante in Docker!")
		println("")
		println("⚡️ Simple workflow:")
		println("• Attach to the container: docker exec -it frizzante-start sh")
		println("• Run environment in container:")
		println("    • Dev environment: make dev")
		println("    • Prod environment: make build")
		println("    • To run the app: ./.gen/bin/app")
		println("• Run prod via docker:")
		println("    • Build image: docker build --target frizzante_prod -t my-app:prod .")
		println("    • Run image: docker run -p 8080:8080 my-app:prod")
		println("    • Via docker compose: docker compose -f compose.yaml -f compose.prod.yaml up -d --build")
		println("🎉 Enjoy!!")
		println("")
		Info("For more info: https://razshare.github.io/frizzante-docs/guides/get-started/")
	}

	<-end
	Success("Bye!")
}
