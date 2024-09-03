# Animesaturn Downloader in Go

Questa utility permette di scaricare anime dal famoso sito Animesaturn e salvarli in formato .mp4 sul computer. 

## Come buildare l'utility
Per creare l'eseguibile, dopo aver installato correttamente go, eseguire nel terminale i seguenti comandi:
```bash
go mod download
go build
```

L'eseguibile può poi essere eseguito con
```bash
./main.exe
```

## Piani futuri:
Possibilità di inserire i paramentri mentre si runna il comando, Es. 
```bash
./main.exe -u https://animesaturn.me/Sword-art-Online -i 1 -f 12 -d ./SAO -n SwordArtOnline_Ep
```