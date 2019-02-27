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

package cmd

import (
	"fmt"
	"github.com/marcoberardelli/youGO"
	"github.com/spf13/cobra"
	"os"
)

var downloader youGO.Downloader
var formatter youGO.Formatter
var path string


var rootCmd = &cobra.Command{
	Use:   "youGO",
	Short: "Download audio from youtube videos",
	Long: `youGO is a tool for downloading audio from youtube playlists and videos.
It also tries to put the correct metadata tags for title and artists.`,
}

func init() {
	var err error
	defaultPath, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rootCmd.PersistentFlags().StringVarP(&path, "path", "p", defaultPath, "Path of the download folder")	
}


func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}