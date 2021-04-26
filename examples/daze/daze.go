// Daze could join an offline-mode server as client.
// Just standing there and do nothing. Automatically reborn after five seconds of death.
//
// BUG(Tnze): Kick by Disconnect: Time Out
package main

import (
	"errors"
	"flag"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/VacPlusPlus/go-mc/bot"
	"github.com/VacPlusPlus/go-mc/bot/basic"
	"github.com/VacPlusPlus/go-mc/chat"
	_ "github.com/VacPlusPlus/go-mc/data/lang/zh-cn"
)

var address = flag.String("address", "127.0.0.1", "The server address")
var client *bot.Client
var player *basic.Player

func main() {
	flag.Parse()
	client = bot.NewClient()
	client.Auth.Name = "Daze"
	player = basic.NewPlayer(client, basic.DefaultSettings)
	basic.EventsListener{
		GameStart:  onGameStart,
		ChatMsg:    onChatMsg,
		Disconnect: onDisconnect,
		Death:      onDeath,
	}.Attach(client)

	//Login
	err := client.JoinServer(*address)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Login success")

	//JoinGame
	for {
		if err = client.HandleGame(); err == nil {
			panic("HandleGame never return nil")
		}

		if err2 := new(bot.PacketHandlerError); errors.As(err, err2) {
			if err := new(DisconnectErr); errors.As(err2, err) {
				log.Print("Disconnect: ", err.Reason)
				return
			} else {
				// print and ignore the error
				log.Print(err2)
			}
		} else {
			log.Fatal(err)
		}
	}
}

func onDeath() error {
	log.Println("Died and Respawned")
	// If we exclude Respawn(...) then the player won't press the "Respawn" button upon death
	go func() {
		time.Sleep(time.Second * 5)
		err := player.Respawn()
		if err != nil {
			log.Print(err)
		}
	}()
	return nil
}

func onGameStart() error {
	log.Println("Game start")
	return nil //if err isn't nil, HandleGame() will return it.
}

func onChatMsg(c chat.Message, _ byte, _ uuid.UUID) error {
	log.Println("Chat:", c.ClearString()) // output chat message without any format code (like color or bold)
	return nil
}

type DisconnectErr struct {
	Reason chat.Message
}

func (d DisconnectErr) Error() string {
	return "disconnect: " + d.Reason.String()
}

func onDisconnect(reason chat.Message) error {
	// return a error value so that we can stop main loop
	return DisconnectErr{Reason: reason}
}
