package responses

import "api-file/main/src/enums"

// FileProgress struct for file progress response.
// Used to send file progress to client with a websocket.
type FileProgress struct {
	App      string         `json:"app"`
	Type     enums.FileType `json:"type"`
	Filename string         `json:"filename"`
	Progress float64        `json:"progress"`
}

// SetFileProgress sets the file progress response.
func (f *FileProgress) SetFileProgress(app string, fileType enums.FileType, filename string, progress float64) {
	f.App = app
	f.Type = fileType
	f.Filename = filename
	f.Progress = progress
}
