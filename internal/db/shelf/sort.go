package shelf

import (
	"github.com/dkaman/recordbaux/internal/db/record"
)

type sortFunc func(i, j *record.Entity) bool

var sortRegistry = map[string]sortFunc {
	"alpha_by_artist": AlphaByArtist,
}

func AlphaByArtist(i, j *record.Entity) bool {
	artistIName := i.Release.BasicInfo.Artists[0].Name
	artistJName := j.Release.BasicInfo.Artists[0].Name
	return artistIName < artistJName
}
