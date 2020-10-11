package main

import (
	"github.com/subosito/gotenv"

	"os"
)

func main() {
	_ = gotenv.Load(".env")
	token := os.Getenv("token")

	client := NewClient(token)
	err := client.SaveToFile("report.csv", 100)
	if err != nil {
		println(err.Error())
	}

}

