package services

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/dkaman/recordbaux/internal/db"
	"github.com/dkaman/recordbaux/internal/db/shelf"
)

type shelfDB db.Repository[*shelf.Entity]

type ShelfService struct {
	logger  *slog.Logger
	shelves shelfDB
}

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
	Err error
}

func NewShelfService(repo shelfDB, log *slog.Logger) *ShelfService {
	logger := log.WithGroup("shelfservice")
	return &ShelfService{
		logger:  logger,
		shelves: repo,
	}
}

func (s *ShelfService) GetAllShelvesCmd() tea.Cmd {
	return func() tea.Msg {
		ss, err := s.shelves.All()
		if err != nil {
			s.logger.Error("repo error",
				slog.String("error", err.Error()),
			)
			return ShelvesLoadedMsg{Shelves: nil, Err: err}
		}

		return ShelvesLoadedMsg{Shelves: ss, Err: err}
	}
}

func (s *ShelfService) GetShelfCmd(id uint) tea.Cmd {
	return func() tea.Msg {
		shlf, err := s.shelves.Get(id)
		return ShelfLoadedMsg{Shelf: shlf, Err: err}
	}
}

func (s *ShelfService) SaveShelfCmd(e *shelf.Entity) tea.Cmd {
	return func() tea.Msg {
		err := s.shelves.Save(e)
		if err != nil {
			s.logger.Error("error saving shelf to repo",
				slog.String("error", err.Error()),
			)
		}

		return ShelfSavedMsg{Err: err}
	}
}

func (s *ShelfService) DeleteShelfCmd(id uint) tea.Cmd {
	return func() tea.Msg {
		err := s.shelves.Delete(id)
		if err != nil {
			s.logger.Error("error deleting shelf",
				slog.String("error", err.Error()),
				slog.Any("id", id),
			)

		}
		return ShelfDeletedMsg{Err: err}
	}
}
