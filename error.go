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



type ErrorServiceCreation struct {
	message string
}

func (err *ErrorServiceCreation) Error() string {
	return err.message
}

func NewErrorServiceCreation(message string) *ErrorServiceCreation {
	return &ErrorServiceCreation{message: message}
}




type ErrorWrongPlaylistID struct {
	message string
}

func (err *ErrorWrongPlaylistID) Error() string{
	return err.message
}

func NewErrorWrongPlaylistId(message string) *ErrorWrongPlaylistID {
	return &ErrorWrongPlaylistID{message: message}
}




type ErrorProblematicName struct {
	message string
}

func (err *ErrorProblematicName) Error() string{
	return err.message
}

func NewErrorProblematicName(message string) *ErrorProblematicName {
	return &ErrorProblematicName{message: message}
}