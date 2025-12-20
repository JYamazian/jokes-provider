package controllers

import (
	"jokes-provider/config"
	"jokes-provider/services"
	"path"

	"github.com/gofiber/fiber/v2"
)

// JokeController handles joke-related endpoints
type JokeController struct {
	jokeService *services.JokeService
}

// NewJokeController creates a new JokeController instance
func NewJokeController() *JokeController {
	return &JokeController{
		jokeService: services.NewJokeService(),
	}
}

// GetRandomJoke godoc
// @Summary      Get a random joke
// @Description  Returns a random joke from the jokes database. Supports caching.
// @Tags         jokes
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Joke  "Random joke object with id and joke fields"
// @Failure      500  {object}  map[string]string  "Failed to retrieve joke"
// @Router       /api/v1/jokes/random [get]
func (ctrl *JokeController) GetRandomJoke(c *fiber.Ctx) error {
	cacheKey := path.Base(c.Path())

	joke, err := ctrl.jokeService.GetRandomJoke(c, cacheKey)
	if err != nil {
		config.LogError(c, "Error retrieving random joke", "error", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve joke",
		})
	}

	return c.Status(fiber.StatusOK).JSON(joke)
}
