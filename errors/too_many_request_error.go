package errors

import "github.com/gofiber/fiber/v2"

// func NewTooManyRequestsError() CustomError {
// 	return &tooManyRequestsError{}
// }

type tooManyRequestsError struct{}

func (t *tooManyRequestsError) Code() int {
	return fiber.StatusTooManyRequests
}

func (t *tooManyRequestsError) Error() string {
	return "You have requested too many in a single time-frame! Please wait another minute!"
}
