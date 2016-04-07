package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	plurgo "github.com/clsung/plurgo/plurkgo"
)

var credPath = flag.String("config", "config.json", "Path to configuration file containing the application's credentials.")

func main() {
	flag.Parse()
	plurkOAuth, err := plurgo.ReadCredentials(*credPath)
	if err != nil {
		log.Fatalf("Error reading credential, %v", err)
	}
	accessToken, authorized, err := plurgo.GetAccessToken(plurkOAuth)

	if authorized {
		bytes, err := json.MarshalIndent(plurkOAuth, "", "  ")
		if err != nil {
			log.Fatalf("failed to store credential: %v", err)
		}
		err = ioutil.WriteFile(*credPath, bytes, 0700)
		if err != nil {
			log.Fatal("failed to write credential: %v", err)
		}
	}
	result, err := plurgo.CallAPI(accessToken, "/APP/Profile/getOwnProfile", map[string]string{})
	if err != nil {
		log.Fatalf("failed: %v", err)
	}
	fmt.Println(string(result))
	var data = map[string]string{}
	data["content"] = "Test plurkAdd from plurkgo"
	data["qualifier"] = "shares"
	data["lang"] = "ja"
	result, err = plurgo.CallAPI(accessToken, "/APP/Timeline/plurkAdd", data)
	if err != nil {
		log.Fatalf("failed: %v", err)
	}
	fmt.Println(string(result))
}
