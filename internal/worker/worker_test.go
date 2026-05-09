package worker_test

import (
	"errors"
	"testing"

	"github.com/gambadaniele1-hue/ticketing-mail/internal/queue"
	"github.com/gambadaniele1-hue/ticketing-mail/internal/worker"
)

type FakeSender struct {
	SimulaErrore bool
}

func (f *FakeSender) Send(job queue.MailJob) error {
	if f.SimulaErrore {
		return errors.New("errore SMTP simulato")
	}

	return nil
}

type FakeQueuer struct {
	JobInCoda []queue.MailJob
	JobInDLQ  []queue.MailJob
}

func (f *FakeQueuer) Requeue(job queue.MailJob) error {
	f.JobInCoda = append(f.JobInCoda, job)
	return nil
}

func (f *FakeQueuer) MoveToDLQ(job queue.MailJob) error {
	f.JobInDLQ = append(f.JobInDLQ, job)
	return nil
}

// ── Test ──────────────────────────────────────────────────────────────────────

func TestInvioRiuscito(t *testing.T) {
	sender := &FakeSender{SimulaErrore: false}
	queuer := &FakeQueuer{}

	w := worker.New(sender, queuer)

	job := queue.MailJob{
		To:      "utente@email.com",
		Subject: "Test",
		HTML:    "<p>Ciao</p>",
		Text:    "Ciao",
	}

	err := w.Process(job)

	if err != nil {
		t.Errorf("errore inattesto %v", err)
	}

	if len(queuer.JobInCoda) != 0 {
		t.Errorf("non dovevano esserci requeue, ce ne sono %d", len(queuer.JobInCoda))
	}
	if len(queuer.JobInDLQ) != 0 {
		t.Errorf("non dovevano esserci job in DLQ, ce ne sono %d", len(queuer.JobInDLQ))
	}
}

func TestRetryDopoErrore(t *testing.T) {
	sender := &FakeSender{SimulaErrore: true}
	queuer := &FakeQueuer{}
	w := worker.New(sender, queuer)

	job := queue.MailJob{
		To:       "utente@email.com",
		Subject:  "Test",
		HTML:     "<p>Ciao</p>",
		Text:     "Ciao",
		Attempts: 0,
	}

	err := w.Process(job)

	if err == nil {
		t.Error("avrebbe dovuto restituire un errore")
	}
	if len(queuer.JobInCoda) != 1 {
		t.Errorf("atteso 1 requeue, ottenuto %d", len(queuer.JobInCoda))
	}
	if queuer.JobInCoda[0].Attempts != 1 {
		t.Errorf("Attempts dovrebbe essere 1, è %d", queuer.JobInCoda[0].Attempts)
	}
	if len(queuer.JobInDLQ) != 0 {
		t.Errorf("non doveva andare in DLQ, ce ne sono %d", len(queuer.JobInDLQ))
	}
}

func TestDLQDopoTreErrori(t *testing.T) {
	sender := &FakeSender{SimulaErrore: true}
	queuer := &FakeQueuer{}
	w := worker.New(sender, queuer)

	job := queue.MailJob{
		To:       "utente@email.com",
		Subject:  "Test",
		HTML:     "<p>Ciao</p>",
		Text:     "Ciao",
		Attempts: 2,
	}

	err := w.Process(job)

	if err == nil {
		t.Error("avrebbe dovuto restituire un errore")
	}
	if len(queuer.JobInDLQ) != 1 {
		t.Errorf("atteso 1 job in DLQ, ottenuto %d", len(queuer.JobInDLQ))
	}
	if len(queuer.JobInCoda) != 0 {
		t.Errorf("non doveva fare requeue, ne ha fatti %d", len(queuer.JobInCoda))
	}
}

func TestJobNonValido(t *testing.T) {
	sender := &FakeSender{}
	queuer := &FakeQueuer{}
	w := worker.New(sender, queuer)

	job := queue.MailJob{}

	err := w.Process(job)

	if err == nil {
		t.Error("un job vuoto dovrebbe restituire errore")
	}
}
