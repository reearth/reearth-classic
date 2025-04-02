package id

import (
	"errors"

	"github.com/reearth/reearthx/idx"
)

var ErrInvalidID = idx.ErrInvalidID
var ErrInvalidSceneID = errors.New("invalid scene ID")
var ErrInvalidWorkspaceID = errors.New("invalid workspace ID")
var ErrInvalidRootLayerID = errors.New("invalid rootLayer ID")
