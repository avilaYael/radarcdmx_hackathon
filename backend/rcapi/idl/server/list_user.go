package server

import (
	"context"
	"fmt"
	//"encoding/json"

	usermodule "github.com/mklfarha/radarcdmx/backend/rcapi/core/module/user"
	"github.com/mklfarha/radarcdmx/backend/rcapi/core/module/user/types"
	pb "github.com/mklfarha/radarcdmx/backend/rcapi/idl/gen"

	pbmapper "github.com/mklfarha/radarcdmx/backend/rcapi/idl/mapper"

	"go.einride.tech/aip/filtering"
	"go.einride.tech/aip/ordering"
	"go.einride.tech/aip/pagination"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func BuildListUserRequest(ctx context.Context, request *pb.ListUserRequest) (types.ListRequest, *pagination.PageToken, error) {
	err := validatePageSizeForListUser(request)
	if err != nil {
		return types.ListRequest{}, nil, err
	}

	// Use pagination.PageToken for offset-based page tokens.
	pageToken, err := pagination.ParsePageToken(request)
	if err != nil {
		return types.ListRequest{}, nil, status.Errorf(codes.InvalidArgument, "invalid page token")
	}

	// parse filters
	declarations := userDeclarations()
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

func (s *server) ListUser(ctx context.Context, request *pb.ListUserRequest) (*pb.ListUserResponse, error) {

	var err error
	req, pageToken, err := BuildListUserRequest(ctx, request)
	if err != nil {

		return nil, err
	}

	// Query the storage.
	var result types.ListResponse
	if request.GetSkipCache() {
		result, err = s.core.User().List(ctx, req, usermodule.WithSkipCache())
	} else {
		result, err = s.core.User().List(ctx, req)
	}
	if err != nil {

		return nil, err
	}

	// Build the response.
	response := &pb.ListUserResponse{
		User: pbmapper.UserSliceToProto(result.User),
	}
	// Set the next page token.
	if result.HasNextPage {
		response.NextPageToken = pageToken.Next(request).String()
	}

	// Respond.
	return response, nil
}

func validatePageSizeForListUser(request *pb.ListUserRequest) error {
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

func userDeclarations() *filtering.Declarations {
	declarations, err := filtering.NewDeclarations(
		filtering.DeclareStandardFunctions(),
		// boolean values
		filtering.DeclareIdent("true", filtering.TypeBool),
		filtering.DeclareIdent("false", filtering.TypeBool),

		//user

		filtering.DeclareIdent("uuid", filtering.TypeString),

		filtering.DeclareIdent("name", filtering.TypeString),

		filtering.DeclareIdent("lastname", filtering.TypeString),

		filtering.DeclareIdent("email", filtering.TypeString),

		filtering.DeclareIdent("password", filtering.TypeString),

		filtering.DeclareEnumIdent("status", pb.UserStatus(0).Type()),

		filtering.DeclareIdent("updated_at", filtering.TypeTimestamp),

		filtering.DeclareIdent("created_by", filtering.TypeString),

		filtering.DeclareIdent("updated_by", filtering.TypeString),

		filtering.DeclareIdent("created_at", filtering.TypeTimestamp),
	)
	if err != nil {
		fmt.Printf("error creating declarions:%v", err)
	}
	return declarations
}
