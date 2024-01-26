package rest

import (
	"context"
	"errors"
	"strconv"

	"github.com/TcMits/aipstr"
	"github.com/TcMits/vnprovince"
	"github.com/TcMits/vnprovince-vercel/rest/proto"
	"github.com/alecthomas/participle/v2"
	"go.einride.tech/aip/pagination"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var stopIteration = errors.New("stop iteration")

func getIDFromDivision(d *vnprovince.Division) int64 {
	return int64(d.ProvinceCode + d.DistrictCode + d.WardCode)
}

type VNProvinceService struct {
	proto.UnimplementedVNProvinceServiceServer

	divisionDeclaration *aipstr.Declaration[selector]
	provinceDeclaration *aipstr.Declaration[selector]
	districtDeclaration *aipstr.Declaration[selector]
	wardDeclaration     *aipstr.Declaration[selector]
	parser              *participle.Parser[aipstr.Filter]
}

func NewVNProvinceService() *VNProvinceService {
	s := &VNProvinceService{
		divisionDeclaration: getFilterDeclaration(),
		parser:              aipstr.NewFilterParser(),
	}

	s.provinceDeclaration = aipstr.NewDeclaration(
		aipstr.WithColumns(
			aipstr.NewColumn("province_code", aipstr.Filterable[selector]()),
			aipstr.NewColumn("province_name", aipstr.Filterable[selector]()),
		),
		aipstr.WithOperatorFuncs(getBasicOperator()...),
	)

	s.districtDeclaration = aipstr.NewDeclaration(
		aipstr.WithColumns(
			aipstr.NewColumn("district_code", aipstr.Filterable[selector]()),
			aipstr.NewColumn("district_name", aipstr.Filterable[selector]()),
		),
		aipstr.WithOperatorFuncs(getBasicOperator()...),
	)

	s.wardDeclaration = aipstr.NewDeclaration(
		aipstr.WithColumns(
			aipstr.NewColumn("ward_code", aipstr.Filterable[selector]()),
			aipstr.NewColumn("ward_name", aipstr.Filterable[selector]()),
		),
		aipstr.WithOperatorFuncs(getBasicOperator()...),
	)

	return s
}

func apiDivisionFromDivision(dst *proto.Division, src *vnprovince.Division) {
	dst.Id = int32(getIDFromDivision(src))
	dst.Name = proto.DivisionResourceName{DivisionId: strconv.FormatInt(getIDFromDivision(src), 10)}.String()
	dst.ProvinceCode = int32(src.ProvinceCode)
	dst.DistrictCode = int32(src.DistrictCode)
	dst.WardCode = int32(src.WardCode)
	dst.ProvinceName = src.ProvinceName
	dst.DistrictName = src.DistrictName
	dst.WardName = src.WardName
}

func (s *VNProvinceService) ListDivisions(ctx context.Context, req *proto.ListDivisionsRequest) (*proto.ListDivisionsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	filter, err := getWhereClauseFromRequest(req, s.parser, s.divisionDeclaration)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	pageToken, _ := pagination.ParsePageToken(req)
	pageSize := req.GetPageSize()
	resp := proto.ListDivisionsResponse{Divisions: make([]*proto.Division, 0, pageSize)}
	index := 0
	toOffset := int(int32(pageToken.Offset) + pageSize)
	offset := int(pageToken.Offset)

	if err := vnprovince.EachDivision(func(d vnprovince.Division) error {
		if !filter(&d) {
			return nil
		}

		if index >= offset && index < toOffset {
			division := proto.Division{}
			apiDivisionFromDivision(&division, &d)
			resp.Divisions = append(resp.Divisions, &division)
		}

		if index >= toOffset {
			resp.NextPageToken = pageToken.Next(req).String()
			return stopIteration
		}

		index += 1
		return nil
	}); err != nil && !errors.Is(err, stopIteration) {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &resp, nil
}

func (s *VNProvinceService) GetDivision(ctx context.Context, req *proto.GetDivisionRequest) (*proto.Division, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	rn := proto.DivisionResourceName{}
	if err := errors.Join(rn.UnmarshalString(req.GetName()), rn.Validate()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	expectedId, _ := strconv.ParseInt(rn.DivisionId, 10, 64)
	resp := proto.Division{}
	if err := vnprovince.EachDivision(func(d vnprovince.Division) error {
		if getIDFromDivision(&d) == expectedId {
			apiDivisionFromDivision(&resp, &d)
			return stopIteration
		}

		return nil
	}); err != nil && !errors.Is(err, stopIteration) {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if resp.Id == 0 {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &resp, nil
}

func apiProvinceFromDivision(dst *proto.Province, src *vnprovince.Division) {
	dst.Id = int32(src.ProvinceCode)
	dst.Name = proto.ProvinceResourceName{ProvinceId: strconv.FormatInt(int64(src.ProvinceCode), 10)}.String()
	dst.ProvinceCode = int32(src.ProvinceCode)
	dst.ProvinceName = src.ProvinceName
}

func (s *VNProvinceService) ListProvinces(ctx context.Context, req *proto.ListProvincesRequest) (*proto.ListProvincesResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	filter, err := getWhereClauseFromRequest(req, s.parser, s.provinceDeclaration)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	pageToken, _ := pagination.ParsePageToken(req)
	pageSize := req.GetPageSize()
	resp := proto.ListProvincesResponse{Provinces: make([]*proto.Province, 0, pageSize)}
	index := 0
	toOffset := int(int32(pageToken.Offset) + pageSize)
	offset := int(pageToken.Offset)
	var previousProvinceCode int64 = 0

	if err := vnprovince.EachDivision(func(d vnprovince.Division) error {
		if !filter(&d) || previousProvinceCode == d.ProvinceCode {
			return nil
		}

		previousProvinceCode = d.ProvinceCode
		if index >= offset && index < toOffset {
			province := proto.Province{}
			apiProvinceFromDivision(&province, &d)
			resp.Provinces = append(resp.Provinces, &province)
		}

		if index >= toOffset {
			resp.NextPageToken = pageToken.Next(req).String()
			return stopIteration
		}

		index += 1
		return nil
	}); err != nil && !errors.Is(err, stopIteration) {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &resp, nil
}

func (s *VNProvinceService) GetProvince(ctx context.Context, req *proto.GetProvinceRequest) (*proto.Province, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	rn := proto.ProvinceResourceName{}
	if err := errors.Join(rn.UnmarshalString(req.GetName()), rn.Validate()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	expectedId, _ := strconv.ParseInt(rn.ProvinceId, 10, 64)
	resp := proto.Province{}
	if err := vnprovince.EachDivision(func(d vnprovince.Division) error {
		if d.ProvinceCode == expectedId {
			apiProvinceFromDivision(&resp, &d)
			return stopIteration
		}

		return nil
	}); err != nil && !errors.Is(err, stopIteration) {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if resp.Id == 0 {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &resp, nil
}

func apiDistrictFromDivision(dst *proto.District, src *vnprovince.Division) {
	dst.Id = int32(src.DistrictCode)
	dst.Name = proto.DistrictResourceName{
		ProvinceId: strconv.FormatInt(int64(src.ProvinceCode), 10),
		DistrictId: strconv.FormatInt(int64(src.DistrictCode), 10),
	}.String()
	dst.DistrictCode = int32(src.DistrictCode)
	dst.DistrictName = src.DistrictName
}

func (s *VNProvinceService) ListDistricts(ctx context.Context, req *proto.ListDistrictsRequest) (*proto.ListDistrictsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	parentRN := proto.ProvinceResourceName{}
	if err := errors.Join(parentRN.UnmarshalString(req.GetParent()), parentRN.Validate()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	filter, err := getWhereClauseFromRequest(req, s.parser, s.districtDeclaration)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// update filter based on parentRN
	if !parentRN.ContainsWildcard() {
		parentID, _ := strconv.ParseInt(parentRN.ProvinceId, 10, 64)
		newF, _ := eqFieldWithInt("province_code", parentID)
		filter = andSelector(filter, newF)
	}

	// update filter to group by district
	var previousDistrictCode int64 = 0
	filter = andSelector(filter, func(d *vnprovince.Division) bool {
		if previousDistrictCode == d.DistrictCode {
			return false
		}

		previousDistrictCode = d.DistrictCode
		return true
	})

	pageToken, _ := pagination.ParsePageToken(req)
	pageSize := req.GetPageSize()
	resp := proto.ListDistrictsResponse{Districts: make([]*proto.District, 0, pageSize)}
	index := 0
	toOffset := int(int32(pageToken.Offset) + pageSize)
	offset := int(pageToken.Offset)

	if err := vnprovince.EachDivision(func(d vnprovince.Division) error {
		if !filter(&d) {
			return nil
		}

		previousDistrictCode = d.DistrictCode
		if index >= offset && index < toOffset {
			district := proto.District{}
			apiDistrictFromDivision(&district, &d)
			resp.Districts = append(resp.Districts, &district)
		}

		if index >= toOffset {
			resp.NextPageToken = pageToken.Next(req).String()
			return stopIteration
		}

		index += 1
		return nil
	}); err != nil && !errors.Is(err, stopIteration) {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &resp, nil
}

func (s *VNProvinceService) GetDistrict(ctx context.Context, req *proto.GetDistrictRequest) (*proto.District, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	rn := proto.DistrictResourceName{}
	if err := errors.Join(rn.UnmarshalString(req.GetName()), rn.Validate()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	expectedProvinceId, _ := strconv.ParseInt(rn.ProvinceId, 10, 64)
	expectedDistrictId, _ := strconv.ParseInt(rn.DistrictId, 10, 64)
	resp := proto.District{}
	if err := vnprovince.EachDivision(func(d vnprovince.Division) error {
		if d.ProvinceCode == expectedProvinceId && d.DistrictCode == expectedDistrictId {
			apiDistrictFromDivision(&resp, &d)
			return stopIteration
		}

		return nil
	}); err != nil && !errors.Is(err, stopIteration) {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if resp.Id == 0 {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &resp, nil
}

func apiWardFromDivision(dst *proto.Ward, src *vnprovince.Division) {
	dst.Id = int32(src.WardCode)
	dst.Name = proto.WardResourceName{
		ProvinceId: strconv.FormatInt(int64(src.ProvinceCode), 10),
		DistrictId: strconv.FormatInt(int64(src.DistrictCode), 10),
		WardId:     strconv.FormatInt(int64(src.WardCode), 10),
	}.String()
	dst.WardCode = int32(src.WardCode)
	dst.WardName = src.WardName
}

func (s *VNProvinceService) ListWards(ctx context.Context, req *proto.ListWardsRequest) (*proto.ListWardsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	parentRN := proto.DistrictResourceName{}
	if err := errors.Join(parentRN.UnmarshalString(req.GetParent()), parentRN.Validate()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	filter, err := getWhereClauseFromRequest(req, s.parser, s.wardDeclaration)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// update filter based on parentRN
	if parentRN.ProvinceId != "-" {
		provinceID, _ := strconv.ParseInt(parentRN.ProvinceId, 10, 64)
		newF, _ := eqFieldWithInt("province_code", provinceID)
		filter = andSelector(filter, newF)
	}

	if parentRN.DistrictId != "-" {
		districtID, _ := strconv.ParseInt(parentRN.DistrictId, 10, 64)
		newF, _ := eqFieldWithInt("district_code", districtID)
		filter = andSelector(filter, newF)
	}

	// remove if ward is not available
	filter = andSelector(filter, func(d *vnprovince.Division) bool {
		return d.WardCode != 0
	})

	pageToken, _ := pagination.ParsePageToken(req)
	pageSize := req.GetPageSize()
	resp := proto.ListWardsResponse{Wards: make([]*proto.Ward, 0, pageSize)}
	index := 0
	toOffset := int(int32(pageToken.Offset) + pageSize)
	offset := int(pageToken.Offset)

	if err := vnprovince.EachDivision(func(d vnprovince.Division) error {
		if !filter(&d) {
			return nil
		}

		if index >= offset && index < toOffset {
			ward := proto.Ward{}
			apiWardFromDivision(&ward, &d)
			resp.Wards = append(resp.Wards, &ward)
		}

		if index >= toOffset {
			resp.NextPageToken = pageToken.Next(req).String()
			return stopIteration
		}

		index += 1
		return nil
	}); err != nil && !errors.Is(err, stopIteration) {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &resp, nil
}

func (s *VNProvinceService) GetWard(ctx context.Context, req *proto.GetWardRequest) (*proto.Ward, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	rn := proto.WardResourceName{}
	if err := errors.Join(rn.UnmarshalString(req.GetName()), rn.Validate()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	expectedProvinceId, _ := strconv.ParseInt(rn.ProvinceId, 10, 64)
	expectedDistrictId, _ := strconv.ParseInt(rn.DistrictId, 10, 64)
	expectedWardId, _ := strconv.ParseInt(rn.WardId, 10, 64)
	resp := proto.Ward{}
	if err := vnprovince.EachDivision(func(d vnprovince.Division) error {
		if d.ProvinceCode == expectedProvinceId && d.DistrictCode == expectedDistrictId && d.WardCode == expectedWardId {
			apiWardFromDivision(&resp, &d)
			return stopIteration
		}

		return nil
	}); err != nil && !errors.Is(err, stopIteration) {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if resp.Id == 0 {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &resp, nil
}
