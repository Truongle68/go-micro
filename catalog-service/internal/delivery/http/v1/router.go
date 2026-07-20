package v1

import (
	"catalog-service/internal/delivery/http/middleware"
	"catalog-service/internal/usecase"

	"github.com/TruongLe68/go-micro/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type V1 struct {
	product  usecase.ProductUsecase
	category usecase.CategoryUsecase
	l        logger.Interface
	v        *validator.Validate
}

type Dependencies struct {
	Product   usecase.ProductUsecase
	Category  usecase.CategoryUsecase
	Logger    logger.Interface
	JWTSecret string
}

func NewRoutes(apiV1Group *gin.RouterGroup, deps Dependencies) {
	r := &V1{
		product:  deps.Product,
		category: deps.Category,
		l:        deps.Logger,
		v:        validator.New(),
	}

	authMid := middleware.AuthMiddleware(deps.JWTSecret)
	adminMid := middleware.RoleMiddleware("admin")
	userOrAdminMid := middleware.RoleMiddleware("user", "admin")

	products := apiV1Group.Group("/products")
	{
		products.POST("/create-product", authMid, adminMid, r.CreateProduct)
		products.GET("/get-product/:id", authMid, userOrAdminMid, r.GetProduct)
		products.POST("/update-product/:id", authMid, adminMid, r.UpdateProduct)
		products.POST("/delete-product/:id", authMid, adminMid, r.DeleteProduct)
		products.GET("/get-product-list", authMid, userOrAdminMid, r.ListProducts)
		products.GET("/get-products-by-category/:id", authMid, userOrAdminMid, r.GetProductsByCategory)
	}

	categories := apiV1Group.Group("/categories")
	{
		categories.POST("/create-category", authMid, adminMid, r.CreateCategory)
		categories.GET("/get-category/:id", authMid, userOrAdminMid, r.GetCategory)
		categories.GET("/get-category-children/:id", authMid, userOrAdminMid, r.GetCategoryChildren)
		categories.POST("/update-category/:id", authMid, adminMid, r.UpdateCategory)
		categories.POST("/delete-category/:id", authMid, adminMid, r.DeleteCategory)
		categories.GET("/get-category-list", authMid, userOrAdminMid, r.ListCategories)
	}
}
