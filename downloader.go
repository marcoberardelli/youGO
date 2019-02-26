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
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// Downloader contains the path where to store the files and a reference to the youtube service used to retreive video information.
type Downloader struct{
	Path string
	YtService *youtube.Service
}


// NewDownloader initializes a new instance of Downloader.
// If the folders songs and formatted don't exist they will be created.
func NewDownloader(path string) (Downloader, error) {

	client := &http.Client{
		Transport: &transport.APIKey{Key: YouTubeAPIKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		return Downloader{}, NewErrorServiceCreation("Impossible to create youtube service")
	}

	// Creating the songs folder, if it not exists
	path = filepath.Join(path, "songs")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}
	// Creating the formatted folder, if it not exists
	formattedPath := filepath.Join(path, "formatted")
	if _, err := os.Stat(formattedPath); os.IsNotExist(err) {
		os.MkdirAll(formattedPath, os.ModePerm)
	}

	downloader := Downloader{
		YtService: service,
		Path: path,
	}

	return downloader, nil
}


func (d *Downloader) download(video VideoData, path string) error {
	path = filepath.Join(path, "%(title)s.%(ext)s")
	cmd := exec.Command("youtube-dl", "-x", "-f", "bestaudio", "-o", path, "https://www.youtube.com/watch?v="+video.VideoID)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		return errors.New(errb.String())
	}
	return nil
}


func (d *Downloader) downloadFromPlaylist(item *youtube.PlaylistItem, tFormatter TitleFormatter, toFormat bool) error {

	// Checking if the video wasn't deleted
	if item.Status == nil {
		return errors.New("Unable to get the video")
	}
	if item.Status.PrivacyStatus == "private" {
		return errors.New("Private video")
	}

	video := VideoData{
		VideoID: item.ContentDetails.VideoId,
		Title: sanitize(item.Snippet.Title),
	}
	
	fmt.Println("Downloading " + video.Title)
	d.download(video, d.Path)

	if toFormat {
		song := tFormatter.FormatTitle(video)
		// Creating a new Formatter, used to add the metadata to the file, if you wrote your own implementation of TitleFormatter.
		formatter, ok := tFormatter.(Formatter)
		if !ok {
			formatter = Formatter{}
		}
		err := formatter.FormatFile(song, d.Path)
		if err != nil {
			return err
		}
	}
	return nil
}


// 
func(d *Downloader) DownloadPlaylistAndFormat(playlistID string, tFormatter TitleFormatter) {

	call := d.YtService.PlaylistItems.List("snippet,contentDetails,status")
	call.MaxResults(50)
	call.PlaylistId(playlistID)
	response, err := call.Do()
	if err != nil {
		fmt.Println(err)
	}
	
	for _, item := range response.Items {
		d.downloadFromPlaylist(item, tFormatter, true)
	}

	// The YouTube API returns just a limited number of videos from a playlist, organized in "pages".
	// Looping until we reach the last page.
	for response.NextPageToken != ""  {
		call = d.YtService.PlaylistItems.List("snippet,contentDetails,status")
		call.PageToken(response.NextPageToken)
		call.MaxResults(50)
		call.PlaylistId(playlistID)
		response, err = call.Do()
		if err != nil {
			fmt.Println(err)
			continue
		}
		
		for _, item := range response.Items {
			err := d.downloadFromPlaylist(item, tFormatter, true)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}



func (d *Downloader) DownloadVideoAndFormat(videoID string, tFormatter TitleFormatter) error {

	call := d.YtService.Videos.List("snippet,status")
	call = call.Id(videoID)
	response, err := call.Do()
	if err != nil {
		return err
	}
	if response.Items[0] == nil {
		return errors.New("Unable to get the video")
	}
	if response.Items[0].Status.PrivacyStatus == "private" {
		return errors.New("Private video")
	}


	video := VideoData{
		VideoID: response.Items[0].Id,
		Title: response.Items[0].Snippet.Title,
	}
	
	fmt.Println("Downloading " + video.Title)
	d.download(video, d.Path)

	song := tFormatter.FormatTitle(video)
	// Creating a new Formatter, used to add the metadata to the file, if you wrote your own implementation of TitleFormatter.
	formatter, ok := tFormatter.(Formatter)
	if !ok {
		formatter = Formatter{}
	}
	formatter.FormatFile(song, d.Path)

	return nil
}