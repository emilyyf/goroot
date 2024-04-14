package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/daviseidel/xlib"
	"log"
	"net"
	"time"
)

type RustResponse struct {
	Ok []byte `json:"Ok"`
}

type Artist struct {
	Name string `json:"name"`
}

type Album struct {
	Name string `json:"name"`
}

type Item struct {
	Name    string   `json:"name"`
	Album   Album    `json:"album"`
	Artists []Artist `json:"artists"`
}

type SpotifyData struct {
	Item Item `json:"item"`
}

func get_data(conn *net.UDPConn, data chan SpotifyData) {
	var buffer bytes.Buffer
	buf := make([]byte, 4096)
	var a RustResponse
	var b SpotifyData

	for {
		_, err := conn.Write([]byte("{\"Get\": {\"Key\": \"Playback\"}}"))
		conn.SetReadDeadline(time.Now().Add(time.Millisecond * 500))
		if err != nil {
			log.Fatal(err)
		}
		for {
			n, err := conn.Read(buf)
			if err != nil {
				data <- SpotifyData{}
				goto reset
			}
			buffer.Write(buf[:n])
			if n == 0 {
				break
			}
		}

		err = json.Unmarshal(buffer.Bytes(), &a)
		if err != nil {
			log.Fatalf("could not marshal json: %s\n", err)
		}

		err = json.Unmarshal(a.Ok, &b)
		if err != nil {
			log.Fatal(err)
		}

		buffer.Reset()

		data <- b

	reset:
		time.Sleep(time.Millisecond * 500)
	}
}

func main() {
	disp := xlib.XOpenDisplay(0)
	root := xlib.XDefaultRootWindow(disp)
	data_chan := make(chan SpotifyData)
	addr, err := net.ResolveUDPAddr("udp", "localhost:8080")
	var music string

	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	go get_data(conn, data_chan)

	for {
		select {
		case data := <-data_chan:
			if data.Item.Name != "" {
				music = fmt.Sprintf("Playing: %s by %s (%s)", data.Item.Name, data.Item.Artists[0].Name, data.Item.Album.Name)
			} else {
				music = "Not playing"
			}
		default:
		}

		date := time.Now().Format("Mon Jan 2 2006 15:04:05.000000")

		str := fmt.Sprintf("%s | %s", music, date)

		xlib.XStoreName(disp, root, str)
		xlib.XFlush(disp)
		time.Sleep(time.Millisecond * (1000 / 60))
	}
}
