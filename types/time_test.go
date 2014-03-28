package types_test

import (
	"time"

	_ "github.com/wfreeman/cq"
	"github.com/wfreeman/cq/types"
	. "launchpad.net/gocheck"
)

func (s *TypesSuite) TestScanTime(c *C) {
	stmt := prepareTest("with {0} as test return test")
	rows, err := stmt.Query(1395967804 * 1000)
	c.Assert(err, IsNil)

	rows.Next()
	var test types.NullTime
	err = rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test.Valid, Equals, true)
	c.Assert(test.Time, DeepEquals, time.Unix(0, 1395967804*1000*1000000))
}
