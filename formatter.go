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
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	s "strings"
)

// TitleFormatter is the interface that contains the FormatTitle method.
// FormatTitle tries to obtains the artist and title of the song by parsing the title of the video.
type TitleFormatter interface {
	FormatTitle(VideoData) SongData
}

// FileFormatter is the interface that contains the FormatFile method.
// FormatFile edits the metadata tags of the file downloaded.
type FileFormatter interface {
	FormatFile(SongData)
}


// VideoData represents the youtube video you want to download.
// It contains the title of the video and its video ID.
type VideoData struct {
	Title string
	
	VideoID string
}


// SongData contains all the information about the song downloaded.
type SongData struct {
	Artist string

	Title string

	// CorrectedName is set to true if the title of the video doesn't contain any peculiar characters.
	CorrectedName bool

	Video VideoData
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


// FormatTitle retreives the title and the artist by parsing the video title.
// If the name of the video is diffucult to parse, FormatTitle returns SongData with empty Artist and Title parameters
func (f Formatter) FormatTitle(video VideoData) SongData {

	if s.Count(video.Title, " - ") != 1 {
		song := SongData{
			Video: video,
			CorrectedName: false,
		}
		return song
	}


	/*
	// TODO: upgrade the algorithm used to parse the title
	r := make([]string, len(f.ArtistDelimiters)*2)
	for _, delimiter := range f.ArtistDelimiters {
		r = append(r, delimiter)
		r = append(r, ", ")
	}
	replacer := s.NewReplacer(r...)
	*/
	// [\p{L}-&.]
	// \[.*\]\s?
	
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
	// TODO: change to map and assign "" to "?" and "_" to "|"
	cchar := []string{"?","|"}
	for _, c := range cchar {
		title = s.Replace(title, c, "", -1)
	}
	return title
}


// FormatFile edits the file downloaded by setting the artist and title metadata.
func (f Formatter) FormatFile(song SongData, path string) error {

	if !song.CorrectedName {
		return nil
	}

	// youtube-dl saves the file without considering special character, such as "?".
	//title := sanitize(song.Video.Title)
	//fmt.Println("Title:  " + title)

	filePath := filepath.Join(path, song.Video.Title)
	formattedPath := filepath.Join(path, "formatted")
	
	
	matches, err := filepath.Glob(filePath + ".*")
	if err != nil {
		return err
	}

	if len(matches) == 0 {
		return errors.New("Impossible to open: " + song.Video.Title)
	}

	filename := s.Replace(matches[0], path, "", 1)
	formattedFilepath := filepath.Join(formattedPath, filename)

	// Calling ffmpeg to edit the metadata tags
	cmd := exec.Command("ffmpeg","-y", "-i", matches[0], "-map", "0", "-c", "copy", "-metadata", fmt.Sprintf(`title="%s"`, song.Title), "-metadata", fmt.Sprintf(`artist="%s"`, song.Artist), formattedFilepath)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err = cmd.Run()
	if err != nil {
		return errors.New(errb.String())
	}
	
	// Removing old file
	err = os.Remove(filepath.Join(path, filename))
	if err != nil {
		return err
	}

	return nil
}
