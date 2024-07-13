package main

import (
	"essh/cmd/internal/session"
	"log"
)

func main() {
	sessions := []session.Session{
		session.New("Session1", "User1", "host1.example.com", "password1", 22),
		session.New("Session2", "User2", "host2.example.com", "password2", 2222),
		session.New("Session3", "User3", "host3.example.com", "password3", 3333),
	}

	err := session.SaveToFile(sessions)
	if err != nil {
		log.Println(err)
		return
	}
}
