package auth

import (
	"net/http"
	httpErr "patika-ecommerce/internal/httpErrors"
	"patika-ecommerce/internal/model"
	"patika-ecommerce/pkg/config"
	mw "patika-ecommerce/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type userHandler struct {
	userRepo *UserRepository
}

func NewUserHandler(r *gin.RouterGroup, cfg *config.Config, userRepo *UserRepository) {
	handler := &userHandler{userRepo: userRepo}

	r.Use(mw.AuthenticationMiddleware(cfg.JWTConfig.SecretKey))
	r.Use(mw.AdminMiddleware())
	r.GET("/:id", handler.getUser)
	r.PUT("/:id", handler.updateUser)
	r.DELETE("/:id", handler.deleteUser)
}

func (u *userHandler) deleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := u.userRepo.DeleteUser(id); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "User deleted"})
}

func (u *userHandler) updateUser(c *gin.Context) {

	id := c.Param("id")
	idx, err := uuid.Parse(id)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	userBody := model.User{
		Base: model.Base{ID: idx},
	}
	if err := c.ShouldBindJSON(&userBody); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	user, err := u.userRepo.UpdateUser(&userBody)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, user)
}

func (u *userHandler) getUser(c *gin.Context) {
	user := &model.User{}
	if err := c.ShouldBindUri(&user); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	id := c.Param("id")
	user, err := u.userRepo.GetUser(id)

	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, user)
}
