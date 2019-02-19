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
	"os"
	"github.com/spf13/cobra"
	"fmt"

	"github.com/marcoberardelli/youGO"
)

var downloader *youGO.Downloader

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "youGO",
	Short: "Download audio from youtube videos",
	Long: `youGO is a tool for downloading audio as mp3 from youtube playlists and videos.
	
It also tries to put the correct title and artist mp3 tags`,
}

func init() {
	var err error
	downloader, err = youGO.NewDownloader()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not initialize the program: %v\n", err)
		os.Exit(1)
	}

	
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}