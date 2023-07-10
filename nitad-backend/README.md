## Prerequisites

- Language - Golang 1.17
- Database - MongoDB
- Cache - Redis
- Bucket - GCP Cloud Storage
- Server - GCP App Engine

## Environment setup

1. Create a `.env` file in the root directory of the project
2. Copy the contents of `.env.template` to `.env`
3. Update the values in `.env` with your own values

- `APP_ENV` = `development` or `production`
- `MONGO_URI` = your mongodb uri e.g. `mongodb://localhost:27017`
- `MONGO_DB` = your mongodb database name
- `REDIS_ENDPOINT` = your redis endpoint, e.g. `localhost:6379`
- `REDIS_PASSWORD` = your redis password

- For these 4 `GCP_SERVICE_ACCOUNT`, `GCP_WORKLOAD_IDENTITY_PROVIDER` `GCP_PROJECTID`, `GCP_BUCKETNAME`
  you need to create your project on google cloud console then you will get all of this information.

## Run

`go run main.go`

## Deployment

To deploy on to google cloud and cloud storage
need to set service account permission to these 4 following roles in [google cloud console](https://cloud.google.com/free?utm_source=google&utm_medium=cpc&utm_campaign=japac-TH-all-en-dr-BKWS-all-super-trial-EXA-dr-1605216&utm_content=text-ad-none-none-DEV_c-CRE_602292303537-ADGP_Hybrid%20%7C%20BKWS%20-%20EXA%20%7C%20Txt%20~%20GCP%20~%20General_Business%20Services%20-%20google%20cloud%20console-KWID_43700071562405490-aud-1596662389894%3Akwd-55675752867&userloc_1012728-network_g&utm_term=KW_google%20cloud%20console&gclsrc=ds&gclsrc=ds&gclid=COL0_O_ppv0CFUPb1AodpLQKXw)

1. App Enine Admin
2. Cloud Build Editor
3. Service Account User
4. Storage Object Admin
