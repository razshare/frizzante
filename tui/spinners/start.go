//go:build !no_spinners

package spinners

func Start(spin *Spinner) {
	spin.Done = make(chan struct{}, 1)
	_, _ = spin.Program.Run()
	spin.Done <- struct{}{}
	return
}
