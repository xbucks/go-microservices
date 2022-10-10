package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/utils"
	dtos2 "github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/dtos"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/features/getting_products"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/shared"
	"net/http"
)

type getProductsEndpoint struct {
	*shared.ProductEndpointBase[shared.InfrastructureConfiguration]
}

func NewGetProductsEndpoint(productEndpointBase *shared.ProductEndpointBase[shared.InfrastructureConfiguration]) *getProductsEndpoint {
	return &getProductsEndpoint{productEndpointBase}
}

func (ep *getProductsEndpoint) MapRoute() {
	ep.ProductsGroup.GET("", ep.getAllProducts())
}

// GetAllProducts
// @Tags Products
// @Summary Get all product
// @Description Get all products
// @Accept json
// @Produce json
// @Param getProductsRequestDto query dtos.GetProductsRequestDto false "GetProductsRequestDto"
// @Success 200 dtos.GetProductsResponseDto
// @Router /api/v1/products [get]
func (ep *getProductsEndpoint) getAllProducts() echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx := c.Request().Context()

		listQuery, err := utils.GetListQueryFromCtx(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		request := &dtos2.GetProductsRequestDto{ListQuery: listQuery}
		if err := c.Bind(request); err != nil {
			ep.Configuration.Log.Warn("Bind", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		query := &getting_products.GetProducts{ListQuery: request.ListQuery}

		queryResult, err := mediatr.Send[*getting_products.GetProducts, *dtos2.GetProductsResponseDto](ctx, query)

		if err != nil {
			ep.Configuration.Log.Warnf("GetProducts", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusOK, queryResult)
	}
}