package spinner

func Stop(s *Spinner) {
	s.Program.Quit()
	<-s.Done
}
