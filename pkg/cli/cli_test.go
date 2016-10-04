package cli_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/supergiant/supergiant/pkg/cli"
	"github.com/supergiant/supergiant/pkg/client"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/test/fake_client"
	cli_lib "github.com/urfave/cli"

	. "github.com/smartystreets/goconvey/convey"
)

// TODO we could use this everywhere
func idInt64(n int64) *int64 { return &n }

func TestCLIRun(t *testing.T) {
	Convey("CLI works correctly", t, func() {
		table := []struct {
			// Input
			command []string
			stdin   string
			// Mocks
			// Expectations
			clientCommandCalled string
			clientCommandArgs   []interface{}
		}{
			// CloudAccounts List
			{
				command:             []string{"supergiant", "cloud_accounts", "list"},
				clientCommandCalled: "CloudAccounts.List",
				clientCommandArgs: []interface{}{
					&model.CloudAccountList{
						BaseList: model.BaseList{
							Filters: map[string][]string{},
						},
					},
				},
			},
			// CloudAccounts Create
			{
				command: []string{"supergiant", "cloud_accounts", "create", "-f", "-"},
				stdin: `{
          "name": "test"
        }`,
				clientCommandCalled: "CloudAccounts.Create",
				clientCommandArgs: []interface{}{
					&model.CloudAccount{
						Name: "test",
					},
				},
			},
			// CloudAccounts Get
			{
				command:             []string{"supergiant", "cloud_accounts", "get", "--id=1"},
				clientCommandCalled: "CloudAccounts.Get",
				clientCommandArgs: []interface{}{
					idInt64(1),
					&model.CloudAccount{},
				},
			},
			// CloudAccounts Update
			{
				command: []string{"supergiant", "cloud_accounts", "update", "--id=1", "-f", "-"},
				stdin: `{
          "name": "test"
        }`,
				clientCommandCalled: "CloudAccounts.Update",
				clientCommandArgs: []interface{}{
					idInt64(1),
					&model.CloudAccount{
						Name: "test",
					},
				},
			},
			// CloudAccounts Delete
			{
				command:             []string{"supergiant", "cloud_accounts", "delete", "--id=1"},
				clientCommandCalled: "CloudAccounts.Delete",
				clientCommandArgs: []interface{}{
					idInt64(1),
					&model.CloudAccount{},
				},
			},
		}

		for _, item := range table {

			var clientCommandCalled string
			var clientCommandArgs []interface{}

			clientFn := func(_ *cli_lib.Context) *client.Client {
				return &client.Client{
					Sessions: &fake_client.Sessions{
						Collection: fake_client.Collection{
							ListFn: func(list model.List) error {
								clientCommandCalled = "Sessions.List"
								clientCommandArgs = []interface{}{list}
								return nil
							},
							CreateFn: func(m model.Model) error {
								clientCommandCalled = "Sessions.Create"
								clientCommandArgs = []interface{}{m}
								return nil
							},
							GetFn: func(id interface{}, m model.Model) error {
								clientCommandCalled = "Sessions.Get"
								clientCommandArgs = []interface{}{id, m}
								return nil
							},
							UpdateFn: func(id interface{}, m model.Model) error {
								clientCommandCalled = "Sessions.Update"
								clientCommandArgs = []interface{}{id, m}
								return nil
							},
							DeleteFn: func(id interface{}, m model.Model) error {
								clientCommandCalled = "Sessions.Delete"
								clientCommandArgs = []interface{}{id, m}
								return nil
							},
						},
					},
					Users: &fake_client.Users{
						Collection: fake_client.Collection{
							ListFn: func(list model.List) error {
								clientCommandCalled = "Users.List"
								clientCommandArgs = []interface{}{list}
								return nil
							},
							CreateFn: func(m model.Model) error {
								clientCommandCalled = "Users.Create"
								clientCommandArgs = []interface{}{m}
								return nil
							},
							GetFn: func(id interface{}, m model.Model) error {
								clientCommandCalled = "Users.Get"
								clientCommandArgs = []interface{}{id, m}
								return nil
							},
							UpdateFn: func(id interface{}, m model.Model) error {
								clientCommandCalled = "Users.Update"
								clientCommandArgs = []interface{}{id, m}
								return nil
							},
							DeleteFn: func(id interface{}, m model.Model) error {
								clientCommandCalled = "Users.Delete"
								clientCommandArgs = []interface{}{id, m}
								return nil
							},
						},
						RegenerateAPITokenFn: func(id interface{}, m *model.User) error {
							clientCommandCalled = "Users.RegenerateAPIToken"
							clientCommandArgs = []interface{}{id, m}
							return nil
						},
					},
					CloudAccounts: &fake_client.CloudAccounts{
						Collection: fake_client.Collection{
							ListFn: func(list model.List) error {
								clientCommandCalled = "CloudAccounts.List"
								clientCommandArgs = []interface{}{list}
								return nil
							},
							CreateFn: func(m model.Model) error {
								clientCommandCalled = "CloudAccounts.Create"
								clientCommandArgs = []interface{}{m}
								return nil
							},
							GetFn: func(id interface{}, m model.Model) error {
								clientCommandCalled = "CloudAccounts.Get"
								clientCommandArgs = []interface{}{id, m}
								return nil
							},
							UpdateFn: func(id interface{}, m model.Model) error {
								clientCommandCalled = "CloudAccounts.Update"
								clientCommandArgs = []interface{}{id, m}
								return nil
							},
							DeleteFn: func(id interface{}, m model.Model) error {
								clientCommandCalled = "CloudAccounts.Delete"
								clientCommandArgs = []interface{}{id, m}
								return nil
							},
						},
					},
					Kubes: &fake_client.Kubes{
						Collection: fake_client.Collection{
							ListFn: func(list model.List) error {
								clientCommandCalled = "Kubes.List"
								clientCommandArgs = []interface{}{list}
								return nil
							},
							CreateFn: func(m model.Model) error {
								clientCommandCalled = "Kubes.Create"
								clientCommandArgs = []interface{}{m}
								return nil
							},
							GetFn: func(id interface{}, m model.Model) error {
								clientCommandCalled = "Kubes.Get"
								clientCommandArgs = []interface{}{id, m}
								return nil
							},
							UpdateFn: func(id interface{}, m model.Model) error {
								clientCommandCalled = "Kubes.Update"
								clientCommandArgs = []interface{}{id, m}
								return nil
							},
							DeleteFn: func(id interface{}, m model.Model) error {
								clientCommandCalled = "Kubes.Delete"
								clientCommandArgs = []interface{}{id, m}
								return nil
							},
						},
					},
					KubeResources: &fake_client.KubeResources{
						Collection: fake_client.Collection{
							ListFn: func(list model.List) error {
								clientCommandCalled = "KubeResources.List"
								clientCommandArgs = []interface{}{list}
								return nil
							},
							CreateFn: func(m model.Model) error {
								clientCommandCalled = "KubeResources.Create"
								clientCommandArgs = []interface{}{m}
								return nil
							},
							GetFn: func(id interface{}, m model.Model) error {
								clientCommandCalled = "KubeResources.Get"
								clientCommandArgs = []interface{}{id, m}
								return nil
							},
							UpdateFn: func(id interface{}, m model.Model) error {
								clientCommandCalled = "KubeResources.Update"
								clientCommandArgs = []interface{}{id, m}
								return nil
							},
							DeleteFn: func(id interface{}, m model.Model) error {
								clientCommandCalled = "KubeResources.Delete"
								clientCommandArgs = []interface{}{id, m}
								return nil
							},
						},
						StartFn: func(id *int64, m *model.KubeResource) error {
							clientCommandCalled = "Users.Start"
							clientCommandArgs = []interface{}{id, m}
							return nil
						},
						StopFn: func(id *int64, m *model.KubeResource) error {
							clientCommandCalled = "Users.Stop"
							clientCommandArgs = []interface{}{id, m}
							return nil
						},
					},
					Volumes: &fake_client.Volumes{
						Collection: fake_client.Collection{
							ListFn: func(list model.List) error {
								clientCommandCalled = "Volumes.List"
								clientCommandArgs = []interface{}{list}
								return nil
							},
							CreateFn: func(m model.Model) error {
								clientCommandCalled = "Volumes.Create"
								clientCommandArgs = []interface{}{m}
								return nil
							},
							GetFn: func(id interface{}, m model.Model) error {
								clientCommandCalled = "Volumes.Get"
								clientCommandArgs = []interface{}{id, m}
								return nil
							},
							UpdateFn: func(id interface{}, m model.Model) error {
								clientCommandCalled = "Volumes.Update"
								clientCommandArgs = []interface{}{id, m}
								return nil
							},
							DeleteFn: func(id interface{}, m model.Model) error {
								clientCommandCalled = "Volumes.Delete"
								clientCommandArgs = []interface{}{id, m}
								return nil
							},
						},
					},
					Entrypoints: &fake_client.Entrypoints{
						Collection: fake_client.Collection{
							ListFn: func(list model.List) error {
								clientCommandCalled = "Entrypoints.List"
								clientCommandArgs = []interface{}{list}
								return nil
							},
							CreateFn: func(m model.Model) error {
								clientCommandCalled = "Entrypoints.Create"
								clientCommandArgs = []interface{}{m}
								return nil
							},
							GetFn: func(id interface{}, m model.Model) error {
								clientCommandCalled = "Entrypoints.Get"
								clientCommandArgs = []interface{}{id, m}
								return nil
							},
							UpdateFn: func(id interface{}, m model.Model) error {
								clientCommandCalled = "Entrypoints.Update"
								clientCommandArgs = []interface{}{id, m}
								return nil
							},
							DeleteFn: func(id interface{}, m model.Model) error {
								clientCommandCalled = "Entrypoints.Delete"
								clientCommandArgs = []interface{}{id, m}
								return nil
							},
						},
					},
					EntrypointListeners: &fake_client.EntrypointListeners{
						Collection: fake_client.Collection{
							ListFn: func(list model.List) error {
								clientCommandCalled = "EntrypointListeners.List"
								clientCommandArgs = []interface{}{list}
								return nil
							},
							CreateFn: func(m model.Model) error {
								clientCommandCalled = "EntrypointListeners.Create"
								clientCommandArgs = []interface{}{m}
								return nil
							},
							GetFn: func(id interface{}, m model.Model) error {
								clientCommandCalled = "EntrypointListeners.Get"
								clientCommandArgs = []interface{}{id, m}
								return nil
							},
							UpdateFn: func(id interface{}, m model.Model) error {
								clientCommandCalled = "EntrypointListeners.Update"
								clientCommandArgs = []interface{}{id, m}
								return nil
							},
							DeleteFn: func(id interface{}, m model.Model) error {
								clientCommandCalled = "EntrypointListeners.Delete"
								clientCommandArgs = []interface{}{id, m}
								return nil
							},
						},
					},
					Nodes: &fake_client.Nodes{
						Collection: fake_client.Collection{
							ListFn: func(list model.List) error {
								clientCommandCalled = "Nodes.List"
								clientCommandArgs = []interface{}{list}
								return nil
							},
							CreateFn: func(m model.Model) error {
								clientCommandCalled = "Nodes.Create"
								clientCommandArgs = []interface{}{m}
								return nil
							},
							GetFn: func(id interface{}, m model.Model) error {
								clientCommandCalled = "Nodes.Get"
								clientCommandArgs = []interface{}{id, m}
								return nil
							},
							UpdateFn: func(id interface{}, m model.Model) error {
								clientCommandCalled = "Nodes.Update"
								clientCommandArgs = []interface{}{id, m}
								return nil
							},
							DeleteFn: func(id interface{}, m model.Model) error {
								clientCommandCalled = "Nodes.Delete"
								clientCommandArgs = []interface{}{id, m}
								return nil
							},
						},
					},
				}
			}

			// Mock stdin
			file, _ := ioutil.TempFile(os.TempDir(), "stdin")
			defer os.Remove(file.Name())
			file.WriteString(item.stdin)
			file.Seek(0, os.SEEK_SET)

			cli.New(clientFn, file, "unversioned").Run(item.command)

			So(clientCommandCalled, ShouldResemble, item.clientCommandCalled)
			So(clientCommandArgs, ShouldResemble, item.clientCommandArgs)
		}
	})
}
