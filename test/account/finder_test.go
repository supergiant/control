package account

import (
	"os"
	"testing"

	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/coreos/etcd/clientv3"
	"github.com/gorilla/mux"
	"github.com/pborman/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/supergiant/supergiant/pkg/account"
	"github.com/supergiant/supergiant/pkg/clouds"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/storage"
	"github.com/supergiant/supergiant/pkg/testutils/assert"
)

const defaultETCDHost = "http://127.0.0.1:2379"

var defaultConfig clientv3.Config

func init() {
	assert.MustRunETCD(defaultETCDHost)
	defaultConfig = clientv3.Config{
		Endpoints: []string{defaultETCDHost},
	}
}

func TestFindDigitalOceanRegions(t *testing.T) {
	doToken := os.Getenv(clouds.EnvDigitalOceanAccessToken)
	require.NotEmpty(t, doToken)

	logrus.SetLevel(logrus.DebugLevel)

	kv := storage.NewETCDRepository(defaultConfig)
	accounts := account.NewService(account.DefaultStoragePrefix, kv)

	accName := uuid.New()
	err := accounts.Create(context.Background(), &model.CloudAccount{
		Name:     accName,
		Provider: clouds.DigitalOcean,
		Credentials: map[string]string{
			clouds.DigitalOceanAccessToken: doToken,
		},
	})
	require.NoError(t, err)

	router := mux.NewRouter()
	handler := account.NewHandler(accounts)
	handler.Register(router)

	req := httptest.NewRequest(http.MethodGet, "/accounts/"+accName+"/regions", nil)
	req = mux.SetURLVars(req, map[string]string{
		"accountName": accName,
	})
	rr := httptest.NewRecorder()

	handler.GetRegions(rr, req)
	require.NoError(t, err)

	regionsJSON, err := ioutil.ReadAll(rr.Body)
	require.NoError(t, err)
	t.Log(string(regionsJSON))

	require.Equal(t, rr.Code, http.StatusOK)

	rs := &account.RegionSizes{}
	err = json.Unmarshal(regionsJSON, &rs)
	require.NoError(t, err)

	require.True(t, len(rs.Regions) > 0)
	require.True(t, len(rs.Sizes) > 0)
	require.Equal(t, clouds.DigitalOcean, rs.Provider)
}
