package spinner

func Start(s *Spinner) (err error) {
	s.Done = make(chan bool, 1)
	go func() {
		_, err = s.Program.Run()
		close(s.Done)
	}()
	return
}
