package main

import (
	"os"
)

const (
	header = `                 _                 _____       _                    
     /\         (_)               / ____|     | |                   
    /  \   _ __  _ _ __ ___   ___| (___   __ _| |_ _   _ _ __ _ __  
   / /\ \ | '_ \| | '_ ` + "`" + ` _ \ / _ \\___ \ / _` + "`" + ` | __| | | | '__| '_ \ 
  / ____ \| | | | | | | | | |  __/____) | (_| | |_| |_| | |  | | | |
 /_/___ \_\_| |_|_|_| |_| |_|\___|_____/ \__,_|\__|\__,_|_|  |_| |_|
 |  __ \                    | |               | |                   
 | |  | | _____      ___ __ | | ___   __ _  __| | ___ _ __          
 | |  | |/ _ \ \ /\ / / '_ \| |/ _ \ / _` + "`" + ` |/ _` + "`" + ` |/ _ \ '__|         
 | |__| | (_) \ V  V /| | | | | (_) | (_| | (_| |  __/ |            
 |_____/ \___/ \_/\_/ |_| |_|_|\___/ \__,_|\__,_|\___|_|           

AnimesaturnDownloader è una utility per scaricare gli anime dal sito AnimeSaturn.
Scritto in Go da Marco Simone.


`
)

func main() {
	initRoot()
	subcmd := parseRoot(os.Args[1:])
	runRoot(subcmd)
}
