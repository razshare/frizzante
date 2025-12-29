package actions

func Generate(options GenerateOptions) (err error) {
	generation := options.Generation
	if generation == ":pick" {
		generation = ""
	}

	return
}
