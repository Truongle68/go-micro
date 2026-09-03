package v1

import (
	"github.com/GoProOrg/core-go-pkg/jwtmanager"
	redismanager "github.com/GoProOrg/core-go-pkg/redismanager/identity"
	"github.com/TruongLe68/go-micro/pkg/ginmw"
	"github.com/TruongLe68/go-micro/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type V1 struct {
	supplier SupplierUC
	po       PurchaseOrderUC
	l        logger.Interface
	v        *validator.Validate
}

type Dependencies struct {
	Supplier      SupplierUC
	PurchaseOrder PurchaseOrderUC
	Verifier      jwtmanager.JWTManager
	Cache         redismanager.BlacklistCacher
	Logger        logger.Interface
}

func NewDependencies(supplier SupplierUC, po PurchaseOrderUC, verifier jwtmanager.JWTManager, cache redismanager.BlacklistCacher, l logger.Interface) *Dependencies {
	return &Dependencies{
		Supplier:      supplier,
		PurchaseOrder: po,
		Verifier:      verifier,
		Cache:         cache,
		Logger:        l,
	}
}

func NewRoutes(apiV1Group *gin.RouterGroup, deps *Dependencies) {
	r := &V1{
		supplier: deps.Supplier,
		po:       deps.PurchaseOrder,
		l:        deps.Logger,
		v:        validator.New(),
	}

	authMW := ginmw.Auth(deps.Verifier, deps.Cache)
	adminMW := ginmw.Role(ginmw.AdminRole)

	// Supplier routes
	suppGroup := apiV1Group.Group("/suppliers")
	suppGroup.Use(authMW, adminMW)
	suppGroup.POST("", r.createSupplier)
	suppGroup.GET("", r.listSuppliers)
	suppGroup.GET("/:id", r.getSupplier)
	suppGroup.PUT("/:id", r.updateSupplier)
	suppGroup.POST("/:id/deactivate", r.deactivateSupplier)
	suppGroup.POST("/:id/reactivate", r.reactivateSupplier)

	// Purchase order routes
	poGroup := apiV1Group.Group("/purchase-orders")
	poGroup.Use(authMW, adminMW)
	poGroup.POST("", r.createPurchaseOrder)
	poGroup.GET("", r.listPurchaseOrders)
	poGroup.GET("/:id", r.getPurchaseOrder)
	poGroup.POST("/:id/confirm", r.confirmPurchaseOrder)
	poGroup.POST("/:id/receive", r.receiveLine)
	poGroup.POST("/:id/cancel", r.cancelPurchaseOrder)
}
