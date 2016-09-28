package model

import "time"

type SessionList struct {
	BaseList
	Items []*Session `json:"items"`
}

type Session struct {
	ID        string    `json:"id"`
	UserID    *int64    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`

	User *User `json:"user"`
}

func (m *Session) Description() string {
	return "Session " + m.ID
}

// Session is not a model persisted in the DB, so we implement model interface.
// TODO there's probably a cleaner way.

func (m *Session) GetID() interface{} {
	return m.ID
}

func (m *Session) GetUUID() string {
	return ""
}

func (m *Session) SetUUID() {
}

func (m *Session) SetActionStatus(status *ActionStatus) {
}

func (m *Session) SetPassiveStatus() {
}
