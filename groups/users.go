package groups

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/jfkgettingshot/roblox/internal"
)

type (
	SortOrder string
)

const (
	SortOrderAsc  SortOrder = "Asc"
	SortOrderDesc SortOrder = "Desc"

	GroupUserApi = "https://groups.roblox.com/v1/groups/%d/users"
)

type UserRole struct {
	User User `json:"user"`
	Role Role `json:"role"`
}

type Role struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Rank int64  `json:"rank"`
}

type User struct {
	HasVerifiedBadge bool   `json:"hasVerifiedBadge"`
	UserID           int64  `json:"userId"`
	Username         string `json:"username"`
	DisplayName      string `json:"displayName"`
}

func formatGroupUserApi(groupId int64, limit int32, cursor string, sortOrder SortOrder) (string, error) {
	parsedUrl, err := url.Parse(fmt.Sprintf(GroupUserApi, groupId))
	if err != nil {
		return "", fmt.Errorf("groups/users.go/formatGroupUserApi failed to parse url: %w", err)
	}

	query := parsedUrl.Query()
	query.Set("limit", fmt.Sprint(limit))
	query.Set("sortOrder", string(sortOrder))

	if cursor != "" {
		query.Set("cursor", cursor)
	}

	parsedUrl.RawQuery = query.Encode()

	return parsedUrl.String(), nil
}

func GetGroupUsersCurs(groupId int64, limit int32, cursor string, sortOrder SortOrder) (*internal.CursorResponse[UserRole], error) {
	apiUrl, err := formatGroupUserApi(groupId, limit, cursor, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("groups/users.go/GetGroupUsersCurs failed to format group user api url: %w", err)
	}

	request, err := http.Get(apiUrl)
	if err != nil {
		return nil, fmt.Errorf("groups/users.go/GetGroupUsersCurs failed to get group users: %w", err)
	}

	cursorResponse, err := internal.ReadCursorResponse[UserRole](request)
	if err != nil {
		return nil, fmt.Errorf("groups/users.go/GetGroupUsersCurs failed to read cursor response: %w", err)
	}

	return cursorResponse, nil
}

// WARNING WILL BE SLOW ON LARGE GROUPS
// GetGroupUsers returns a list of users in a group
// even on a failure if some users are returned
// the function will return the users and the error
func GetGroupUsers(groupId int64, sortOrder SortOrder) ([]UserRole, error) {
	cursor := ""
	users := make([]UserRole, 0)

	for {
		response, err := GetGroupUsersCurs(groupId, 100, cursor, sortOrder)
		cursor, err := internal.CursorHandler(&users, response, err)
		if err != nil {
			return users, fmt.Errorf("groups/users.go/GetGroupUsers failed to get group users: %w", err)
		}
		if len(cursor) == 0 {
			break
		}
	}

	return users, nil
}
