package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"github.com/joho/godotenv"
	migrate "github.com/govel-golang-migration/govel-golang-migration"
)

func main() {
	flag.Parse()

	method := flag.Arg(0)

	dispatch(method)
}

func dispatch(method string) {
	fmt.Printf("dispatch %s method \n", method)

	switch method {
	case "install":
		fmt.Println("install ...")

		loadEnv()

		migrate.Install(os.Getenv("DB_HOST"))
	case "make":
		fmt.Println("make")

		fileName := strings.Trim(flag.Arg(1), "")

		if fileName == "" {
			fmt.Println("file name is required")

			return
		}

		migrate.Make(fileName)
	case "migrate":
		fmt.Println("migrate")

		loadEnv()

		migrate.Migrate(os.Getenv("DB_HOST"))
	case "rollback":
		fmt.Println("rollback")

		stage, err := strconv.Atoi(strings.Trim(flag.Arg(1), ""))

		if err != nil {
			fmt.Println("stage is required and must be numeric")

			return
		}

		loadEnv()

		migrate.Rollback(stage, os.Getenv("DB_HOST"))
	default:
		fmt.Println("method not support")
	}
}

func loadEnv() {
	_, filename, _, _ := runtime.Caller(0)

	envFile, _ := filepath.Abs(filepath.Join(filename, "../../.env"))

	if err := godotenv.Load(envFile); err != nil {
		panic(err)
	}
}
