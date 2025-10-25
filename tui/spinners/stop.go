//go:build !no_spinners

package spinners

func Stop(spin *Spinner) {
	spin.Program.Quit()
	<-spin.Done
}
