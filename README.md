# Animesaturn Downloader in Go
<p align="center">
    <a href="https://github.com/MrRainbow0704/animesaturnDownloaderGo/releases/latest/"><img alt="Release" src="https://img.shields.io/github/v/release/MrRainbow0704/animesaturnDownloaderGo"></a>
    <a href="https://www.gnu.org/licenses/"><img alt="License" src="https://img.shields.io/github/license/MrRainbow0704/animesaturnDownloaderGo"></a>
    <a href="https://github.com/MrRainbow0704/animesaturnDownloaderGo/releases/latest/"><img alt="GitHub" src="https://img.shields.io/github/downloads/MrRainbow0704/animesaturnDownloaderGo/total"></a>
    <br />
    <a href="https://go.dev/"><img alt="Go" src="https://img.shields.io/badge/Go-00ADD8.svg?&logo=go&logoColor=white"></a>
    <a href="https://www.typescriptlang.org/"><img alt="TypeScript" src="https://img.shields.io/badge/TypeScript-3178C6?logo=typescript&logoColor=fff"></a>
    <a href="https://wails.io/"><img alt="Wails" src="https://img.shields.io/badge/Wails-df0000.svg?logo=wails&logoColor=white"></a>
    <a href="https://svelte.dev/"><img alt="Svelte" src="https://img.shields.io/badge/Svelte-f1413d.svg?logo=svelte&logoColor=white"></a>
</p>

Questa utility permette di scaricare anime dal famoso sito Animesaturn e salvarli in formato .mp4 sul computer. Contiene una versione CLI e una con interfaccia grafica.

## Build

Per creare l'eseguibile, dopo aver installato correttamente go e make, eseguire nel terminale uno dei seguenti comandi:

```console
# Esegue la build di tutte le versioni
foo@bar:~/animesaturnDownloaderGo$ make

# Esegue la build per windows
foo@bar:~/animesaturnDownloaderGo$ make win

# Esegue la build per linux
foo@bar:~/animesaturnDownloaderGo$ make linux

# Esegue la build per mac
foo@bar:~/animesaturnDownloaderGo$ make mac
```
Al termine della build, l'eseguibile sarà disponibile all'interno della cartella `./bin` con nome `animesaturn-downloader-(VERSIONE)-(PIATTAFORMA)`. <br>
**Nota bene: È possibile fare la build per Linux da un dispositivo Windows, a patto che si usi il WSL (Windows Subsystem for Linux).** <br>
**Inoltre non è possibile eseguire la build per mac su hardware non mac. [Prendetevela con Apple.](https://github.com/wailsapp/wails/issues/1041#issuecomment-2492133624)**

## Utilizzo

### CLI

Per ottenere informazioni a proposito della CLI si può usare il seguente comando:

```console
foo@bar:~/animesaturnDownloaderGo$ ./bin/animesaturn-downloader -h
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

```console
foo@bar:~/animesaturnDownloaderGo$ ./bin/animesaturn-downloader -u https://your-url-here/anime -f 1 -l 12 -d ./my-anime -n MyAnime_ -w 3
```

Questo comando invoca l'eseguibile con i seguenti parametri:

-   url: https[]()://your-url-here/anime
-   primo episodio: 1
-   ultimo episodio: 12
-   cartella output: ./my-anime
-   nome dei file: MyAnime\_
-   worker da usare: 3

NB: Il prgoramma aggiunge "i.mp4" alla fine di ogni file con i uguale al numero dell'episodio scaricato.

### Applicazione Grafica

L'eseguibile per l'applicazione grafica è reperibile in `./bin/animesaturn-downloader-(VERSIONE)-(PIATTAFORMA)-gui` e presenta una interfaccia grafica realizzata con svelte e wails.
