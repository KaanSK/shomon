package thehive

type Alert struct {
	Id           string       `json:"_id,omitempty"`
	Type         string       `json:"type"`
	Source       string       `json:"source"`
	SourceRef    string       `json:"sourceRef"`
	Title        string       `json:"title"`
	Description  string       `json:"description"`
	ExternalLink string       `json:"externalLink"`
	Tags         []string     `json:"tags"`
	Observables  []Observable `json:"observables"`
}

type Observable struct {
	DataType string `json:"dataType"`
	Data     string `json:"data"`
}
