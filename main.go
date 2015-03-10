package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	irc "github.com/fluffle/goirc/client"
	"github.com/mattaitchison/envconfig"
)

var nick = envconfig.String("irc_nick", "ircslackline", "nick to use in irc")
var server = envconfig.String("irc_server", "irc.freenode.net:7000", "IRC server to connect to")
var useSSL = envconfig.Bool("irc_ssl", true, "use ssl to connect to server")
var channel = envconfig.String("irc_channel", "#lanciv", "IRC channel to connect to")
var slackChan = envconfig.String("slack_channel", "#lanciv_ircline", "Slack channel to link to")
var slackToken = envconfig.String("slack_token", "", "slack API token")

func main() {
	quit := make(chan bool)

	// api := slack.New(slackToken)
	// api.SetDebug(true)

	// Or, create a config and fiddle with it first:
	cfg := irc.NewConfig(nick, nick, "Gliderbot!")
	cfg.SSL = useSSL
	cfg.Server = server
	cfg.Version = "Gliderbot!"

	client := irc.Client(cfg)

	// Add handlers to do things here!
	// e.g. join a channel on connect.
	client.HandleFunc("connected",
		func(conn *irc.Conn, line *irc.Line) {
			log.Println("Connected to IRC server", server)
			conn.Join(channel)
			// conn.EnableStateTracking()
		})

	client.HandleFunc("privmsg",
		func(conn *irc.Conn, line *irc.Line) {

			// url := fmt.Sprintf("https://slack.com/api/chat.postMessage")
			res, err := http.PostForm("https://slack.com/api/chat.postMessage",
				url.Values{"token": {slackToken}, "channel": {slackChan}, "text": {line.Text()}, "username": {line.Nick}})
			if err != nil {
				log.Fatal(err)
			}

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Fatal(err)
			}
			defer res.Body.Close()

			log.Printf("%s", body)

			// log.Println(line.Nick)

			// log.Println(line.Text())
		})
	// And a signal on disconnect

	if err := client.Connect(); err != nil {
		log.Fatal("IRC connection failed:", err)
	}

	// Wait for disconnect
	<-quit
}
