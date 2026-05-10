# 📧 Ticketing Mail

### Microservizio Mail — Go

> Microservizio generico per l'invio di email transazionali. Consuma job da una coda Redis e spedisce tramite SMTP, senza conoscere il contenuto delle mail.

---

## 📌 Panoramica

**Ticketing Mail** è un servizio indipendente con un unico scopo: ricevere job dalla coda Redis pubblicati dal backend Laravel e inviare le email. Tutta la logica del contenuto (template, dati, lingua) rimane in Laravel — il microservizio si occupa esclusivamente di spedire.

---

## 🛠️ Stack

| Componente | Versione                 |
| ---------- | ------------------------ |
| Go         | 1.21+                    |
| Redis      | 7.x                      |
| SMTP       | — (Mailtrap in sviluppo) |

---

## ⚙️ Installazione

```bash
# 1. Clona la repository
git clone https://github.com/gambadaniele1-hue/ticketing-mail.git
cd ticketing-mail

# 2. Installa le dipendenze
go mod download

# 3. Copia il file di configurazione
cp .env.example .env

# 4. Configura le variabili nel file .env
REDIS_HOST=127.0.0.1
REDIS_PORT=6379
REDIS_QUEUE=mail:queue
REDIS_DLQ=mail:dlq

SMTP_HOST=smtp.mailtrap.io
SMTP_PORT=587
SMTP_USER=
SMTP_PASSWORD=
SMTP_FROM=noreply@ticketing.com

# 5. Avvia il servizio
go run main.go
```

---

## 📨 Struttura del Job

Il microservizio consuma job in formato JSON dalla coda Redis:

```json
{
  "to": "utente@email.com",
  "subject": "Oggetto della mail",
  "html": "<html>...contenuto HTML renderizzato da Laravel...</html>",
  "text": "Versione plain text fallback"
}
```

---

## 🔄 Flusso di elaborazione

```
Laravel (pubblica job)
        │
        ▼
Redis — mail:queue
        │
        ▼
Go — consuma job
        │
        ├── OK  → job completato
        │
        └── FAIL → rimette in coda
                    (max 3 tentativi)
                         │
                         └── se >= 3 → mail:dlq
                                       (dead letter queue)
```

---

## 🧪 Testing

Il progetto è sviluppato con approccio **TDD (Test Driven Development)**.

```bash
# Esegui tutti i test
go test ./...

# Esegui i test con output dettagliato
go test -v ./...
```

---

## 📦 Repository collegate

| Repository                                                              | Descrizione                               |
| ----------------------------------------------------------------------- | ----------------------------------------- |
| [`ticketing-api`](https://github.com/gambadaniele1-hue/ticketing-api)   | Backend Laravel — pubblica i job su Redis |
| [`ticketing-app`](https://github.com/gambadaniele1-hue/ticketing-app)   | Frontend Lovable                          |
| [`ticketing-docs`](https://github.com/gambadaniele1-hue/ticketing-docs) | Documentazione completa                   |

---

## 👤 Autore

Progetto realizzato come elaborato di quinta superiore — Informatica.

---

_Mail Service v1.0 — Go_
