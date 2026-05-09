package queue_test

import (
	"encoding/json"
	"testing"

	"github.com/gambadaniele1-hue/ticketing-mail/internal/queue"
)

func TestMailJobSerialization(t *testing.T) {
	rawjson := `{
		"to": "gamba@gmail.com",
		"subject": "oggeto della mail",
		"html": "<h1>Titolo</h1>",
		"text": "testo in chiaro"
	}`
	var job queue.MailJob

	err := json.Unmarshal([]byte(rawjson), &job)

	if err != nil {
		t.Fatalf("errore inatteso: %v", err)
	}

	if job.To != "gamba@gmail.com" {
		t.Errorf("job.To: atteso %q, ricevuto %q", "gamba@gmail.com", job.To)
	}

	if job.Subject != "oggeto della mail" {
		t.Errorf("job.Subjectg: atteso %q, ottenuto %q", "oggeto della mail", job.Subject)
	}

	if job.HTML != "<h1>Titolo</h1>" {
		t.Errorf("job.HTML: atteso %q, ottenuto %q", "<h1>Titolo</h1>", job.HTML)
	}

	if job.Text != "testo in chiaro" {
		t.Errorf("job.Text: atteso %q, ottenuto %q", "testo in chiaro", job.Text)
	}
}

func TestMailJobValidation(t *testing.T) {
	jobVouto := queue.MailJob{}

	if jobVouto.IsValid() {
		t.Error("un job vuoto non dovrebbe essere valido")
	}

	jobValido := queue.MailJob{
		To:      "gamba@email",
		Subject: "Oggetto",
		HTML:    "<h1>Titolo</h1>",
		Text:    "testo in chiaro",
	}

	if !jobValido.IsValid() {
		t.Error("un job completo dovrebbe essere valido")
	}
}
