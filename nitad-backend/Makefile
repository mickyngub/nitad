dev:
	air run main.go


test:
	APP_ENV=test go test -v ./.../.../test/...

mock-build:
	mockgen -package subcategory_test \
	-destination api/setup/mock_gcp.go \
	github.com/birdglove2/nitad-backend/gcp Uploader