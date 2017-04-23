package main

import (
	"fmt"
	"log"

	"github.com/kyleterry/quasar"
)

func addHandler(match quasar.Result, msg quasar.Message, com quasar.Communication) {
	if err := com.Send(fmt.Sprintf("Hello, %s!", msg.Nick), msg); err != nil {
		log.Print(err)
	}
}

func removeHandler(match quasar.Result, msg quasar.Message, com quasar.Communication) {
	if err := com.Send(fmt.Sprintf("Hello, %s!", msg.Nick), msg); err != nil {
		log.Print(err)
	}
}

func listHandler(match quasar.Result, msg quasar.Message, com quasar.Communication) {
	if err := com.Send(fmt.Sprintf("Hello, %s!", msg.Nick), msg); err != nil {
		log.Print(err)
	}
}
