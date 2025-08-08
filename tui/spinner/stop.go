package spinner

import "fmt"

func Stop(manager *Spinner) {
	manager.Program.Quit()
	<-manager.Done
	fmt.Print("\r\033[K")
}
