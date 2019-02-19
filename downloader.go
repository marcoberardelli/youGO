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
	"net/http"
	"log"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
	"runtime"
	"fmt"
	"os"
	"github.com/mitchellh/go-homedir"
	"path/filepath"

	"os/exec"
	"bytes"
)

type Downloader struct{
	YtService *youtube.Service
}


func NewDownloader() (*Downloader, error) {

	client := &http.Client{
		Transport: &transport.APIKey{Key: YouTubeAPIKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		return nil, NewErrorServiceCreation("Impossible to create youtube service")
	}

	// TODO: convert using filepath.Join()
	// Also
	var dir, errorDir string
	if runtime.GOOS == "windows" {
		dir, _ = homedir.Dir()
		dir = dir + "\\Desktop\\songs"
		errorDir = dir + "\\error"
	} else {
		dir, _ = homedir.Dir()
		dir = dir + "/songs"
		errorDir = dir + "/error"
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
		os.MkdirAll(errorDir, os.ModePerm)
	}



	// formatter, err := NewFormatter(" & ", " x ", " ft. ", " feat ")
	// if err != nil {
	// 	return nil, err
	// }

	downloader := &Downloader{
		YtService: service,
	}

	return downloader, nil
}


func(d *Downloader) DownloadPlaylistAndFormat(playlistID string, tFormatter TitleFormatter) {

	call := d.YtService.PlaylistItems.List("snippet,contentDetails,status")
	call.MaxResults(50)
	call.PlaylistId(playlistID)
	response, err := call.Do()
	if err != nil {
		fmt.Println(err)
	}
	
	for _, item := range response.Items {
		if item.Status == nil {
			continue
		}
		if item.Status.PrivacyStatus == "public" {
			video := VideoData{
				VideoID: item.ContentDetails.VideoId,
				Title: item.Snippet.Title,
			}
			
			song := tFormatter.FormatTitle(video)
			fmt.Println("Downloading " + item.Snippet.Title)
			d.download(video)
			formatter, ok := tFormatter.(Formatter)
			if !ok {
				formatter = Formatter{}
			}
			formatter.FormatFile(song)
		}
	}

	for response.NextPageToken != ""  {
		call = d.YtService.PlaylistItems.List("snippet,contentDetails,status")
		call.PageToken(response.NextPageToken)
		call.MaxResults(50)
		call.PlaylistId(playlistID)
		response, err = call.Do()
		if err != nil {
			fmt.Println(err)
		}
		
		for _, item := range response.Items {
			
			if item.Status == nil {
				continue
			}
			if item.Status.PrivacyStatus == "public" {
				video := VideoData{
					VideoID: item.ContentDetails.VideoId,
					Title: item.Snippet.Title,
				}
				
				song := tFormatter.FormatTitle(video)
				fmt.Println("Downloading " + item.Snippet.Title)
				d.download(video)
				formatter, ok := tFormatter.(Formatter)
				if !ok {
					formatter = Formatter{}
				}
				formatter.FormatFile(song)
			}
		}
	}
}

func (d *Downloader) download (video VideoData) {
	path := filepath.Join("songs", "error", "%(title)s.%(ext)s")
	cmd := exec.Command("youtube-dl", "-x", "-f", "bestaudio", "-o", path, "https://www.youtube.com/watch?v="+video.VideoID)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		log.Printf(errb.String())
	}
}

func (d *Downloader) DownloadMp3(videoID string) {

	/*
	call := d.YtService.Videos.List("snippet,status")
	call = call.Id(videoID)
	response, err := call.Do()
	if err != nil {
		fmt.Println(err)
	}
	if response.Items[0].Status.PrivacyStatus == "public" {
		songInfo := &SongInfo{
			VideoID: response.Items[0].Id,
			Title:  response.Items[0].Snippet.Title,
		}
		fmt.Println("Downloading " + response.Items[0].Snippet.Title)
	}
	*/
}