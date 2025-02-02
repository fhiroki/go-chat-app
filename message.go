package main

import "time"

type message struct {
	Name      string
	Email     string
	Message   string
	AvatarURL string
	When      time.Time
}
