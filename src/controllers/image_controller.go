package controllers

import (
	"api-file/main/src/dto/requests"
	"api-file/main/src/dto/responses"
	"api-file/main/src/errors"
	"api-file/main/src/models"
	"api-file/main/src/services"
	upload "api-file/main/src/utils"
	errorutil "github.com/ArnoldPMolenaar/api-utils/errors"
	"github.com/ArnoldPMolenaar/api-utils/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/h2non/bimg"
	"os"
)

func CreateImage(c *fiber.Ctx) error {
	// Parse the request.
	request := requests.CreateImage{}
	if err := c.BodyParser(&request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate image fields.
	validate := utils.NewValidator()
	if err := validate.Struct(request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, utils.ValidatorErrors(err))
	}

	// Check if the storage path exists.
	storagePath, err := services.GetStoragePath(request.AppStoragePathID)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err)
	} else if storagePath.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.StoragePathExists, "Storage path does not exist.")
	}

	// Extract the extension from the image.
	filename, extension, err := upload.GetExtensionFromFilename(request.Name)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.ParseFilename, err)
	}

	// Check if the image is available.
	if available, err := services.IsImageAvailable(request.AppStoragePathID, filename, extension); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err)
	} else if available {
		return errorutil.Response(c, fiber.StatusConflict, errors.ImageExists, "Image already")
	}

	// Convert data to bytes.
	mimeType, base64Data, err := upload.GetMimeTypeAndBase64(request.Data)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.ParseBase64, err)
	} else if isValid := upload.IsValidImage(mimeType); !isValid {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.ImageTypeInvalid, "Invalid image.")
	}
	data, err := upload.Base64ToBytes(base64Data)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.ParseBase64, err)
	}

	// Upload the image.
	width, height, err := uploadImage(storagePath, request.FolderID, filename, data)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errors.UploadImage, err)
	}

	// Create the image.
	image, err := services.CreateImage(request.FolderID, filename, extension, mimeType, len(data), width, height, request.Description)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err)
	}

	// Return the image.
	response := responses.Image{}
	response.SetImage(&image)

	return c.JSON(response)
}

func uploadImage(appStoragePath *models.AppStoragePath, folderID uint, filename string, data []byte) (int, int, error) {
	path := appStoragePath.Path
	if folderPath, err := services.GetFolderPath(appStoragePath.ID, folderID); err != nil {
		return 0, 0, err
	} else {
		path += folderPath
	}

	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return 0, 0, err
	}

	img := bimg.NewImage(data)
	size, err := img.Size()
	if err != nil {
		return 0, 0, err
	}

	width := size.Width
	height := size.Height
	chunks := upload.ChunkBytes(data)

	file, err := os.OpenFile(path+filename, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	var seeker int64
	for i, chunk := range chunks {
		_, err := file.WriteAt(chunk, seeker)
		if err != nil {
			return 0, 0, err
		}
		seeker += int64(len(chunk))
		// TODO: Write process to websocket connection.
		log.Debug("Processed: %d", (i*100)/len(chunks))
	}

	return width, height, err
}
