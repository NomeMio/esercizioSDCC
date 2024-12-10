# esercizioSDCC

Questo progetto implementa un sistema distribuito per eseguire operazioni MapReduce utilizzando Docker e Go.

## Requisiti

- **Docker** e **Docker Compose** installati.
- **Go** (versione 1.19 o successiva) per modifiche e test locali.

---

## Come avviare il sistema

### 1. Eseguire con Docker Compose

Nella directory principale del progetto, eseguire:

```bash
docker-compose up --build
```

---

## Modifica del numero di worker

### 1. Aggiornare il comando del master

Nel file `docker-compose.yml`, modificare la parte `command` della sezione `master` per aggiungere o rimuovere indirizzi worker.

**Esempio:**

```yaml
command: >
  ./master
  -a
  worker1:8000,worker2:8000,worker3:8000
  -n
  1000
  -m
  40
```

Per aggiungere un worker, includere un nuovo indirizzo come `worker4:8000`.

### 2. Aggiungere una nuova sezione per il worker

Aggiungere la configurazione del worker nella sezione `services` di `docker-compose.yml`.

**Esempio:**

```yaml
worker4:
  build:
    context: ./worker
  networks:
    - internal_network
  ports:
    - "8004:8000"
  environment:
    - WORKER_NAME=worker4
```

---

## Arrestare il sistema

Per fermare e rimuovere tutti i container, eseguire:

```bash
docker-compose down
```

---

## Ispezione dei dati generati

### 1. Ispezionare il volume Docker

Per visualizzare i dettagli del volume:

```bash
docker volume inspect eserciziosdcc_app-data
```

### 2. Visualizzare i file nel volume

- Recuperare il percorso `Mountpoint` dall'output precedente.
- Eseguire:

```bash
sudo ls -l <Mountpoint>
```

---

## Struttura del codice

### File principali

- **`master/main.go`**: Contiene la logica principale per orchestrare i worker.
- **`master/utilis`**: Libreria di funzioni utili per la gestione di file e calcoli distribuiti.
- **`worker/`**: Codice e configurazione per i container dei worker.

---

## Parametri principali

- `-a`: Indirizzi dei worker, separati da virgola (es. `worker1:8000,worker2:8000`).
- `-n`: Numero di interi generati casualmente.
- `-m`: Valore massimo degli interi generati.

---

## Descrizione del flusso

### Master

1. Genera un file di numeri casuali.
2. Suddivide i dati e li assegna ai worker.
3. Raccoglie i risultati delle fasi Map e Reduce dai worker.

### Worker

1. Esegue la fase **Map** per identificare il valore minimo e massimo.
2. Esegue la fase **Reduce** per consolidare i dati.

I risultati finali vengono salvati in un file specificato in `configuration.FILE_NAME_REPLAY`.

---

## Debug e log

Per abilitare il logging dettagliato, assicurarsi che `log.SetFlags` sia configurato correttamente nel codice del master.

