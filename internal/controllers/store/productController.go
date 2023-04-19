package store

import (
	"net/http"

	"github.com/Rafii271/UTS-EAI/config"
	"github.com/Rafii271/UTS-EAI/domain/entity"
	"github.com/Rafii271/UTS-EAI/internal/dto"

	"github.com/gofiber/fiber/v2"
)

func (s *StoreImplementation) CreateProduct(c *fiber.Ctx) error {
	req := dto.ProductRequest{}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	store := c.Locals("store").(entity.StoreToken)

	// check if category name exists
	var category entity.Category

	err := config.DB.Where("name = ?", req.CategoryName).First(&category).Error

	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "Category name does not exist",
		})
	}

	// check if category name is exist in database
	categoryName := entity.Category{}

	err = config.DB.Where("name = ?", req.CategoryName).First(&categoryName).Error

	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "Category name does not exist",
		})
	}

	newProduct := entity.Product{
		Name:         req.Name,
		Price:        req.Price,
		Stock:        req.Stock,
		Description:  req.Description,
		Image:        req.Image,
		CategoryName: req.CategoryName,
		StoreID:      store.StoreID,
	}

	err = config.DB.Create(&newProduct).Error

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	// preload store name
	err = config.DB.Preload("Store").Where("id = ?", newProduct.ID).First(&newProduct).Error

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	// preload category name
	err = config.DB.Preload("Category").Where("id = ?", newProduct.ID).First(&newProduct).Error

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	response := dto.ProductResponse{
		ID:          newProduct.ID,
		Name:        newProduct.Name,
		Price:       newProduct.Price,
		Stock:       newProduct.Stock,
		Description: newProduct.Description,
		Image:       newProduct.Image,
		Category:    newProduct.Category.Name,
		StoreName:   newProduct.Store.Name,
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "success",
		"data":    response,
	})

}

func (s *StoreImplementation) DetailsProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	product := entity.Product{}

	err := config.DB.Preload("Store").Where("id = ?", id).First(&product).Error

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	response := dto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Price:       product.Price,
		Stock:       product.Stock,
		Description: product.Description,
		Image:       product.Image,
		StoreName:   product.Store.Name,
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    response,
	})
}

func (s *StoreImplementation) UpdateProduct(c *fiber.Ctx) error {
	req := dto.ProductUpdateRequest{}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	store := c.Locals("store").(entity.StoreToken)

	// check if category name exists
	var category entity.Category

	err := config.DB.Where("name = ?", req.CategoryName).First(&category).Error

	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "Category name does not exist",
		})
	}

	// check if category name is exist in database
	categoryName := entity.Category{}

	err = config.DB.Where("name = ?", req.CategoryName).First(&categoryName).Error

	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "Category name does not exist",
		})
	}

	id := c.Params("id")

	updateProduct := entity.Product{}

	// find the product with the given id and store.StoreID
	err = config.DB.Where("id = ? AND store_id = ?", id, store.StoreID).First(&updateProduct).Error

	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "Product not found",
		})
	}

	// update the product fields with the new values
	updateProduct.Name = req.Name
	updateProduct.Price = req.Price
	updateProduct.Stock = req.Stock
	updateProduct.Description = req.Description
	updateProduct.Image = req.Image
	updateProduct.CategoryName = req.CategoryName

	// save the changes to the database
	err = config.DB.Save(&updateProduct).Error

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	// preload store and category name
	err = config.DB.Preload("Store").Preload("Category").Where("id = ?", updateProduct.ID).First(&updateProduct).Error

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	response := dto.ProductResponse{
		ID:          updateProduct.ID,
		Name:        updateProduct.Name,
		Price:       updateProduct.Price,
		Stock:       updateProduct.Stock,
		Description: updateProduct.Description,
		Image:       updateProduct.Image,
		Category:    updateProduct.Category.Name,
		StoreName:   updateProduct.Store.Name,
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    response,
	})
}

func (s *StoreImplementation) DeleteProduct(c *fiber.Ctx) error {
    id := c.Params("id")
    store := c.Locals("store").(entity.StoreToken)

    // check if product exists and belongs to the store
    product := entity.Product{}
    err := config.DB.Where("id = ? AND store_id = ?", id, store.StoreID).First(&product).Error
    if err != nil {
        return c.Status(http.StatusNotFound).JSON(fiber.Map{
            "message": "Product not found or does not belong to the store",
        })
    }

    // delete product
    err = config.DB.Delete(&product).Error
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(err.Error())
    }

    return c.Status(http.StatusOK).JSON(fiber.Map{
        "message": "Product deleted successfully",
    })
}

