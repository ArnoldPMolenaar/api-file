package controllers

import (
	"api-file/main/src/dto/requests"
	"api-file/main/src/dto/responses"
	"api-file/main/src/errors"
	"api-file/main/src/models"
	"api-file/main/src/services"
	errorutil "github.com/ArnoldPMolenaar/api-utils/errors"
	"github.com/ArnoldPMolenaar/api-utils/utils"
	"github.com/gofiber/fiber/v2"
)

// GetFolder func to get a folder.
func GetFolder(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id, err := utils.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Find the folder.
	folder, folders, err := services.GetFolder(id, true)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if folder == nil || folder.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.FolderExists, "Folder does not exist.")
	}

	// Return the storage path.
	response := responses.FolderPreload{}
	response.SetFolderPreload(folder, folders)

	return c.JSON(response)
}

// CreateFolder func to create a folder.
func CreateFolder(c *fiber.Ctx) error {
	var err error

	// Parse the request.
	request := requests.CreateFolder{}
	if err := c.BodyParser(&request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate folder fields.
	validate := utils.NewValidator()
	if err := validate.Struct(request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, utils.ValidatorErrors(err))
	}

	// Check if the folder is available.
	available := false
	if request.ParentFolderID != nil {
		available, err = services.IsFolderAvailable(request.AppStoragePathID, request.Name, *request.ParentFolderID)
	} else {
		available, err = services.IsFolderAvailable(request.AppStoragePathID, request.Name)
	}
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if available {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.FolderExists, "Folder already exists.")
	}

	// Create the folder.
	var folder *models.Folder
	if request.ParentFolderID != nil {
		folder, err = services.CreateFolder(request.AppStoragePathID, request.Name, request.Color, *request.ParentFolderID)
	} else {
		folder, err = services.CreateFolder(request.AppStoragePathID, request.Name, request.Color)
	}
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Return the folder.
	response := responses.Folder{}
	response.SetFolder(folder)

	return c.JSON(response)
}

// UpdateFolder func to update a folder.
func UpdateFolder(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id, err := utils.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Parse the request.
	request := requests.UpdateFolder{}
	if err := c.BodyParser(&request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate storage fields.
	validate := utils.NewValidator()
	if err := validate.Struct(request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, utils.ValidatorErrors(err))
	}

	// Find the folder.
	folder, _, err := services.GetFolder(id)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if folder == nil || folder.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.FolderExists, "Folder does not exist.")
	}

	// Check if the folder is available.
	available := false
	if request.ParentFolderID != nil {
		available, err = services.IsFolderAvailable(request.AppStoragePathID, request.Name, *request.ParentFolderID)
	} else {
		available, err = services.IsFolderAvailable(request.AppStoragePathID, request.Name)
	}
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if available {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.FolderExists, "Folder already exists.")
	}

	// Update the folder.
	if folder, err = services.UpdateFolder(folder, request.Name, request.Color); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Return the folder.
	response := responses.Folder{}
	response.SetFolder(folder)

	return c.JSON(response)
}

// DeleteFolder func to delete a folder.
func DeleteFolder(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id, err := utils.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Find the folder.
	folder, _, err := services.GetFolder(id)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if folder == nil || folder.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.FolderExists, "Folder does not exist.")
	}

	// Delete the folder.
	if err := services.DeleteFolder(folder); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// RestoreFolder func to restore a folder.
func RestoreFolder(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id, err := utils.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Find the folder.
	if deleted, err := services.IsFolderDeleted(id); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !deleted {
		return errorutil.Response(c, fiber.StatusNotFound, errors.FolderExists, "Folder does not exist.")
	}

	// Restore the folder.
	if err := services.RestoreFolder(id); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
