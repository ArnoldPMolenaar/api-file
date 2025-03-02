package errors

// Define error codes as constants.
const (
	AppExists            = "appExists"
	StoragePathExists    = "storagePathExists"
	StoragePathAvailable = "storagePathAvailable"
	FolderExists         = "folderExists"
	ImageExists          = "imageExists"
	ImageTypeInvalid     = "imageTypeInvalid"
	ParseBase64          = "parseBase64"
	ParseFilename        = "parseFilename"
	UploadImage          = "uploadImage"
	ConvertImage         = "convertImage"
	CodeInvalid          = "codeInvalid"
	CodeExists           = "codeExists"
	IORead               = "ioRead"
	DocumentExist        = "documentExist"
	DocumentTypeInvalid  = "documentTypeInvalid"
	UploadDocument       = "uploadDocument"
	// Add more error codes as needed.
)
