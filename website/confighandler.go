package website

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Callbackuri  string
	Clientid     string
	Clientsecret string
}

var config Config

func ConfigInit() {
	//Check if file exists anton_config.json
	if _, err := os.Stat("./partyqueue.json"); os.IsNotExist(err) {
		log.Println("partyqueue.json did not exist and have been created. Please fill in the fields")
		file, _ := json.MarshalIndent(config, "", " ")

		_ = ioutil.WriteFile("./partyqueue.json", file, 0644)
		os.Exit(1)
	}

	//Load config.
	file, err := os.Open("./partyqueue.json")
	if err != nil {
		log.Println(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Println(err)
	}

	log.Println(config)
}
