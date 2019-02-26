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
	"github.com/spf13/cobra"
	"fmt"
	"os"
	"github.com/marcoberardelli/youGO"
)

var formatter youGO.Formatter

// playlistCmd represents the playlist command
var playlistCmd = &cobra.Command{
	Use:   "playlist",
	Short: "Download YouTube playlist",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
		  return errors.New("Missing playlist link/ID")
		}
		return nil
	  },

	RunE: func(cmd *cobra.Command, args []string) error {
		
		formatter, err := youGO.NewFormatter(" & ", " x ", " ft. ", " feat ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not initialize the program: %v\n", err)
			os.Exit(1)
		}
		downloader.DownloadPlaylistAndFormat(args[0], formatter)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(playlistCmd)

	

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// playlistCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// playlistCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
