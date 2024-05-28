package main

import (
	"encoding/json"
	"log"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func handleRequests(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	review := admissionv1.AdmissionReview{
		Request: &admissionv1.AdmissionRequest{},
	}

	err := decoder.Decode(&review)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	review.Response = &admissionv1.AdmissionResponse{
		UID:     review.Request.UID,
		Allowed: false, //TODO: Replace this with actual authorization logic
		Result: &metav1.Status{
			Message: "Unauthorized",
		},
	}

	response, err := json.Marshal(review)
	if err != nil {
		http.Error(w, "Error creating response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func main() {
	http.HandleFunc("/validate", handleRequests)
	log.Fatal(http.ListenAndServeTLS(":443", "/etc/webhook/certs/tls.crt", "/etc/webhook/certs/tls.key", nil))
}
