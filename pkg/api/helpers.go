package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
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

const defaultListLimit = "50"

func handleList(core *core.Core, r *http.Request, m model.Model, listPtr interface{}) (resp *Response, err error) {
	q := r.URL.Query()

	query, limit, offset := buildQuery(q, m), q.Get("limit"), q.Get("offset")
	if limit == "" {
		limit = defaultListLimit
	}

	if err := listModels(core, m, listPtr, query, limit, offset); err != nil {
		return nil, err
	}

	return &Response{
		http.StatusOK,
		listPtr,
	}, nil
}

func handleKubeList(core *core.Core, r *http.Request) (resp *Response, err error) {
	m, listPtr := new(model.Kube), new(model.KubeList)
	q := r.URL.Query()

	query, limit, offset := buildQuery(q, m), q.Get("limit"), q.Get("offset")
	if limit == "" {
		limit = defaultListLimit
	}

	// get kube models
	if err := listModels(core, m, listPtr, query, limit, offset); err != nil {
		return nil, err
	}

	// populate kube models with nodes and helm releases
	for _, k := range listPtr.Items {
		nodes, err := listKubeNodes(core, k.Name)
		if err != nil {
			return nil, err
		}
		k.Nodes = nodes

		releases, err := listKubeReleases(core, k.Name)
		if err != nil {
			return nil, err
		}
		k.HelmReleases = releases
	}

	return &Response{
		http.StatusOK,
		listPtr,
	}, nil
}

func listKubeNodes(core *core.Core, kname string) ([]*model.Node, error) {
	var query string
	if kname != "" {
		query = fmt.Sprintf(`kube_name = '%s'`, kname)
	}

	list := new(model.NodeList)
	if err := listModels(core, new(model.Node), list, query, nil, nil); err != nil {
		return nil, err
	}
	return list.Items, nil
}

func listKubeReleases(core *core.Core, kname string) ([]*model.HelmRelease, error) {
	var query string
	if kname != "" {
		query = fmt.Sprintf(`kube_name = '%s'`, kname)
	}

	list := new(model.HelmReleaseList)
	if err := listModels(core, new(model.HelmRelease), list, query, nil, nil); err != nil {
		return nil, err
	}
	return list.Items, nil
}

func listModels(core *core.Core, m model.Model, listPtr, filter, limit, offset interface{}) error {
	listValue := reflect.ValueOf(listPtr).Elem()
	slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(m)), 0, 0)
	items := listValue.FieldByName("Items")
	items.Set(slice)

	scope := core.DB.Where(filter).Limit(limit).Offset(offset)

	return scope.Find(items.Addr().Interface())
}

func buildQuery(q url.Values, m model.Model) string {
	var andQueries []string
	for _, field := range model.RootFieldJSONNames(m) {

		// ?filter.name=this&filter.name=that
		filterValues := q["filter."+field]

		var orQueries []string
		for _, val := range filterValues {
			orQueries = append(orQueries, fmt.Sprintf("%s = '%s'", field, val))
		}

		if len(orQueries) > 0 {
			andQueries = append(andQueries, "("+strings.Join(orQueries, " OR ")+")")
		}
	}
	andQuery := strings.Join(andQueries, " AND ")

	return andQuery
}
