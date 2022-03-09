package subcategory

import (
	"context"
	"log"
	"mime/multipart"

	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/gcp"
	"github.com/birdglove2/nitad-backend/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	ListSubcategory(ctx context.Context) ([]*Subcategory, errors.CustomError)
	ListUnsetSubcategory(ctx context.Context) ([]*Subcategory, errors.CustomError)
	GetSubcategoryById(ctx context.Context, id string) (*Subcategory, errors.CustomError)
	AddSubcategory(ctx context.Context, subcategoryDTO *SubcategoryDTO) (*Subcategory, errors.CustomError)
	EditSubcategory(ctx *fiber.Ctx, subcate *SubcategoryDTO) (*Subcategory, errors.CustomError)
	DeleteSubcategory(ctx context.Context, id string) errors.CustomError

	InsertToCategory(ctx context.Context, subcate *Subcategory, categoryId primitive.ObjectID) (*Subcategory, errors.CustomError)
	HandleUpdateImage(ctx *fiber.Ctx, oldImageURL string, newImageFile *multipart.FileHeader) (string, errors.CustomError)
	// FindByIds2(ctx context.Context, sids []string) ([]Subcategory, []primitive.ObjectID, errors.CustomError)
	FindByIds3(ctx context.Context, sids []string) ([]Subcategory, []primitive.ObjectID, errors.CustomError)
}

type subcategoryService struct {
	repository Repository
	gcpService gcp.Uploader
}

func NewService(repository Repository, gcpService gcp.Uploader) Service {
	return &subcategoryService{repository, gcpService}

}

func (s *subcategoryService) ListSubcategory(ctx context.Context) ([]*Subcategory, errors.CustomError) {
	subcates, err := s.repository.ListSubcategory(ctx)
	if err != nil {
		return nil, err
	}

	for _, subcate := range subcates {
		subcate.Image = gcp.GetURL(subcate.Image, collectionName)
	}

	return subcates, nil
}

func (s *subcategoryService) ListUnsetSubcategory(ctx context.Context) ([]*Subcategory, errors.CustomError) {
	subcates, err := s.repository.ListUnsetSubcategory(ctx)
	if err != nil {
		return nil, err
	}

	for _, subcate := range subcates {
		subcate.Image = gcp.GetURL(subcate.Image, collectionName)
	}
	return subcates, nil
}

func (s *subcategoryService) GetSubcategoryById(ctx context.Context, id string) (*Subcategory, errors.CustomError) {
	oid, err := database.ExtractOID(id)
	if err != nil {
		return nil, err
	}
	subcate, err := s.repository.GetSubcategoryById(ctx, oid)
	if err != nil {
		return nil, err
	}

	subcate.Image = gcp.GetURL(subcate.Image, collectionName)

	return subcate, nil
}

func (s *subcategoryService) AddSubcategory(ctx context.Context, subcategoryDTO *SubcategoryDTO) (*Subcategory, errors.CustomError) {
	addedSubcategory := new(Subcategory)
	imageFilename, err := s.gcpService.UploadFile(ctx, subcategoryDTO.Image, collectionName)
	if err != nil {
		return addedSubcategory, err
	}
	addedSubcategory.Image = imageFilename

	err = utils.CopyStruct(subcategoryDTO, addedSubcategory)
	if err != nil {
		return addedSubcategory, err
	}

	addedSubcate, err := s.repository.AddSubcategory(ctx, addedSubcategory)
	if err != nil {
		s.gcpService.DeleteFile(ctx, imageFilename)
		return addedSubcategory, err
	}

	return addedSubcate, nil
}

//TODO fix logic oldSubcate and newSubcate
func (s *subcategoryService) EditSubcategory(ctx *fiber.Ctx, subcateDTO *SubcategoryDTO) (*Subcategory, errors.CustomError) {
	editedSubcate := new(Subcategory)
	oldSubcate, err := s.repository.GetSubcategoryById(ctx.Context(), subcateDTO.ID)
	if err != nil {
		return editedSubcate, err
	}

	editedImageFilename, err := s.HandleUpdateImage(ctx, oldSubcate.Image, subcateDTO.Image)
	if err != nil {
		return editedSubcate, err
	}

	err = utils.CopyStruct(subcateDTO, editedSubcate)
	if err != nil {
		return editedSubcate, err
	}

	editedSubcate.Image = editedImageFilename
	editedSubcate.CategoryId = subcateDTO.CategoryId

	editedSubcate, err = s.repository.EditSubcategory(ctx.Context(), editedSubcate)
	if err != nil {
		return editedSubcate, err
	}
	return editedSubcate, err
}

func (s *subcategoryService) DeleteSubcategory(ctx context.Context, id string) errors.CustomError {
	oid, err := database.ExtractOID(id)
	if err != nil {
		return err
	}
	subcate, err := s.repository.GetSubcategoryById(ctx, oid)
	if err != nil {
		return err
	}

	if subcate.CategoryId != primitive.NilObjectID {
		return errors.NewBadRequestError("Unable to delete subcategory that is still in categeoryId " + subcate.CategoryId.Hex())
	}

	s.gcpService.DeleteFile(ctx, subcate.Image)

	return s.repository.DeleteSubcategory(ctx, subcate.ID)

}

func (p *subcategoryService) HandleUpdateImage(ctx *fiber.Ctx, oldImageURL string, newImageFile *multipart.FileHeader) (string, errors.CustomError) {
	if newImageFile == nil {
		return oldImageURL, nil
	}

	log.Println("new", oldImageURL)

	p.gcpService.DeleteFile(ctx.Context(), oldImageURL)
	newUploadImageURL, err := p.gcpService.UploadFile(ctx.Context(), newImageFile, collectionName)
	if err != nil {
		p.gcpService.DeleteFile(ctx.Context(), newUploadImageURL)
		return oldImageURL, err
	}
	return newUploadImageURL, nil
}

func (s *subcategoryService) FindByIds3(ctx context.Context, sids []string) ([]Subcategory, []primitive.ObjectID, errors.CustomError) {
	var subcategories []Subcategory
	var objectIDs []primitive.ObjectID

	for _, sid := range sids {
		subcate, err := s.GetSubcategoryById(ctx, sid)
		if err != nil {
			return subcategories, objectIDs, err
		}
		objectIDs = append(objectIDs, subcate.ID)
		subcategories = append(subcategories, *subcate)
	}

	return subcategories, objectIDs, nil
}

func (s *subcategoryService) InsertToCategory(ctx context.Context, subcate *Subcategory, categoryId primitive.ObjectID) (*Subcategory, errors.CustomError) {
	return s.repository.InsertToCategory(ctx, subcate, categoryId)
}
