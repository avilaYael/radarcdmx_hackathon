package contacto

import (
	"encoding/json"
	"log"

	"github.com/guregu/null/v6"

	"github.com/mklfarha/radarcdmx/backend/rcapi/entity/mapper"
)

type Contacto struct {
	Telefono null.String `json:"telefono"`
	Correo   null.String `json:"correo"`
	SitioWeb null.String `json:"sitio_web"`
}

func (e Contacto) String() string {
	res, _ := json.Marshal(e)
	return string(res)
}

func (e Contacto) PrimaryKeyValues() []string {
	return []string{}
}

func ContactoFromJSON(data json.RawMessage) Contacto {
	entity := Contacto{}
	if data == nil {
		return entity
	}
	if len(data) == 0 {
		return entity
	}

	if err := json.Unmarshal(data, &entity); err != nil {
		if err2 := mapper.FlexibleUnmarshal(data, &entity); err2 != nil {
			log.Printf("flexible unmarshal error ContactoFromJSON: %v\n", err2)
		}
	}
	return entity
}

func ContactoSliceFromJSON(data json.RawMessage) []Contacto {
	entity := []Contacto{}
	if data == nil {
		return entity
	}
	if len(data) == 0 {
		return entity
	}

	if err := json.Unmarshal(data, &entity); err != nil {
		entity = []Contacto{}
		var rawSlice []json.RawMessage
		if err2 := json.Unmarshal(data, &rawSlice); err2 == nil {
			for _, raw := range rawSlice {
				item := Contacto{}
				if err3 := mapper.FlexibleUnmarshal(raw, &item); err3 != nil {
					log.Printf("flexible unmarshal error ContactoSliceFromJSON item: %v\n", err3)
				}
				entity = append(entity, item)
			}
		}
	}
	return entity
}

func (e Contacto) ToJSON() json.RawMessage {
	res, err := json.Marshal(e)
	if err != nil {
		log.Printf("marshal error ContactoToJSON: %v\n", err)
	}
	return res
}

func ContactoToJSON(e Contacto) json.RawMessage {
	res, err := json.Marshal(e)
	if err != nil {
		log.Printf("marshal error ContactoToJSON: %v\n", err)
	}
	return res
}

func ContactoSliceToJSON(e []Contacto) json.RawMessage {
	if e == nil {
		return json.RawMessage{}
	}
	res, err := json.Marshal(e)
	if err != nil {
		log.Printf("marshal error ContactoSliceToJSON: %v\n", err)
	}
	return res
}
