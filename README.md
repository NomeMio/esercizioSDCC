# esercizioSDCC


per runnare docker compose up --build nella cartella esercizioSDCC

per aggiungere o togliere workers bisogna modificare due cose nel file yaml:
- 1 aggiungere/togleire nome:porto nella parte comand di master

- 2 aggiungere tipo (  worker2:
  build:
  context: ./worker
  networks:
  - internal_network
  ports:
  - "8002:8000"
  environment:
  - WORKER_NAME=worker2) nell sezione sotto gli altri workere
