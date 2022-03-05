package subcategory_test

import (
	"mime/multipart"
	"testing"

	"github.com/birdglove2/nitad-backend/api/subcategory"
	"github.com/birdglove2/nitad-backend/utils"
	"github.com/stretchr/testify/require"
)

func createDummySubcategory(t *testing.T) subcategory.Subcategory {
	dummySubcate := subcategory.Subcategory{
		Title: "dummy subcate title",
		Image: "dummy subcate image url",
	}

	adddedSubcategory, err := subcategory.Add(&dummySubcate)
	require.Equal(t, err, nil)
	require.Equal(t, dummySubcate.Title, adddedSubcategory.Title)
	require.Equal(t, dummySubcate.Image, adddedSubcategory.Image)
	require.NotEqual(t, nil, adddedSubcategory.ID)
	return *adddedSubcategory
}

func randomImages(n int) []*multipart.FileHeader {
	results := make([]*multipart.FileHeader, n)
	for i := 0; i < n; i++ {
		results[i] = &multipart.FileHeader{
			Filename: utils.RandomString(5) + ".jpg",
		}
	}
	return results
}
