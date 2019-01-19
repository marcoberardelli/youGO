// Copyright 2019 Marco Berardelli
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

package main

import(
	"os"
	"fmt"
	"log"
)

func printUsage() {

}

func main() {

	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		printUsage()
		return
	}

	if len(os.Args) < 3 {
		fmt.Println("Invalid arguments")
		printUsage()
		return
	}

	d, err := NewDownloader()
	if err != nil {
		log.Fatal(err)
	}
	
	if os.Args[1] == "-p" || os.Args[1] == "--playlist" {
		d.DownloadPlaylist(os.Args[2])
	} else if os.Args[1] == "-v" || os.Args[1] == "-video" {
		d.DownloadMp3(os.Args[2])
	} else {
		printUsage()
	}
}
