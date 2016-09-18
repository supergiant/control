package api

import (
	"testing"

	"github.com/supergiant/supergiant/pkg/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestListFiltering(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("Filtering on API list operations works correctly", t, func() {
		requestor := createAdmin(srv.Core)
		sg := srv.Core.NewAPIClient("token", requestor.APIToken)

		// Create a bunch of Users (because it's easy)
		users := []*model.User{
			{
				Username: "user1",
				Password: "password",
			},
			{
				Username: "user2",
				Password: "password",
			},
			{
				Username: "admin1",
				Password: "password",
				Role:     "admin",
			},
		}
		for _, user := range users {
			sg.Users.Create(user)
		}

		table := []struct {
			filters          map[string][]string
			usernamesMatched []string
		}{
			// Empty
			{
				nil,
				[]string{"user1", "user2", "admin1"},
			},
			// Empty
			{
				map[string][]string{"username": []string{}},
				[]string{"user1", "user2", "admin1"},
			},
			// Single key, single value
			{
				map[string][]string{"username": []string{"user1"}},
				[]string{"user1"},
			},
			// Single key, single invalid value
			{
				map[string][]string{"username": []string{"userbutt"}},
				nil,
			},
			// Single key, multiple values (OR)
			{
				map[string][]string{"username": []string{"user1", "user2"}},
				[]string{"user1", "user2"},
			},
			// Multiple keys, single values (AND)
			{
				map[string][]string{"username": []string{"user1"}, "role": []string{"admin"}},
				nil,
			},
			// Multiple keys, single values (AND)
			{
				map[string][]string{"username": []string{"admin1"}, "role": []string{"admin"}},
				[]string{"admin1"},
			},
			// Multiple keys, multiple values (AND of ORs)
			{
				map[string][]string{"username": []string{"admin1", "user2"}, "role": []string{"user"}},
				[]string{"user2"},
			},
			// Multiple keys, multiple values (AND of ORs)
			{
				map[string][]string{"username": []string{"admin1", "user2"}, "role": []string{"user", "admin"}},
				[]string{"admin1", "user2"},
			},
		}

		for _, item := range table {
			list := &model.UserList{
				BaseList: model.BaseList{
					Filters: item.filters,
				},
			}
			err := sg.Users.List(list)

			So(err, ShouldBeNil)

			var usernamesMatched []string
			for _, user := range list.Items {

				// Don't include the requestor in the User list (this is specific to
				// Users, since it requires a User to do anything).
				if user.Username == requestor.Username {
					continue
				}

				usernamesMatched = append(usernamesMatched, user.Username)
			}

			So(usernamesMatched, ShouldResemble, item.usernamesMatched)
		}
	})
}

func TestListPagination(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("Pagination on API list operations works correctly", t, func() {
		requestor := createAdmin(srv.Core)
		sg := srv.Core.NewAPIClient("token", requestor.APIToken)

		// Create 2 more users
		sg.Users.Create(&model.User{Username: "user1", Password: "password"})
		sg.Users.Create(&model.User{Username: "user2", Password: "password"})

		table := []struct {
			limit             int64
			offset            int64
			expectedTotal     int64
			expectedUsernames []string
		}{
			// We currently do not accept 0 limits, though this should probably change.
			{
				limit:             0, // so this defaults to 25
				offset:            0,
				expectedTotal:     3,
				expectedUsernames: []string{"bossman", "user1", "user2"},
			},
			{
				limit:             1,
				offset:            0,
				expectedTotal:     3,
				expectedUsernames: []string{"bossman"},
			},
			{
				limit:             2,
				offset:            0,
				expectedTotal:     3,
				expectedUsernames: []string{"bossman", "user1"},
			},
			{
				limit:             25,
				offset:            1,
				expectedTotal:     3,
				expectedUsernames: []string{"user1", "user2"},
			},
			{
				limit:             25,
				offset:            2,
				expectedTotal:     3,
				expectedUsernames: []string{"user2"},
			},
			{
				limit:             25,
				offset:            3,
				expectedTotal:     3,
				expectedUsernames: nil,
			},
		}

		for _, item := range table {
			list := &model.UserList{
				BaseList: model.BaseList{
					Limit:  item.limit,
					Offset: item.offset,
				},
			}
			err := sg.Users.List(list)

			So(err, ShouldBeNil)

			var usernamesMatched []string
			for _, user := range list.Items {
				usernamesMatched = append(usernamesMatched, user.Username)
			}

			So(list.Total, ShouldEqual, item.expectedTotal)
			So(usernamesMatched, ShouldResemble, item.expectedUsernames)
		}
	})
}
