package api

import (
	"testing"

	"github.com/supergiant/supergiant/pkg/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestValidToken(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("Given a User and Client with the user's token", t, func() {
		user := &model.User{
			Username: "tester",
			Password: "boofar26",
		}
		srv.Core.Users.Create(user)

		sg := srv.Core.NewAPIClient("token", user.APIToken)

		Convey("When a request is made", func() {
			err := sg.Users.Get(user.ID, user)

			Convey("There should be no error", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestImproperToken(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("Given a Client with an improperly formatted token", t, func() {
		sg := srv.Core.NewAPIClient("tokin'?", "yah mon")

		Convey("When a request is made", func() {
			err := sg.Nodes.Delete(1, new(model.Node))

			Convey("There should be a 401 error", func() {
				So(err.(*model.Error).Status, ShouldEqual, 401)
			})
		})
	})
}

func TestInvalidToken(t *testing.T) {
	srv := newTestServer()
	go srv.Start()
	defer srv.Stop()

	Convey("Given a Client with an invalid token", t, func() {
		sg := srv.Core.NewAPIClient("token", "ThisFormattedCorrectlyCuz32Chars")

		Convey("When a request is made", func() {
			err := sg.Nodes.Delete(1, new(model.Node))

			Convey("There should be a 401 error", func() {
				So(err.(*model.Error).Status, ShouldEqual, 401)
			})
		})
	})
}
