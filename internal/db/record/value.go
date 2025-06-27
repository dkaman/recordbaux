package record

import (

	discogs "github.com/dkaman/discogs-golang"
)

type Entity struct {
	Release discogs.ReleaseInstance
}
