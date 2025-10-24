package spinners

func Stop(spin *Spinner) {
	spin.Program.Quit()
	<-spin.Done
}
