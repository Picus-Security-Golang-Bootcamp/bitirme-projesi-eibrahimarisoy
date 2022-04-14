package role

// TODO swagger implementation
import (
	"net/http"
	"patika-ecommerce/internal/model"
	"patika-ecommerce/pkg/config"

	mw "patika-ecommerce/pkg/middleware"

	httpErr "patika-ecommerce/internal/httpErrors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type roleHandler struct {
	roleRepo *RoleRepository
}

func NewRoleHandler(r *gin.RouterGroup, cfg *config.Config, roleRepo *RoleRepository) {
	handler := &roleHandler{
		roleRepo: roleRepo,
	}

	r.Use(mw.AuthenticationMiddleware(cfg.JWTConfig.SecretKey))
	r.Use(mw.AdminMiddleware())
	r.POST("", handler.createRole)
	r.GET("/", handler.getRoles)
	r.GET("/:id", handler.getRole)
	r.PUT("/:id", handler.updateRole)
	r.DELETE("/:id", handler.deleteRole)
}

// createRole creates a new role
func (r *roleHandler) createRole(c *gin.Context) {
	role := &model.Role{}
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	if err := r.roleRepo.InsertRole(role); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, role)
}

// getRoles returns all roles
func (r *roleHandler) getRoles(c *gin.Context) {
	roles, err := r.roleRepo.GetRoles()
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, roles)
}

// deleteRole deletes a role
func (r *roleHandler) deleteRole(c *gin.Context) {
	id := c.Param("id")
	if err := r.roleRepo.DeleteRole(id); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Role deleted"})
}

// updateRole updates a role
func (r *roleHandler) updateRole(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}
	role := &model.Role{
		Base: model.Base{ID: id},
	}

	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	if err := r.roleRepo.UpdateRole(role); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, role)
}

// getRole returns a role
func (r *roleHandler) getRole(c *gin.Context) {
	id := c.Param("id")
	role, err := r.roleRepo.GetRole(id)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, role)
}
