package users

import (
	"context"
	"net/http"
	"quiz_platform/internal/handler/repository"
	"quiz_platform/internal/models"
	"quiz_platform/internal/utility"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserWithRoles struct {
	models.User
	Roles []*models.Role
}

type RoleWithSelected struct {
	models.Role
	Selected bool
}

func UsersListGetHandler(c *gin.Context) {
	ctx := context.Background()
	users, err := repository.UserRepositoryInstance.GetAllUsers(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userIds := make([]int32, len(users))
	userIdMap := make(map[int32]*UserWithRoles)
	usersWithRoles := make([]UserWithRoles, len(users))
	for i, u := range users {
		userIds[i] = u.Id
		usersWithRoles[i] = UserWithRoles{
			User:  *u,
			Roles: make([]*models.Role, 0),
		}
		userIdMap[u.Id] = &usersWithRoles[i]
	}

	ids, roles, err := repository.UserRepositoryInstance.GetUsersRoles(ctx, userIds)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, id := range ids {
		userIdMap[id].Roles = append(userIdMap[id].Roles, roles[i])
	}

	baseHInterface, _ := c.Get("BaseH")
	baseH, _ := baseHInterface.(*gin.H)
	c.HTML(http.StatusOK, "users_list.html", utility.MergeMaps(*baseH, gin.H{
		"title": "Users",
		"users": usersWithRoles}))
}

func UserEditFormGetHandler(c *gin.Context) {
	i, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := int32(i)

	ctx := context.Background()
	user, err := repository.UserRepositoryInstance.GetUserById(ctx, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userWithRoles := UserWithRoles{
		User:  *user,
		Roles: make([]*models.Role, 0),
	}

	_, userWithRoles.Roles, err = repository.UserRepositoryInstance.
		GetUsersRoles(ctx, []int32{user.Id})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	selectedRoles := make(map[int32]struct{})
	for _, v := range userWithRoles.Roles {
		selectedRoles[v.Id] = struct{}{}
	}

	allRoles, err := repository.RoleRepositoryInstance.GetAllRoles(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	allRoleDTOs := make([]*RoleWithSelected, len(allRoles))
	for i, v := range allRoles {
		_, wasSelected := selectedRoles[v.Id]
		allRoleDTOs[i] = &RoleWithSelected{
			Role:     *v,
			Selected: wasSelected,
		}
	}

	baseHInterface, _ := c.Get("BaseH")
	baseH, _ := baseHInterface.(*gin.H)
	c.HTML(http.StatusOK, "users_form.html", utility.MergeMaps(*baseH, gin.H{
		"title": "Edit user",
		"user":  userWithRoles,
		"roles": allRoleDTOs}))
}

func UserEditFormPostHandler(c *gin.Context) {
	i, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId := int32(i)

	username := c.PostForm("username")
	if len(username) > 100 || len(username) < 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid username"})
		return
	}

	selectedRoles := c.PostFormArray("roles")
	selectedRoleIds := make([]int32, len(selectedRoles))
	for i, v := range selectedRoles {
		roleId, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		selectedRoleIds[i] = int32(roleId)
	}

	ctx := context.Background()
	err = repository.TransactionManager.Run(ctx, func(ctx context.Context) error {

		err = repository.UserRepositoryInstance.UpdateUserName(ctx, userId, username)
		if err != nil {
			return err
		}

		err = repository.RoleRepositoryInstance.RemoveUserRoles(ctx, userId)
		if err != nil {
			return err
		}

		err = repository.RoleRepositoryInstance.AddUserRoles(ctx, userId, selectedRoleIds)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/users")
}

func UserDeletePostHandler(c *gin.Context) {
	i, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := int32(i)

	ctx := context.Background()
	err = repository.UserRepositoryInstance.DeleteUser(ctx, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/users")
}
