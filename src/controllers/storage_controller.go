package controllers

import (
	"api-file/main/src/database"
	"api-file/main/src/dto/requests"
	"api-file/main/src/dto/responses"
	"api-file/main/src/errors"
	"api-file/main/src/models"
	"api-file/main/src/services"
	errorutil "github.com/ArnoldPMolenaar/api-utils/errors"
	"github.com/ArnoldPMolenaar/api-utils/pagination"
	"github.com/ArnoldPMolenaar/api-utils/utils"
	"github.com/gofiber/fiber/v2"
)

// GetStoragePaths func to get all storage paths for the app.
func GetStoragePaths(c *fiber.Ctx) error {
	storagePaths := make([]models.AppStoragePath, 0)
	values := c.Request().URI().QueryArgs()
	allowedColumns := map[string]bool{
		"id":   true,
		"app":  true,
		"path": true,
	}

	queryFunc := pagination.Query(values, allowedColumns)
	sortFunc := pagination.Sort(values, allowedColumns)
	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}
	limit := c.QueryInt("limit", 10)
	if limit < 1 {
		limit = 10
	}
	offset := pagination.Offset(page, limit)

	db := database.Pg.Scopes(queryFunc, sortFunc).
		Limit(limit).
		Offset(offset).
		Find(&storagePaths)
	if db.Error != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, db.Error.Error())
	}

	total := int64(0)
	database.Pg.Scopes(queryFunc).
		Model(&models.AppStoragePath{}).
		Count(&total)
	pageCount := pagination.Count(int(total), limit)

	paginationModel := pagination.CreatePaginationModel(limit, page, pageCount, int(total), toStoragePathPagination(storagePaths))

	return c.Status(fiber.StatusOK).JSON(paginationModel)
}

// GetStoragePath func to get a storage path for the app.
func GetStoragePath(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id, err := utils.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Find the storage path.
	storagePath, err := services.GetStoragePath(id)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if storagePath == nil || storagePath.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.StoragePathExists, "Storage path does not exist.")
	}

	// Return the storage path.
	response := responses.AppStoragePath{}
	response.SetAppStoragePath(storagePath)

	return c.JSON(response)
}

// CreateStoragePath func to create a storage path for the app.
func CreateStoragePath(c *fiber.Ctx) error {
	// Parse the request.
	request := requests.CreateAppStoragePath{}
	if err := c.BodyParser(&request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate storage fields.
	validate := utils.NewValidator()
	if err := validate.Struct(request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, utils.ValidatorErrors(err))
	}

	// Check if app exists.
	if available, err := services.IsAppAvailable(request.App); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !available {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.AppExists, "AppName does not exist.")
	}

	// Check if storage path exists.
	if available, err := services.IsStorageAvailable(request.App, request.Path); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if available {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.StoragePathAvailable, "Storage path already available.")
	}

	// Create the storage path.
	storagePath, err := services.CreateStoragePath(request.App, request.Path)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Return the storage path.
	response := responses.AppStoragePath{}
	response.SetAppStoragePath(storagePath)

	return c.JSON(response)
}

// UpdateStoragePath func to update a storage path for the app.
func UpdateStoragePath(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id, err := utils.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Parse the request.
	request := requests.UpdateAppStoragePath{}
	if err := c.BodyParser(&request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate storage fields.
	validate := utils.NewValidator()
	if err := validate.Struct(request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, utils.ValidatorErrors(err))
	}

	// Check if app exists.
	if available, err := services.IsAppAvailable(request.App); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !available {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.AppExists, "AppName does not exist.")
	}

	// Check if storage path exists.
	if available, err := services.IsStorageAvailable(request.App, request.Path); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if available {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.StoragePathAvailable, "Storage path already available.")
	}

	// Find the storage path.
	storagePath, err := services.GetStoragePath(id)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if storagePath == nil || storagePath.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.StoragePathExists, "Storage path does not exist.")
	}

	// Update the storage path.
	storagePath, err = services.UpdateStoragePath(storagePath, request.App, request.Path)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Return the storage path.
	response := responses.AppStoragePath{}
	response.SetAppStoragePath(storagePath)

	return c.JSON(response)
}

// toStoragePathPagination func to convert the storage paths to a response struct.
func toStoragePathPagination(storagePaths []models.AppStoragePath) []responses.AppStoragePath {
	result := make([]responses.AppStoragePath, len(storagePaths))

	for i := range storagePaths {
		response := responses.AppStoragePath{}
		response.SetAppStoragePath(&storagePaths[i])
		result[i] = response
	}

	return result
}
