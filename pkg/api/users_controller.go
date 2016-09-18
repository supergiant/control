package api

import (
	"fmt"
	"net/http"

	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

type errorForbidden struct {
	user *model.User
}

func (err *errorForbidden) Error() string {
	return fmt.Sprintf("User %d cannot perform this operation", *err.user.ID)
}

func ensureAdmin(user *model.User) error {
	if user.Role != model.UserRoleAdmin {
		return &errorForbidden{user}
	}
	return nil
}

func ensureSameUser(id *int64, user *model.User) error {
	if *id != *user.ID {
		return &errorForbidden{user}
	}
	return nil
}

//------------------------------------------------------------------------------

func ListUsers(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	// Admin can see everyone
	if user.Role == model.UserRoleAdmin {
		return handleList(core, r, new(model.User), new(model.UserList))
	}

	list := &model.UserList{
		Items: []*model.User{user},
		BaseList: model.BaseList{
			Limit: 1,
			Total: 1,
		},
	}

	// Only show User themself if not admin
	return &Response{
		http.StatusOK,
		list,
	}, nil
}

func CreateUser(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	if err := ensureAdmin(user); err != nil {
		return nil, err
	}

	item := new(model.User)
	if err := decodeBodyInto(r, item); err != nil {
		return nil, err
	}
	if err := core.Users.Create(item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusCreated)
}

func UpdateUser(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}

	// Ensure the requester is this User, or an Admin
	if err := ensureSameUser(id, user); err != nil {
		if err = ensureAdmin(user); err != nil {
			return nil, err
		}
	}

	item := new(model.User)
	if err := decodeBodyInto(r, item); err != nil {
		return nil, err
	}

	// Only admins can change User roles
	if user.Role != model.UserRoleAdmin {
		item.Role = ""
	}

	if err := core.Users.Update(id, new(model.User), item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusAccepted)
}

func GetUser(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.User)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}

	// Ensure the requester is this User, or an Admin
	if err := ensureSameUser(id, user); err != nil {
		if err = ensureAdmin(user); err != nil {
			return nil, err
		}
	}

	if err := core.Users.Get(id, item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusOK)
}

func DeleteUser(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.User)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}

	// Ensure the requester is this User, or an Admin
	if err := ensureSameUser(id, user); err != nil {
		if err = ensureAdmin(user); err != nil {
			return nil, err
		}
	}

	if err := core.Users.Delete(id, item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusAccepted)
}

func RegenerateUserAPIToken(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.User)
	id, err := parseID(r)
	if err != nil {
		return nil, err
	}

	// Ensure the requester is this User, or an Admin
	if err := ensureSameUser(id, user); err != nil {
		if err = ensureAdmin(user); err != nil {
			return nil, err
		}
	}

	if err := core.Users.RegenerateAPIToken(id, item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusAccepted)
}
