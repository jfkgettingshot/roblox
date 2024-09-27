package groups_test

import (
	"testing"

	"github.com/jfkgettingshot/roblox/groups"
)

func TestGetGroupUsersCurs(t *testing.T) {
	// Dont test GetGroupUsers since it takes a long ass time
	// User group id 1 since people will join it for being
	// the first group for whatever reason

	users, err := groups.GetGroupUsersCurs(1, 100, "", groups.SortOrderAsc)
	if err != nil {
		t.Errorf("groups.GetGroupUsersCurs(1, 100, \"\", groups.SortOrderAsc) returned an error: %v", err)
		return
	}

	if len(users.Data) == 0 {
		t.Errorf("groups.GetGroupUsersCurs(1, 100, \"\", groups.SortOrderAsc) returned 0 users")
		return
	}
}
