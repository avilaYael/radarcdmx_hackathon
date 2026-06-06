package enum

import (
	"encoding/json"
	"fmt"
	"github.com/guregu/null/v6"
)

//go:generate go run github.com/dmarkham/enumer -type=EstablecimientoStatus -json
type EstablecimientoStatus int64

const (
	ESTABLECIMIENTO_STATUS_INVALID = iota
)

func (e EstablecimientoStatus) ToInt64() int64 {
	return int64(e)
}

func (e EstablecimientoStatus) ToNullInt() null.Int {
	return null.NewInt(int64(e), true)
}

func EstablecimientoStatusFromString(in string) EstablecimientoStatus {
	switch in {
	case "invalid":
		return ESTABLECIMIENTO_STATUS_INVALID
	}
	return ESTABLECIMIENTO_STATUS_INVALID
}

func EstablecimientoStatusFromPointerString(in *string) EstablecimientoStatus {
	if in == nil {
		return ESTABLECIMIENTO_STATUS_INVALID
	}
	return EstablecimientoStatusFromString(*in)
}

func (e EstablecimientoStatus) String() string {
	switch e {
	case ESTABLECIMIENTO_STATUS_INVALID:
		return "invalid"
	}

	return "invalid"
}

func (e EstablecimientoStatus) StringPtr() *string {
	val := e.String()
	return &val
}

func EstablecimientoStatusSliceToJSON(in []EstablecimientoStatus) json.RawMessage {
	res := make([]int64, len(in))
	for i, e := range in {
		res[i] = int64(e)
	}
	jr, err := json.Marshal(res)
	if err != nil {
		fmt.Printf("error marshaling EstablecimientoStatus slice to json: %v", err)
		return json.RawMessage{}
	}
	return jr
}

func JSONToEstablecimientoStatusSlice(in json.RawMessage) []EstablecimientoStatus {
	res := []int64{}
	err := json.Unmarshal(in, &res)
	if err != nil {
		fmt.Printf("error unmarshaling EstablecimientoStatus slice to int slice: %v", err)
		return nil
	}
	if len(res) == 0 {
		return nil
	}
	finalRes := []EstablecimientoStatus{}
	for _, r := range res {
		finalRes = append(finalRes, EstablecimientoStatus(r))
	}
	return finalRes
}
