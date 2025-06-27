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
	artistIName := i.Release.BasicInfo.Artists[0].Name
	artistJName := j.Release.BasicInfo.Artists[0].Name
	return artistIName < artistJName
}
