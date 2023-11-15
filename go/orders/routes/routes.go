package routes

import (
	"log"

	"github.com/amirul-zafrin/event-ticketing/orders.git/database"
	"github.com/amirul-zafrin/event-ticketing/orders.git/models"
	"github.com/amirul-zafrin/event-ticketing/orders.git/services"
	"github.com/amirul-zafrin/event-ticketing/orders.git/utilities"
	"github.com/gofiber/fiber/v2"
)

func CreateOrder(c *fiber.Ctx) error {
	var order models.Orders

	if err := c.BodyParser(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": err.Error()})
	}

	if err := database.Database.Db.Create(&order).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failed", "message": err})
	}
	c.Locals("order", order)
	return c.Next()
	// return c.Status(200).SendString("OK!")
}

func UpdateOrder(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": "Please make sure id is integer"})
	}

	var order models.Orders
	if err := utilities.FindOrder(id, &order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": err.Error()})
	}

	if err := database.Database.Db.Updates(&order); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failed", "message": err})
	}
	c.Locals("order", order)
	return c.Next()
}

func UpdateEventSeatRMQ(c *fiber.Ctx) error {
	order := c.Locals("order").(models.Orders)
	// detail, err := json.Marshal(order.Details)
	// if err != nil {
	// 	log.Fatalf("Error during json marshall: %s", err)
	// }
	log.Println("Connecting to RMQ")
	rmq := services.ConnectRMQ()
	updateReq := services.LockingRequest{EventID: int(order.EventID), IsPaid: order.IsPaid, Details: order.Details}
	log.Println("Trying to publish!")
	err := rmq.PublishUpdateSeatStatus(&updateReq, "seat_requests")
	if err != nil {
		log.Println("Failed to publish seat update request: %s", err)
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": order})
}
