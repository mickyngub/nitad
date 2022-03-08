package subcategory

import (
	"context"
	"mime/multipart"

	"github.com/birdglove2/nitad-backend/database"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/gcp"
	"github.com/birdglove2/nitad-backend/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	ListSubcategory(ctx context.Context) ([]Subcategory, errors.CustomError)
	ListUnsetSubcategory(ctx context.Context) ([]Subcategory, errors.CustomError)
	GetSubcategoryById(ctx context.Context, id string) (*Subcategory, errors.CustomError)
	AddSubcategory(ctx context.Context, subcategoryDTO *SubcategoryDTO) (*Subcategory, errors.CustomError)
	EditSubcategory(ctx *fiber.Ctx, subcate *SubcategoryDTO) (*Subcategory, errors.CustomError)
	DeleteSubcategory(ctx context.Context, id string) errors.CustomError

	InsertToCategory(ctx *fiber.Ctx, subcate *Subcategory, categoryId primitive.ObjectID) (*Subcategory, errors.CustomError)
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

func (s *subcategoryService) ListSubcategory(ctx context.Context) ([]Subcategory, errors.CustomError) {
	return s.repository.ListSubcategory(ctx)
}

func (s *subcategoryService) ListUnsetSubcategory(ctx context.Context) ([]Subcategory, errors.CustomError) {
	return s.repository.ListUnsetSubcategory(ctx)
}

func (s *subcategoryService) GetSubcategoryById(ctx context.Context, id string) (*Subcategory, errors.CustomError) {
	oid, err := database.ExtractOID(id)
	if err != nil {
		return nil, err
	}
	return s.repository.GetSubcategoryById(ctx, oid)
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
		s.gcpService.DeleteFile(ctx, imageFilename, collectionName)
		return addedSubcategory, err
	}

	return addedSubcate, nil
}

//TODO fix logic oldSubcate and newSubcate
func (s *subcategoryService) EditSubcategory(ctx *fiber.Ctx, subcateDTO *SubcategoryDTO) (*Subcategory, errors.CustomError) {
	editedSubcate := new(Subcategory)
	oldSubcate, err := s.GetSubcategoryById(ctx.Context(), subcateDTO.ID.Hex())
	if err != nil {
		return editedSubcate, err
	}

	imageURL, err := s.HandleUpdateImage(ctx, oldSubcate.Image, subcateDTO.Image)
	if err != nil {
		return editedSubcate, err
	}

	err = utils.CopyStruct(subcateDTO, editedSubcate)
	if err != nil {
		return editedSubcate, err
	}

	editedSubcate.Image = imageURL
	// if subcateDTO.CategoryId != primitive.NilObjectID {
	editedSubcate.CategoryId = subcateDTO.CategoryId
	// } else {
	// editedSubcate.CategoryId = oldSubcate.CategoryId
	// }

	editedSubcate, err = s.repository.EditSubcategory(ctx.Context(), editedSubcate)
	if err != nil {
		return editedSubcate, err
	}
	return editedSubcate, err
}

func (s *subcategoryService) DeleteSubcategory(ctx context.Context, id string) errors.CustomError {
	subcate, err := s.GetSubcategoryById(ctx, id)
	if err != nil {
		return err
	}

	if subcate.CategoryId != primitive.NilObjectID {
		return errors.NewBadRequestError("Unable to delete subcategory that is still in categeoryId " + subcate.CategoryId.Hex())
	}

	s.gcpService.DeleteFile(ctx, subcate.Image, collectionName)

	return s.repository.DeleteSubcategory(ctx, subcate.ID)

}

func (p *subcategoryService) HandleUpdateImage(ctx *fiber.Ctx, oldImageURL string, newImageFile *multipart.FileHeader) (string, errors.CustomError) {
	if newImageFile == nil {
		return oldImageURL, nil
	}

	p.gcpService.DeleteFile(ctx.Context(), oldImageURL, collectionName)
	newUploadImageURL, err := p.gcpService.UploadFile(ctx.Context(), newImageFile, collectionName)
	if err != nil {
		p.gcpService.DeleteFile(ctx.Context(), newUploadImageURL, collectionName)
		return oldImageURL, err
	}
	return newUploadImageURL, nil
}

// func (s *subcategoryService) HandleUpdateImage(ctx *fiber.Ctx, subcate *Subcategory) (*Subcategory, errors.CustomError) {
// 	oldSubcategory, err := s.repository.GetSubcategoryById(ctx.Context(), subcate.ID)
// 	if err != nil {
// 		return subcate, err
// 	}

// 	files, err := utils.ExtractUpdatedFiles(ctx, "image")
// 	if err != nil {
// 		return subcate, err
// 	}

// 	subcate.Image = oldSubcategory.Image
// 	// if there is file passed, delete the old one and upload a new one
// 	if len(files) > 0 {
// 		newUploadFilename, err := collections_helper.HandleUpdateSingleFile(s.gcpService, ctx.Context(), files[0], subcate.Image, collectionName)
// 		if err != nil {
// 			return subcate, err
// 		}
// 		// if upload success, pass the url to the subcategory struct
// 		subcate.Image = newUploadFilename
// 	}

// 	subcate.CreatedAt = oldSubcategory.CreatedAt
// 	return subcate, nil
// }

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

// func (s *subcategoryService) FindByIds2(ctx context.Context, sids []string) ([]Subcategory, []primitive.ObjectID, errors.CustomError) {
// 	var objectIds []primitive.ObjectID
// 	var subcategories []Subcategory

// 	sids = utils.RemoveDuplicateIds(sids)

// 	for _, sid := range sids {
// 		oid, err := utils.IsValidObjectId(sid)
// 		if err != nil {
// 			return subcategories, objectIds, err
// 		}

// 		subcate, err := s.repository.GetSubcategoryById(ctx, oid)
// 		if err != nil {
// 			return subcategories, objectIds, err
// 		}
// 		objectIds = append(objectIds, oid)
// 		subcategories = append(subcategories, *subcate)
// 	}

// 	return subcategories, objectIds, nil
// }

func (s *subcategoryService) InsertToCategory(ctx *fiber.Ctx, subcate *Subcategory, categoryId primitive.ObjectID) (*Subcategory, errors.CustomError) {
	return s.repository.InsertToCategory(ctx.Context(), subcate, categoryId)
}
