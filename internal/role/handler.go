package role

// TODO swagger implementation
import (
	"patika-ecommerce/internal/model"
	"patika-ecommerce/pkg/config"

	mw "patika-ecommerce/pkg/middleware"

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
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := r.roleRepo.InsertRole(role); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, role)
}

// getRoles returns all roles
func (r *roleHandler) getRoles(c *gin.Context) {
	roles, err := r.roleRepo.GetRoles()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, roles)
}

// deleteRole deletes a role
func (r *roleHandler) deleteRole(c *gin.Context) {
	id := c.Param("id")
	if err := r.roleRepo.DeleteRole(id); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(204, gin.H{"message": "Role deleted"})
}

// updateRole updates a role
func (r *roleHandler) updateRole(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}
	role := &model.Role{
		Base: model.Base{ID: id},
	}

	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := r.roleRepo.UpdateRole(role); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, role)
}

// getRole returns a role
func (r *roleHandler) getRole(c *gin.Context) {
	id := c.Param("id")
	role, err := r.roleRepo.GetRole(id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, role)
}
