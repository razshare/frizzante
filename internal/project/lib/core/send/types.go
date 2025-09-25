package send

type FileOrElseConfig struct {
	UseDisk bool
	OrElse  func()
}
