package establecimiento

import (
	"encoding/json"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/guregu/null/v6"
	"github.com/mklfarha/radarcdmx/backend/rcapi/entity/contacto"
	"github.com/mklfarha/radarcdmx/backend/rcapi/entity/ubicacion"
	"time"

	"github.com/mklfarha/radarcdmx/backend/rcapi/entity/mapper"
)

type Establecimiento struct {
	UUID            uuid.UUID           `json:"uuid"`
	IdDenue         null.Int64          `json:"id_denue"`
	Clee            null.String         `json:"clee"`
	Nombre          null.String         `json:"nombre"`
	RazonSocial     null.String         `json:"razon_social"`
	PerOcu          null.String         `json:"per_ocu"`
	CodigoActividad null.Int64          `json:"codigo_actividad"`
	NombreActividad null.String         `json:"nombre_actividad"`
	UsoDeSuelo      null.String         `json:"uso_de_suelo"`
	ClaveCatastral  null.String         `json:"clave_catastral"`
	Contacto        contacto.Contacto   `json:"contacto"`
	Ubicacion       ubicacion.Ubicacion `json:"ubicacion"`
	FechaAlta       null.Time           `json:"fecha_alta"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
}

func (e Establecimiento) String() string {
	res, _ := json.Marshal(e)
	return string(res)
}

func (e Establecimiento) PrimaryKeyValues() []string {
	return []string{
		e.UUID.String(),
	}
}

func EstablecimientoFromJSON(data json.RawMessage) Establecimiento {
	entity := Establecimiento{}
	if data == nil {
		return entity
	}
	if len(data) == 0 {
		return entity
	}

	if err := json.Unmarshal(data, &entity); err != nil {
		if err2 := mapper.FlexibleUnmarshal(data, &entity); err2 != nil {
			fmt.Printf("flexible unmarshal error EstablecimientoFromJSON: %v\n", err2)
		}
	}
	return entity
}

func EstablecimientoSliceFromJSON(data json.RawMessage) []Establecimiento {
	entity := []Establecimiento{}
	if data == nil {
		return entity
	}
	if len(data) == 0 {
		return entity
	}

	if err := json.Unmarshal(data, &entity); err != nil {
		entity = []Establecimiento{}
		var rawSlice []json.RawMessage
		if err2 := json.Unmarshal(data, &rawSlice); err2 == nil {
			for _, raw := range rawSlice {
				item := Establecimiento{}
				if err3 := mapper.FlexibleUnmarshal(raw, &item); err3 != nil {
					fmt.Printf("flexible unmarshal error EstablecimientoSliceFromJSON item: %v\n", err3)
				}
				entity = append(entity, item)
			}
		}
	}
	return entity
}

func (e Establecimiento) ToJSON() json.RawMessage {
	res, err := json.Marshal(e)
	if err != nil {
		fmt.Printf("marshal error EstablecimientoToJSON: %v\n", err)
	}
	return res
}

func EstablecimientoToJSON(e Establecimiento) json.RawMessage {
	res, err := json.Marshal(e)
	if err != nil {
		fmt.Printf("marshal error EstablecimientoToJSON: %v\n", err)
	}
	return res
}

func EstablecimientoSliceToJSON(e []Establecimiento) json.RawMessage {
	if e == nil {
		return json.RawMessage{}
	}
	res, err := json.Marshal(e)
	if err != nil {
		fmt.Printf("marshal error EstablecimientoSliceToJSON: %v\n", err)
	}
	return res
}
