package core

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/supergiant/supergiant/pkg/client"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/util"
)

const (
	SessionCookieName = "supergiant_session"
	sessionTTL        = 3 * time.Hour
)

var (
	ErrorBadLogin = errors.New("Invalid credentials")
)

type SessionsInterface interface {
	Client(id string) *client.Client
	List() []*model.Session
	Create(*model.Session) error
	Get(id string, m *model.Session) error
	Delete(id string) error
}

type Sessions struct {
	core     *Core
	sessions *SafeMap // map[session-id]*Session
	clients  *SafeMap // map[session-id]*Client (for reusing http clients)
}

func NewSessions(core *Core) *Sessions {
	return &Sessions{
		core:     core,
		sessions: NewSafeMap(core),
		clients:  NewSafeMap(core),
	}
}

//------------------------------------------------------------------------------

type SessionExpirer struct {
	core *Core
}

func (s *SessionExpirer) Perform() error {
	for _, session := range s.core.Sessions.List() {
		if time.Since(session.CreatedAt) > sessionTTL {
			if err := s.core.Sessions.Delete(session.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

//------------------------------------------------------------------------------

// Each Session reuses a single Client instance, and this method fetches that.
func (c *Sessions) Client(id string) *client.Client {
	if ci := c.clients.Get(id); ci != nil {
		return ci.(*client.Client)
	}
	return nil
}

//------------------------------------------------------------------------------

func (c *Sessions) List() (items []*model.Session) {
	items = make([]*model.Session, 0)
	for _, si := range c.sessions.List() {
		items = append(items, si.(*model.Session))
	}
	return
}

func (c *Sessions) Create(m *model.Session) error {
	// Find User by username
	if err := c.core.DB.Where("username = ?", m.User.Username).First(m.User); err != nil {
		return ErrorBadLogin
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword(m.User.EncryptedPassword, []byte(m.User.Password)); err != nil {
		return ErrorBadLogin
	}

	// Build Session (and set to user-passed value)
	*m = model.Session{
		ID:        util.RandomString(32),
		UserID:    m.User.ID,
		CreatedAt: time.Now(),
	}

	// Create Session
	c.sessions.Put(m.Description(), m.ID, m)

	// Create Client
	client := c.core.APIClient("session", m.ID)
	c.clients.Put("Client of "+m.Description(), m.ID, client)

	return nil
}

func (c *Sessions) Get(id string, m *model.Session) error {
	si := c.sessions.Get(id)
	if si == nil {
		return fmt.Errorf("Could not find session %s", id)
	}
	session := si.(*model.Session)

	session.User = new(model.User)

	// Load User
	if err := c.core.Users.Get(session.UserID, session.User); err != nil {
		return err
	}

	// NOTE there's no unmarshalling here, we just do this for consistency with
	// other collections
	*m = *session

	return nil
}

func (c *Sessions) Delete(id string) error {
	m := &model.Session{ID: id}
	c.sessions.Delete(m.Description(), id)
	c.clients.Delete("Client of "+m.Description(), id)
	return nil
}
