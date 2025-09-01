package spinner

func Stop(spin *Spinner) {
	spin.Program.Quit()
	<-spin.Done
}
