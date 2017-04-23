package main

import (
	"flag"
	"log"

	"github.com/kyleterry/quasar"
)

const helpText = `Tenyks Feeds:
    tenyks: list feeds
    tenyks: delete feed https://example.com/feeds/thing.rss`
const description = "Add RSS feeds and let tenyks check them for new posts. New posts will be sent to the channel that added the feed."
const version = "1.0.0"

func matchHello(msg quasar.Message) (quasar.Result, error) {
	res := make(quasar.Result)
	if msg.Payload != "hello" {
		return nil, quasar.ErrNoMatch
	}
	return res, nil
}

var matchAdd = quasar.NewRegexMatcher("^add feed (?P<url>(http[s]?)://(.+))$")
var matchRemove = quasar.NewRegexMatcher("^remove feed (?P<url>(http[s]?)://(.+))$")
var matchList = quasar.NewRegexMatcher("^list feeds$")

func main() {
	sendAddr := flag.String("send-addr", "tcp://localhost:61124", "Address to send messages to")
	recvAddr := flag.String("recv-addr", "tcp://localhost:61123", "Address to receive messages from")
	flag.Parse()

	config := &quasar.Config{
		Name:    "Tenyks Feeds",
		Version: "1.0",
		Service: quasar.ServiceConfig{
			SendAddr: sendAddr,
			RecvAddr: recvAddr,
		},
	}

	service := quasar.NewService(config)
	service.HelpText = helpText
	service.Description = description

	service.Handle(quasar.MsgHandler{
		MatcherFunc:  quasar.MatcherFunc(matchAdd),
		DirectOnly:   true,
		HelpText:     "tenyks: add feed https://example.com/feeds/thing.rss",
		MatchHandler: quasar.HandlerFunc(addHandler),
	})

	service.Handle(quasar.MsgHandler{
		MatcherFunc:  quasar.MatcherFunc(matchRemove),
		DirectOnly:   true,
		HelpText:     "tenyks: remove feed https://example.com/feeds/thing.rss",
		MatchHandler: quasar.HandlerFunc(removeHandler),
	})

	service.Handle(quasar.MsgHandler{
		MatcherFunc:  quasar.MatcherFunc(matchList),
		DirectOnly:   true,
		HelpText:     "tenyks: list feeds",
		MatchHandler: quasar.HandlerFunc(listHandler),
	})

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
