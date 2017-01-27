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
	// TODO we can probably consolidate all same error codes (would need to be in
	// model if that's where we keep the immutability check on fields).
	if _, ok := err.(*core.ErrorMissingRequiredParent); ok {
		return 422
	}
	if _, ok := err.(*core.ErrorValidationFailed); ok {
		return 422
	}
	if _, ok := err.(*model.ErrorChangedImmutableField); ok {
		return 422
	}
	return 500
}

const logViewBytesize int64 = 4096

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

func decodeBodyInto(r *http.Request, item model.Model) error {
	if err := json.NewDecoder(r.Body).Decode(item); err != nil {
		return &bodyDecodingError{err}
	}
	model.ZeroReadonlyFields(item)
	return nil
}

func itemResponse(core *core.Core, item model.Model, status int) (*Response, error) {
	core.SetResourceActionStatus(item)
	item.SetPassiveStatus()
	return &Response{status, item}, nil
}

const defaultListLimit = 25

func handleList(core *core.Core, r *http.Request, m model.Model, listPtr interface{}) (resp *Response, err error) {
	listValue := reflect.ValueOf(listPtr).Elem()

	slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(m)), 0, 0)
	items := listValue.FieldByName("Items")
	items.Set(slice)

	qstr := r.URL.Query()

	var andQueries []string
	for _, field := range model.RootFieldJSONNames(m) {

		// ?filter.name=this&filter.name=that
		filterValues := qstr["filter."+field]

		var orQueries []string
		for _, val := range filterValues {
			orQueries = append(orQueries, fmt.Sprintf("%s = '%s'", field, val))
		}

		if len(orQueries) > 0 {
			andQueries = append(andQueries, "("+strings.Join(orQueries, " OR ")+")")
		}
	}
	andQuery := strings.Join(andQueries, " AND ")

	baseScope := core.DB
	if andQuery != "" {
		baseScope = baseScope.Where(andQuery)
	}

	// BaseList
	pagination := model.BaseList{}

	if err := baseScope.Model(m).Count(&pagination.Total); err != nil {
		return nil, err
	}
	offsetParam := qstr.Get("offset")
	limitParam := qstr.Get("limit")

	pagination.Limit = defaultListLimit
	if limitParam != "" {
		if pagination.Limit, err = strconv.ParseInt(limitParam, 10, 64); err != nil {
			return nil, err
		}
	}

	if offsetParam != "" {
		if pagination.Offset, err = strconv.ParseInt(offsetParam, 10, 64); err != nil {
			return nil, err
		}
	}

	// TODO we may want to actually allow 0 limits here, and instead use pointers
	// to int64, because limit 0 will still return total count.
	scope := baseScope
	if pagination.Limit != 0 {
		scope = scope.Limit(pagination.Limit)
	}
	scope = scope.Offset(pagination.Offset)

	if err := scope.Find(items.Addr().Interface()); err != nil {
		return nil, err
	}

	for i := 0; i < items.Len(); i++ {
		item := items.Index(i).Interface().(model.Model)
		core.SetResourceActionStatus(item)
		item.SetPassiveStatus()
	}

	// Yeah... kinda nasty
	listValue.FieldByName("BaseList").Set(reflect.ValueOf(pagination))

	return &Response{
		http.StatusOK,
		listPtr,
	}, nil
}
