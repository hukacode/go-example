package main

import (
	"os"
	"runtime"

	"github.com/atotto/clipboard"
	"github.com/ncruces/zenity"
)

func main() {
	content, err := clipboard.ReadAll()
	check(err)

	if content != "" {
		inboxFilePath := ""

		if runtime.GOOS == "windows" {
			inboxFilePath = "C:\\Dropbox\\_huka\\digital-garden\\content\\note\\gtd\\inbox.md"
		} else if runtime.GOOS == "linux" {
			inboxFilePath = "/home/huka/Dropbox/_huka/digital-garden/content/note/gtd/inbox.md"
		}

		file, err := os.OpenFile(inboxFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		check(err)

		defer file.Close()

		if _, err = file.WriteString("\n\n" + content); err != nil {
			panic(err)
		}

		zenity.Notify("Sent to inbox!")
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
