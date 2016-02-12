//Based on the script written by /u/redditgoogle @ https://www.reddit.com/r/Music/comments/4597e6/soundcloud_could_be_forced_to_close_after_44m/czw8q9q
package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var lock chan bool

var maxParallel = flag.Int("maxparallel", 2, "Defines the maximum parrallel artists to download")
var silent = flag.Bool("silent", false, "Defines if the youtube-dl output is to be printed")

func main() {
	flag.Parse()
	var artists []string
	list, err := os.Open("artistlist.txt")
	if err != nil {
		log.Println("Could not open artist list file : artistlist.txt")
		return
	}
	defer list.Close()

	scanner := bufio.NewScanner(list)

	var total = 0
	for scanner.Scan() {
		total++
		artist := scanner.Text()
		artists = append(artists, artist)
	}
	log.Println("downloading " + strconv.Itoa(total))
	lock = make(chan bool, *maxParallel)
	for _, artist := range artists {
		go fetch(artist)
	}
	time.Sleep(100) // Wait for goroutines to lock
	for len(lock) != 0 {

	}
}

func fetch(artist string) {
	lock <- true
	log.Println("Running for " + artist)
	cmd := exec.Command("youtube-dl.exe", "soundcloud.com/"+artist+"/tracks", "-o "+artist+"/%(title)s.%(ext)s", "--add-metadata", "--write-description", "--no-progress")
	if !*silent {
		cmd.Stdout = NewLogWriter(artist)
	}
	cmd.Stderr = os.Stderr
	cmd.Run()
	log.Println("Finished " + artist)
	<-lock
}

type LogWriter struct {
	artist        string
	progress, max int
}

func NewLogWriter(artist string) *LogWriter {
	lw := &LogWriter{}
	lw.artist = artist
	return lw
}

func (lw LogWriter) Write(p []byte) (n int, err error) {
	log.Println(lw.artist + ": " + string(p))
	// I tried matching the youtube-dl progress but it's harder than I thought to get it right
	/*match, _ := regexp.Match(".*?(\\[download\\])( )(Downloading)( )(video)( )(\\d+)( )(of)( )(\\d+).*?", p)
	if match {
		log.Println(lw.artist + ": " + string(p))
	}*/
	return len(p), nil
}
