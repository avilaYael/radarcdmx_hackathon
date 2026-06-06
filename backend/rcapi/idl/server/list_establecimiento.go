package server

import (
	"context"
	"fmt"
	//"encoding/json"

	establecimientomodule "github.com/mklfarha/radarcdmx/backend/rcapi/core/module/establecimiento"
	"github.com/mklfarha/radarcdmx/backend/rcapi/core/module/establecimiento/types"
	pb "github.com/mklfarha/radarcdmx/backend/rcapi/idl/gen"

	pbmapper "github.com/mklfarha/radarcdmx/backend/rcapi/idl/mapper"

	"go.einride.tech/aip/filtering"
	"go.einride.tech/aip/ordering"
	"go.einride.tech/aip/pagination"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func BuildListEstablecimientoRequest(ctx context.Context, request *pb.ListEstablecimientoRequest) (types.ListRequest, *pagination.PageToken, error) {
	err := validatePageSizeForListEstablecimiento(request)
	if err != nil {
		return types.ListRequest{}, nil, err
	}

	// Use pagination.PageToken for offset-based page tokens.
	pageToken, err := pagination.ParsePageToken(request)
	if err != nil {
		return types.ListRequest{}, nil, status.Errorf(codes.InvalidArgument, "invalid page token")
	}

	// parse filters
	declarations := establecimientoDeclarations()
	filter, err := filtering.ParseFilter(request, declarations)
	if err != nil {
		return types.ListRequest{}, nil, fmt.Errorf("error parsing filter: %w", err)
	}

	/* // enable for debugging
	if filter.CheckedExpr != nil {
		b, _ := json.Marshal(filter.CheckedExpr.Expr)
		fmt.Printf("filtering: %v \n", string(b))
	}
	*/

	orderBy, err := ordering.ParseOrderBy(request)
	if err != nil {
		return types.ListRequest{}, nil, fmt.Errorf("error parsing order by: %w", err)
	}

	/* // enable for debugging
	if orderBy.Fields != nil {
		b, _ := json.Marshal(orderBy.Fields)
		fmt.Printf("ordering: %v \n", string(b))
	}
	*/

	return types.ListRequest{
		Offset:                pageToken.Offset,
		PageSize:              request.GetPageSize(),
		Filter:                filter,
		FilteringDeclarations: declarations,
		OrderBy:               orderBy,
		IncludeFields:         request.GetIncludeFields(),
		ExcludeFields:         request.GetExcludeFields(),
	}, &pageToken, nil
}

func (s *server) ListEstablecimiento(ctx context.Context, request *pb.ListEstablecimientoRequest) (*pb.ListEstablecimientoResponse, error) {

	var err error
	req, pageToken, err := BuildListEstablecimientoRequest(ctx, request)
	if err != nil {

		return nil, err
	}

	// Query the storage.
	var result types.ListResponse
	if request.GetSkipCache() {
		result, err = s.core.Establecimiento().List(ctx, req, establecimientomodule.WithSkipCache())
	} else {
		result, err = s.core.Establecimiento().List(ctx, req)
	}
	if err != nil {

		return nil, err
	}

	// Build the response.
	response := &pb.ListEstablecimientoResponse{
		Establecimiento: pbmapper.EstablecimientoSliceToProto(result.Establecimiento),
	}
	// Set the next page token.
	if result.HasNextPage {
		response.NextPageToken = pageToken.Next(request).String()
	}

	// Respond.
	return response, nil
}

func validatePageSizeForListEstablecimiento(request *pb.ListEstablecimientoRequest) error {
	// Handle request constraints.
	const (
		defaultPageSize = 10
	)
	switch {
	case request.PageSize < 0:
		return status.Errorf(codes.InvalidArgument, "page size is negative")
	case request.PageSize == 0:
		request.PageSize = defaultPageSize
	}
	return nil
}

func establecimientoDeclarations() *filtering.Declarations {
	declarations, err := filtering.NewDeclarations(
		filtering.DeclareStandardFunctions(),
		// boolean values
		filtering.DeclareIdent("true", filtering.TypeBool),
		filtering.DeclareIdent("false", filtering.TypeBool),

		//contacto

		filtering.DeclareIdent("contacto.telefono", filtering.TypeString),

		filtering.DeclareIdent("contacto.correo", filtering.TypeString),

		filtering.DeclareIdent("contacto.sitio_web", filtering.TypeString),

		//ubicacion

		filtering.DeclareIdent("ubicacion.entidad", filtering.TypeString),

		filtering.DeclareIdent("ubicacion.municipio", filtering.TypeString),

		filtering.DeclareIdent("ubicacion.localidad", filtering.TypeString),

		filtering.DeclareIdent("ubicacion.manzana", filtering.TypeInt),

		filtering.DeclareIdent("ubicacion.codigo_postal", filtering.TypeString),

		filtering.DeclareIdent("ubicacion.calle", filtering.TypeString),

		filtering.DeclareIdent("ubicacion.num_ext", filtering.TypeString),

		filtering.DeclareIdent("ubicacion.num_int", filtering.TypeString),

		filtering.DeclareIdent("ubicacion.latitud", filtering.TypeFloat),

		filtering.DeclareIdent("ubicacion.longitud", filtering.TypeFloat),

		//establecimiento

		filtering.DeclareIdent("uuid", filtering.TypeString),

		filtering.DeclareIdent("id_denue", filtering.TypeInt),

		filtering.DeclareIdent("clee", filtering.TypeString),

		filtering.DeclareIdent("nombre", filtering.TypeString),

		filtering.DeclareIdent("razon_social", filtering.TypeString),

		filtering.DeclareIdent("per_ocu", filtering.TypeString),

		filtering.DeclareIdent("codigo_actividad", filtering.TypeInt),

		filtering.DeclareIdent("nombre_actividad", filtering.TypeString),

		filtering.DeclareIdent("uso_de_suelo", filtering.TypeString),

		filtering.DeclareIdent("clave_catastral", filtering.TypeString),

		filtering.DeclareIdent("fecha_alta", filtering.TypeTimestamp),

		filtering.DeclareIdent("created_at", filtering.TypeTimestamp),

		filtering.DeclareIdent("updated_at", filtering.TypeTimestamp),
	)
	if err != nil {
		fmt.Printf("error creating declarions:%v", err)
	}
	return declarations
}
