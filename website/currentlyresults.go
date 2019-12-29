package website

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
