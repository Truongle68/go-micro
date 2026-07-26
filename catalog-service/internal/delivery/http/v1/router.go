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
		products.POST("", authMid, adminMid, r.CreateProduct)
		products.GET("/search", r.SearchProducts)
		products.GET("/:id", r.GetProduct)
		products.PUT("/:id", authMid, adminMid, r.UpdateProduct)
		products.DELETE("/:id", authMid, adminMid, r.DeleteProduct)
		products.GET("", authMid, r.ListProducts)
	}

	categories := apiV1Group.Group("/categories")
	{
		categories.POST("", authMid, adminMid, r.CreateCategory)
		categories.GET("/:id", r.GetCategory)
		categories.GET("/:id/children", r.GetCategoryChildren)
		categories.PUT("/:id", authMid, adminMid, r.UpdateCategory)
		categories.DELETE("/:id", authMid, adminMid, r.DeleteCategory)
		categories.GET("", r.ListCategories)
		categories.GET("/:id/products", r.GetProductsByCategory)
	}
}
