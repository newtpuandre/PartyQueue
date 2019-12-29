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
