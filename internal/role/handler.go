package role

import (
	"github.com/gin-gonic/gin"
)

type roleHandler struct {
	roleRepo *RoleRepository
}

func NewRoleHandler(r *gin.RouterGroup, roleRepo *RoleRepository) {
	handler := &roleHandler{roleRepo: roleRepo}

	r.POST("", handler.createRole)
	r.GET("/", handler.getRoles)
	r.GET("/:id", handler.getRole)
	r.PUT("/:id", handler.updateRole)
	r.DELETE("/:id", handler.deleteRole)
}

func (r *roleHandler) createRole(c *gin.Context) {

}

func (r *roleHandler) getRoles(c *gin.Context) {

}

func (r *roleHandler) deleteRole(c *gin.Context) {

}

func (r *roleHandler) updateRole(c *gin.Context) {

}

func (r *roleHandler) getRole(c *gin.Context) {}
