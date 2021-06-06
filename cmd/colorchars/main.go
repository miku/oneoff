package main

import (
	"bufio"
	"io"
	"log"
	"os"

	"github.com/fatih/color"
)

var colors = []*color.Color{
	color.New(color.FgCyan),
	color.New(color.FgRed),
	color.New(color.FgGreen),
	color.New(color.FgYellow),
	color.New(color.FgBlue),
	color.New(color.FgMagenta),
	color.New(color.FgWhite),
	color.New(color.FgHiCyan),
	color.New(color.FgHiRed),
	color.New(color.FgHiGreen),
	color.New(color.FgHiYellow),
	color.New(color.FgHiBlue),
	color.New(color.FgHiMagenta),
	color.New(color.FgHiWhite),
	color.New(color.BgCyan),
	color.New(color.BgRed),
	color.New(color.BgGreen),
	color.New(color.BgYellow),
	color.New(color.BgBlue),
	color.New(color.BgMagenta),
	color.New(color.BgWhite),
	color.New(color.BgHiCyan),
	color.New(color.BgHiRed),
	color.New(color.BgHiGreen),
	color.New(color.BgHiYellow),
	color.New(color.BgHiBlue),
	color.New(color.BgHiMagenta),
	color.New(color.BgHiWhite),
}

func main() {
	var br = bufio.NewReader(os.Stdin)
	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		for _, c := range line {
			colors[int(c)%len(colors)].Printf("%c", c)
		}
	}

}
