package cmds

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dkaman/recordbaux/internal/db"
	"github.com/dkaman/recordbaux/internal/db/shelf"
)

type ShelvesLoadedMsg struct {
	Shelves []*shelf.Entity
	Err     error
}

type ShelfLoadedMsg struct {
	Shelf *shelf.Entity
	Err   error
}

type ShelfSavedMsg struct {
	Err error
}

type ShelfDeletedMsg struct {
	Err  error
}

type shelfDB db.Repository[*shelf.Entity]

func GetAllShelvesCmd(repo shelfDB, logger *slog.Logger) tea.Cmd {
	l := logger.WithGroup("getallshelvescmd")
	return func() tea.Msg {
		ss, err := repo.All()
		if err != nil {
			l.Error("repo error",
				slog.String("error", err.Error()),
			)
			return ShelvesLoadedMsg{Shelves: nil, Err: err}
		}

		return ShelvesLoadedMsg{Shelves: ss, Err: err}
	}
}

func GetShelfCmd(repo shelfDB, id uint, logger *slog.Logger) tea.Cmd {
	l := logger.WithGroup("getshelfcmd")
	return func() tea.Msg {
		s, err := repo.Get(id)
		l.Info("s",
			"s", s,
		)
		return ShelfLoadedMsg{Shelf: s, Err: err}
	}
}

func SaveShelfCmd(repo shelfDB, e *shelf.Entity, logger *slog.Logger) tea.Cmd {
	l := logger.WithGroup("saveshelfcmd")
	return func() tea.Msg {
		err := repo.Save(e)
		if err != nil {
			l.Error("error saving shelf to repo",
				slog.String("error", err.Error()),
			)
		}

		return ShelfSavedMsg{Err: err}
	}
}

func DeleteShelfCmd(repo shelfDB, id uint, logger *slog.Logger) tea.Cmd {
	l := logger.WithGroup("deleteshelfcmd")
	return func() tea.Msg {
		err := repo.Delete(id)
		if err != nil {
			l.Error("error deleting shelf",
				slog.String("error", err.Error()),
				slog.Any("id", id),
			)

		}
		return ShelfDeletedMsg{Err: err}
	}
}
