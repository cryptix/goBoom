package goBoom

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestReverse(t *testing.T) {
	Convey("reverting a simple string", t, func() {
		So(reverse("test 123"), ShouldEqual, "321 tset")
	})
}
