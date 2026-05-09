package queue

type Mailjob struct {
	To       string `json:"to"`
	Subject  string `json:"subject"`
	HTML     string `json:"html"`
	Text     string `json:"text"`
	Attempts int    `json:"attempts"`
}

func (j *Mailjob) IsValid() bool {
	return j.To != "" && j.Subject != "" && (j.HTML != "" || j.Text != "")
}
