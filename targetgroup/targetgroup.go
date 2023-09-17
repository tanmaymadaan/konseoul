package targetgroup

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"io"
	"konseoul/common"
	"log"
	"net/http"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type Handler struct {
	db        *gorm.DB
	customMux *common.CustomServeMux
}

type TargetGroup struct {
	ID          uuid.UUID `gorm:"primaryKey"`
	Name        string
	LBAlgorithm string `gorm:"column:load_balancing_algorithm"`
}

type CreateTargetGroupDTO struct {
	Name        string `json:"name" validate:"required,ascii"`
	LBAlgorithm string `json:"lbAlgorithm" validate:"required,eq=round-robin"`
}

func NewTargetGroupHandler(db *gorm.DB, customMux *common.CustomServeMux) *Handler {
	tgh := &Handler{
		db:        db,
		customMux: customMux,
	}

	baseUri := "/v1/target-group"

	// Register specific route handlers
	tgh.customMux.HandleFunc(baseUri+"/create", tgh.createTargetGroupHandler)

	return tgh
}

func (tgh *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tgh.customMux.ServeHTTP(w, r)
}

func requestBodyValidator(reqBody io.ReadCloser, v any) error {
	body, err := io.ReadAll(reqBody)
	if err != nil {
		return errors.New("bad request")
	}

	if err := json.Unmarshal(body, &v); err != nil {
		return errors.New("bad request - invalid json")
	}

	err = validate.Struct(v)
	//validationErrors := err.(validator.ValidationErrors)
	if err != nil {
		return err
	}

	return nil
}

func (tgh *Handler) createTargetGroupHandler(w http.ResponseWriter, r *http.Request) {
	var createTargetGroupDTO CreateTargetGroupDTO
	err := requestBodyValidator(r.Body, &createTargetGroupDTO)
	if err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}

	res, err := tgh.createTargetGroup(&createTargetGroupDTO)
	if err != nil {
		http.Error(w, "Something went wrong"+err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResp, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error happened in JSON marshal. Err: %s\n", err)
	}

	_, err = w.Write(jsonResp)
	if err != nil {
		log.Printf("Error happened in JSON marshal. Err: %s", err)
	}
}

func (tgh *Handler) createTargetGroup(createTargetGroupDTO *CreateTargetGroupDTO) (common.StandardResponse, error) {
	var existingTargetGroups []TargetGroup
	dbRes := tgh.db.Where(TargetGroup{Name: createTargetGroupDTO.Name}).Find(&existingTargetGroups)
	if dbRes.Error != nil {
		panic(dbRes.Error.Error())
	}

	if len(existingTargetGroups) > 0 {
		return common.StandardResponse{StatusCode: http.StatusUnprocessableEntity, Data: nil, Message: "Target group with same name exists already"}, nil
	}

	var newTargetGroup = TargetGroup{ID: uuid.New(), Name: createTargetGroupDTO.Name, LBAlgorithm: createTargetGroupDTO.LBAlgorithm}
	tgh.db.Create(&newTargetGroup)

	return common.StandardResponse{StatusCode: http.StatusOK, Data: map[string]interface{}{
		"id":          newTargetGroup.ID,
		"name":        newTargetGroup.Name,
		"lbAlgorithm": newTargetGroup.LBAlgorithm,
	}, Message: "Success"}, nil
}
