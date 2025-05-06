# Animesaturn Downloader in Go
Questa utility permette di scaricare anime dal famoso sito Animesaturn e salvarli in formato .mp4 sul computer. Contiene una versione CLI e una con interfaccia grafica.

## Build
Per creare l'eseguibile, dopo aver installato correttamente go, eseguire nel terminale il seguente comando:
- Per Linux/MacOS: 
```console
foo@bar:~/animesaturnDownloaderGo$ ./scripts/build.sh
```
- Per Windows:
```powershell
PS ~/animesaturnDownloaderGo> .\scripts\build.ps1
```

## Utilizzo
### CLI
Per ottenere informazioni a proposito della CLI si può usare il seguente comando:
- Per Linux/MacOS:
```console
foo@bar:~/animesaturnDownloaderGo$ ./bin/animesaturn-downloader -h
```
- Per Windows:
```powershell
PS ~/animesaturnDownloaderGo> .\bin\animesaturn-downloader.exe -h
```
Che produrrà un output simile a questo:
```console
AnimesaturnDownloader è una utility per scaricare gli anime dal sito AnimeSaturn.
Scritto in Go da Marco Simone.


Questa schermata di aiuto è divisa in più parti, usa "animesaturn-downloader <sottocomando> -h" per vedere la schermata di aiuto per il sottocomando specifico.

I sottocomandi disponibili sono:
  download              Scarica gli episodi di un anime
  search                Cerca un anime per nome

Utilizzo: animesaturn-downloader <sottocomando> [opzioni]

Flag globali:
  -h, --help            stampa le informazioni di aiuto
  -v, --verbose         stampa altre informazioni di debug
  -V, --version         stampa la versione del programma e termina il programma
```

Un esempio di comando può essere:
- Per Linux/MacOS:
```console
foo@bar:~/animesaturnDownloaderGo$ ./bin/animesaturn-downloader -u https://your-url-here/anime -f 1 -l 12 -d ./my-anime -n MyAnime_ -w 3
```
- Per Windows:
```powershell
PS ~/animesaturnDownloaderGo> .\bin\animesaturn-downloader.exe -u https://your-url-here/anime -f 1 -l 12 -d ./my-anime -n MyAnime_ -w 3
```
Questo comando invoca l'eseguibile con i seguenti parametri:
- url: https[]()://your-url-here/anime
- primo episodio: 1
- ultimo episodio: 12
- cartella output: ./my-anime
- nome dei file: MyAnime\_
- worker da usare: 3
  
NB: Il prgoramma aggiunge "i.mp4" alla fine di ogni file con i uguale al numero dell'episodio scaricato.

### Applicazione Grafica
L'eseguibile per l'applicazione grafica è reperibile in ./bin/animesaturn-downloader-gui.exe