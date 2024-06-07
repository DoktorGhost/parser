package usecase

import "parser/internal/storage"

type UseCaseParser struct {
	storage storage.Repo
}

func NewUseCaseParser(storage storage.Repo) *UseCaseParser {
	return &UseCaseParser{storage: storage}
}

func (p *UseCaseParser) Add(name string, card storage.Card) error {
	err := p.storage.Add(name, card)
	return err
}
