# Animesaturn Downloader in Go

Questa utility permette di scaricare anime dal famoso sito Animesaturn e salvarli in formato .mp4 sul computer.

## Come buildare l'utility

Per creare l'eseguibile, dopo aver installato correttamente go, eseguire nel terminale i seguenti comandi:

```console
foo@bar:~/animesaturnDownloaderGo$ go mod download
foo@bar:~/animesaturnDownloaderGo$ go build
```

L'eseguibile può poi essere eseguito in due modi:

1. Invocando l'eseguibile con inserimento dei parametri manuale

```console

foo@bar:~/animesaturnDownloaderGo$ ./main.exe
Inserisci il link alla pagina dell'anime: https://your-url-here/anime
Inserisci il primo episodio da scaricare: 1
Inserisci l'ultimo episodio da scaricare: 12
Inserisci il percorso dove salvare i file [Vuoto per: "Percorso corrente"]: ./my-anime
Inserisci il nome per i file: MyAnime_
```

2. O, come CLI

```console
foo@bar:~/animesaturnDownloaderGo$ ./main.exe -u https://your-url-here/anime -f 1 -l 12 -d ./my-anime -n MyAnime_
```

Entrambi i comandi hanno questo risultato:
Invocare l'eseguibile con i seguenti parametri:

-   url: https[]()://your-url-here/anime
-   primo episodio: 1
-   ultimo episodio: 12
-   cartella output: ./my-anime
-   nome dei file: MyAnime\_

NB: Il prgoramma aggiunge "i.mp4" alla fine di ogni file con i uguale al numero dell'episodio scaricato.
