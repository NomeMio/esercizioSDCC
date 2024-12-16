# esercizioSDCC

Questo progetto implementa un sistema distribuito per eseguire operazioni MapReduce utilizzando Docker e Go.

## Requisiti

- **Docker** e **Docker Compose** installati.
- **Go** versione 1.23.4 o successiva.

---

## Come avviare il sistema

### 1. Eseguire con Docker Compose

Nella directory principale del progetto, eseguire:

```bash
docker-compose up --build
```

---

## Modificare il numero di workers 


### 1. Aggiungere/Rimuovere il servizio relativo al worker

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


### 2. Aggiornare il comando del master

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

In seguito aggiungere la dipendeza da parte del master del nuovo worker.

```yaml
    depends_on:
      - worker1
      - worker2
      - worker3
      - worker4
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

## Parametri principali

### Master flags
- `-a`: Indirizzi dei worker, separati da virgola (es. `worker1:8000,worker2:8000`).
- `-n`: Numero di interi generati casualmente.
- `-m`: Valore massimo degli interi generati.
### Worker flags
- `-p`: Indica la porta su cui ascoltare.

---

## Descrizione del flusso

### Master

1. Genera un file di numeri casuali.
2. Suddivide i dati in shard e li assegna ai worker.
3. Sincronizza le fasi di map e reduce in modo che la seconda inizi solo dopo che la prima fase sia conclusa per tutti i worker.

### Worker

1. Esegue la fase **Map** e aspetta lo shuffle delle chiavi.
2. Esegue la fase **Reduce** per consolidare i dati.

I risultati finali vengono salvati in un file specificato in `configuration.FILE_NAME_REPLAY`.

---

