package api

import (
	"testing"

	"github.com/supergiant/supergiant/pkg/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUsersList(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	user, admin := createUserAndAdmin(srv.Core)

	Convey("Given a user and an admin", t, func() {

		Convey("When the user Lists Users", func() {
			sg := srv.Core.APIClient("token", user.APIToken)
			users := new(model.UserList)
			err := sg.Users.List(users)

			Convey("They should see only themself", func() {
				So(err, ShouldBeNil)
				So(*users.Items[0].ID, ShouldEqual, *user.ID)
			})
		})

		Convey("When the admin Lists Users", func() {
			sg := srv.Core.APIClient("token", admin.APIToken)
			users := new(model.UserList)
			err := sg.Users.List(users)

			Convey("They should see everyone", func() {
				So(err, ShouldBeNil)
				So(users.Total, ShouldEqual, 2)
			})
		})
	})
}

func TestUsersCreate(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	user, admin := createUserAndAdmin(srv.Core)

	newUser := &model.User{
		Username: "my-new-freind",
		Password: "password",
	}

	Convey("Given a user and an admin", t, func() {

		Convey("When the user Creates a User", func() {
			sg := srv.Core.APIClient("token", user.APIToken)
			err := sg.Users.Create(newUser)

			Convey("They should receive a 403 Forbidden error", func() {
				So(err.(*model.Error).Status, ShouldEqual, 403)
			})
		})

		Convey("When the admin Creates a User", func() {
			sg := srv.Core.APIClient("token", admin.APIToken)
			err := sg.Users.Create(newUser)

			Convey("They should be no error", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestUsersGet(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	user, admin := createUserAndAdmin(srv.Core)

	Convey("Given a user and an admin", t, func() {

		Convey("When the user Gets another User", func() {
			sg := srv.Core.APIClient("token", user.APIToken)
			err := sg.Users.Get(admin.ID, admin)

			Convey("They should receive a 403 Forbidden error", func() {
				So(err.(*model.Error).Status, ShouldEqual, 403)
			})
		})

		Convey("When the user Gets themself", func() {
			sg := srv.Core.APIClient("token", user.APIToken)
			err := sg.Users.Get(user.ID, user)

			Convey("There should be no error", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When the admin Gets another User", func() {
			sg := srv.Core.APIClient("token", admin.APIToken)
			err := sg.Users.Get(user.ID, user)

			Convey("There should be no error", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When the admin Gets themself", func() {
			sg := srv.Core.APIClient("token", admin.APIToken)
			err := sg.Users.Get(admin.ID, admin)

			Convey("There should be no error", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestUsersUpdate(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	user, admin := createUserAndAdmin(srv.Core)

	srv.Core.Users.Get(user.ID, user)
	origUserPass := string(user.EncryptedPassword)

	srv.Core.Users.Get(admin.ID, admin)
	origAdminPass := string(admin.EncryptedPassword)

	Convey("Given a user and an admin", t, func() {

		Convey("When the user Updates another User", func() {
			sg := srv.Core.APIClient("token", user.APIToken)
			err := sg.Users.Update(admin.ID, &model.User{Password: "new-password"})

			Convey("They should receive a 403 Forbidden error", func() {
				So(err.(*model.Error).Status, ShouldEqual, 403)
			})
		})

		Convey("When the user Updates themself", func() {
			sg := srv.Core.APIClient("token", user.APIToken)
			err := sg.Users.Update(user.ID, &model.User{Password: "new-password"})

			reloadedUser := new(model.User)
			srv.Core.Users.Get(user.ID, reloadedUser)
			newUserPass := string(reloadedUser.EncryptedPassword)

			Convey("The Update should be successful", func() {
				So(err, ShouldBeNil)
				So(newUserPass, ShouldNotBeEmpty)
				So(newUserPass, ShouldNotEqual, origUserPass)
			})
		})

		Convey("When the user tries to Update their role", func() {
			sg := srv.Core.APIClient("token", user.APIToken)
			updatedUser := *user
			updatedUser.Password = "password" // just required for update, no change
			updatedUser.Role = "admin"
			err := sg.Users.Update(user.ID, &updatedUser)

			Convey("The Update should be succesful, but the role unchanged", func() {
				So(err, ShouldBeNil)
				So(updatedUser.Role, ShouldEqual, "user")
			})
		})

		Convey("When the admin Updates another User", func() {
			sg := srv.Core.APIClient("token", admin.APIToken)
			err := sg.Users.Update(user.ID, &model.User{Password: "new-password"})

			reloadedUser := new(model.User)
			srv.Core.Users.Get(user.ID, reloadedUser)
			newUserPass := string(reloadedUser.EncryptedPassword)

			Convey("The Update should be successful", func() {
				So(err, ShouldBeNil)
				So(newUserPass, ShouldNotBeEmpty)
				So(newUserPass, ShouldNotEqual, origUserPass)
			})
		})

		Convey("When the admin Updates themself", func() {
			sg := srv.Core.APIClient("token", admin.APIToken)
			err := sg.Users.Update(admin.ID, &model.User{Password: "new-password"})

			reloadedAdmin := new(model.User)
			srv.Core.Users.Get(admin.ID, reloadedAdmin)
			newAdminPass := string(reloadedAdmin.EncryptedPassword)

			Convey("The Update should be successful", func() {
				So(err, ShouldBeNil)
				So(newAdminPass, ShouldNotBeEmpty)
				So(newAdminPass, ShouldNotEqual, origAdminPass)
			})
		})

		Convey("When the admin updates a User role", func() {
			sg := srv.Core.APIClient("token", admin.APIToken)
			updatedUser := *user
			updatedUser.Password = "password" // just required for update, no change
			updatedUser.Role = "admin"
			err := sg.Users.Update(user.ID, &updatedUser)

			Convey("The Update should be succesful, and the role changed", func() {
				So(err, ShouldBeNil)
				So(updatedUser.Role, ShouldEqual, "admin")
			})
		})
	})
}

func TestUsersRegenerateAPIToken(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	user, admin := createUserAndAdmin(srv.Core)

	Convey("Given a user and an admin", t, func() {

		Convey("When the user calls RegenerateAPIToken on another User", func() {
			sg := srv.Core.APIClient("token", user.APIToken)
			err := sg.Users.RegenerateAPIToken(admin.ID, new(model.User))

			Convey("They should receive a 403 Forbidden error", func() {
				So(err.(*model.Error).Status, ShouldEqual, 403)
			})
		})

		Convey("When the user calls RegenerateAPIToken on themself", func() {
			sg := srv.Core.APIClient("token", user.APIToken)
			reloadedUser := new(model.User)
			err := sg.Users.RegenerateAPIToken(user.ID, reloadedUser)

			Convey("The token should be regenerated", func() {
				So(err, ShouldBeNil)
				So(reloadedUser.APIToken, ShouldNotEqual, user.APIToken)
			})
		})

		Convey("When the admin calls RegenerateAPIToken on another User", func() {
			sg := srv.Core.APIClient("token", admin.APIToken)
			reloadedUser := new(model.User)
			err := sg.Users.RegenerateAPIToken(user.ID, reloadedUser)

			Convey("The token should be regenerated", func() {
				So(err, ShouldBeNil)
				So(reloadedUser.APIToken, ShouldNotEqual, user.APIToken)
			})
		})

		Convey("When the admin calls RegenerateAPIToken on themself", func() {
			sg := srv.Core.APIClient("token", admin.APIToken)
			reloadedAdmin := new(model.User)
			err := sg.Users.RegenerateAPIToken(admin.ID, reloadedAdmin)

			Convey("The token should be regenerated", func() {
				So(err, ShouldBeNil)
				So(reloadedAdmin.APIToken, ShouldNotEqual, admin.APIToken)
			})
		})
	})
}

func TestUsersDelete(t *testing.T) {
	Convey("Given a user and an admin", t, func() {

		// NOTE we do this each assertion since we're deleting
		srv := newTestServer()
		go srv.Start()
		defer srv.Stop()

		observer := &model.User{
			Username: "observer",
			Password: "password",
			Role:     "admin",
		}
		srv.Core.Users.Create(observer)
		observerClient := srv.Core.APIClient("token", observer.APIToken)

		user, admin := createUserAndAdmin(srv.Core)

		Convey("When the user Deletes another User", func() {
			sg := srv.Core.APIClient("token", user.APIToken)
			err := sg.Users.Delete(admin.ID, admin)

			Convey("They should receive a 403 Forbidden error", func() {
				So(err.(*model.Error).Status, ShouldEqual, 403)
			})
		})

		Convey("When the user Deletes themself", func() {
			sg := srv.Core.APIClient("token", user.APIToken)
			deleteErr := sg.Users.Delete(user.ID, user)
			getErr := observerClient.Users.Get(user.ID, new(model.User))

			Convey("The User should be Deleted", func() {
				So(deleteErr, ShouldBeNil)
				So(getErr.(*model.Error).Status, ShouldEqual, 404)
			})
		})

		Convey("When the admin Deletes another User", func() {
			sg := srv.Core.APIClient("token", admin.APIToken)
			deleteErr := sg.Users.Delete(user.ID, user)
			getErr := observerClient.Users.Get(user.ID, new(model.User))

			Convey("The User should be Deleted", func() {
				So(deleteErr, ShouldBeNil)
				So(getErr.(*model.Error).Status, ShouldEqual, 404)
			})
		})

		Convey("When the admin Deletes themself", func() {
			sg := srv.Core.APIClient("token", admin.APIToken)
			deleteErr := sg.Users.Delete(admin.ID, admin)
			getErr := observerClient.Users.Get(admin.ID, new(model.User))

			Convey("The User should be Deleted", func() {
				So(deleteErr, ShouldBeNil)
				So(getErr.(*model.Error).Status, ShouldEqual, 404)
			})
		})
	})
}
