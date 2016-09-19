package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
)

type Response struct {
	Status int
	Object interface{}
}

//------------------------------------------------------------------------------

type bodyDecodingError struct { // status bad request
	err error
}

func (e *bodyDecodingError) Error() string {
	return "Error decoding JSON body: " + e.err.Error()
}

var (
	errorUnauthorized  = errors.New("Unauthorized")
	errorBadAuthHeader = errors.New("Improperly formatted Authorization header")
)

//------------------------------------------------------------------------------

func errorHTTPStatus(err error) int {
	if _, ok := err.(*bodyDecodingError); ok {
		return 400
	}
	if err == core.ErrorBadLogin {
		return 400
	}
	if err == errorUnauthorized || err == errorBadAuthHeader {
		return 401
	}
	if _, ok := err.(*errorForbidden); ok {
		return 403
	}
	if err == gorm.ErrRecordNotFound {
		return 404
	}
	return 500
}

const logViewBytesize int64 = 2048

func logHandler(core *core.Core) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if user := loadUser(core, w, r); user == nil {
			return
		}

		if core.LogPath == "" {
			msg := "No log file configured!\nCreate file and provide path to --log-file at startup.\n"
			w.Write([]byte(msg))
			return
		}

		file, err := os.Open(core.LogPath)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		stat, err := os.Stat(core.LogPath)
		if err != nil {
			panic(err)
		}

		fileBytesize := stat.Size()

		var offset int64
		var bufferSize int64

		if fileBytesize < logViewBytesize {
			offset = 0
			bufferSize = fileBytesize
		} else {
			offset = fileBytesize - logViewBytesize
			bufferSize = logViewBytesize
		}

		buf := make([]byte, bufferSize)

		if _, err := file.ReadAt(buf, offset); err != nil {
			panic(err)
		}
		w.Write(buf)
	}
}

func loadUser(core *core.Core, w http.ResponseWriter, r *http.Request) *model.User {
	auth := r.Header.Get("Authorization")
	tokenMatch := regexp.MustCompile(`^SGAPI (token|session)="([A-Za-z0-9]{32})"$`).FindStringSubmatch(auth)

	if len(tokenMatch) != 3 {
		respond(w, nil, errorBadAuthHeader)
		return nil
	}

	switch tokenMatch[1] {
	case "token":
		user := new(model.User)
		if err := core.DB.Where("api_token = ?", tokenMatch[2]).First(user); err != nil {
			respond(w, nil, errorUnauthorized)
			return nil
		}

		return user

	case "session":
		session := new(model.Session)
		if err := core.Sessions.Get(tokenMatch[2], session); err != nil {
			respond(w, nil, errorUnauthorized)
			return nil
		}

		return session.User
	}

	respond(w, nil, errorBadAuthHeader)
	return nil
}

func respond(w http.ResponseWriter, resp *Response, err error) {
	if err != nil {
		status := errorHTTPStatus(err)
		resp = &Response{
			Status: status,
			Object: &model.Error{
				Status:  status,
				Message: err.Error(),
			},
		}
	}
	body, marshalErr := json.MarshalIndent(resp.Object, "", "  ")
	if marshalErr != nil {
		panic(marshalErr)
	}
	w.WriteHeader(resp.Status)
	w.Write(append(body, []byte{10}...)) // add line break (without string conversion)
}

func openHandler(c *core.Core, fn func(*core.Core, *http.Request) (*Response, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := fn(c, r)
		respond(w, resp, err)
	}
}

func restrictedHandler(core *core.Core, fn func(*core.Core, *model.User, *http.Request) (*Response, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := loadUser(core, w, r)
		if user == nil {
			return
		}
		resp, err := fn(core, user, r)
		respond(w, resp, err)
	}
}

//------------------------------------------------------------------------------

func parseID(r *http.Request) (*int64, error) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return nil, err
	}
	id64 := int64(id)
	return &id64, nil
}

func parseIncludes(r *http.Request) (includes []string) {
	if includesVal := r.URL.Query().Get("includes"); includesVal != "" {
		includes = strings.Split(includesVal, " ")
	}
	return
}

func decodeBodyInto(r *http.Request, item model.Model) error {
	if err := json.NewDecoder(r.Body).Decode(item); err != nil {
		return &bodyDecodingError{err}
	}
	model.ZeroReadonlyFields(item)
	return nil
}

func itemResponse(core *core.Core, item model.Model, status int) (*Response, error) {
	core.SetResourceActionStatus(item)
	return &Response{status, item}, nil
}

func handleList(core *core.Core, r *http.Request, m model.Model) (*Response, error) {
	slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(m)), 0, 0)
	itemsPtr := reflect.New(slice.Type())
	items := itemsPtr.Elem()
	items.Set(slice)

	qstr := r.URL.Query()

	var andQueries []string
	for _, field := range model.IndexedFields(m) {
		if val := qstr.Get(field.JSONName); val != "" {
			if field.Kind == reflect.Int64 {
				andQueries = append(andQueries, fmt.Sprintf("%s = %s", field.JSONName, val))
			} else { // string
				andQueries = append(andQueries, fmt.Sprintf("%s = '%s'", field.JSONName, val))
			}
		}
	}
	andQuery := strings.Join(andQueries, " AND ")

	scope := core.DB
	if andQuery != "" {
		scope = scope.Where(andQuery)
	}

	if err := scope.Find(itemsPtr.Interface()); err != nil {
		return nil, err
	}

	for i := 0; i < items.Len(); i++ {
		core.SetResourceActionStatus(items.Index(i).Interface().(model.Model))
	}

	return &Response{
		http.StatusOK,
		items.Interface(),
	}, nil
}
