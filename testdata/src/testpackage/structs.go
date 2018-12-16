package vrchat

import (
	"context"
	"time"
)

// User represents user's in-game information
type User struct {
	AcceptedTOSVersion             *int      `json:"acceptedTOSVersion"`
	AllowAvatarCopying             *bool     `json:"allowAvatarCopying"`
	CurrentAvatar                  *string   `json:"currentAvatar"`
	CurrentAvatarAssetURL          *string   `json:"currentAvatarAssetUrl"`
	CurrentAvatarImageURL          *string   `json:"currentAvatarImageUrl"`
	CurrentAvatarThumbnailImageURL *string   `json:"currentAvatarThumbnailImageUrl"`
	DeveloperType                  *string   `json:"developerType"`
	DisplayName                    *string   `json:"displayName"`
	EmailVerified                  *bool     `json:"emailVerified"`
	FriendGroupNames               []string  `json:"friendGroupNames"`
	FriendKey                      *int      `json:"friendKey"`
	Friends                        []*User   `json:"friends"`
	HasBirthday                    *bool     `json:"hasBirthday"`
	HasEmail                       *bool     `json:"hasEmail"`
	HasLoggedInFromClient          *bool     `json:"hasLoggedInFromClient"`
	HasPendingEmail                *bool     `json:"hasPendingEmail"`
	HomeLocation                   *string   `json:"homeLocation"`
	ID                             *string   `json:"id"`
	IsFriend                       *bool     `json:"isFriend"`
	LastLogin                      *string   `json:"last_login"`
	ObfuscatedEmail                *string   `json:"obfuscatedEmail"`
	ObfuscatedPendingEmail         *string   `json:"obfuscatedPendingEmail"`
	PastDisplayNames               []*string `json:"pastDisplayNames"`
	Status                         *string   `json:"status"`
	StatusDescription              *string   `json:"statusDescription"`
	SteamDetails                   struct{}  `json:"steamDetails"`
	Tags                           []string  `json:"tags"`
	Unsubscribe                    *bool     `json:"unsubscribe"`
	Username                       *string   `json:"username"`
}

type credential struct {
	username string `json:"username"`
	password string `json:"password"`
}

type service struct {
}

// AvatarService handles communication with the avatar related
// methods of the VRChat API.
type AvatarService service

// Avatar represents 3D avatar resource.
type Avatar struct {
	ID             *string   `json:"id"`
	Name           *string   `json:"name"`
	Description    *string   `json:"description"`
	AuthorID       *string   `json:"authorId"`
	AuthorName     *string   `json:"authorName"`
	Tags           []*string `json:"tags"`
	AssetURL       *string   `json:"assetUrl"`
	AssetURLObject struct {
	} `json:"assetUrlObject"`
	ImageURL          *string `json:"imageUrl"`
	ThumbnailImageURL *string `json:"thumbnailImageUrl"`
	ReleaseStatus     *string `json:"releaseStatus"`
	Version           *int    `json:"version"`
	Featured          *bool   `json:"featured"`
	UnityPackages     []struct {
		ID              *string `json:"id"`
		AssetURL        *string `json:"assetUrl"`
		UnityVersion    *string `json:"unityVersion"`
		UnitySortNumber *int    `json:"unitySortNumber"`
		AssetVersion    *int    `json:"assetVersion"`
		Platform        *string `json:"platform"`
	} `json:"unityPackages"`
	UnityPackageUpdated   *bool   `json:"unityPackageUpdated"`
	UnityPackageURL       *string `json:"unityPackageUrl"`
	UnityPackageURLObject struct {
	} `json:"unityPackageUrlObject"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (s *AvatarService) GetByID(ctx context.Context, id string) error {
	return nil
}

func (s *AvatarService) Choose(ctx context.Context, id string) error {
	return nil
}

// AvatarListOptions specifies the optional parameters to the
// Avatar.List method.
type AvatarListOptions struct {
	// Order specifies the direction to sort avatars. Possible values are: ascending, descending.
	Order string `url:"order,omitempty"`

	// User specifies the avatar's owner. Possible values are: me, friends.
	User string `url:"user,omitempty"`

	// Sort specifies how to sort avatars. Possible values are: created, updated, order, _created_at, _updated_at.
	Sort string `url:"sort,omitempty"`

	// ReleaseStatus specifies the status how the avatar is released. Possible values are: public, private, hidden, all.
	ReleaseStatus string `url:"releaseStatus,omitempty"`
}

func (s *AvatarService) List(ctx context.Context, opt *AvatarListOptions) error {
	return nil
}
