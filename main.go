package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

var (
	token string
)

func init() {

	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	fmt.Println("\nBay Bay Maszkaro!")
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Printf("%v on %v: %v\n", m.Author.Username, m.ChannelID, m.Content)
	if m.Author.ID == s.State.User.ID {
		return
	}

	if len(m.Content) <= 1 || m.Content[0:2] != "->" {
		return
	}

	command := strings.Split(m.Content, " ")

	switch command[0] {
	case "->throw":
		throw(s, m)
	case "->testimage":
		testImage(s, m)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("->Gotowe %v", m.Author.Username))
	default:
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("->NieWkurwiajMnie %v", m.Author.Username))
	}
}

func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{255, 255, 0, 255}
    point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}

func testImage(s *discordgo.Session, m *discordgo.MessageCreate) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	addLabel(img, 20, 30, "Hello Go")

	f, err := os.Create(fmt.Sprintf("/tmp/discordgo/welcome_%v.png", m.Author.Username))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	// var f io.Writer
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
	// file := discordgo.File{
	// 	Name:        fmt.Sprintf("%v", m.Author.Username),
	// 	ContentType: "image/png",
	// 	Reader:      f,
	// }

	// message := discordgo.MessageSend{
	// 	File: &file,
	// }
	_, err = s.ChannelFileSend(m.ChannelID, "asd", f)
	if err!=nil {
		panic(err)
	}
	// return f
	// foo := bufio.NewWriter(f)
	// return foo
}
