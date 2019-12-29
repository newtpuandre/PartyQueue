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
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

var (
	loggedIn    = false
	urlLogin    string
	auth        spotify.Authenticator
	ch          = make(chan *spotify.Client)
	finalToken  Token
	queue       []QueueSong
	isplaying   = true
	currPlaying CurrentlyPlaying
)

type searchForm struct {
	SearchString string
}

//Token stores the spotify api token
type Token struct {
	Token   string
	Expiry  time.Time
	Refresh string
}

//SpotifyTokenRefresh are used to refresh the spotify token.
type SpotifyTokenRefresh struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	Scope       string `json:"scope"`
}

//SetupAuth configures the spotify lib
func SetupAuth() {
	auth = spotify.NewAuthenticator(config.Callbackuri, spotify.ScopeUserReadCurrentlyPlaying, spotify.ScopeUserReadPlaybackState, spotify.ScopeUserModifyPlaybackState)
	auth.SetAuthInfo(config.Clientid, config.Clientsecret)
}

//404 Route
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "website/templates/404.html")
}

func homeHandle(w http.ResponseWriter, r *http.Request) {
	var data WebsiteStruct
	data.IsLoggedIn = loggedIn
	keys, ok := r.URL.Query()["search"]

	if !ok || len(keys[0]) < 1 { //Viewing home
		tmpl, err := template.ParseFiles("website/templates/home.html")
		if err != nil {
			log.Println(err)
		}
		if isplaying {
			data.IsPlaying = "true"
		} else {
			data.IsPlaying = ""
		}
		tmpl.Execute(w, data)
		return
	} else { //Searching
		key := keys[0]
		data.SearchResults = spotifySearch(key)
		if isplaying {
			data.IsPlaying = "true"
		} else {
			data.IsPlaying = ""
		}

	}

	tmpl := template.Must(template.ParseFiles("website/templates/home.html"))
	tmpl.Execute(w, data)
}

func callback(w http.ResponseWriter, r *http.Request) {
	var tok *oauth2.Token
	var err error

	tok, err = auth.Token("abc123", r)
	if err != nil {
		fmt.Println(err)
	}
	var newToken Token
	newToken.Token = tok.AccessToken
	newToken.Expiry = tok.Expiry
	newToken.Refresh = tok.RefreshToken

	finalToken = newToken

	w.Header().Set("Content-Type", "application/json")

	loggedIn = true
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func login(w http.ResponseWriter, r *http.Request) {
	urlLogin = auth.AuthURL("abc123")
	http.Redirect(w, r, urlLogin, http.StatusSeeOther)
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

func queueadd(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var newQueue QueueSong

	newQueue.SongId = vars["id"]

	queue = append(queue, newQueue)
	fmt.Println(queue)

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func pauseQueue(w http.ResponseWriter, r *http.Request) {
	isplaying = !isplaying
	fmt.Println("Queue state", isplaying)

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func queueRemove(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var slicevalue int

	for i, value := range queue {
		if value.SongId == vars["id"] {
			slicevalue = i
			break
		}
	}

	copy(queue[slicevalue:], queue[slicevalue+1:])
	queue = queue[:len(queue)-1]

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func queueSkip(w http.ResponseWriter, r *http.Request) {
	nextSong()
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func queueDown(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var slicevalue int

	for i, value := range queue {
		if value.SongId == vars["id"] {
			slicevalue = i
			break
		}
	}

	tempStorage := queue[slicevalue+1]
	queue[slicevalue+1] = queue[slicevalue]
	queue[slicevalue] = tempStorage

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func queueUp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var slicevalue int

	for i, value := range queue {
		if value.SongId == vars["id"] {
			slicevalue = i
			break
		}
	}

	tempStorage := queue[slicevalue-1]
	queue[slicevalue-1] = queue[slicevalue]
	queue[slicevalue] = tempStorage

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func showqueue(w http.ResponseWriter, r *http.Request) {

	var queueData QueueStruct
	if len(queue) >= 1 {
		for _, value := range queue {
			result := resolveSpotifyIds(value.SongId)
			convertResult := convertSongResult(result)
			queueData.Songs = append(queueData.Songs, convertResult)
		}
	}

	queueData.Size = len(queueData.Songs) - 1
	if finalToken.Token != "" {
		queueData.CurrPlay = currPlaying
		queueData.CurrPlayImage = currPlaying.Item.Album.Images[0].URL
		queueData.CurrPlayArtist = currPlaying.Item.Artists[0].Name
	}

	tmpl := template.Must(template.ParseFiles("website/templates/queue.html"))
	tmpl.Execute(w, queueData)
}

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
