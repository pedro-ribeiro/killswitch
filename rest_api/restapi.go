package restapi

import (
	"encoding/json"
	"killswitch/features"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func BindAPI(port string, store features.FeatureCrudifier, failed chan bool) {
	log.Println("Starting REST API")

	router := mux.NewRouter()
	router.HandleFunc("/features/{key}", func(w http.ResponseWriter, r *http.Request) { featureGetter(store, w, r) })

	go func() {
		log.Fatal(http.ListenAndServe(":"+port, router))
		failed <- true
	}()
}

type GenericErrorResponse struct {
	Err string `json:"error"`
}

type featureGetterRequest struct {
	key string
}

type FeatureGetterResponse struct {
	Key         string `json:"key"`
	Description string `json:"description"`
	IsActive    bool   `json:"isActive"`
}

func featureGetter(store features.FeatureCrudifier, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	featureKey := vars["key"]

	response, err := getFeatureFromStore(store, &featureGetterRequest{featureKey})

	if err != nil {
		log.Printf("Could not retrieve feature with key '%s'\n", featureKey)
		answerWithError(w, err.Error())
		return
	}

	answerWithOk(w, response)
}

func getFeatureFromStore(store features.FeatureCrudifier, request *featureGetterRequest) (FeatureGetterResponse, error) {
	feature, err := store.GetFeature(request.key)

	if err != nil {
		return FeatureGetterResponse{}, err //FIXME: should I be returning a reference here? I'd rather have it be a VO
	}

	return FeatureGetterResponse{
		Key:         feature.Key,
		Description: feature.Description,
		IsActive:    feature.IsActive,
	}, nil
}

func answerWithOk(w http.ResponseWriter, okResponse interface{}) {
	w.WriteHeader(http.StatusOK)

	result, err := json.Marshal(okResponse)

	if err != nil {
		answerWithError(w, err.Error())
		return
	}

	w.Write(result)
}

func answerWithError(w http.ResponseWriter, errMessage string) {
	w.WriteHeader(http.StatusInternalServerError)

	result, err := json.Marshal(GenericErrorResponse{errMessage})
	if err != nil {
		result, _ = json.Marshal(GenericErrorResponse{"unknown error"})
	}

	w.Write(result)
}
