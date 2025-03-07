package errors

// Define error codes as constants.
const (
	AppExists            = "appExists"
	StoragePathExists    = "storagePathExists"
	StoragePathAvailable = "storagePathAvailable"
	StoragePathFull      = "storagePathFull"
	FolderExists         = "folderExists"
	ImageExists          = "imageExists"
	ImageTypeInvalid     = "imageTypeInvalid"
	ParseBase64          = "parseBase64"
	ParseFilename        = "parseFilename"
	DeleteImage          = "deleteImage"
	UploadImage          = "uploadImage"
	ConvertImage         = "convertImage"
	CodeInvalid          = "codeInvalid"
	CodeExists           = "codeExists"
	DocumentExist        = "documentExist"
	DocumentTypeInvalid  = "documentTypeInvalid"
	UploadDocument       = "uploadDocument"
	DeleteDocument       = "deleteDocument"
	// Add more error codes as needed.
)
