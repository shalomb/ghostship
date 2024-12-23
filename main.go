package main

import (
	gitstatus "github.com/shalomb/ghostship/gitstatus"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Printf("gitstatus: %+v", gitstatus.Status())
}
