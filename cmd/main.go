package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/redis/go-redis/v9"

	"github.com/gambadaniele1-hue/ticketing-mail/internal/config"
	"github.com/gambadaniele1-hue/ticketing-mail/internal/mail"
	"github.com/gambadaniele1-hue/ticketing-mail/internal/queue"
	"github.com/gambadaniele1-hue/ticketing-mail/internal/worker"
)

func main() {
	log.Println("[main] avvio ticketing-mail...")

	// 1. Carica configurazione
	cfg := config.Load()

	// 2. Connessione Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("[main] impossibile connettersi a Redis: %v", err)
	}
	log.Println("[main] connesso a Redis:", cfg.RedisAddr)

	// 3. Inizializza i pezzi
	redisQueuer := queue.NewRedisQueuer(redisClient)
	smtpSender := mail.NewSMTPSender(mail.SMTPConfig{
		Host:     cfg.SMTPHost,
		Port:     cfg.SMTPPort,
		Username: cfg.SMTPUsername,
		Password: cfg.SMTPPassword,
		From:     cfg.SMTPFrom,
	})
	w := worker.New(smtpSender, redisQueuer)

	// 4. Avvia il loop in una goroutine
	go startWorkerLoop(ctx, redisQueuer, w)

	// 5. Aspetta il segnale di shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("[main] shutdown in corso...")
	cancel()
	log.Println("[main] shutdown completato")
}

func startWorkerLoop(ctx context.Context, redisQueuer *queue.RedisQueuer, w *worker.Worker) {
	log.Printf("[main] in ascolto sulla coda %s...", queue.MainQueue)
	for {
		job, err := redisQueuer.Pop(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Println("[main] shutdown completato")
				return
			}
			log.Printf("[main] errore lettura coda: %v", err)
			continue
		}
		if job == nil {
			continue
		}

		log.Printf("[main] job ricevuto per: %s", job.To)
		if err := w.Process(*job); err != nil {
			log.Printf("[main] errore processo: %v", err)
		}
	}
}
