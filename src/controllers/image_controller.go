package controllers

import (
	"api-file/main/src/dto/requests"
	"api-file/main/src/dto/responses"
	"api-file/main/src/enums"
	"api-file/main/src/errors"
	"api-file/main/src/models"
	"api-file/main/src/services"
	upload "api-file/main/src/utils"
	"fmt"
	errorutil "github.com/ArnoldPMolenaar/api-utils/errors"
	"github.com/ArnoldPMolenaar/api-utils/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/h2non/bimg"
	"os"
)

// GetImage method to get the image by ID.
func GetImage(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id, err := utils.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Get the image.
	image, err := services.GetImageById(id, true)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err)
	} else if image.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.ImageExists, "Image does not exist.")
	}

	// Return the image.
	response := responses.Image{}
	response.SetImage(&image)

	return c.JSON(response)
}

// GetImageFile method to get the image file by ID.
func GetImageFile(c *fiber.Ctx) error {
	// Get the ID and size from the URL.
	id, err := utils.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Try to get image from cache.
	var filePath string
	filePath, err = services.GetImageFromCache(id)

	if filePath == "" || err != nil {
		// Get the image.
		image, err := services.GetImage(id)
		if err != nil {
			return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err)
		} else if image.ID == 0 {
			return errorutil.Response(c, fiber.StatusNotFound, errors.ImageExists, "Image does not exist.")
		}

		// Construct the file path.
		path, err := services.GetPath(&image.Folder.AppStoragePath, image.FolderID)
		if err != nil {
			return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err)
		}
		filePath = fmt.Sprintf("%s%s.%s", path, image.Name, image.Extension)
		_ = services.SaveImageToCache(image.ID, filePath)
	}

	// Send the file as a response.
	return c.SendFile(filePath)
}

// GetImageFileSize method to get the image file by ID.
func GetImageFileSize(c *fiber.Ctx) error {
	// Get the ID and size from the URL.
	size := enums.Size(c.Params("size"))
	id, err := utils.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Try to get image from cache.
	var filePath string
	filePath, err = services.GetImageFromCache(id, size.String())

	if filePath == "" || err != nil {
		// Get the image size.
		imageSize, err := services.GetImageSizeById(id, size)
		if err != nil {
			return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err)
		} else if imageSize.ID == 0 {
			return errorutil.Response(c, fiber.StatusNotFound, errors.ImageExists, "Image does not exist.")
		}

		// Construct the file path.
		path, err := services.GetPath(&imageSize.Image.Folder.AppStoragePath, imageSize.Image.FolderID)
		if err != nil {
			return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err)
		}
		filePath := fmt.Sprintf("%s%s-%s.webp", path, imageSize.Image.Name, size)
		_ = services.SaveImageToCache(imageSize.Image.ID, filePath, size.String())
	}

	// Send the file as a response.
	return c.SendFile(filePath)
}

// CreateImage method to create an image.
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

	// Check if the storage path is full.
	if available, err := services.IsStorageSpaceAvailable(request.AppStoragePathID); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err)
	} else if !available {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.StoragePathFull, "Storage path is full.")
	}

	// Extract the extension from the image.
	filename, extension, err := upload.GetExtensionFromFilename(request.Name)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.ParseFilename, err)
	}

	// Check if the image is available.
	if available, err := services.IsImageAvailable(request.FolderID, filename, extension); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err)
	} else if available {
		return errorutil.Response(c, fiber.StatusConflict, errors.ImageExists, "Image already exists.")
	}

	// Convert data to bytes.
	mimeType, base64Data, err := upload.GetMimeTypeAndBase64(request.Data)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.ParseBase64, err)
	} else if isValid := upload.IsValidImage(mimeType); !isValid {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.ImageTypeInvalid, fmt.Sprintf("Invalid image for %s.", mimeType))
	}
	data, err := upload.Base64ToBytes(base64Data)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.ParseBase64, fmt.Sprintf("Error while decoding bytes. Amount of correct parsed bytes: %d", err))
	}

	// Upload the image.
	progress := 100.0
	if !request.IsNotResizable {
		progress = 100.0 / 7
	}

	fileProgress := responses.FileProgress{}
	fileProgress.SetFileProgress(storagePath.AppName, enums.Image, request.Name, 0.0)

	width, height, err := uploadImage(storagePath, request.FolderID, request.Name, data, progress, &fileProgress)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errors.UploadImage, err)
	}

	// Create web size images.
	var imageSizes []models.ImageSize
	if !request.IsNotResizable {
		if imageSizes, err = convertAndUploadImages(storagePath, request.FolderID, filename, data, request.Quality, progress, &fileProgress); err != nil {
			return errorutil.Response(c, fiber.StatusInternalServerError, errors.ConvertImage, err)
		}
	}

	// Create the image.
	image, err := services.CreateImage(request.FolderID, filename, extension, mimeType, len(data), width, height, request.Description, imageSizes)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err)
	}

	// Return the image.
	response := responses.Image{}
	response.SetImage(&image)

	return c.JSON(response)
}

// UpdateImage method to update the image fields like description.
func UpdateImage(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id, err := utils.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Parse the request.
	request := requests.UpdateImage{}
	if err := c.BodyParser(&request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate image fields.
	validate := utils.NewValidator()
	if err := validate.Struct(request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, utils.ValidatorErrors(err))
	}

	// Check if the image exists.
	image, err := services.GetImageById(id, true)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err)
	} else if image.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.ImageExists, "Image does not exist.")
	}

	// Check if the image data has been modified since it was last fetched.
	if request.UpdatedAt.Unix() < image.UpdatedAt.Unix() {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.OutOfSync, "Data is out of sync.")
	}

	var filename *string
	var extension *string
	var mimeType *string
	var size *int
	var width *int
	var height *int
	var imageSizes *[]models.ImageSize

	if request.Name != nil && request.Data != nil {
		// Delete the old image.
		if err := deleteImage(&image); err != nil {
			return errorutil.Response(c, fiber.StatusInternalServerError, errors.DeleteImage, err)
		}

		// Extract the extension from the image.
		parsedFilename, parsedExtension, err := upload.GetExtensionFromFilename(*request.Name)
		if err != nil {
			return errorutil.Response(c, fiber.StatusBadRequest, errors.ParseFilename, err)
		}
		filename = &parsedFilename
		extension = &parsedExtension

		// Convert data to bytes.
		mimeType, base64Data, err := upload.GetMimeTypeAndBase64(*request.Data)
		if err != nil {
			return errorutil.Response(c, fiber.StatusBadRequest, errors.ParseBase64, err)
		} else if isValid := upload.IsValidImage(mimeType); !isValid {
			return errorutil.Response(c, fiber.StatusBadRequest, errors.ImageTypeInvalid, fmt.Sprintf("Invalid image for %s.", mimeType))
		}
		data, err := upload.Base64ToBytes(base64Data)
		if err != nil {
			return errorutil.Response(c, fiber.StatusBadRequest, errors.ParseBase64, fmt.Sprintf("Error while decoding bytes. Amount of correct parsed bytes: %d", err))
		}
		dataLen := len(data)
		size = &dataLen

		// Upload the image.
		progress := 100.0
		if request.IsNotResizable == nil || !*request.IsNotResizable {
			progress = 100.0 / 7
		}

		fileProgress := responses.FileProgress{}
		fileProgress.SetFileProgress(image.Folder.AppStoragePath.AppName, enums.Image, *request.Name, 0.0)

		imageWidth, imageHeight, err := uploadImage(&image.Folder.AppStoragePath, image.FolderID, *request.Name, data, progress, &fileProgress)
		if err != nil {
			return errorutil.Response(c, fiber.StatusInternalServerError, errors.UploadImage, err)
		}
		width = &imageWidth
		height = &imageHeight

		// Create web size images.
		if request.IsNotResizable == nil || !*request.IsNotResizable {
			if createdImageSizes, err := convertAndUploadImages(&image.Folder.AppStoragePath, image.FolderID, *filename, data, *request.Quality, progress, &fileProgress); err != nil {
				return errorutil.Response(c, fiber.StatusInternalServerError, errors.ConvertImage, err)
			} else {
				imageSizes = &createdImageSizes
			}
		}
	}

	// Update the image.
	image, err = services.UpdateImage(&image, filename, extension, mimeType, size, width, height, request.Description, imageSizes)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err)
	}

	// Return the image.
	response := responses.Image{}
	response.SetImage(&image)

	return c.JSON(response)
}

// DeleteImage func to delete an image.
func DeleteImage(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id, err := utils.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Find the image.
	image, err := services.GetImageById(id, true)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if image.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.ImageExists, "Image does not exist.")
	}

	// Delete the image.
	if err := services.DeleteImage(&image); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// DeleteImageHard func to delete an image for ever.
func DeleteImageHard(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id, err := utils.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Find the image.
	image, err := services.GetImageById(id, true)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if image.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.ImageExists, "Image does not exist.")
	}

	// Delete the image.
	if err := services.DeleteImage(&image, true); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Delete the image from the storage path.
	if err := deleteImage(&image); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.InternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// RestoreImage func to restore a image.
func RestoreImage(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id, err := utils.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Find the image.
	if deleted, err := services.IsImageDeleted(id); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !deleted {
		return errorutil.Response(c, fiber.StatusNotFound, errors.ImageExists, "Image does not exist.")
	}

	// Restore the image.
	if err := services.RestoreImage(id); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// Upload the image to the storage path.
func uploadImage(appStoragePath *models.AppStoragePath, folderID uint, filename string, data []byte, progress float64, fileProgress *responses.FileProgress) (width, height int, err error) {
	path, err := services.GetPath(appStoragePath, folderID)
	if err != nil {
		return 0, 0, err
	}

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return 0, 0, err
	}

	img := bimg.NewImage(data)
	size, err := img.Size()
	if err != nil {
		return 0, 0, err
	}

	width = size.Width
	height = size.Height
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
		percentage := float64(i) * 100.0 / float64(len(chunks))
		if progress == 100.0 {
			fileProgress.Progress = percentage
			BroadcastProgress(fileProgress)
		} else {
			fileProgress.Progress = progress * percentage / 100.0
			BroadcastProgress(fileProgress)
		}
	}

	fileProgress.Progress = progress
	BroadcastProgress(fileProgress)

	return width, height, err
}

// Convert and upload the images to the storage path.
func convertAndUploadImages(appStoragePath *models.AppStoragePath, folderID uint, filename string, data []byte, quality int, progress float64, fileProgress *responses.FileProgress) ([]models.ImageSize, error) {
	var imageSizes []models.ImageSize
	sizes := map[enums.Size]int{
		enums.XS:  600,
		enums.SM:  960,
		enums.MD:  1280,
		enums.LG:  1920,
		enums.XL:  2560,
		enums.XXL: 3840,
	}
	path, err := services.GetPath(appStoragePath, folderID)
	if err != nil {
		return imageSizes, err
	}
	originalSize, err := bimg.NewImage(data).Size()
	if err != nil {
		return imageSizes, err
	}

	var amountOfImagesToCreate int8
	for _, width := range sizes {
		if originalSize.Width > width {
			amountOfImagesToCreate++
		}
	}

	calculatedProgress := (100.0 - progress) / float64(amountOfImagesToCreate)
	var currentImage int8
	for size, width := range sizes {
		currentImage++
		filenameSize := fmt.Sprintf("%s-%s.webp", filename, size)

		if originalSize.Width <= width {
			continue
		}

		newHeight := originalSize.Height * width / originalSize.Width
		resized, err := bimg.NewImage(data).Resize(width, newHeight)
		if err != nil {
			return imageSizes, err
		}

		s, err := bimg.NewImage(resized).Size()
		if err != nil {
			return imageSizes, err
		}

		converted, err := bimg.NewImage(resized).Convert(bimg.WEBP)
		if err != nil {
			return imageSizes, err
		}

		processed, err := bimg.NewImage(converted).Process(bimg.Options{Quality: quality})
		if err != nil {
			return imageSizes, err
		}

		err = bimg.Write(path+filenameSize, processed)
		if err != nil {
			return imageSizes, err
		}

		imageSizes = append(imageSizes, models.ImageSize{
			Size:   size,
			Width:  s.Width,
			Height: s.Height,
		})

		fileProgress.Progress = progress + calculatedProgress*float64(currentImage)
		BroadcastProgress(fileProgress)
	}

	return imageSizes, nil
}

// Delete the image from the storage path.
func deleteImage(image *models.Image) error {
	path, err := services.GetPath(&image.Folder.AppStoragePath, image.FolderID)
	if err != nil {
		return err
	}

	if err := os.Remove(fmt.Sprintf("%s%s.%s", path, image.Name, image.Extension)); err != nil {
		return err
	}

	for i := range image.ImageSizes {
		if err := os.Remove(fmt.Sprintf("%s%s-%s.webp", path, image.Name, image.ImageSizes[i].Size.String())); err != nil {
			return err
		}
	}

	return nil
}
