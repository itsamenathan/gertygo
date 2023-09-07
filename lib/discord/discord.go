package discord

import (
	"fmt"
	"log"
	"main/lib/finnhub"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/imroc/req/v3"
	"github.com/tidwall/gjson"
)

var (
	session *discordgo.Session
	botID   string
)

func Init(token string) {
	var err error
	session, err = discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	session.AddHandler(messageCreate)
	session.Identify.Intents = discordgo.IntentsGuildMessages

	err = session.Open()
	if err != nil {
		fmt.Println("error opening Discord session,", err)
		return
	}

	user, err := session.User("@me")
	if err != nil {
		fmt.Println("Error fetching bot user:", err)
		return
	}
	botID = user.ID
}

func Close() {
	err := session.Close()
	if err != nil {
		fmt.Println("error closing Discord session,", err)
	}
}

func send(channelID string, message string) {
	_, err := session.ChannelMessageSend(channelID, message)
	if err != nil {
		log.Println(err)
	}
}
func sendEmbed(channelID string, message discordgo.MessageEmbed) {
	_, err := session.ChannelMessageSendEmbed(channelID, &message)
	if err != nil {
		log.Println(err)
	}
}

func handleStockCommand(m *discordgo.MessageCreate, args []string) {
	if len(args) < 2 {
		return
	}

	// Get the stock quote
	stockSymbol := args[2]
	quote, err := finnhub.GetStockQuote(stockSymbol)
	if err != nil {
		log.Printf("Error getting stock quote for %s: %+v", stockSymbol, err)
		return
	}
	// Format the message
	var (
		color              int
		arrow              string
		priceCurrent       = quote.GetC()
		priceChange        = quote.GetD()
		priceChangePercent = quote.GetDp()
	)
	if priceChange < 0 {
		color = 0xFF0000 //red
		arrow = "🔻"
	} else {
		color = 0x00FF00 //green
		arrow = "🔺"
	}
	embed := discordgo.MessageEmbed{
		Title:       strings.ToUpper(stockSymbol),
		Description: fmt.Sprintf("$%.2f %s $%.2f (%.2f%%)", priceCurrent, arrow, priceChange, priceChangePercent),
		Color:       color,
	}

	// Send the message
	sendEmbed(m.ChannelID, embed)
}

func handleCatFactCommand(m *discordgo.MessageCreate) {
	client := req.C()
	resp, err := client.R().Get("https://catfact.ninja/fact")
	if err != nil {
		log.Println(err)
		return
	}
	value := gjson.Get(resp.String(), "fact")
	embed := discordgo.MessageEmbed{
		Title:       "Cat Fact",
		Description: value.String(),
	}
	sendEmbed(m.ChannelID, embed)
}

func handleCatPicCommand(m *discordgo.MessageCreate) {
	client := req.C()
	resp, err := client.R().Get("https://api.thecatapi.com/v1/images/search")
	if err != nil {
		log.Println(err)
		return
	}
	image := discordgo.MessageEmbedImage{
		URL: gjson.Get(resp.String(), "0.url").String(),
	}
	embed := discordgo.MessageEmbed{
		Title:       "Cat Pic",
		Description: "Look at this majestic beast!",
		Image:       &image,
	}
	sendEmbed(m.ChannelID, embed)
}

func handleDogPicCommand(m *discordgo.MessageCreate) {
	client := req.C()
	resp, err := client.R().Get("https://api.thedogapi.com/v1/images/search")
	if err != nil {
		log.Println(err)
		return
	}
	image := discordgo.MessageEmbedImage{
		URL: gjson.Get(resp.String(), "0.url").String(),
	}
	embed := discordgo.MessageEmbed{
		Title:       "Doc Pic",
		Description: "Look at this majestic beast!",
		Image:       &image,
	}
	sendEmbed(m.ChannelID, embed)
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == botID {
		return
	}

	log.Printf("message: %s\n", m.Content)

	// Ignore all messages that don't mention me at first
	prefix := fmt.Sprintf("<@%s>", botID)
	if !strings.HasPrefix(m.Content, prefix) {
		return
	}

	// Parse message into array of string
	args := strings.Split(m.Content, " ")
	if len(args) < 1 {
		return
	}

	command := args[1]

	switch command {
	case "ping":
		send(m.ChannelID, "Pong!")
	case "stock":
		handleStockCommand(m, args)
	case "catfact":
		handleCatFactCommand(m)
	case "catpic":
		handleCatPicCommand(m)
	case "dogpic":
		handleDogPicCommand(m)
	}
}
