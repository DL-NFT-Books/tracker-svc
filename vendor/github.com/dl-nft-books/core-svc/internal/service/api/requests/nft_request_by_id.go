package requests

import (
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

func NewNftRequestByIdRequest(r *http.Request) (int64, error) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return 0, errors.Wrap(err, "failed to get id from the url path")
	}

	return int64(id), nil
}
