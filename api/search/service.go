package search

import (
	"github.com/birdglove2/nitad-backend/api/category"
	"github.com/birdglove2/nitad-backend/api/project"
	"github.com/birdglove2/nitad-backend/errors"
)

func SearchAll() (Search, errors.CustomError) {
	search := Search{}

	categorySearch, err := category.SearchAll()
	if err != nil {
		return search, err
	}

	projectSearch, err := project.SearchAll()
	if err != nil {
		return search, err
	}

	search.Category = categorySearch
	search.Project = projectSearch

	return search, nil
}
