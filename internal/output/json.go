package output

import (
	"encoding/json"
	"io"

	"github.com/RakhaYandra/repos-sync/internal/model"
)

func JSON(w io.Writer, r model.Report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
