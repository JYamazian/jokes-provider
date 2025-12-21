package wrapper

import (
	"jokes-provider/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func shouldSkipCache(c *fiber.Ctx) bool {
	cacheControl := c.Get(utils.HeaderCacheControl)
	return strings.Contains(strings.ToLower(cacheControl), utils.CacheControlNoCache)
}
