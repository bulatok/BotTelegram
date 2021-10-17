package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
)

const tok = "1819449359:AAG2jsR1H1tgxeCcooa-OfNd6zmxmqQQuoQ"
type Origin struct{
	Name string
	url string
}
type Hero struct{
	//Id int `json:"id"`
	Name string `json:"name"`
	ImageURL string `json:"image"`
	Status string `json:"status"`
	Gender string `json:"gender"`
	Origin Origin `json:"origin"`
	//Episods []string `json:"episods"`
}


func getHero(heroURL string) *Hero{
	clint := &http.Client{
		Timeout: time.Second * 100,
	}
	resp, err := clint.Get(heroURL)
	if err != nil{
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil{
		log.Fatal(err)
	}

	newHero := &Hero{}
	//var newHero interface{}
	err = json.Unmarshal(body, &newHero)
	if err != nil{
		log.Fatal(err)
	}
	return newHero
}
func PrettyHero(h *Hero) string{
	res := fmt.Sprintf("name: %s\nstatus: %s\ngender: %s\nstatus: %s\nfrom: %s\n",
		h.Name, h.Status, h.Gender, h.Status, h.Origin.Name)
	return res
}
func getNum(s string) (int, string){
	for _, v := range s{
		if v < '0' || v > '9'{
			return -1, "ты еблан?"
		}
	}
	num, _ := strconv.Atoi(s)
	if num < 1 || num > 672{
		return -1, "ты еблан?"
	}
	return num, "OK"
}

var urlBase = `https://rickandmortyapi.com/api/character/`
func main(){
	// for metrics
	f, err := os.OpenFile("metrics", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil{
		log.Println(err)
	}
	defer f.Close()

	// bot implementation
	b, err := tb.NewBot(tb.Settings{
		Token: tok,
		Poller: &tb.LongPoller{Timeout: 10*time.Second},
	})
	if err != nil{
		log.Fatal(err)
	}


	b.Handle("/start", func(m *tb.Message){
		//btnRand := tb.ReplyMarkup{}
		b.Send(m.Sender,`Please enter number (1-672) to see a hero from "rick and morty"`)
	})

	b.Handle(tb.OnText, func(m *tb.Message){
		var s string
		//log.Println("Handle1")
		s = fmt.Sprintf(`"%s" %s `, m.Text, m.Sender.FirstName)

		n, err := getNum(m.Text)
		if n == -1{
			s += fmt.Sprintln(m.Text)
			b.Send(m.Sender, err)
		}else{
			rick := getHero(urlBase + m.Text)
			ph := &tb.Photo{File:tb.FromURL(rick.ImageURL)}	// getting a photo of hero
			s += fmt.Sprintln(time.Now().Format("02-01-2006 15:04:05"))
			b.Send(m.Sender, ph)
			b.Send(m.Sender, PrettyHero(rick))
		}
		fmt.Println(s)
		f.WriteString(s)
	})
	b.Start()
}