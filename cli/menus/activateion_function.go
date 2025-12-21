package menus

type ActivationFunction func(query string) (value string, active bool)
