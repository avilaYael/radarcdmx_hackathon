package handlers

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

//go:embed sectores.json
var sectoresJSON []byte

type sectorItem struct {
	CodigoSector []string `json:"codigo_sector"`
	Sector       string   `json:"sector"`
}

// ListSectoresHandler serves the static sector catalog from sectores.json.
type ListSectoresHandler struct {
	logger *log.Logger
	items  []sectorItem
}

// NewListSectoresHandler constructs ListSectoresHandler and validates JSON at startup.
func NewListSectoresHandler(logger *log.Logger) (*ListSectoresHandler, error) {
	items := make([]sectorItem, 0)
	if err := json.Unmarshal(sectoresJSON, &items); err != nil {
		return nil, fmt.Errorf("invalid sectores.json: %w", err)
	}

	return &ListSectoresHandler{logger: logger, items: items}, nil
}

func (h *ListSectoresHandler) Method() string  { return http.MethodGet }
func (h *ListSectoresHandler) Pattern() string { return "/api/sectores" }

func (h *ListSectoresHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(h.items); err != nil {
		h.logger.Printf("encode sectores response: %v", err)
	}
}
