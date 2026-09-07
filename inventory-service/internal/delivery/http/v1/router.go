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
	supplier  SupplierUC
	po        PurchaseOrderUC
	warehouse WarehouseUC
	stock     StockUC
	l         logger.Interface
	v         *validator.Validate
}

type Dependencies struct {
	Supplier      SupplierUC
	PurchaseOrder PurchaseOrderUC
	Warehouse     WarehouseUC
	Stock         StockUC
	Verifier      jwtmanager.JWTManager
	Cache         redismanager.BlacklistCacher
	Logger        logger.Interface
}

func NewDependencies(supplier SupplierUC, po PurchaseOrderUC, w WarehouseUC, stock StockUC, verifier jwtmanager.JWTManager, cache redismanager.BlacklistCacher, l logger.Interface) *Dependencies {
	return &Dependencies{
		Supplier:      supplier,
		PurchaseOrder: po,
		Warehouse:     w,
		Stock:         stock,
		Verifier:      verifier,
		Cache:         cache,
		Logger:        l,
	}
}

func NewRoutes(apiV1Group *gin.RouterGroup, deps *Dependencies) {
	r := &V1{
		supplier:  deps.Supplier,
		po:        deps.PurchaseOrder,
		warehouse: deps.Warehouse,
		stock:     deps.Stock,
		l:         deps.Logger,
		v:         validator.New(),
	}

	authMW := ginmw.Auth(deps.Verifier, deps.Cache)
	adminMW := ginmw.Role(ginmw.AdminRole)

	// Public stock routes (Front Store)
	apiV1Group.GET("/stock/availability", r.checkAvailability)
	apiV1Group.GET("/stock/availability/:sku", r.getSKUAvailability)

	// Stock management routes (Dashboard)
	stockGroup := apiV1Group.Group("/stock")
	stockGroup.Use(authMW, adminMW)
	stockGroup.GET("/levels", r.listStockLevels)
	stockGroup.GET("/levels/:id", r.getStockLevel)
	stockGroup.POST("/adjust", r.adjustStock)
	stockGroup.PUT("/levels/:id/thresholds", r.updateThresholds)
	stockGroup.GET("/summary", r.getStockSummary)
	stockGroup.GET("/movements", r.listStockMovements)

	// Warehouse routes
	whGroup := apiV1Group.Group("/warehouses")
	whGroup.Use(authMW, adminMW)
	whGroup.GET("", r.listWarehouses)

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
