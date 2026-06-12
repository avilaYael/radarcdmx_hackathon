package ubicacion

import (
	"encoding/json"
	"log"

	"github.com/guregu/null/v6"

	"github.com/mklfarha/radarcdmx/backend/rcapi/entity/mapper"
)

type Ubicacion struct {
	Entidad      null.String `json:"entidad"`
	Municipio    null.String `json:"municipio"`
	Localidad    null.String `json:"localidad"`
	Manzana      null.Int64  `json:"manzana"`
	CodigoPostal null.String `json:"codigo_postal"`
	Calle        null.String `json:"calle"`
	NumExt       null.String `json:"num_ext"`
	NumInt       null.String `json:"num_int"`
	Latitud      null.Float  `json:"latitud"`
	Longitud     null.Float  `json:"longitud"`
}

func (e Ubicacion) String() string {
	res, _ := json.Marshal(e)
	return string(res)
}

func (e Ubicacion) PrimaryKeyValues() []string {
	return []string{}
}

func UbicacionFromJSON(data json.RawMessage) Ubicacion {
	entity := Ubicacion{}
	if data == nil {
		return entity
	}
	if len(data) == 0 {
		return entity
	}

	if err := json.Unmarshal(data, &entity); err != nil {
		if err2 := mapper.FlexibleUnmarshal(data, &entity); err2 != nil {
			log.Printf("flexible unmarshal error UbicacionFromJSON: %v\n", err2)
		}
	}
	return entity
}

func UbicacionSliceFromJSON(data json.RawMessage) []Ubicacion {
	entity := []Ubicacion{}
	if data == nil {
		return entity
	}
	if len(data) == 0 {
		return entity
	}

	if err := json.Unmarshal(data, &entity); err != nil {
		entity = []Ubicacion{}
		var rawSlice []json.RawMessage
		if err2 := json.Unmarshal(data, &rawSlice); err2 == nil {
			for _, raw := range rawSlice {
				item := Ubicacion{}
				if err3 := mapper.FlexibleUnmarshal(raw, &item); err3 != nil {
					log.Printf("flexible unmarshal error UbicacionSliceFromJSON item: %v\n", err3)
				}
				entity = append(entity, item)
			}
		}
	}
	return entity
}

func (e Ubicacion) ToJSON() json.RawMessage {
	res, err := json.Marshal(e)
	if err != nil {
		log.Printf("marshal error UbicacionToJSON: %v\n", err)
	}
	return res
}

func UbicacionToJSON(e Ubicacion) json.RawMessage {
	res, err := json.Marshal(e)
	if err != nil {
		log.Printf("marshal error UbicacionToJSON: %v\n", err)
	}
	return res
}

func UbicacionSliceToJSON(e []Ubicacion) json.RawMessage {
	if e == nil {
		return json.RawMessage{}
	}
	res, err := json.Marshal(e)
	if err != nil {
		log.Printf("marshal error UbicacionSliceToJSON: %v\n", err)
	}
	return res
}
