package main

type shorty struct {
	URL           *string `json:"url,omitempty"`
	ShortCode     *string `json:"shortcode,omitempty"`
	StartDate     *string `json:"startDate,omitempty"`
	RedirectCount *uint64 `json:"redirectCount,omitempty"`
	LastSeenDate  *string `json:"lastSeenDate,omitempty"`
}

type shortyError struct {
	Msg string `json:"message"`
}
