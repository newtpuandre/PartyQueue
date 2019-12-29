package website

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func QueueHandler() {

	for {
		if isplaying {
			playing := spotifyGetCurrent()
			currPlaying = playing
			if len(queue) >= 1 && playing.ProgressMS+3000 >= playing.Item.DurationMS {
				nextSong()
			}
			time.Sleep(2 * time.Second)
		}
		if finalToken.Token != "" {
			time := strings.Split(time.Until(finalToken.Expiry).String(), "m")
			timeInt, err := strconv.Atoi(time[0])
			if err == nil {
				if timeInt < 10 {
					updateSpotifyToken(finalToken)
				}
			} else {
				updateSpotifyToken(finalToken)
			}
		}
	}
}

func nextSong() {
	x, a := queue[0], queue[1:]
	queue = a
	spotifyPlaySong(x.SongId)
	currPlaying = spotifyGetCurrent()
	time.Sleep(200)
}

func spotifyPlaySong(songId string) {
	fmt.Println(songId + " Song id !")
	var song PlaySongJson
	song.uris = append(song.uris, songId)

	var bearer = "Bearer " + finalToken.Token
	searchUrl := "https://api.spotify.com/v1/me/player/play"

	body := fmt.Sprintf("{\"uris\":[\"spotify:track:" + songId + "\"]}")
	fmt.Println(body)

	req, err := http.NewRequest("PUT", searchUrl, strings.NewReader(body))
	if err != nil {
		log.Println(err)
	}

	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	log.Println(bodyString)
}

func spotifyGetCurrent() CurrentlyPlaying {
	var bearer = "Bearer " + finalToken.Token
	currentPlayingurl := "https://api.spotify.com/v1/me/player"

	req, err := http.NewRequest("GET", currentPlayingurl, nil)
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}

	var currentlyPlaying CurrentlyPlaying
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&currentlyPlaying)
	if err != nil {
		panic(err)
	}

	//fmt.Println(currentlyPlaying)
	return currentlyPlaying

}

func resolveSpotifyIds(id string) SongLookup {
	var bearer = "Bearer " + finalToken.Token
	currentPlayingurl := "https://api.spotify.com/v1/tracks/"

	req, err := http.NewRequest("GET", currentPlayingurl+id, nil)
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}

	var resolvedSong SongLookup
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&resolvedSong)
	if err != nil {
		panic(err)
	}

	//fmt.Println(currentlyPlaying)
	return resolvedSong
}

func convertResults(result SearchResult) []ConvertedResults {
	var results []ConvertedResults

	for _, item := range result.Tracks.Items {
		var newResult ConvertedResults
		newResult.Songname = item.Name
		newResult.Songid = item.ID
		newResult.SongLength = item.DurationMS / 1000
		newResult.AlbumArt = item.Album.Images[0].URL
		newResult.Artist = item.Artists[0].Name
		results = append(results, newResult)
	}

	return results

}

func convertSongResult(result SongLookup) ConvertedResults {

	var newResult ConvertedResults
	newResult.Songname = result.Name
	newResult.Songid = result.ID
	newResult.SongLength = result.DurationMS / 1000
	newResult.AlbumArt = result.Album.Images[0].URL
	newResult.Artist = result.Artists[0].Name

	return newResult

}

func updateSpotifyToken(token Token) {
	apiUrl := "https://accounts.spotify.com"
	resource := "/api/token"
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", token.Refresh)

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String()

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	r.SetBasicAuth(config.Clientid, config.Clientsecret)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, _ := client.Do(r)

	var spot SpotifyTokenRefresh
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&spot)
	if err != nil {
		panic(err)
	}

	token.Token = spot.AccessToken
	token.Expiry = time.Now().Local().Add(time.Hour)
}

func spotifySearch(search string) []ConvertedResults {
	var bearer = "Bearer " + finalToken.Token
	searchUrl := "https://api.spotify.com/v1/search?q="
	typeUrl := "&type=track%2Cartist"

	req, err := http.NewRequest("GET", searchUrl+url.QueryEscape(search)+typeUrl, nil)
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}

	var results SearchResult
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&results)
	if err != nil {
		panic(err)
	}

	return convertResults(results)

}
