package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/skip2/go-qrcode"
	"github.com/muesli/termenv"
)

func main() {
	var args = os.Args[1:]

	var level = qrcode.Medium
	if len(args) > 1 {
		if args[0] == "--low" {
			level = qrcode.Low
			args = args[1:]
		} else if args[0] == "--medium" {
			level = qrcode.Medium
			args = args[1:]
		} else if args[0] == "--high" {
			level = qrcode.High
			args = args[1:]
		} else if args[0] == "--highest" {
			level = qrcode.Highest
			args = args[1:]
		}
	}

	if len(args) > 1 || (len(args) == 1 && args[0] == "--help") {
		fmt.Printf("Generate & print unicode QR codes on the command line.\n")
		fmt.Printf("Usage: %s [--low|--medium|--high|--highest] [text]\n", os.Args[0])
		fmt.Printf("If no text is given, read from STDIN.\n")
		os.Exit(1)
	}

	var content = ""
	if len(args) < 1 {
		data, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println(err)
			return
		}
		content = strings.TrimSuffix(string(data), "\n")
	} else {
		content = args[0]
	}

	qr, err := qrcode.New(content, level)
	if err != nil {
		fmt.Println(err)
		return
	}

	output := bytes.NewBuffer([]byte{})

	bitmap := qr.Bitmap()
	height := len(bitmap)
	width := len(bitmap[0])

	var line *bytes.Buffer
	for y := 1; y < height - 2; y += 2 {
		line = bytes.NewBuffer([]byte{})
		for x := 1; x < width - 1; x++ {
			if bitmap[y][x] && bitmap[y+1][x] {
				line.WriteRune('█')
			} else if bitmap[y][x] && !bitmap[y+1][x] {
				line.WriteRune('▀')
			} else if !bitmap[y][x] && bitmap[y+1][x] {
				line.WriteRune('▄')
			} else if !bitmap[y][x] && !bitmap[y+1][x] {
				line.WriteRune(' ')
			}
		}
		output.WriteString(termenv.String(line.String()).Reverse().String())
		output.WriteByte('\n')
	}
	output.WriteString(strings.Repeat("▀", width - 2))
	output.WriteByte('\n')
	output.WriteTo(os.Stdout)
}
