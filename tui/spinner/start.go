package spinner

func Start(spin *Spinner) {
	spin.Done = make(chan bool, 1)
	_, _ = spin.Program.Run()
	spin.Done <- true
	return
}
