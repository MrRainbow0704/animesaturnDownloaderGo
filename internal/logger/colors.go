package logger

type Color string

const (
	ColorReset           Color = "\033[0m"
	ColorFgBlack         Color = "\033[30m"
	ColorFgRed           Color = "\033[31m"
	ColorFgGreen         Color = "\033[32m"
	ColorFgYellow        Color = "\033[33m"
	ColorFgBlue          Color = "\033[34m"
	ColorFgMagenta       Color = "\033[35m"
	ColorFgCyan          Color = "\033[36m"
	ColorFgWhite         Color = "\033[37m"
	ColorFgDefault       Color = "\033[39m"
	ColorFgBlackBright   Color = "\033[90m"
	ColorFgRedBright     Color = "\033[91m"
	ColorFgGreenBright   Color = "\033[92m"
	ColorFgYellowBright  Color = "\033[93m"
	ColorFgBlueBright    Color = "\033[94m"
	ColorFgMagentaBright Color = "\033[95m"
	ColorFgCyanBright    Color = "\033[96m"
	ColorFgWhiteBright   Color = "\033[97m"
	ColorBgBlack         Color = "\033[40m"
	ColorBgRed           Color = "\033[41m"
	ColorBgGreen         Color = "\033[42m"
	ColorBgYellow        Color = "\033[43m"
	ColorBgBlue          Color = "\033[44m"
	ColorBgMagenta       Color = "\033[45m"
	ColorBgCyan          Color = "\033[46m"
	ColorBgWhite         Color = "\033[47m"
	ColorBgDefault       Color = "\033[49m"
	ColorBgBlackBright   Color = "\033[100m"
	ColorBgRedBright     Color = "\033[101m"
	ColorBgGreenBright   Color = "\033[102m"
	ColorBgYellowBright  Color = "\033[103m"
	ColorBgBlueBright    Color = "\033[104m"
	ColorBgMagentaBright Color = "\033[105m"
	ColorBgCyanBright    Color = "\033[106m"
	ColorBgWhiteBright   Color = "\033[107m"
)

// Restituisce una stringa avvolta nei colori selezionati.
func Colorize(s string, c ...Color) string {
	return ColorizeNoReset(s, c...) + string(ColorReset)
}

func ColorizeNoReset(s string, c ...Color) string {
	var cs Color
	for _, cx := range c {
		cs = cs + cx
	}
	return string(cs) + s
}
