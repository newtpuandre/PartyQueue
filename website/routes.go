package website

import (
	"fmt"
	"log"
	"net/http"
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
