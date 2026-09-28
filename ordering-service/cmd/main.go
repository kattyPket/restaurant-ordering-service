package main // ⬅️ แก้จาก "cmd" เป็น "main"

import (
	"log"
	"ordering-service/configs"
	"ordering-service/internal/db"
	"ordering-service/internal/menus"
	"ordering-service/internal/middleware"
	order "ordering-service/internal/orders"
	"ordering-service/internal/table" // ⬅️ เพิ่ม
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func main() {
	configs.LoadConfig()

	sessionTTL := time.Duration(configs.Cfg.Session.ExpirationHours) * time.Hour
	if sessionTTL <= 0 {
		log.Fatal("session.expiration_hours must be greater than zero")
	}

	db.ConnectDB()

	if err := db.DB.AutoMigrate(
		&table.DiningTable{},
		&table.TableSession{},
		&menus.MenuCategory{},
		&menus.MenuItem{},
		&order.Order{},
		&order.OrderItem{},
	); err != nil {
		log.Fatalf("❌ Failed to migrate table_sessions: %v", err)
	}

	if err := db.DB.Exec(`
    CREATE UNIQUE INDEX IF NOT EXISTS idx_table_sessions_one_active_per_table
    ON table_sessions (table_id)
    WHERE closed_at IS NULL;
`).Error; err != nil {
		log.Fatalf("failed to create active-session index: %v", err)
	}

	// 3. Initialize Fiber App
	app := fiber.New()
	app.Use(logger.New())

	app.Get("/ping", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Ordering Service is running! 🚀",
		})
	})

	// --- table module wiring ---  ⬅️ เพิ่มทั้งบล็อกนี้
	tableRepo := table.NewRepository(db.DB) // ปรับ db.DB ให้ตรงกับชื่อ exported *gorm.DB ตัวจริงใน internal/db
	tableService := table.NewService(tableRepo, sessionTTL)
	tableHandler := table.NewHandler(tableService)

	menuRepo := menus.NewRepository(db.DB)
	menuService := menus.NewService(menuRepo)
	menuHandler := menus.NewHandler(menuService)

	orderRepo := order.NewRepository(db.DB)
	orderService := order.NewService(orderRepo, menuRepo)
	orderHandler := order.NewHandler(orderService)

	api := app.Group("/api")
	table.SetupTableRoutes(api, tableHandler)
	app.Post("/internal/events/table-checkout-completed", tableHandler.HandleCheckoutCompleted)
	menus.SetupRoutes(app, menuHandler)
	order.RegisterRoutes(api, orderHandler, middleware.TableSessionGuard(tableService))
	// --- end table module wiring ---

	// 4. Start Server on port 8081
	log.Printf("🚀 Starting Ordering Service on port %s", configs.Cfg.App.Port)
	if err := app.Listen(configs.Cfg.App.Port); err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}

//go run cmd/main.go
