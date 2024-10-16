package usecase

import (
	"github.com/teris-io/shortid"
)

type Usecase struct {
	repo Repo
}

type Repo interface {
	AddToDB(id, recievedLink string) error
	GetFromDB(long_link string) (string, error)
}

func NewUsecase(repo Repo) *Usecase {
	return &Usecase{
		repo: repo,
	}
}

func (uc *Usecase) GetLink(link string) (string, error) {
	return uc.repo.GetFromDB(link)
}

func (uc *Usecase) AddLink(id, link string) error {
	return uc.repo.AddToDB(id, link)
}

func (uc *Usecase) CreateId() (string, error) {
	return shortid.Generate()
}
