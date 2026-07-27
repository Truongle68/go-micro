package v1

import (
	"catalog-service/internal/delivery/http/middleware"
	"catalog-service/internal/domain"
	"catalog-service/internal/usecase"
	"catalog-service/pkg/redis"

	"github.com/GoProOrg/core-go-pkg/jwtmanager"
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
	Product  usecase.ProductUsecase
	Category usecase.CategoryUsecase
	Logger   logger.Interface
	Verifier jwtmanager.JWTManager
	Cache    redis.IdentityCacher
}

func NewRoutes(apiV1Group *gin.RouterGroup, deps Dependencies) {
	r := &V1{
		product:  deps.Product,
		category: deps.Category,
		l:        deps.Logger,
		v:        validator.New(),
	}

	authMid := middleware.Auth(deps.Verifier, deps.Cache)
	adminMid := middleware.Role(string(domain.UserRoleAdmin))

	products := apiV1Group.Group("/products")
	{
		products.POST("", authMid, adminMid, r.createProduct)
		products.GET("/search", r.searchProducts)
		products.GET("/:id", r.getProduct)
		products.PUT("/:id", authMid, adminMid, r.updateProduct)
		products.DELETE("/:id", authMid, adminMid, r.deleteProduct)
		products.GET("", authMid, r.listProducts)
	}

	categories := apiV1Group.Group("/categories")
	{
		categories.POST("", authMid, adminMid, r.createCategory)
		categories.GET("/:id", r.getCategory)
		categories.GET("/:id/children", r.getCategoryChildren)
		categories.PUT("/:id", authMid, adminMid, r.updateCategory)
		categories.DELETE("/:id", authMid, adminMid, r.deleteCategory)
		categories.GET("", r.listCategories)
		categories.GET("/:id/products", r.getProductsByCategory)
	}
}
