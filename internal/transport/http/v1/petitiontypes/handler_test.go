package petitiontypes

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	domain "github.com/ParkPawapon/request-api/internal/domain/petitiontype"
	apperrors "github.com/ParkPawapon/request-api/internal/pkg/errors"
	"github.com/ParkPawapon/request-api/internal/transport/http/response"
	"github.com/gin-gonic/gin"
)

func TestListPetitionTypesReturnsDataEnvelope(t *testing.T) {
	router := testRouter(NewHandler(fakeListUseCase{
		items: []domain.PetitionType{
			{ID: 2, Name: "ขอลาออก"},
			{ID: 1, Name: "คำร้องทั่วไป"},
		},
	}))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/petition-types", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var body response.Envelope[[]PetitionTypeResponse]
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !body.OK || body.Data == nil {
		t.Fatalf("expected success envelope, got %+v", body)
	}
	if len(*body.Data) != 2 {
		t.Fatalf("expected 2 petition types, got %d", len(*body.Data))
	}
	if (*body.Data)[0].PetitionTypeID != 2 || (*body.Data)[0].PetitionTypeName != "ขอลาออก" {
		t.Fatalf("unexpected first item: %+v", (*body.Data)[0])
	}
}

func TestListPetitionTypesNormalizesUseCaseError(t *testing.T) {
	router := testRouter(NewHandler(fakeListUseCase{
		err: apperrors.Internal("Internal Server Error", errors.New("postgres password leaked detail")),
	}))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/petition-types", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	var body response.Envelope[[]PetitionTypeResponse]
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.OK || body.Error == nil {
		t.Fatalf("expected error envelope, got %+v", body)
	}
	if body.Error.Message != "Internal Server Error" {
		t.Fatalf("unexpected error message %q", body.Error.Message)
	}
}

func TestListPetitionTypesFailsClosedWithoutUseCase(t *testing.T) {
	router := testRouter(NewHandler(nil))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/petition-types", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status %d, got %d", http.StatusServiceUnavailable, w.Code)
	}
}

type fakeListUseCase struct {
	items []domain.PetitionType
	err   error
}

func (u fakeListUseCase) Execute(context.Context) ([]domain.PetitionType, error) {
	if u.err != nil {
		return nil, u.err
	}
	return u.items, nil
}

func testRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	group := router.Group("/v1")
	RegisterRoutes(group, handler)
	return router
}
