package cmds

type TransitionToMainMenuMsg struct {}

type TransitionToLoadedShelfMsg struct {
	ShelfID uint
}

type TransitionToLoadedPlaylistMsg struct {
	PlaylistID uint
}

type TransitionToLoadedBinMsg struct {
	BinID uint
}

type TransitionToCreatePlaylistMsg struct {
	ShelfIDs []uint
}
