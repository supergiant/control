package account

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/pborman/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/supergiant/control/pkg/account"
	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/storage/file"
)

func TestFindDigitalOceanRegions(t *testing.T) {
	doToken := os.Getenv(clouds.EnvDigitalOceanAccessToken)
	if doToken == "" {
		t.SkipNow()
	}

	logrus.SetLevel(logrus.DebugLevel)

	s, err := file.NewFileRepository(fmt.Sprintf("/tmp/sg-storage-%d", time.Now().UnixNano()))
	require.Nil(t, err, "setup file storage provider")

	accounts := account.NewService(account.DefaultStoragePrefix, s)

	accName := uuid.New()
	err = accounts.Create(context.Background(), &model.CloudAccount{
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
