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
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

var (
	token string
	// db    *sql.DB
	err error
)

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	// db, err = sql.Open("mysql", "username:password@tcp(127.0.0.1:3306)/test")
	// if err != nil {
	// 	panic(err.Error())
	// }
	// defer db.Close()

	// db.SetConnMaxLifetime(time.Minute * 3)
	// db.SetMaxOpenConns(10)
	// db.SetMaxIdleConns(10)

	// insert, err := db.Query("INSERT INTO test VALUES ( 2, 'TEST' )")
	// if err != nil {
	// 	panic(err.Error())
	// }
	// defer insert.Close()

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Initial handlers
	dg.AddHandler(ready)

	// Handlers for Chat
	dg.AddHandler(messageCreate)
	dg.AddHandler(messageLogging)
	dg.AddHandler(messagePingPrevent)

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

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

func ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateStatus(1, "IDKCloud.com")
}
func messageLogging(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Printf("%v on %v: %v\n", m.Author.Username, m.ChannelID, m.Content)
}

func messagePingPrevent(s *discordgo.Session, m *discordgo.MessageCreate) {
	// userGroups, err := s.UserGuilds(100, "", "")
	// if err != nil {
	// 	panic(err)
	// }
	// for i, v := range userGroups {
	// 	fmt.Println(i, v.ID, v.Name)
	// }

	fmt.Println(m.Member.Roles)
	fmt.Println(m.GuildID)
	guild, err := s.State.Guild(m.GuildID)
	if err != nil {
		panic(err)
	}

	fmt.Println(guild.Name)

	for _, v := range guild.Roles {
		fmt.Println(v)
	}

	// insert, err := db.Query(
	// 	"IF NOT EXISTS (SELECT 1 FROM ChatPing WHERE user_id = ? OR group_id = ?) BEGIN INSERT ChatPing(author,user_id,group_id,is_user) VALUES (?,?,?,?) END;",
	// 	m.Author.ID,
	// 	m.Author.)
	// if err != nil {
	// 	panic(err.Error())
	// }
	// // be careful deferring Queries if you are using transactions
	// defer insert.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if len(m.Content) <= 1 || m.Content[0:2] != "->" {
		return
	}

	command := strings.Split(m.Content, " ")

	switch command[0] {
	case "->jebnijBasemSynu":
		s.ChannelMessageSend(m.ChannelID, "@everyone")
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
	color := color.RGBA{255, 0, 0, 255}
	point := fixed.Point26_6{
		X: fixed.Int26_6(x * 64),
		Y: fixed.Int26_6(y * 64),
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}

func testImage(s *discordgo.Session, m *discordgo.MessageCreate) {
	img := image.NewRGBA(image.Rect(0, 0, 640, 360))
	welcomeText := fmt.Sprintf("Witaj %v", m.Author.Username)
	addLabel(img, 10, 30, welcomeText)
	addLabel(img, 10, 45, "***** ***")

	fileName := fmt.Sprintf("/tmp/discordgo/welcome_%v.png", m.Author.Username)
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}

	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
	f.Close()

	f, err = os.Open(fileName)
	if err != nil {
		panic(err)
	}

	_, err = s.ChannelFileSend(m.ChannelID, "asd.png", f)
	if err != nil {
		panic(err)
	}

	f.Close()

	err = os.Remove(fileName)
	if err != nil {
		panic(err)
	}
}
