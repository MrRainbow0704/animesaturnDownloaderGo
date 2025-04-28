package logger

type color uint

func (c color) String() string {
	return colorsCodes[c]
}

var colorsCodes = []string{
	"\033[0m",
	"\033[30m",
	"\033[31m",
	"\033[32m",
	"\033[33m",
	"\033[34m",
	"\033[35m",
	"\033[36m",
	"\033[37m",
	"\033[39m",
	"\033[90m",
	"\033[91m",
	"\033[92m",
	"\033[93m",
	"\033[94m",
	"\033[95m",
	"\033[96m",
	"\033[97m",
	"\033[40m",
	"\033[41m",
	"\033[42m",
	"\033[43m",
	"\033[44m",
	"\033[45m",
	"\033[46m",
	"\033[47m",
	"\033[49m",
	"\033[100m",
	"\033[101m",
	"\033[102m",
	"\033[103m",
	"\033[104m",
	"\033[105m",
	"\033[106m",
	"\033[107m",
}

const (
	reset color = iota
	fgBlack
	fgRed
	fgGreen
	fgYellow
	fgBlue
	fgMagenta
	fgCyan
	fgWhite
	fgDefault
	fgBlackBright
	fgRedBright
	fgGreenBright
	fgYellowBright
	fgBlueBright
	fgMagentaBright
	fgCyanBright
	fgWhiteBright
	bgBlack
	bgRed
	bgGreen
	bgYellow
	bgBlue
	bgMagenta
	bgCyan
	bgWhite
	bgDefault
	bgBlackBright
	bgRedBright
	bgGreenBright
	bgYellowBright
	bgBlueBright
	bgMagentaBright
	bgCyanBright
	bgWhiteBright
)

// Restituisce una stringa avvolta nei colori selezionati.
func Colorize(s string, c ...color) string {
	cs := ""
	for _, cx := range c {
		cs = cs + cx.String()
	}
	return cs + s + reset.String()
}
