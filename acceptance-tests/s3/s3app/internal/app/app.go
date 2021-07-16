package app

import (
	"github.com/aws/aws-sdk-go/aws"
	awscredentials "github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/aws/aws-sdk-go/aws/session"
	"log"
	"net/http"
	"s3app/internal/credentials"

	"github.com/gorilla/mux"
)

func App(creds *credentials.S3Service) *mux.Router {
	config := aws.NewConfig().
		WithCredentials(
			awscredentials.NewStaticCredentials(
			creds.AccessKeyId,
			creds.AccessKeySecret,
			"")).
		WithRegion(creds.Region)

	session, err := session.NewSession(config)
	if err != nil {
		panic(err)
	}
	r := mux.NewRouter()

	r.HandleFunc("/", aliveness).Methods("HEAD", "GET")
	r.HandleFunc("/{file_name}", handleUpload(session, creds)).Methods("POST")
	r.HandleFunc("/{file_name}", handleDownload(session, creds)).Methods("GET")
	r.HandleFunc("/{file_name}", handleDelete(session, creds)).Methods("DELETE")

	return r
}

func aliveness(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handled aliveness test.")
	w.WriteHeader(http.StatusNoContent)
}