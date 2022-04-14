package auth

import (
	"net/http"
	"patika-ecommerce/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type userHandler struct {
	userRepo *UserRepository
}

func NewUserHandler(r *gin.RouterGroup, userRepo *UserRepository) {
	handler := &userHandler{userRepo: userRepo}

	r.POST("/", handler.createUser)
	r.GET("/", handler.getUsers)
	r.GET("/:id", handler.getUser)
	r.PUT("/:id", handler.updateUser)
	r.DELETE("/:id", handler.deleteUser)
}

func (u *userHandler) createUser(c *gin.Context) {
	var userBody model.User
	if err := c.ShouldBindJSON(&userBody); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := u.userRepo.InsertUser(&userBody)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, user)
}

func (u *userHandler) getUsers(c *gin.Context) {
	users, err := u.userRepo.GetAll()

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (u *userHandler) deleteUser(c *gin.Context) {
	user := &model.User{}
	if err := c.ShouldBindUri(&user); err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}

	id := c.Param("id")
	if err := u.userRepo.DeleteUser(id); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(204, gin.H{"message": "User deleted"})
}

func (u *userHandler) updateUser(c *gin.Context) {
	modelUser := &model.User{}
	if err := c.ShouldBindUri(&modelUser); err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}

	id := c.Param("id")
	idx, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}

	userBody := model.User{
		Base: model.Base{ID: idx},
	}
	if err := c.ShouldBindJSON(&userBody); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := u.userRepo.UpdateUser(&userBody)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, user)
}

func (u *userHandler) getUser(c *gin.Context) {
	user := &model.User{}
	if err := c.ShouldBindUri(&user); err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}

	id := c.Param("id")
	user, err := u.userRepo.GetUser(id)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, user)
}
