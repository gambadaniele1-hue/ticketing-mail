package queue

type MailJob struct {
	To       string `json:"to"`
	Subject  string `json:"subject"`
	HTML     string `json:"html"`
	Text     string `json:"text"`
	Attempts int    `json:"attempts"`
}

func (j *MailJob) IsValid() bool {
	return j.To != "" && j.Subject != "" && (j.HTML != "" || j.Text != "")
}
