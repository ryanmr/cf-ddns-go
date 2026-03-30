package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"

	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

type Payload struct {
	Username  string  `json:"username"`
	AvatarURL string  `json:"avatar_url"`
	Content   string  `json:"content"`
	Embeds    []Embed `json:"embeds"`
}

type Embed struct {
	Author      Author  `json:"author"`
	Title       string  `json:"title"`
	URL         string  `json:"url"`
	Description string  `json:"description"`
	Color       int     `json:"color"`
	Fields      []Field `json:"fields"`
	Footer      Footer  `json:"footer"`
}

type Author struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	IconURL string `json:"icon_url"`
}

type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type Footer struct {
	Text string `json:"text"`
}

func handleDiscordWebhook(newIP string) {
	url := getDiscordWebhookUrl()

	payload := getPayload(newIP)

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal discord payload")
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error().Err(err).Msg("Failed to create discord request")
		return
	}
	req.Header.Set("Content-Type", "application/json")

	log.Info().Msg("Sending discord webhook request...")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send discord webhook")
		return
	}
	defer resp.Body.Close()

	log.Info().Msg("Sent discord webhook request")

	if resp.StatusCode != http.StatusNoContent {
		log.Warn().Msg("Discord complained about this webhook and did not respond with 204")
	}
}

func getDiscordWebhookUrl() string {
	value, ok := os.LookupEnv("DISCORD_WEBHOOK_URL")
	if !ok {
		log.Fatal().Msg("DISCORD_WEBHOOK_URL was not set")
	}
	return value
}

func getPayload(newIP string) Payload {
	randomIndex := rand.Intn(len(quips))
	randomSilly := quips[randomIndex]

	return Payload{
		Username:  "Ryanbot",
		AvatarURL: "https://i.imgur.com/4M34hi2.png",
		Content:   "👋 There's been a network infrastructure update.\n\nFor details on this script, visit [ryanmr/cf-ddns-go](https://github.com/ryanmr/cf-ddns-go).",
		Embeds: []Embed{
			{
				Author: Author{
					Name:    "Server",
					URL:     fmt.Sprintf("https://ryanrampersad.com/?s=server&cloudflare&ip=%s&idx=%d&t=%d", newIP, randomIndex, time.Now().Unix()),
					IconURL: "https://i.imgur.com/4M34hi2.png",
				},
				Title:       "IP Updated",
				URL:         fmt.Sprintf("https://ryanrampersad.com/?s=ip&cloudflare&ip=%s&idx=%d&t=%d", newIP, randomIndex, time.Now().Unix()),
				Description: fmt.Sprintf("New IP is %s", newIP),
				Color:       15258703,
				Fields:      []Field{},
				Footer: Footer{
					Text: randomSilly,
				},
			},
		},
	}
}

var quips []string = []string{
	"Wow, how unusual. Did the power go out?",
	"It's Minnesota. It's probably another storm.",
	"It's Minnesota. It's probably another drought.",
	"Is your router running? You better go catch it.",
	"Try as you might, but you will never escape Cloudflare.",
	"Hello, is this thing on?",
	"One wrong move and you'll never connect to your precious server again...",
	"Finally, that old ip has been getting so stale.",
	"You were supposed to cut the red wire!",
	"Turns out, it's just random sometimes.",
	"For a second there, I thought Google discontinued the Internet.",
	"For a second there, I thought Microsoft extinguished the Internet.",
	"For a second there, I thought Apple walled off the Internet.",
	"Introducing our brand new revolutionary product, network connectivity. Again.",
	"No, seriously, it's updated so you should go look.",
	"Sometimes when there's an update, I just don't tell anyone anyway.",
	"REDACTED. CLASSIFIED.",
	"Going, going, gone!",
	"Have you tried turning the entire internet off and on again?",
	"Plot twist: the IP was inside us all along.",
	"This is your ISP's way of keeping things exciting.",
	"New IP just dropped. Limited edition.",
	"I've seen more stable connections on a tin can telephone.",
	"Somewhere, a DHCP lease just expired and chose violence.",
	"Don't worry, I wrote it down on a sticky note this time.",
	"The packets must flow.",
	"Another day, another dollar, another IP address.",
	"I blame cosmic rays.",
	"This wouldn't happen if we all just used IPv6.",
	"Your ISP called. They said 'lol.'",
	"Breaking news: local server changes address, neighbors unaffected.",
	"In this economy? Even IP addresses can't afford to stay put.",
	"The DNS gods demand a sacrifice.",
	"Fun fact: this is the 4,294,967,296th coolest IP address possible.",
	"I'm not saying it was aliens, but it was definitely your modem.",
	"Achievement unlocked: IP Address Rotation Survivor.",
	"If anyone asks, we were always at this address.",
	"Roses are red, violets are blue, your IP changed, and your DNS did too.",
}
