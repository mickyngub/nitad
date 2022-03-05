package subcategory

import (
	"context"
	"mime/multipart"

	"github.com/birdglove2/nitad-backend/errors"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	ListSubcategory(ctx context.Context) ([]Subcategory, errors.CustomError)
	GetSubcategoryById(ctx context.Context, oid primitive.ObjectID) (*Subcategory, errors.CustomError)
	AddSubcategory(ctx context.Context, subcate *Subcategory) (*Subcategory, errors.CustomError)
	EditSubcategory(ctx context.Context, subcate *Subcategory) (*Subcategory, errors.CustomError)
	DeleteSubcategory(ctx context.Context, oid primitive.ObjectID) errors.CustomError
}

type Service interface {
	ListSubcategory(ctx context.Context) ([]Subcategory, errors.CustomError)
	GetSubcategoryById(ctx context.Context, oid primitive.ObjectID) (*Subcategory, errors.CustomError)
	AddSubcategory(ctx context.Context, files []*multipart.FileHeader, subcate *Subcategory) (*Subcategory, errors.CustomError)
	EditSubcategory(ctx *fiber.Ctx, subcate *Subcategory) (*Subcategory, errors.CustomError)
	DeleteSubcategory(ctx context.Context, oid primitive.ObjectID) errors.CustomError

	HandleUpdateImage(ctx *fiber.Ctx, subcate *Subcategory) (*Subcategory, errors.CustomError)
	FindByIds2(ctx context.Context, sids []string) ([]Subcategory, []primitive.ObjectID, errors.CustomError)
}
