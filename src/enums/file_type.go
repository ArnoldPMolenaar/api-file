package enums

type FileType string

const (
	Image    FileType = "image"
	Document FileType = "document"
)

func (t FileType) String() string {
	return string(t)
}
