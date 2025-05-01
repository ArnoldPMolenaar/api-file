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
	"os"
)

// GetDocument method to get a document by its ID.
func GetDocument(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id, err := utils.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Get the document.
	document, err := services.GetDocumentById(id)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err)
	} else if document.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.DocumentExist, "Document does not exist.")
	}

	// Return the document.
	response := responses.Document{}
	response.SetDocument(&document)

	return c.JSON(response)
}

// GetDocumentFile method to get the document file by ID.
func GetDocumentFile(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id, err := utils.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Get the document.
	document, err := services.GetDocumentById(id)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err)
	} else if document.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.DocumentExist, "Document does not exist.")
	}

	// Construct the file path.
	path, err := services.GetPath(&document.Folder.AppStoragePath, document.FolderID)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err)
	}
	filePath := fmt.Sprintf("%s%s.%s", path, document.Name, document.Extension)

	// Send the file as a response.
	return c.SendFile(filePath)
}

// CreateDocument method to create an document.
func CreateDocument(c *fiber.Ctx) error {
	// Parse the request.
	request := requests.CreateDocument{}
	if err := c.BodyParser(&request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate document fields.
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

	// Extract the extension from the document.
	filename, extension, err := upload.GetExtensionFromFilename(request.Name)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.ParseFilename, err)
	}

	// Check if the document is available.
	if available, err := services.IsDocumentAvailable(request.FolderID, filename, extension); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err)
	} else if available {
		return errorutil.Response(c, fiber.StatusConflict, errors.DocumentExist, "Document already exists.")
	}

	// Convert data to bytes.
	mimeType, base64Data, err := upload.GetMimeTypeAndBase64(request.Data)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.ParseBase64, err)
	} else if isValid := upload.IsValidDocument(mimeType); !isValid {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.DocumentTypeInvalid, fmt.Sprintf("Invalid document for %s.", mimeType))
	}
	data, err := upload.Base64ToBytes(base64Data)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.ParseBase64, fmt.Sprintf("Error while decoding bytes. Amount of correct parsed bytes: %d", err))
	}

	// Upload the document.
	fileProgress := responses.FileProgress{}
	fileProgress.SetFileProgress(storagePath.AppName, enums.Document, request.Name, 0.0)

	err = uploadDocument(storagePath, request.FolderID, request.Name, data, &fileProgress)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errors.UploadDocument, err)
	}

	// Create the document.
	document, err := services.CreateDocument(request.FolderID, filename, extension, mimeType, len(data))
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err)
	}

	// Return the document.
	response := responses.Document{}
	response.SetDocument(&document)

	return c.JSON(response)
}

// UpdateDocument method to update a document.
func UpdateDocument(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id, err := utils.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Parse the request.
	request := requests.UpdateDocument{}
	if err := c.BodyParser(&request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate document fields.
	validate := utils.NewValidator()
	if err := validate.Struct(request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, utils.ValidatorErrors(err))
	}

	// Get the document.
	document, err := services.GetDocumentById(id)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err)
	} else if document.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.DocumentExist, "Document does not exist.")
	}

	// Check if the document data has been modified since it was last fetched.
	if request.UpdatedAt.Unix() < document.UpdatedAt.Unix() {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.OutOfSync, "Data is out of sync.")
	}

	// Extract the extension from the document.
	filename, extension, err := upload.GetExtensionFromFilename(request.Name)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.ParseFilename, err)
	}

	// Convert data to bytes.
	mimeType, base64Data, err := upload.GetMimeTypeAndBase64(request.Data)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.ParseBase64, err)
	} else if isValid := upload.IsValidDocument(mimeType); !isValid {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.DocumentTypeInvalid, fmt.Sprintf("Invalid document for %s.", mimeType))
	}
	data, err := upload.Base64ToBytes(base64Data)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.ParseBase64, fmt.Sprintf("Error while decoding bytes. Amount of correct parsed bytes: %d", err))
	}

	// Delete existing document.
	if err := deleteDocument(&document); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errors.DeleteDocument, err)
	}

	// Upload the document.
	fileProgress := responses.FileProgress{}
	fileProgress.SetFileProgress(document.Folder.AppStoragePath.AppName, enums.Document, request.Name, 0.0)

	err = uploadDocument(&document.Folder.AppStoragePath, document.FolderID, request.Name, data, &fileProgress)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errors.UploadDocument, err)
	}

	// Update the document.
	document, err = services.UpdateDocument(&document, filename, extension, mimeType, len(data))
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err)
	}

	// Return the document.
	response := responses.Document{}
	response.SetDocument(&document)

	return c.JSON(response)
}

// DeleteDocument func to delete an document.
func DeleteDocument(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id, err := utils.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Find the document.
	document, err := services.GetDocumentById(id)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if document.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.DocumentExist, "Document does not exist.")
	}

	// Delete the document.
	if err := services.DeleteDocument(&document); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// DeleteDocument func to delete an document.
func DeleteDocumentHard(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id, err := utils.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Find the document.
	document, err := services.GetDocumentById(id)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if document.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.DocumentExist, "Document does not exist.")
	}

	// Delete the document.
	if err := services.DeleteDocument(&document, true); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Delete the document from the storage.
	if err := deleteDocument(&document); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.InternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// RestoreDocument func to restore a document.
func RestoreDocument(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id, err := utils.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Find the document.
	if deleted, err := services.IsDocumentDeleted(id); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !deleted {
		return errorutil.Response(c, fiber.StatusNotFound, errors.DocumentExist, "Document does not exist.")
	}

	// Restore the document.
	if err := services.RestoreDocument(id); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// Upload the document to the storage path.
func uploadDocument(appStoragePath *models.AppStoragePath, folderID uint, filename string, data []byte, fileProgress *responses.FileProgress) error {
	path, err := services.GetPath(appStoragePath, folderID)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	chunks := upload.ChunkBytes(data)
	file, err := os.OpenFile(path+filename, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	var seeker int64
	for i, chunk := range chunks {
		_, err := file.WriteAt(chunk, seeker)
		if err != nil {
			return err
		}
		seeker += int64(len(chunk))
		fileProgress.Progress = float64(i) * 100.0 / float64(len(chunks))
		BroadcastProgress(fileProgress)
	}

	fileProgress.Progress = 100.0
	BroadcastProgress(fileProgress)

	return err
}

// Delete the document from the storage path.
func deleteDocument(document *models.Document) error {
	path, err := services.GetPath(&document.Folder.AppStoragePath, document.FolderID)
	if err != nil {
		return err
	}

	if err := os.Remove(fmt.Sprintf("%s%s.%s", path, document.Name, document.Extension)); err != nil {
		return err
	}

	return nil
}
