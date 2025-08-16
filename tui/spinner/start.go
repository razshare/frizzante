package spinner

func Start(s *Spinner) {
	s.Done = make(chan bool, 1)
	_, _ = s.Program.Run()
	s.Done <- true
	return
}
