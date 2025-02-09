package responses

import (
	"api-file/main/src/models"
	"time"
)

type FolderPreload struct {
	ID               uint       `json:"id"`
	AppStoragePathID uint       `json:"app_storage_path_id"`
	Name             string     `json:"name"`
	Color            string     `json:"color"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
	Folders          []Folder   `json:"folders"`
	Images           []Image    `json:"images"`
	Documents        []Document `json:"documents"`
}

func (f *FolderPreload) SetFolderPreload(folder *models.Folder, folders []*models.Folder) {
	f.ID = folder.ID
	f.AppStoragePathID = folder.AppStoragePathID
	f.Name = folder.Name
	f.Color = folder.Color
	f.CreatedAt = folder.CreatedAt
	f.UpdatedAt = folder.UpdatedAt

	f.Folders = make([]Folder, len(folders))
	for i := range folders {
		f.Folders[i] = Folder{}
		f.Folders[i].SetFolder(folders[i])
	}

	f.Images = make([]Image, len(folder.Images))
	for i := range folder.Images {
		f.Images[i] = Image{}
		f.Images[i].SetImage(&folder.Images[i])
	}

	f.Documents = make([]Document, len(folder.Documents))
	for i := range folder.Documents {
		f.Documents[i] = Document{}
		f.Documents[i].SetDocument(&folder.Documents[i])
	}
}
