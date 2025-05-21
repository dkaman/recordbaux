package physical

type sortFunc func(i, j *Record) bool

func AlphaByArtist(i, j *Record) bool {
	artistIName := i.Release.BasicInfo.Artists[0].Name
	artistJName := j.Release.BasicInfo.Artists[0].Name
	return artistIName < artistJName
}
