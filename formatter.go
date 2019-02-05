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
	s "strings"
	"os"
	"fmt"
	"regexp"
	"log"
	"github.com/bogem/id3v2"
)

const token = "---"

// SongInfo is used to store the title and the artist of the song and it also has the path where the mp3 should be saved.
// If the title contains any character complicated to parse the Error var will be updated to true.
type SongInfo struct {

	Title string

	// If there are more than an artist, each one will be separeted by a comma.
	Artist string

	// Path where the file will be saved.
	Path string

	// Error will be true if there are characters that cause any problems while parsing the title/artist.
	Error bool

	VideoID string
}

// A Formatter has the job of manipulating the title of the youtube video to extract the title and the artist of the song.
type Formatter struct{

	// Slice that contains all the characters for identifying if the song has more than an artist.
	ArtistDelimiters []string

	// A Regexp is a compiled regular expression.
	// The regular expression will remove all non-alphanumeric characters (including spaces, see bug) 
	Regexp *regexp.Regexp

	FilesNotDone []*SongInfo
}

// NewFormatter inizializes a Formatter with "[^a-zA-Z0-9]+" as compiled regular expression and with delimiters passed as arguments.
// Then it returns the pointer to the struct
func NewFormatter(delimiters ...string) (*Formatter, error){
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
    if err != nil {
        return nil, err
    }
    
	return &Formatter{ArtistDelimiters: delimiters, Regexp: reg}, nil
}

func (f *Formatter) FormatFolder(folderPath string) error{

	return nil
}

// FormatMp3 takes a pointer to SongInfo as parameter 
func (f *Formatter) FormatMp3(songInfo *SongInfo) error {

	correctName := s.Replace(songInfo.Path, token, " ", -1)
	err := os.Rename(songInfo.Path, correctName)
	if err != nil {
		f.FilesNotDone = append(f.FilesNotDone, songInfo)
		fmt.Println("Redownloading the song")
		return err
	}

	if songInfo.Error {
		//wg.Done()
		return nil
	}

	tag, err := id3v2.Open(correctName, id3v2.Options{Parse: true})
	if err != nil {
 		fmt.Println("Error while opening mp3 file: " + err.Error())
 	}
	tag.SetArtist(songInfo.Artist)
	tag.SetTitle(songInfo.Title)
	if err = tag.Save(); err != nil {
		log.Fatal("Error while saving a tag: ", err)
	}
	tag.Close()
	
	//wg.Done()
	return nil
}

func (f *Formatter) FormatTitle(title string) (*SongInfo, error) {

	r := make([]string, len(f.ArtistDelimiters)*2)
	for _, delimiter := range f.ArtistDelimiters {
		r = append(r, delimiter)
		r = append(r, ", ")
	}
	replacer := s.NewReplacer(r...)
	
	if s.Count(title, " - ") != 1 {
		
		noSpace := s.Replace(title, " ", token, -1)
		songInfo := &SongInfo{Title: f.Regexp.ReplaceAllString(noSpace, token), Artist: "", Error: true}
		return songInfo, NewErrorProblematicName("Name has difficult characters to understand, saving in Formatter.FilesNotDone")
	}
	
	nameSplitted := s.Split(title, " - ")
	filename := f.Regexp.ReplaceAllString(nameSplitted[1], " ")
	return &SongInfo{Title: s.Replace(filename, " ", token, -1), Artist: replacer.Replace(nameSplitted[0])}, nil
}

