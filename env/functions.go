package env

import "github.com/joho/godotenv"

// LoadDotenv loads environments variables from a dotenv file.
//
// If you call LoadDotenv without any args it will default to loading .env in the current path.
//
// LoadDotenv will not override an env variable that already exists.
func LoadDotenv(fileNames ...string) error {
	return godotenv.Load(fileNames...)
}
