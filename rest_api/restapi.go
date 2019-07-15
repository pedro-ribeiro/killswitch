package restapi

import (
	"encoding/json"
	"fmt"
	"killswitch/features"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func BindAPI(port string, store features.FeatureStore, failed chan bool) {
	log.Println("Starting REST API")

	router := mux.NewRouter()
	router.HandleFunc("/features/{key}", func(w http.ResponseWriter, r *http.Request) { featureGetter(store, w, r) }).Methods("GET")
	router.HandleFunc("/features", func(w http.ResponseWriter, r *http.Request) { allFeaturesGetter(store, w, r) }).Methods("GET")
	router.HandleFunc("/features", func(w http.ResponseWriter, r *http.Request) { featureUpserter(store, w, r) }).Methods("PUT")

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

type FeatureResponse struct {
	Key         string `json:"key"`
	Description string `json:"description"`
	IsActive    bool   `json:"isActive"`
}

type AllFeaturesGetterResponse struct {
	Results map[string]FeatureResponse `json:"results"`
}

type FeatureUpsertRequest struct {
	Key         string `json:"key"`
	Description string `json:"description"`
	IsActive    bool   `json:"isActive"`
}

func allFeaturesGetter(store features.FeatureStore, w http.ResponseWriter, r *http.Request) {
	response, err := getAllFeaturesFromStore(store)

	if err != nil {
		fmt.Errorf("Could not retrieve all features")
		answerWithError(w, err)
		return
	}

	answerWithOk(w, response)
}

func featureGetter(store features.FeatureStore, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	featureKey := vars["key"]

	response, err := getFeatureFromStore(store, &FeatureGetterRequest{featureKey})

	if err != nil {
		fmt.Errorf("Could not retrieve feature with key '%s'\n", featureKey)
		answerWithError(w, err)
		return
	}

	answerWithOk(w, response)
}

func featureUpserter(store features.FeatureStore, w http.ResponseWriter, r *http.Request) {
	upsertRequest := FeatureUpsertRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&upsertRequest)

	if err != nil {
		fmt.Errorf("Could not unmarshal request: %s", err)
		answerWithError(w, err)
		return
	}

	response, err := upsertFeatureInStore(store, &upsertRequest)

	if err != nil {
		fmt.Errorf("Could not upsert feature: %s", err)
		answerWithError(w, err)
		return
	}

	answerWithOk(w, response)
}

func getFeatureFromStore(store features.FeatureStore, request *FeatureGetterRequest) (FeatureResponse, error) {
	feature, err := store.GetFeatureByKey(request.key)

	if err != nil {
		if _, ok := err.(features.NotFoundError); ok {
			return FeatureResponse{}, RestNotFoundError{fmt.Sprintf("Did not find feature with key '%s'", request.key)} //FIXME: should I be returning a reference here? I'd rather have it be a VO
		}
	}

	return responseFromFeature(feature), nil
}

func getAllFeaturesFromStore(store features.FeatureStore) (AllFeaturesGetterResponse, error) {
	values, err := store.GetAllFeatures()

	if err != nil {
		return AllFeaturesGetterResponse{}, err
	}

	results := make(map[string]FeatureResponse)

	for k, f := range values {
		results[k] = responseFromFeature(f)
	}

	return AllFeaturesGetterResponse{results}, nil
}

func upsertFeatureInStore(store features.FeatureStore, request *FeatureUpsertRequest) (FeatureResponse, error) {
	feature, err := store.UpsertFeature(features.Feature{
		Key:         request.Key,
		Description: request.Description,
		IsActive:    request.IsActive,
	})

	if err != nil {
		return FeatureResponse{}, err
	}

	return responseFromFeature(feature), nil
}

func responseFromFeature(value features.Feature) FeatureResponse {
	return FeatureResponse{
		Key:         value.Key,
		Description: value.Description,
		IsActive:    value.IsActive,
	}
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
