package bin

import (
	"github.com/dkaman/recordbaux/internal/db/record"
)

var (
	SortAlphaByArtist = "alpha_by_artist"
)

type sortFunc func(i, j *record.Entity) bool

var sortRegistry = map[string]sortFunc {
	SortAlphaByArtist: AlphaByArtist,
}

func AlphaByArtist(i, j *record.Entity) bool {
	return i.Artists[0] < j.Artists[0]
}
