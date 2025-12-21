package controllers

import (
	"jokes-provider/config"
	"jokes-provider/helpers"
	"jokes-provider/services"
	"jokes-provider/utils"
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
// @Router       /v1/jokes/random [get]
func (ctrl *JokeController) GetRandomJoke(c *fiber.Ctx) error {
	cacheKey := path.Base(c.Path())

	joke, err := ctrl.jokeService.GetRandomJoke(c, cacheKey)
	if err != nil {
		config.LogError(c, "Error retrieving random joke", utils.JSONKeyError, err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			utils.JSONKeyError: utils.ErrMsgFailedToRetrieve,
		})
	}

	return c.Status(fiber.StatusOK).JSON(joke)
}

// GetJokeByID godoc
// @Summary      Get a joke by ID
// @Description  Returns a specific joke by its ID from the jokes database. Supports caching.
// @Tags         jokes
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Joke ID"
// @Success      200  {object}  models.Joke  "Joke object with id and joke fields"
// @Failure      404  {object}  map[string]string  "Joke not found"
// @Failure      500  {object}  map[string]string  "Failed to retrieve joke"
// @Router       /v1/jokes/{id} [get]
func (ctrl *JokeController) GetJokeByID(c *fiber.Ctx) error {
	jokeID := c.Params(utils.ParamID)

	if jokeID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			utils.JSONKeyError: utils.ErrMsgJokeIDRequired,
		})
	}

	joke, err := ctrl.jokeService.GetJokeByID(c, jokeID)
	if err != nil {
		if err == helpers.ErrJokeNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				utils.JSONKeyError: utils.ErrMsgJokeNotFound,
				utils.JSONKeyID:    jokeID,
			})
		}
		config.LogError(c, "Error retrieving joke by ID", utils.JSONKeyID, jokeID, utils.JSONKeyError, err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			utils.JSONKeyError: utils.ErrMsgFailedToRetrieve,
		})
	}

	return c.Status(fiber.StatusOK).JSON(joke)
}
