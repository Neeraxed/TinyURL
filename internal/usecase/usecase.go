package usecase

import (
	"tinyurl/internal/repository"
)

type Usecase struct {
	st *repository.Storage
}

func NewUsecase(st *repository.Storage) *Usecase {
	return &Usecase{
		st: st,
	}
}

func (uc *Usecase) GetLink(link string) (string, error) {
	return uc.st.GetFromDB(link)
}
func (uc *Usecase) AddLink(id, link string) error {
	return uc.st.AddToDB(id, link)
}
