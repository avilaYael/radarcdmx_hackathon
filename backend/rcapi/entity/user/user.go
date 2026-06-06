package user

import (
	"encoding/json"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/guregu/null/v6"
	"github.com/mklfarha/radarcdmx/backend/rcapi/enum"
	"time"

	"github.com/mklfarha/radarcdmx/backend/rcapi/entity/mapper"
)

type User struct {
	UUID      uuid.UUID       `json:"uuid"`
	Name      null.String     `json:"name"`
	Lastname  null.String     `json:"lastname"`
	Email     string          `json:"email"`
	Password  string          `json:"password"`
	Status    enum.UserStatus `json:"status"`
	UpdatedAt time.Time       `json:"updated_at"`
	CreatedBy uuid.UUID       `json:"created_by"`
	UpdatedBy uuid.UUID       `json:"updated_by"`
	CreatedAt time.Time       `json:"created_at"`
}

func (e User) String() string {
	res, _ := json.Marshal(e)
	return string(res)
}

func (e User) PrimaryKeyValues() []string {
	return []string{
		e.UUID.String(),
	}
}

func UserFromJSON(data json.RawMessage) User {
	entity := User{}
	if data == nil {
		return entity
	}
	if len(data) == 0 {
		return entity
	}

	if err := json.Unmarshal(data, &entity); err != nil {
		if err2 := mapper.FlexibleUnmarshal(data, &entity); err2 != nil {
			fmt.Printf("flexible unmarshal error UserFromJSON: %v\n", err2)
		}
	}
	return entity
}

func UserSliceFromJSON(data json.RawMessage) []User {
	entity := []User{}
	if data == nil {
		return entity
	}
	if len(data) == 0 {
		return entity
	}

	if err := json.Unmarshal(data, &entity); err != nil {
		entity = []User{}
		var rawSlice []json.RawMessage
		if err2 := json.Unmarshal(data, &rawSlice); err2 == nil {
			for _, raw := range rawSlice {
				item := User{}
				if err3 := mapper.FlexibleUnmarshal(raw, &item); err3 != nil {
					fmt.Printf("flexible unmarshal error UserSliceFromJSON item: %v\n", err3)
				}
				entity = append(entity, item)
			}
		}
	}
	return entity
}

func (e User) ToJSON() json.RawMessage {
	res, err := json.Marshal(e)
	if err != nil {
		fmt.Printf("marshal error UserToJSON: %v\n", err)
	}
	return res
}

func UserToJSON(e User) json.RawMessage {
	res, err := json.Marshal(e)
	if err != nil {
		fmt.Printf("marshal error UserToJSON: %v\n", err)
	}
	return res
}

func UserSliceToJSON(e []User) json.RawMessage {
	if e == nil {
		return json.RawMessage{}
	}
	res, err := json.Marshal(e)
	if err != nil {
		fmt.Printf("marshal error UserSliceToJSON: %v\n", err)
	}
	return res
}
