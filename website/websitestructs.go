package website

type WebsiteStruct struct {
	IsPlaying     string
	IsLoggedIn    bool
	SearchResults []ConvertedResults
}

type QueueStruct struct {
	Size           int
	CurrPlay       CurrentlyPlaying
	CurrPlayImage  string
	CurrPlayArtist string
	Songs          []ConvertedResults
}
