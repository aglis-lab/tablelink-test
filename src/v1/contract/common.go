package contract

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"tablelink/src/app"

	"github.com/go-chi/chi/v5"
)

type UserRequest struct {
}

func ValidateParamRequest(ctx context.Context, idKey string) string {
	return chi.URLParamFromCtx(ctx, idKey)
}

func ValidateJSONRequest[
	K UserRequest,
](ctx context.Context, r *http.Request) (K, error) {
	_, span := app.Tracer().Start(ctx, "ValidateRequest")
	defer span.End()

	var payload K
	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		span.RecordError(fmt.Errorf("read request body err: %v", err))
		return payload, err
	}

	if err := json.Unmarshal(bodyByte, &payload); err != nil {
		span.RecordError(fmt.Errorf("unmarshal request body err: %v", err))
		return payload, err
	}

	if err := app.Validator().Struct(payload); err != nil {
		span.RecordError(fmt.Errorf("validate request body err: %v", err))
		return payload, err
	}

	return payload, nil
}
