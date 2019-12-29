package main

import (
	website "PARTYQUEUE/website"
	"sync"
)

var (
	//WebPort is the port the webserver listens to
	WebPort = "80"
)

func main() {
	//Reads config file
	website.ConfigInit()

	//Initiates the spotify auth
	website.SetupAuth()
	wg := &sync.WaitGroup{}

	//Start webserver
	wg.Add(1)
	go func() {
		website.StartServer(WebPort)
	}()

	//Start queue handling
	wg.Add(1)
	go func() {
		website.QueueHandler()
	}()

	//Wait for the group to stop
	wg.Wait()
}
