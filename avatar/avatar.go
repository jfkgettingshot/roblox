package avatar

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	AvatarApi = "https://avatar.roblox.com/v1/users/%d/avatar"
)

func formatAvatarApi(userId int64) string {
	return fmt.Sprintf(AvatarApi, userId)
}

func GetAvatar(userId int64) (*Avatar, error) {
	avatar := &Avatar{}
	response, err := http.Get(formatAvatarApi(userId))
	if err != nil {
		return nil, fmt.Errorf("avatar/avatar.go/GetAvatar failed to get avatar: %w", err)
	}
	defer response.Body.Close()

	jsonDecoder := json.NewDecoder(response.Body)
	err = jsonDecoder.Decode(avatar)
	if err != nil {
		return nil, fmt.Errorf("avatar/avatar.go/GetAvatar failed to decode avatar: %w", err)
	}

	return avatar, nil
}

type Avatar struct {
	Scales              Scales     `json:"scales"`
	PlayerAvatarType    string     `json:"playerAvatarType"`
	BodyColors          BodyColors `json:"bodyColors"`
	Assets              []Asset    `json:"assets"`
	DefaultShirtApplied bool       `json:"defaultShirtApplied"`
	DefaultPantsApplied bool       `json:"defaultPantsApplied"`
	Emotes              []Emote    `json:"emotes"`
}

type Asset struct {
	ID               int64     `json:"id"`
	Name             string    `json:"name"`
	AssetType        AssetType `json:"assetType"`
	CurrentVersionID int64     `json:"currentVersionId"`
	Meta             *Meta     `json:"meta,omitempty"`
}

type AssetType struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Meta struct {
	Order     int64 `json:"order"`
	Puffiness int64 `json:"puffiness"`
	Version   int64 `json:"version"`
}

type BodyColors struct {
	HeadColorID     int64 `json:"headColorId"`
	TorsoColorID    int64 `json:"torsoColorId"`
	RightArmColorID int64 `json:"rightArmColorId"`
	LeftArmColorID  int64 `json:"leftArmColorId"`
	RightLegColorID int64 `json:"rightLegColorId"`
	LeftLegColorID  int64 `json:"leftLegColorId"`
}

type Emote struct {
	AssetID   int64  `json:"assetId"`
	AssetName string `json:"assetName"`
	Position  int64  `json:"position"`
}

// Golang has no double type
// would use float64 but
// i want less ram usage
type Scales struct {
	Height     float32 `json:"height"`
	Width      float32 `json:"width"`
	Head       float32 `json:"head"`
	Depth      float32 `json:"depth"`
	Proportion float32 `json:"proportion"`
	BodyType   float32 `json:"bodyType"`
}
