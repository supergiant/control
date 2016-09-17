package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

func CreateSession(core *core.Core, r *http.Request) (*Response, error) {
	item := new(model.Session)
	if err := decodeBodyInto(r, item); err != nil {
		return nil, err
	}
	if err := core.Sessions.Create(item); err != nil {
		return nil, err
	}
	return itemResponse(core, item, http.StatusCreated)
}

// ---------- separating open (above) from restricted (below) handlers ---------

func GetSession(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.Session)
	id := mux.Vars(r)["id"]
	if err := core.Sessions.Get(id, item); err != nil {
		return nil, err
	}

	// Ensure the requester is this Session, or an Admin
	if err := ensureSameUser(item.UserID, user); err != nil {
		if err = ensureAdmin(user); err != nil {
			return nil, err
		}
	}

	return itemResponse(core, item, http.StatusOK)
}

func ListSessions(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	allSessions := core.Sessions.List()
	var sessions []*model.Session

	if user.Role == model.UserRoleAdmin {
		// Admin can see everyone's session
		sessions = allSessions
	} else {
		// User can see only their own
		sessions = make([]*model.Session, 0)
		for _, session := range allSessions {
			if *session.UserID == *user.ID {
				sessions = append(sessions, session)
				break
			}
		}
	}

	list := &model.SessionList{
		Items: sessions,
		Pagination: model.Pagination{
			Limit: int64(len(sessions)),
			Total: int64(len(sessions)),
		},
	}

	return &Response{
		http.StatusOK,
		list,
	}, nil
}

func DeleteSession(core *core.Core, user *model.User, r *http.Request) (*Response, error) {
	item := new(model.Session)
	id := mux.Vars(r)["id"]
	if err := core.Sessions.Get(id, item); err != nil {
		return nil, err
	}

	// Ensure the requester is this Session, or an Admin
	if err := ensureSameUser(item.UserID, user); err != nil {
		if err = ensureAdmin(user); err != nil {
			return nil, err
		}
	}

	if err := core.Sessions.Delete(id); err != nil {
		return nil, err
	}
	return &Response{
		Status: http.StatusAccepted,
	}, nil
}
