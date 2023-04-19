package routes

import (
	"github.com/Rafii271/UTS-EAI/internal/controllers/global"
	"github.com/Rafii271/UTS-EAI/internal/controllers/store"
	"github.com/Rafii271/UTS-EAI/internal/controllers/user"
	"github.com/Rafii271/UTS-EAI/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	userImplementation := user.UserImplementation{}
	storeImplementation := store.StoreImplementation{}
	globalImplementation := global.GlobalImplementation{}

	api := app.Group("/api")

	// =================== AUTH ===================
	register := api.Group("/register")
	register.Post("", userImplementation.Register)

	login := api.Group("/login")
	login.Post("", userImplementation.Login)

	users := api.Group("/user").Use(middleware.AuthUser(middleware.Config{
		Unauthorized: func(c *fiber.Ctx) error {
			return c.Status(401).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		},
	}))
	users.Post("/logout", userImplementation.Logout)
	users.Get("/profile", userImplementation.ShowProfile)
	// =================== AUTH ===================

	// ==================== Global ====================
	products := api.Group("/products")
	products.Get("/detail", globalImplementation.DetailsProduct)
	products.Get("", globalImplementation.ShowAllProduct)
	categories := api.Group("/categories")
	categories.Get("", globalImplementation.GetAllCategories)
	// ==================== Global ====================

	// =================== STORE ===================
	storeAPI := api.Group("/store").Use(middleware.AuthUser(middleware.Config{
		Unauthorized: func(c *fiber.Ctx) error {
			return c.Status(401).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		},
	}))
	storeAPI.Post("/create", storeImplementation.RegisterStore)

	Store := api.Group("/login-store")
	Store.Post("", storeImplementation.LoginStore)

	// =================== Product ===================
	productAPI := api.Group("/product").Use(middleware.AuthStore(middleware.Config{
		Unauthorized: func(c *fiber.Ctx) error {
			return c.Status(401).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		},
	}))
	productAPI.Post("/create", storeImplementation.CreateProduct)
	productAPI.Get("/shows", storeImplementation.GetStoreProducts)
	productAPI.Put("/update/:id", storeImplementation.UpdateProduct)
	productAPI.Delete("/delete/:id", storeImplementation.DeleteProduct)

}
