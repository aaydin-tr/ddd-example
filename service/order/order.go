package order

import (
	"errors"

	"github.com/aaydin-tr/e-commerce/domain/order"
	entity "github.com/aaydin-tr/e-commerce/entity"
	"github.com/aaydin-tr/e-commerce/valueobject"
	"github.com/google/uuid"
)

var (
	ErrInsufficientStock = errors.New("Insufficient stock")
)

type OrderServiceInterface interface {
	Create(product *entity.Product, orderQuantity int) error
}

type OrderService struct {
	orderRepository order.OrderRepository
}

func NewOrderService(orderRepository order.OrderRepository) *OrderService {
	return &OrderService{orderRepository: orderRepository}
}

func (s *OrderService) Create(product *entity.Product, orderQuantity int) error {
	quantity, err := valueobject.NewQuantity(orderQuantity)
	if err != nil {
		return err
	}

	if product.Stock.Value() < quantity.Value() {
		return ErrInsufficientStock
	}

	err = s.orderRepository.Create(&entity.Order{
		ID:        uuid.New(),
		ProductID: product.ID,
		Quantity:  quantity,
	})

	if err != nil {
		return err
	}

	defer product.DecreaseStock(quantity.Value())
	defer product.IncreaseDemand(quantity.Value())

	if product.Campaign == nil || (product.Campaign != nil && !product.Campaign.IsActive()) {
		return nil
	}

	remainingTargetSaleCount := product.Campaign.RemainingTargetSalesCount(quantity.Value())
	var avaibleStockForCampaign int
	if remainingTargetSaleCount <= 0 {
		avaibleStockForCampaign = product.Campaign.TargetSalesCount.Value() - product.Campaign.TotalSales.Value()
		product.Campaign.Close()
	} else {
		avaibleStockForCampaign = quantity.Value()
	}

	err = product.Campaign.IncreaseTotalSales(avaibleStockForCampaign)
	if err != nil {
		return err
	}
	err = product.Campaign.UpdateAverageItemPrice(product.Price.Value(), avaibleStockForCampaign)
	if err != nil {
		return err
	}

	return nil
}
