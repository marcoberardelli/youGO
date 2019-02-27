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
	"errors"
	"fmt"
	"github.com/marcoberardelli/youGO"
	"github.com/spf13/cobra"
	"os"
)

// videoCmd represents the video command
var videoCmd = &cobra.Command{
	Use:   "video",
	Short: "Download audio from a YouTube video",
	Long: `Download a YouTube video by passing its ID
The downloaded file is saved in the songs folder.`,

	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
		return errors.New("Missing video link/ID")
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		formatter, err := youGO.NewFormatter(" & ", " x ", " ft. ", " feat ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not initialize the program: %v\n", err)
			os.Exit(1)
		}

		downloader, err = youGO.NewDownloader(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not initialize the program: %v\n", err)
			os.Exit(1)
		}

		downloader.DownloadVideoAndFormat(args[0], formatter)

		return nil
	},
	
}

func init() {
	rootCmd.AddCommand(videoCmd)
}
