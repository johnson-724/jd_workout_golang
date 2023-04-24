package file

import "github.com/joho/godotenv"

func LoadConfigAndEnv(){
	loadEnv()
}

func loadEnv() {
	path := AccessFromCurrentDir(".env")

	if err := godotenv.Load(path); err != nil {
		panic(err)
	}
}