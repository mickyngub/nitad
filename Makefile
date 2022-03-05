dev:
	./air run main.go


test-api:
	APP_ENV=test go test -v  ./.../.../test/... 

mock-build:
	mockgen -package subcategory_test \
	-destination api/subcategory/test/mock_gcp_test.go \
	github.com/birdglove2/nitad-backend/gcp Uploader