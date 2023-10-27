package user

import (
	"github.com/gofiber/fiber/v2"
	"google.golang.org/protobuf/encoding/protojson"

	"EchoWave/internal/proto"
)

// handler is the user handler.
// job of the handler is to have controller methods attached to it,
// this way handler can access dependencies and pass them to the controller methods.
type handler struct {
	service proto.UserServiceClient
}

// InitRoutes initializes the user routes.
func InitRoutes(v1Group fiber.Router, service *proto.UserServiceClient) {
	h := &handler{service: *service}

	userRoutes := v1Group.Group("/users")

	userRoutes.Get("/:username", h.getUser)
}

// getUser is the controller method for GET /users/:username.
func (h *handler) getUser(ctx *fiber.Ctx) error {
	username := ctx.Params("username")

	user, err := h.service.GetUserByUsername(ctx.Context(), &proto.GetUserByUsernameRequest{Username: username})
	if err != nil {
		return err
	}

	return ctx.JSON(protojson.Format(user))
}
