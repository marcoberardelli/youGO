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
	"github.com/BrianAllred/goydl"
	"runtime"
	"fmt"
	"path/filepath"
	"os"
	"github.com/mitchellh/go-homedir"
	"sync"
)

var wg sync.WaitGroup

type Downloader struct{
	PathFolder string
	PathErrorFolder string
	YtService *youtube.Service
	FormatterUtil *Formatter
}


func NewDownloader() (*Downloader, error) {

	client := &http.Client{
		Transport: &transport.APIKey{Key: YouTubeAPIKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		return nil, NewErrorServiceCreation("Impossible to create youtube service")
	}

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
	formatter, err := NewFormatter(" & ", " x ", " ft. ", " feat ")
	if err != nil {
		return nil, err
	}
	downloader := &Downloader{
		PathFolder: dir,
		PathErrorFolder: errorDir,
		YtService: service,
		FormatterUtil: formatter,
	}

	return downloader, nil
}


func(d *Downloader) DownloadPlaylist(playlistID string) {

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
			songInfo, err := d.FormatterUtil.FormatTitle(item.Snippet.Title)
			songInfo.VideoID = item.ContentDetails.VideoId
			if err != nil {
				switch err.(type) {
				case *ErrorProblematicName:
					songInfo.Path = d.PathErrorFolder + filepath.FromSlash("/"+songInfo.Title) + ".mp3"
				default:
					fmt.Println(err)
				}
			} else {
				songInfo.Path = d.PathFolder + filepath.FromSlash("/"+songInfo.Title) + ".mp3"
			}
			fmt.Println("Downloading " + item.Snippet.Title)
			d.download(songInfo)
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
				songInfo, err := d.FormatterUtil.FormatTitle(item.Snippet.Title)
				songInfo.VideoID = item.ContentDetails.VideoId
				if err != nil {
					switch err.(type) {
					case *ErrorProblematicName:
						songInfo.Path = d.PathErrorFolder + filepath.FromSlash("/"+songInfo.Title) + ".mp3"
					default:
						fmt.Println(err)
					}
				} else {
					songInfo.Path = d.PathFolder + filepath.FromSlash("/"+songInfo.Title) + ".mp3"
				}
				fmt.Println("Downloading " + item.Snippet.Title)
				d.download(songInfo)
				err = d.download(songInfo);
				for err == nil {
					err = d.download(songInfo)
				}
			}
		}
	}
	wg.Wait()
}


func (d *Downloader) DownloadMp3(videoID string) {

	call := d.YtService.Videos.List("snippet,status")
	call = call.Id(videoID)
	response, err := call.Do()
	if err != nil {
		fmt.Println(err)
	}
	if response.Items[0].Status.PrivacyStatus == "public" {
		songInfo, err := d.FormatterUtil.FormatTitle(response.Items[0].Snippet.Title)
		songInfo.VideoID = videoID
		if err != nil {
			switch err.(type) {
			case *ErrorProblematicName:
				songInfo.Path = d.PathErrorFolder + filepath.FromSlash("/"+songInfo.Title) + ".mp3"
			default:
				fmt.Println(err)
			}
		} else {
			songInfo.Path = d.PathFolder + filepath.FromSlash("/"+songInfo.Title) + ".mp3"
		}
		fmt.Println("Downloading " + response.Items[0].Snippet.Title)

		err = d.download(songInfo);
		for err == nil {
			err = d.download(songInfo)
		}
	}
}


func (d *Downloader) download(songInfo *SongInfo) error {
	y := goydl.NewYoutubeDl()
	
	
	//filename := songInfo.Path + filepath.FromSlash("/"+songInfo.Title) + ".mp3"
	y.Options.Output.Value = songInfo.Path
	y.Options.ExtractAudio.Value = true
	y.Options.Format.Value = "140"
	y.Options.AudioFormat.Value = "mp3"
	command, err := y.Download("https://www.youtube.com/watch?v=" + songInfo.VideoID)
	if err != nil {
		log.Fatal(err)
	}
	command.Wait()
	//wg.Add(1)
	return d.FormatterUtil.FormatMp3(songInfo)
}