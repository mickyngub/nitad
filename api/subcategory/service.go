package subcategory

import (
	"context"

	"github.com/birdglove2/nitad-backend/api/collections_helper"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/gcp"
	"github.com/birdglove2/nitad-backend/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	ListSubcategory(ctx context.Context) ([]Subcategory, errors.CustomError)
	ListUnsetSubcategory(ctx context.Context) ([]Subcategory, errors.CustomError)
	GetSubcategoryById(ctx context.Context, oid primitive.ObjectID) (*Subcategory, errors.CustomError)
	AddSubcategory(ctx context.Context, subcategoryDTO *SubcategoryDTO) (*Subcategory, errors.CustomError)
	EditSubcategory(ctx *fiber.Ctx, subcate *Subcategory) (*Subcategory, errors.CustomError)
	DeleteSubcategory(ctx context.Context, oid primitive.ObjectID) errors.CustomError

	HandleUpdateImage(ctx *fiber.Ctx, subcate *Subcategory) (*Subcategory, errors.CustomError)
	FindByIds2(ctx context.Context, sids []string) ([]Subcategory, []primitive.ObjectID, errors.CustomError)
	FindByIds3(ctx context.Context, osids []primitive.ObjectID) ([]Subcategory, errors.CustomError)
}

type subcategoryService struct {
	repository Repository
	gcpService gcp.Uploader
}

func NewService(repository Repository, gcpService gcp.Uploader) Service {
	return &subcategoryService{repository, gcpService}

}

func (s *subcategoryService) ListSubcategory(ctx context.Context) ([]Subcategory, errors.CustomError) {
	return s.repository.ListSubcategory(ctx)
}

func (s *subcategoryService) ListUnsetSubcategory(ctx context.Context) ([]Subcategory, errors.CustomError) {
	return s.repository.ListUnsetSubcategory(ctx)
}

func (s *subcategoryService) GetSubcategoryById(ctx context.Context, oid primitive.ObjectID) (*Subcategory, errors.CustomError) {
	return s.repository.GetSubcategoryById(ctx, oid)
}

func (s *subcategoryService) AddSubcategory(ctx context.Context, subcategoryDTO *SubcategoryDTO) (*Subcategory, errors.CustomError) {
	addedSubcategory := new(Subcategory)
	imageFilename, err := s.gcpService.UploadFile(ctx, subcategoryDTO.Image, collectionName)
	if err != nil {
		return addedSubcategory, err
	}
	addedSubcategory.Image = imageFilename

	// if subcategoryDTO.CategoryId
	err = utils.CopyStruct(subcategoryDTO, addedSubcategory)
	if err != nil {
		return addedSubcategory, err
	}

	addedSubcate, err := s.repository.AddSubcategory(ctx, addedSubcategory)
	if err != nil {
		s.gcpService.DeleteFile(ctx, imageFilename, collectionName)
		return addedSubcategory, err
	}

	return addedSubcate, nil
}

func (s *subcategoryService) EditSubcategory(ctx *fiber.Ctx, subcate *Subcategory) (*Subcategory, errors.CustomError) {
	subcate, err := s.HandleUpdateImage(ctx, subcate)
	if err != nil {
		return subcate, err
	}

	editedSubcate, err := s.repository.EditSubcategory(ctx.Context(), subcate)
	if err != nil {
		return subcate, err
	}
	return editedSubcate, err

}

func (s *subcategoryService) DeleteSubcategory(ctx context.Context, oid primitive.ObjectID) errors.CustomError {
	subcate, err := s.repository.GetSubcategoryById(ctx, oid)
	if err != nil {
		return err
	}

	if subcate.CategoryId != primitive.NilObjectID {
		return errors.NewBadRequestError("Unable to delete subcategory that is still in categeoryId " + subcate.CategoryId.Hex())
	}

	s.gcpService.DeleteFile(ctx, subcate.Image, collectionName)

	return s.repository.DeleteSubcategory(ctx, oid)

}

func (s *subcategoryService) HandleUpdateImage(ctx *fiber.Ctx, subcate *Subcategory) (*Subcategory, errors.CustomError) {
	oldSubcategory, err := s.repository.GetSubcategoryById(ctx.Context(), subcate.ID)
	if err != nil {
		return subcate, err
	}

	files, err := utils.ExtractUpdatedFiles(ctx, "image")
	if err != nil {
		return subcate, err
	}

	subcate.Image = oldSubcategory.Image
	// if there is file passed, delete the old one and upload a new one
	if len(files) > 0 {
		newUploadFilename, err := collections_helper.HandleUpdateSingleFile(s.gcpService, ctx.Context(), files[0], subcate.Image, collectionName)
		if err != nil {
			return subcate, err
		}
		// if upload success, pass the url to the subcategory struct
		subcate.Image = newUploadFilename
	}

	subcate.CreatedAt = oldSubcategory.CreatedAt
	return subcate, nil
}

func (s *subcategoryService) FindByIds3(ctx context.Context, osids []primitive.ObjectID) ([]Subcategory, errors.CustomError) {
	var subcategories []Subcategory

	for _, osid := range osids {
		subcate, err := s.repository.GetSubcategoryById(ctx, osid)
		if err != nil {
			return subcategories, err
		}
		subcategories = append(subcategories, *subcate)
	}

	return subcategories, nil
}

func (s *subcategoryService) FindByIds2(ctx context.Context, sids []string) ([]Subcategory, []primitive.ObjectID, errors.CustomError) {
	var objectIds []primitive.ObjectID
	var subcategories []Subcategory

	sids = utils.RemoveDuplicateIds(sids)

	for _, sid := range sids {
		oid, err := utils.IsValidObjectId(sid)
		if err != nil {
			return subcategories, objectIds, err
		}

		subcate, err := s.repository.GetSubcategoryById(ctx, oid)
		if err != nil {
			return subcategories, objectIds, err
		}
		objectIds = append(objectIds, oid)
		subcategories = append(subcategories, *subcate)
	}

	return subcategories, objectIds, nil
}
