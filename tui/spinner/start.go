package spinner

import "time"

func Start(manager *Spinner) (err error) {
	manager.Done = make(chan bool)
	go func() {
		_, err = manager.Program.Run()
		close(manager.Done)
	}()
	time.Sleep(100 * time.Millisecond)
	return
}
