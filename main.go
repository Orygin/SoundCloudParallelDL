// Based on the script written by /u/redditgoogle @ https://www.reddit.com/r/Music/comments/4597e6/soundcloud_could_be_forced_to_close_after_44m/czw8q9q
package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"math/rand/v2"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var maxParallel = flag.Int("maxparallel", 1, "Defines the maximum parrallel artists to download")
var silent = flag.Bool("silent", false, "Defines if the youtube-dl output is to be printed")
var cookieFile = flag.String("cookiefile", "", "Defines the cookie file")

func main() {
	flag.Parse()
	var artists []string
	list, err := os.Open(workingDir + "artistlist.txt")
	if err != nil {
		log.Printf("Could not open artist list file : %v", err)
		return
	}

	scanner := bufio.NewScanner(list)

	var total = 0
	for scanner.Scan() {
		total++
		artist := scanner.Text()
		artists = append(artists, artist)
	}
	list.Close()
	log.Println("downloading " + strconv.Itoa(total))

	rand.Shuffle(len(artists), func(i, j int) {
		artists[i], artists[j] = artists[j], artists[i]
	})

	for {
		for _, artist := range artists {
			fetch(artist)
		}
		time.Sleep(time.Hour * 24)
	}
}

func fetch(artist string) {
	log.Println("Running for " + artist)
	args := []string{
		"soundcloud.com/" + artist + "/tracks",
		"-o", workingDir + artist + "/%(artist)s-%(title)s.%(ext)s",
		"-x", "--audio-format", "flac", "--audio-quality", "8",
		/*"--add-metadata", "--write-description", */ "-w", "--no-progress",
		"--sleep-requests", "6",
		//"--extractor-args", "soundcloud:formats=*_aac",
	}
	if *cookieFile != "" {
		args = append(args, "--cookies", *cookieFile)
	}
	log.Printf("args: %s", strings.Join(args, " "))

	ctx, cancel := context.WithCancel(context.Background())

	cmd := exec.CommandContext(ctx, cmdBinName, args...)
	if !*silent {
		cmd.Stdout = NewLogWriter(artist, cancel)
	}
	cmd.Stderr = os.Stderr
	cmd.Cancel = func() error {
		proc, err := os.FindProcess(-cmd.Process.Pid)
		if err != nil {
			log.Printf("Could not find process: %v", err)
			return err
		}
		return proc.Signal(syscall.SIGTERM)
	}
	aggrementCmd(cmd)

	if err := cmd.Run(); err != nil {
		log.Printf("could not run ytdlp: %s", err.Error())
	}
	log.Println("Finished " + artist)
}

type LogWriter struct {
	artist        string
	progress, max int
	stop          context.CancelFunc
}

func NewLogWriter(artist string, stop context.CancelFunc) *LogWriter {
	return &LogWriter{
		artist: artist,
		stop:   stop,
	}
}

func (lw LogWriter) Write(p []byte) (n int, err error) {
	str := string(p)
	log.Println(lw.artist + ": " + string(p))
	if strings.Contains(str, "has already been downloaded") {
		lw.stop()
	}
	// I tried matching the youtube-dl progress but it's harder than I thought to get it right
	/*match, _ := regexp.Match(".*?(\\[download\\])( )(Downloading)( )(video)( )(\\d+)( )(of)( )(\\d+).*?", p)
	if match {
		log.Println(lw.artist + ": " + string(p))
	}*/
	return len(p), nil
}
