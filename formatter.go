// Copyright © 2019 Marco Berardelli marco.berardelli@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package youGO

import(
	"os/exec"
	"fmt"
	"regexp"
	"log"
	"path/filepath"
	s "strings"
	"bytes"
	"os"
)


type TitleFormatter interface {
	FormatTitle(VideoData) SongData
}


// VideoData represents the youtube video you want to download.
// It contains the title of the video and its video ID.
type VideoData struct {
	Title string
	VideoID string
}


// SongData contains all the information about the song downloaded.
type SongData struct {
	Title string
	Artist string
	Video VideoData
	// CorrectedName will be true if the title of the video doesn't contain any peculiar characters.
	CorrectedName bool
}


// Formatter implements TitleFormatter and is used to extract the title and artist of the song from the video title.
// It also update the file with the title and artist metadata
type Formatter struct{

	// Slice that contains all the characters for identifying if the song has more than an artist, such as "ft." and "&".
	ArtistDelimiters []string

	// A Regexp is a compiled regular expression.
	// The regular expression will remove all non-alphanumeric characters.
	Regexp *regexp.Regexp
}


// NewFormatter inizializes a Formatter with "[^a-zA-Z0-9]+" as compiled regular expression and with delimiters passed as arguments.
func NewFormatter(delimiters ...string) (Formatter, error){
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
    if err != nil {
        return Formatter{}, err
    }
	
	return Formatter{ArtistDelimiters: delimiters, Regexp: reg, }, nil
}


func (f Formatter) FormatTitle(video VideoData) SongData {

	if s.Count(video.Title, " - ") != 1 {
		song := SongData{
			Video: video,
			CorrectedName: false,
		}
		return song
	}


	/*
	// TODO:
	r := make([]string, len(f.ArtistDelimiters)*2)
	for _, delimiter := range f.ArtistDelimiters {
		r = append(r, delimiter)
		r = append(r, ", ")
	}
	replacer := s.NewReplacer(r...)
	*/
	
	
	// TODO: 
	titleSplitted := s.Split(video.Title, " - ")
	song := SongData{
		Title: titleSplitted[1],
		Artist: titleSplitted[0],
		Video: video,
		CorrectedName: true,
	}
	
	return song
}

func sanitize(title string) string {
	cchar := []string{"?"}
	for _, c := range cchar {
		title = s.Replace(title, c, "", -1)
	}
	return title
}

func (f Formatter) FormatFile(song SongData) {

	if !song.CorrectedName {
		return
	}

	// youtube-dl saves the file without considering special character, such as ?.
	title := sanitize(song.Video.Title)

	errorPath := filepath.Join("songs", "error")
	matches, err := filepath.Glob(filepath.Join(errorPath, title) + ".*")
    if err != nil {
		fmt.Println(err)
		return
	}

	filename := s.Replace(matches[0], errorPath, "", 1)
	path := filepath.Join("songs", filename)

	cmd := exec.Command("ffmpeg","-y", "-i", filepath.Join(errorPath, filename), "-map", "0", "-c", "copy", "-metadata", fmt.Sprintf(`title="%s"`, song.Title), "-metadata", fmt.Sprintf(`author="%s"`, song.Artist), path)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err = cmd.Run()
	if err != nil {
		log.Printf(errb.String())
	}
	
	err = os.Remove(filepath.Join(errorPath, filename))
	if err != nil {
		log.Printf("Error deleting a copy of " + filename + " in " + errorPath)
	}
}
