package main

import (
	"github.com/johnjones4/golden-hour-bot/lib/runtime"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	go runtime.StartLocalAlertQueuer()
	go runtime.StartLocalDequeuerRuntime()

	runtime.StartLocalServerRuntime()
}
