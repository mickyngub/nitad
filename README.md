### Prerequisites

- Language - Golang 1.17
- Database - MongoDB
- Cache - Redis
- Bucket - GCP Cloud Storage
- Server - GCP App Engine

### Environment setup

- Create a `.env` file in the root directory of the project
- Copy the contents of `.env.example` to `.env`
- Update the values in `.env` with your own values

### Run

`go run main.go`

### Deployment

To deploy on to google cloud and cloud storage
need to set service account permission to these 4 following roles in [google cloud console](https://cloud.google.com/free?utm_source=google&utm_medium=cpc&utm_campaign=japac-TH-all-en-dr-BKWS-all-super-trial-EXA-dr-1605216&utm_content=text-ad-none-none-DEV_c-CRE_602292303537-ADGP_Hybrid%20%7C%20BKWS%20-%20EXA%20%7C%20Txt%20~%20GCP%20~%20General_Business%20Services%20-%20google%20cloud%20console-KWID_43700071562405490-aud-1596662389894%3Akwd-55675752867&userloc_1012728-network_g&utm_term=KW_google%20cloud%20console&gclsrc=ds&gclsrc=ds&gclid=COL0_O_ppv0CFUPb1AodpLQKXw)

1. App Enine Admin
2. Cloud Build Editor
3. Service Account User
4. Storage Object Admin
