package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// ShowAccount godoc
//
//	@Summary		测试 ping, 任何请求方式都可以
//	@Tags			test
//	@Produce		json
//	@Success		200
//	@Router			/api/v1/ping	[get]
func GetPing(c *fiber.Ctx) error {
	err := c.JSON(map[string]string{"ping": "pang"})
	if err != nil {
		log.Errorf("GetPing err: %v", err)
		return err
	}
	return nil
}
