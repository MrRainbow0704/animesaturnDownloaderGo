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

Questa utility permette di scaricare anime dal famoso sito AnimeSaturn e salvarli in formato .mp4 sul computer. Contiene una versione CLI e una con interfaccia grafica.

> [!NOTE]
> Questo programma è distribuito sotto la [terza versione della GNU General Public License](LICENSE.md) (GNP GPLv3). Si è liberi di copiare, modificare e distribuire questo software amche per scopi commerciali a patto che il codice risultante venga distribuito in modo open source e sotto la stessa identica licenza. Questo software viene fornito così com'è e senza alcuna garanzia.

> [!NOTE]
> In caso di bug o errori per favore [aprite un Issue](https://github.com/MrRainbow0704/animesaturnDownloaderGo/issues/new/choose) con tutte le informazioni relative.
> Questo programma è un hobby e non ne traggo guadagno, quindi vi risponderò appena mi sarà possibile.
>
> Inoltre, se possedete abbastanza conoscenze da riuscire a risolvere da soli il vostro problema, potete [aprire una Pull Request](https://github.com/MrRainbow0704/animesaturnDownloaderGo/compare) per contribuire a questo progetto con il vostro codice.

## Build

Per creare l'eseguibile, dopo aver installato correttamente go e make (e npm per la build della versione GUI), eseguire nel terminale uno dei seguenti comandi:

```console
# Esegue la build di tutte le versioni
foo@bar:~/animesaturnDownloaderGo$ make

# Esegue la build per windows
foo@bar:~/animesaturnDownloaderGo$ make win

# Esegue la build per linux
foo@bar:~/animesaturnDownloaderGo$ make linux

# Esegue la build per mac
foo@bar:~/animesaturnDownloaderGo$ make mac

# Esegue la build solo della versione CLI
# per una delle piattaforme precedenti
foo@bar:~/animesaturnDownloaderGo$ make [ARCH]-cli

# Esegue la build solo della versione GUI
# per una delle piattaforme precedenti
foo@bar:~/animesaturnDownloaderGo$ make [ARCH]-gui
```

Al termine della build, gli eseguibili saranno disponibili all'interno della cartella `./bin` con nome `animesaturn-downloader-[VERSIONE]-[PIATTAFORMA]` per la versione CLI e `animesaturn-downloader-[VERSIONE]-[PIATTAFORMA]-gui` per la versione GUI. <br>

> [!TIP]
> È possibile fare la build per Linux da un dispositivo Windows, a patto che si usi il WSL (Windows Subsystem for Linux).

> [!WARNING]
> Non è possibile eseguire la build per mac su hardware non mac. [Prendetevela con Apple](https://github.com/wailsapp/wails/issues/1041#issuecomment-2492133624).

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

L'utility mette a disposizione 2 sottocomandi: [`download`](#sottocomando-download) e [`search`](#sottocomando-search).

#### Sottocomando `download`

Il sottocomando `downlaod` provvede un modo per scaricare anime dal sito. Maggiori informazioni sul funzionamento del comando possono essere ottenute con il seguente comando:

```console
foo@bar:~/animesaturnDownloaderGo$ ./bin/animesaturn-downloader download -h
```

Per scaricare un anime, si può usare il seguente comando:

```console
foo@bar:~/animesaturnDownloaderGo$ ./bin/animesaturn-downloader download -u https://your-url-here/anime -f 1 -l 12 -d ./my-anime -n MyAnime_ -w 3
```

Questo comando invoca l'eseguibile con i seguenti parametri:

- url: https[]()://your-url-here/anime
- primo episodio: 1
- ultimo episodio: 12
- cartella output: ./my-anime
- nome dei file: MyAnime\_
- worker da usare: 3

> [!NOTE]
> Il prgoramma aggiunge  `i.mp4` alla fine di ogni file con `i` uguale al numero dell'episodio scaricato.

> [!NOTE]
> In caso di spazi in qualunque degli argomenti, si raccomanda di avvolgere quell'argomento in virgolette o apici (`""` o `''`) o eseguire l'escape degli spazi (`\ ` per shell POSIX, `` `  `` per powershell, `^ ` per cmd windows) al fine di minimizzare errori di parsing.

> [!CAUTION]
> Questa utility supporta solo gli url presenti in questa lista: [Domini Ufficiali](https://www.animesaturn.me/).
>
> Per quanto altri url potrebbero funzionare non sono supportati e il loro utilizzo potrebbe portare al furto di dati personali da parte di terzi o anche al download di malware. Si ricorda di usare cautela quando si naviga in rete.

#### Sottocomando `search`

Il sottocomando `search` provvede un modo per cercare anime dal sito. Maggiori informazioni sul funzionamento del comando possono essere ottenute con il seguente comando:

```console
foo@bar:~/animesaturnDownloaderGo$ ./bin/animesaturn-downloader search -h
```

Per cercare un anime, si può usare il seguente comando:

```console
foo@bar:~/animesaturnDownloaderGo$ ./bin/animesaturn-downloader search -s Attack
```

Questo comando cercherà tutti gli anime il cui nome contiene la parola "Attack" e ne restituirà titoli e url.

### Applicazione Grafica

L'eseguibile per l'applicazione grafica è reperibile in `./bin/animesaturn-downloader-[VERSIONE]-[PIATTAFORMA]-gui` e presenta una interfaccia grafica realizzata con svelte e wails.
