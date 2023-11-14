package common

type WebFingerResponse struct {
	Subject string           `json:"subject"`
	Links   []WebFingerLinks `json:"links"`
}

type WebFingerLinks struct {
	Rel  string `json:"rel"`
	Type string `json:"type"`
	Href string `json:"href"`
}
