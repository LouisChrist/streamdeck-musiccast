# MusicCast plugin for Streamdeck

Simple plugin to change and monitor the power status of Yamaha MusicCast devices.
Works only on windows.

## Install

Download *musiccast.streamDeckPlugin* from release page and install it by opening the file.

## Build from source

Clone repo and build executable with `go build`.

    git clone git@github.com:LouisChrist/streamdeck-musiccast.git de.louischrist.musiccast.sdPlugin
    cd de.louischrist.musiccast.sdPlugin
    go build

### Manual install

Copy directory to %appdata%\Elgato\StreamDeck\Plugins.

### Create *.streamDeckPlugin* file

Use [DistributionTool](https://developer.elgato.com/documentation/stream-deck/sdk/exporting-your-plugin/)
from the StreamDeck SDK site to build.

    DistributionTool.exe de.louischrist.musiccast.sdPlugin Release
