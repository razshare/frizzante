package spinner

import "fmt"

func Stop(s *Spinner) {
	s.Program.Quit()
	<-s.Done
	fmt.Print("\r\033[K")
}
