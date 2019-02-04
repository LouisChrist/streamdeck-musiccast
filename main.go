package main

import (
	"log"

	"github.com/LouisChrist/streamdeck-musiccast/sdplugin"
)

func main() {
	// file, err := os.OpenFile("C:\\Users\\Louis\\Desktop\\streamdeck.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer file.Close()
	// log.SetOutput(file)

	plugin, err := sdplugin.New(newMusicCastHandler())
	if err != nil {
		log.Fatal(err)
	}
	defer plugin.Close()
	log.Fatal(plugin.Run())
}
