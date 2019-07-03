package jumpcloud

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestAccUserGroup(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	posixName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	gid := acctest.RandIntRange(1, 1000)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccUserGroup(rName, gid, posixName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("jumpcloud_user_group.test_group", "name", rName),
					resource.TestCheckResourceAttr("jumpcloud_user_group.test_group",
						"attributes.posix_groups", fmt.Sprintf("%d:%s", gid, posixName)),
				),
			},
		},
	})
}

func testAccUserGroup(name string, gid int, posixName string) string {
	return fmt.Sprintf(`
		resource "jumpcloud_user_group" "test_group" {
    		name = "%s"
			attributes = {
				posix_groups = "%d:%s"
			}
		}`, name, gid, posixName,
	)
}

func TestResourceUserGroup(t *testing.T) {
	suite.Run(t, new(ResourceUserGroupSuite))
}

type ResourceUserGroupSuite struct {
	suite.Suite
	A              *assert.Assertions
	TestHTTPServer *httptest.Server
}

func (s *ResourceUserGroupSuite) SetupSuite() {
	s.A = assert.New(s.Suite.T())
}

func (s *ResourceUserGroupSuite) TestTrueUserGroupRead() {
	cases := []struct {
		ResponseStatus int
		UserGroupNil   bool
		OK             bool
		ErrorNil       bool
		Payload        []byte
	}{
		{http.StatusNotFound, true, false, true, []byte("irrelevant")},
		{http.StatusOK, false, true, true, []byte("{}")},
	}

	for _, c := range cases {
		testServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(c.ResponseStatus)
			rw.Write(c.Payload)
		}))

		config := &jcapiv2.Configuration{
			BasePath: testServer.URL,
		}

		ug, ok, err := userGroupReadHelper(config, "id")
		s.A.Equal(c.OK, ok)
		s.A.Equal(c.UserGroupNil, ug == nil)
		s.A.Equal(c.ErrorNil, err == nil)
		testServer.Close()
	}
}
