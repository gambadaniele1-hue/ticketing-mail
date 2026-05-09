package worker

import (
	"fmt"
	"log"

	"github.com/gambadaniele1-hue/ticketing-mail/internal/queue"
)

const MaxAttempts = 3

type Sender interface {
	Send(job queue.MailJob) error
}

type Queuer interface {
	Requeue(job queue.MailJob) error
	MoveToDLQ(job queue.MailJob) error
}

type Worker struct {
	sender Sender
	queuer Queuer
}

func New(sender Sender, queuer Queuer) *Worker {
	return &Worker{
		sender: sender,
		queuer: queuer,
	}
}

func (w *Worker) Process(job queue.MailJob) error {
	if !job.IsValid() {
		return fmt.Errorf("job non valido: to=%q subject=%q", job.To, job.Subject)
	}

	err := w.sender.Send(job)
	if err == nil {
		log.Printf("[worker] email inviata a %s", job.To)
		return nil
	}

	job.Attempts++
	log.Printf("[worker] tentativo %d/%d fallito per %s", job.Attempts, MaxAttempts, job.To)

	if job.Attempts >= MaxAttempts {
		log.Printf("[worker] max tentativi raggiunti, sposto in DLQ: %s", job.To)
		if err := w.queuer.MoveToDLQ(job); err != nil {
			return fmt.Errorf("errore DLQ: %w", err)
		}
		return fmt.Errorf("job in DLQ dopo %d tentativi", MaxAttempts)
	}

	if err := w.queuer.Requeue(job); err != nil {
		return fmt.Errorf("errore requeue: %w", err)
	}

	return fmt.Errorf("invio fallito, job rimesso in coda (tentativo %d)", job.Attempts)
}
