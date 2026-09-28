// service.go
package order

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"ordering-service/internal/menus"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrOrderNotFound = errors.New("order not found")

type Service interface {
	PlaceOrder(tableID, sessionID uuid.UUID, req *PlaceOrderRequest) (*OrderResponse, error)
	GetOrderTracking(tableID, sessionID, orderID uuid.UUID) (*Order, error)
}

type orderService struct {
	orderRepo Repository
	menuRepo  menus.Repository // ຕ້ອງ Inject MenuRepo ເຂົ້າມາເພື່ອກວດລາຄາ
}

func NewService(repo Repository, menuRepo menus.Repository) Service {
	return &orderService{orderRepo: repo, menuRepo: menuRepo}
}

func (s *orderService) PlaceOrder(tableID, sessionID uuid.UUID, req *PlaceOrderRequest) (*OrderResponse, error) {
	var totalAmount float64
	var orderItems []OrderItem

	// 1. ດຶງລາຄາຈິງຈາກ DB (ບໍ່ເຊື່ອ Frontend)
	for _, reqItem := range req.Items {
		menuItem, err := s.menuRepo.GetMenuItemByID(reqItem.MenuItemID.String())
		if err != nil || !menuItem.IsAvailable {
			return nil, errors.New("ມີເມນູທີ່ບໍ່ພ້ອມຂາຍ ຫຼື ບໍ່ພົບໃນລະບົບ")
		}

		itemTotal := menuItem.Price * float64(reqItem.Quantity)
		totalAmount += itemTotal

		// ຈັດການ Pointer ຂອງ Note
		note := reqItem.Note
		var notePtr *string
		if note != "" {
			notePtr = &note
		}

		orderItems = append(orderItems, OrderItem{
			ID:         uuid.New(),
			MenuItemID: reqItem.MenuItemID,
			Quantity:   reqItem.Quantity,
			UnitPrice:  menuItem.Price,
			TotalPrice: itemTotal,
			Note:       notePtr,
			Status:     "pending",
		})
	}

	// สร้างเลขที่บิลแบบง่ายๆ เช่น ORD-20260915-142151
	currentTime := time.Now().Format("20060102-150405")
	generatedOrderNumber := fmt.Sprintf("ORD-%s", currentTime)

	// 2. ສ້າງ Order Status 'pending'
	newOrder := Order{
		ID:          uuid.New(),
		OrderNumber: generatedOrderNumber, // 👈 เพิ่มบรรทัดนี้เข้าไป!
		OrderType:   "dine-in",            // 👈 เพิ่มประเภทออร์เดอร์เข้าไปด้วย (เพราะใน Log มันเป็นค่าว่างเหมือนกัน)
		TableID:     &tableID,
		SessionID:   sessionID,
		Subtotal:    totalAmount,
		GrandTotal:  totalAmount,
		Status:      "pending",
	}

	// 3. ບັນທຶກລົງ DB ດ້ວຍ Transaction
	if err := s.orderRepo.CreateOrderWithItems(&newOrder, orderItems); err != nil {
		return nil, err
	}

	// 4. ຍິງ Event ຫາ Backoffice ຜ່ານ Internal Endpoint (ແທນການໃຊ້ Redis/Hub ໂດຍກົງ)
	go s.publishOrderPlacedEvent(newOrder.ID, tableID)

	return &OrderResponse{
		OrderID:     newOrder.ID,
		OrderNumber: newOrder.OrderNumber,
		Status:      newOrder.Status,
		GrandTotal:  newOrder.GrandTotal,
	}, nil
}

// Internal HTTP Call ໄປຫາ backoffice-service
func (s *orderService) publishOrderPlacedEvent(orderID, tableID uuid.UUID) {
	payload := map[string]string{
		"event":    "order.placed",
		"order_id": orderID.String(),
		"table_id": tableID.String(),
	}
	jsonData, _ := json.Marshal(payload)

	// ຍິງໄປຫາ Internal Endpoint ຂອງ Backoffice (:8080)
	resp, err := http.Post("http://localhost:8080/internal/events", "application/json", bytes.NewBuffer(jsonData))
	if err == nil {
		defer resp.Body.Close()
	}
	// ໃນ Layer 2 (Outbox Pattern) ເຮົາຈະປ່ຽນຈຸດນີ້ເປັນການ Insert ລົງຕາຕະລາງ outbox_events ແທນການຍິງ http.Post ໂດຍກົງ
}

// ເພີ່ມຟັງຊັນນີ້ຕໍ່ທ້າຍໄຟລ໌ service.go
func (s *orderService) GetOrderTracking(tableID, sessionID, orderID uuid.UUID) (*Order, error) {
	// ເອີ້ນໃຊ້ repo ດຶງຂໍ້ມູນ order ທີ່ລູກຄ້າສັ່ງ
	order, err := s.orderRepo.GetOrderByIDAndTableIDAndSessionID(orderID, tableID, sessionID)
	if err != nil {
		// ຖ້າຫາບໍ່ພົບ ຫຼື ບໍ່ແມ່ນຂອງໂຕະຕົວເອງ ໃຫ້ return 404 (ErrOrderNotFound)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	return order, nil
}
