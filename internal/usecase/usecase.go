package usecase

import (
	"time"

	"github.com/teris-io/shortid"
	"go.uber.org/zap"
)

type Usecase struct {
	repo Repo
	log  *zap.Logger
}

type Repo interface {
	AddToDB(id, recievedLink string) error
	GetFromDB(long_link string) (string, error)
}

func NewUsecase(repo Repo, log *zap.Logger) *Usecase {
	return &Usecase{
		repo: repo,
		log:  log,
	}
}

func (uc *Usecase) GetLink(id string) (string, error) {
	res, err := uc.repo.GetFromDB(id)
	if err != nil {
		uc.log.Info("Failed to get link from database",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
	} else {
		uc.log.Info("Pulled the link from database",
			zap.String("link", res),
			zap.Time("time", time.Now()),
		)
	}

	return res, err
}

func (uc *Usecase) AddLink(link string) (string, error) {
	id, err := shortid.Generate()
	if err != nil {
		uc.log.Error("Failed to generate link",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
	}
	err = uc.repo.AddToDB(id, link)

	if err != nil {
		uc.log.Error("Failed to add link to database",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
	} else {
		uc.log.Info("Added link to database",
			zap.String("id", id),
			zap.String("link", link),
			zap.Time("time", time.Now()),
		)
	}

	return id, err
}
