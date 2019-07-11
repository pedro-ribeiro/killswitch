package restapi

import (
	"encoding/json"
	"fmt"
	"killswitch/features"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func BindAPI(port string, store features.FeatureCrudifier, failed chan bool) {
	log.Println("Starting REST API")

	router := mux.NewRouter()
	router.HandleFunc("/features/{key}", func(w http.ResponseWriter, r *http.Request) { featureGetter(store, w, r) }).Methods("GET")
	router.HandleFunc("/features", func(w http.ResponseWriter, r *http.Request) { featureUpserter(store, w, r) }).Methods("POST")

	go func() {
		log.Fatal(http.ListenAndServe(":"+port, router))
		failed <- true
	}()
}

type GenericErrorResponse struct {
	Err string `json:"error"`
}

type FeatureGetterRequest struct {
	key string
}

type FeatureGetterResponse struct {
	Key         string `json:"key"`
	Description string `json:"description"`
	IsActive    bool   `json:"isActive"`
}

type FeatureUpsertRequest struct {
	Key         string `json:"key"`
	Description string `json:"description"`
	IsActive    bool   `json:"isActive"`
}

type FeatureUpsertResponse struct {
	Key         string `json:"key"`
	Description string `json:"description"`
	IsActive    bool   `json:"isActive"`
}

func featureGetter(store features.FeatureCrudifier, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	featureKey := vars["key"]

	response, err := getFeatureFromStore(store, &FeatureGetterRequest{featureKey})

	if err != nil {
		log.Printf("Could not retrieve feature with key '%s'\n", featureKey)
		answerWithError(w, err)
		return
	}

	answerWithOk(w, response)
}

func featureUpserter(store features.FeatureCrudifier, w http.ResponseWriter, r *http.Request) {
	upsertRequest := FeatureUpsertRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&upsertRequest)

	if err != nil {
		log.Printf("Could not unmarshal request: %s", err)
		answerWithError(w, err)
		return
	}

	response, err := upsertFeatureInStore(store, &upsertRequest)

	if err != nil {
		log.Printf("Could not upsert feature: %s", err)
		answerWithError(w, err)
		return
	}

	answerWithOk(w, response)
}

func getFeatureFromStore(store features.FeatureCrudifier, request *FeatureGetterRequest) (FeatureGetterResponse, error) {
	feature, err := store.GetFeature(request.key)

	if err != nil {
		if _, ok := err.(features.NotFoundError); ok {
			return FeatureGetterResponse{}, RestNotFoundError{fmt.Sprintf("Did not find feature with key '%s'", request.key)} //FIXME: should I be returning a reference here? I'd rather have it be a VO
		}
	}

	return FeatureGetterResponse{
		Key:         feature.Key,
		Description: feature.Description,
		IsActive:    feature.IsActive,
	}, nil
}

func upsertFeatureInStore(store features.FeatureCrudifier, request *FeatureUpsertRequest) (FeatureUpsertResponse, error) {
	feature, err := store.UpsertFeature(features.Feature{
		Key:         request.Key,
		Description: request.Description,
		IsActive:    request.IsActive,
	})

	if err != nil {
		return FeatureUpsertResponse{}, err
	}

	return FeatureUpsertResponse{
		Key:         feature.Key,
		Description: feature.Description,
		IsActive:    feature.IsActive,
	}, nil
}

func answerWithOk(w http.ResponseWriter, okResponse interface{}) {
	w.WriteHeader(http.StatusOK)

	//FIXME: convert to json.Coder
	result, err := json.Marshal(okResponse)

	if err != nil {
		answerWithError(w, err)
		return
	}

	w.Write(result)
}

func answerWithError(w http.ResponseWriter, err error) {
	if _, ok := err.(RestNotFoundError); ok {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	result, err := json.Marshal(GenericErrorResponse{err.Error()})
	if err != nil {
		result, _ = json.Marshal(GenericErrorResponse{"unknown error"})
	}

	w.Write(result)
}

type RestNotFoundError struct {
	s string
}

func (e RestNotFoundError) Error() string {
	return e.s
}
