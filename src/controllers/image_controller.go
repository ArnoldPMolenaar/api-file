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
		return errorutil.Response(c, fiber.StatusBadRequest, errors.ImageTypeInvalid, "Invalid image.")
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

	width, height, err := uploadImage(storagePath, request.FolderID, request.Name, data, progress)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errors.UploadImage, err)
	}

	// Create web size images.
	var imageSizes []models.ImageSize
	if !request.IsNotResizable {
		if imageSizes, err = convertAndUploadImages(storagePath, request.FolderID, filename, data, request.Quality, progress); err != nil {
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

// Upload the image to the storage path.
func uploadImage(appStoragePath *models.AppStoragePath, folderID uint, filename string, data []byte, progress float64) (int, int, error) {
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
		percentage := float64(i) * 100.0 / float64(len(chunks))
		if progress == 100.0 {
			log.Debugf("Processed: %.2f", percentage)
		} else {
			log.Debugf("Processed: %.2f", progress*percentage/100.0)
		}
	}

	// TODO: Make here one last process message with 100% or `progress` value.
	log.Debugf("Processed: %.2f", progress)

	return width, height, err
}

func convertAndUploadImages(appStoragePath *models.AppStoragePath, folderID uint, filename string, data []byte, quality int, progress float64) ([]models.ImageSize, error) {
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

		// TODO: Write process to websocket connection.
		log.Debugf("Processed: %.2f", progress+calculatedProgress*float64(currentImage))
	}

	return imageSizes, nil
}
