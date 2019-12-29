package website

type SongLookup struct {
	Album            Album        `json:"album"`
	Artists          []Artist     `json:"artists"`
	AvailableMarkets []string     `json:"available_markets"`
	DiscNumber       int64        `json:"disc_number"`
	DurationMS       int64        `json:"duration_ms"`
	Explicit         bool         `json:"explicit"`
	ExternalIDS      ExternalIDS  `json:"external_ids"`
	ExternalUrls     ExternalUrls `json:"external_urls"`
	Href             string       `json:"href"`
	ID               string       `json:"id"`
	IsLocal          bool         `json:"is_local"`
	Name             string       `json:"name"`
	Popularity       int64        `json:"popularity"`
	PreviewURL       string       `json:"preview_url"`
	TrackNumber      int64        `json:"track_number"`
	Type             string       `json:"type"`
	URI              string       `json:"uri"`
}

type SearchResult struct {
	Artists Artists `json:"artists"`
	Tracks  Artists `json:"tracks"`
}

type Album struct {
	AlbumType            string       `json:"album_type"`
	Artists              []Artist     `json:"artists"`
	AvailableMarkets     []string     `json:"available_markets"`
	ExternalUrls         ExternalUrls `json:"external_urls"`
	Href                 string       `json:"href"`
	ID                   string       `json:"id"`
	Images               []Image      `json:"images"`
	Name                 string       `json:"name"`
	ReleaseDate          string       `json:"release_date"`
	ReleaseDatePrecision string       `json:"release_date_precision"`
	TotalTracks          int64        `json:"total_tracks"`
	Type                 string       `json:"type"`
	URI                  string       `json:"uri"`
}

type Artists struct {
	Href     string      `json:"href"`
	Items    []Item      `json:"items"`
	Limit    int64       `json:"limit"`
	Next     *string     `json:"next"`
	Offset   int64       `json:"offset"`
	Previous interface{} `json:"previous"`
	Total    int64       `json:"total"`
}

type Artist struct {
	ExternalUrls ExternalUrls `json:"external_urls"`
	Href         string       `json:"href"`
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Type         string       `json:"type"`
	URI          string       `json:"uri"`
}

type Item struct {
	Album            AlbumClass      `json:"album"`
	Artists          []ArtistElement `json:"artists"`
	AvailableMarkets []string        `json:"available_markets"`
	DiscNumber       int64           `json:"disc_number"`
	DurationMS       int64           `json:"duration_ms"`
	Explicit         bool            `json:"explicit"`
	ExternalIDS      ExternalIDS     `json:"external_ids"`
	ExternalUrls     ExternalUrls    `json:"external_urls"`
	Href             string          `json:"href"`
	ID               string          `json:"id"`
	IsLocal          bool            `json:"is_local"`
	Name             string          `json:"name"`
	Popularity       int64           `json:"popularity"`
	PreviewURL       *string         `json:"preview_url"`
	TrackNumber      int64           `json:"track_number"`
	Type             ItemType        `json:"type"`
	URI              string          `json:"uri"`
}

type AlbumClass struct {
	AlbumType            AlbumTypeEnum        `json:"album_type"`
	Artists              []ArtistElement      `json:"artists"`
	AvailableMarkets     []string             `json:"available_markets"`
	ExternalUrls         ExternalUrls         `json:"external_urls"`
	Href                 string               `json:"href"`
	ID                   string               `json:"id"`
	Images               []Image              `json:"images"`
	Name                 string               `json:"name"`
	ReleaseDate          string               `json:"release_date"`
	ReleaseDatePrecision ReleaseDatePrecision `json:"release_date_precision"`
	TotalTracks          int64                `json:"total_tracks"`
	Type                 AlbumTypeEnum        `json:"type"`
	URI                  string               `json:"uri"`
}

type ArtistElement struct {
	ExternalUrls ExternalUrls `json:"external_urls"`
	Href         string       `json:"href"`
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Type         ArtistType   `json:"type"`
	URI          string       `json:"uri"`
}

type ExternalUrls struct {
	Spotify string `json:"spotify"`
}

type Image struct {
	Height int64  `json:"height"`
	URL    string `json:"url"`
	Width  int64  `json:"width"`
}

type ExternalIDS struct {
	Isrc string `json:"isrc"`
}

type AlbumTypeEnum string

const (
	Compilation AlbumTypeEnum = "compilation"
	Single      AlbumTypeEnum = "single"
)

type ArtistType string

type ReleaseDatePrecision string

const (
	Day  ReleaseDatePrecision = "day"
	Year ReleaseDatePrecision = "year"
)

type ItemType string

const (
	Track ItemType = "track"
)

type PlaySongJson struct {
	uris []string `json:"context_uri"`
}

type CurrentlyPlaying struct {
	Device               Device  `json:"device"`
	ShuffleState         bool    `json:"shuffle_state"`
	RepeatState          string  `json:"repeat_state"`
	Timestamp            int64   `json:"timestamp"`
	Context              Context `json:"context"`
	ProgressMS           int64   `json:"progress_ms"`
	Item                 Item    `json:"item"`
	CurrentlyPlayingType string  `json:"currently_playing_type"`
	Actions              Actions `json:"actions"`
	IsPlaying            bool    `json:"is_playing"`
}

type Actions struct {
	Disallows Disallows `json:"disallows"`
}

type Disallows struct {
	Pausing bool `json:"pausing"`
}

type Context struct {
	ExternalUrls ExternalUrls `json:"external_urls"`
	Href         string       `json:"href"`
	Type         string       `json:"type"`
	URI          string       `json:"uri"`
}

type Device struct {
	ID               string `json:"id"`
	IsActive         bool   `json:"is_active"`
	IsPrivateSession bool   `json:"is_private_session"`
	IsRestricted     bool   `json:"is_restricted"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	VolumePercent    int64  `json:"volume_percent"`
}

type ConvertedResults struct {
	Songname   string
	Songid     string
	SongLength int64
	AlbumArt   string
	Artist     string
}

type QueueSong struct {
	SongId string
}
