package rest

import (
	"context"
	"testing"

	"github.com/TcMits/vnprovince-vercel/rest/proto"
)

func Test_vnProvinceService_ListDivisions(t *testing.T) {
	s := NewVNProvinceService()
	ctx := context.Background()

	req := &proto.ListDivisionsRequest{PageSize: -10}
	if _, err := s.ListDivisions(ctx, req); err == nil {
		t.Fatalf("ListDivisions() error = %v, wantErr %v", err, true)
	}

	req = &proto.ListDivisionsRequest{PageSize: 10}
	resp, err := s.ListDivisions(ctx, req)
	if err != nil {
		t.Fatalf("ListDivisions() error = %v, wantErr %v", err, false)
	}

	if len(resp.Divisions) != 10 {
		t.Fatalf("ListDivisions() error = %v, wantErr %v", err, false)
	}

	if resp.NextPageToken == "" {
		t.Fatalf("ListDivisions() error = %v, wantErr %v", err, false)
	}

	req = &proto.ListDivisionsRequest{PageSize: 10, Filter: "province_code=1 AND district_code=1 AND ward_code=25"}
	resp, err = s.ListDivisions(ctx, req)
	if err != nil {
		t.Fatalf("ListDivisions() error = %v, wantErr %v", err, false)
	}

	if len(resp.Divisions) != 1 {
		t.Fatalf("ListDivisions() error = %v, wantErr %v", err, false)
	}

	if resp.NextPageToken != "" {
		t.Fatalf("ListDivisions() error = %v, wantErr %v", err, false)
	}

	if resp.Divisions[0].Id != 27 {
		t.Fatalf("ListDivisions() error = %v, wantErr %v", err, false)
	}

	req = &proto.ListDivisionsRequest{PageSize: 10, Skip: 10}
	resp, err = s.ListDivisions(ctx, req)
	if err != nil {
		t.Fatalf("ListDivisions() error = %v, wantErr %v", err, false)
	}

	if len(resp.Divisions) != 10 {
		t.Fatalf("ListDivisions() error = %v, wantErr %v", err, false)
	}

	if resp.Divisions[0].Id == 3 {
		t.Fatalf("ListDivisions() error = %v, wantErr %v", err, false)
	}
}

func Test_vnProvinceService_GetDivision(t *testing.T) {
	s := NewVNProvinceService()
	ctx := context.Background()

	req := &proto.GetDivisionRequest{Name: "divisions/1"}
	if _, err := s.GetDivision(ctx, req); err == nil {
		t.Fatalf("GetDivision() error = %v, wantErr %v", err, true)
	}

	req = &proto.GetDivisionRequest{Name: "divisions/3"}
	resp, err := s.GetDivision(ctx, req)
	if err != nil {
		t.Fatalf("GetDivision() error = %v, wantErr %v", err, false)
	}

	if resp.Id != 3 {
		t.Fatalf("GetDivision() error = %v, wantErr %v", err, false)
	}
}

func Test_vnProvinceService_ListProvinces(t *testing.T) {
	s := NewVNProvinceService()
	ctx := context.Background()

	req := &proto.ListProvincesRequest{PageSize: -10}
	if _, err := s.ListProvinces(ctx, req); err == nil {
		t.Fatalf("ListProvinces() error = %v, wantErr %v", err, true)
	}

	req = &proto.ListProvincesRequest{PageSize: 10}
	resp, err := s.ListProvinces(ctx, req)
	if err != nil {
		t.Fatalf("ListProvinces() error = %v, wantErr %v", err, false)
	}

	if len(resp.Provinces) != 10 {
		t.Fatalf("ListProvinces() error = %v, wantErr %v", err, false)
	}

	if resp.NextPageToken == "" {
		t.Fatalf("ListProvinces() error = %v, wantErr %v", err, false)
	}

	req = &proto.ListProvincesRequest{PageSize: 10, Filter: "province_code=1"}
	resp, err = s.ListProvinces(ctx, req)
	if err != nil {
		t.Fatalf("ListProvinces() error = %v, wantErr %v", err, false)
	}

	if len(resp.Provinces) != 1 {
		t.Fatalf("ListProvinces() error = %v, wantErr %v", err, false)
	}

	if resp.NextPageToken != "" {
		t.Fatalf("ListProvinces() error = %v, wantErr %v", err, false)
	}

	if resp.Provinces[0].Id != 1 {
		t.Fatalf("ListProvinces() error = %v, wantErr %v", err, false)
	}

	req = &proto.ListProvincesRequest{PageSize: 10, Skip: 10}
	resp, err = s.ListProvinces(ctx, req)
	if err != nil {
		t.Fatalf("ListProvinces() error = %v, wantErr %v", err, false)
	}

	if len(resp.Provinces) != 10 {
		t.Fatalf("ListProvinces() error = %v, wantErr %v", err, false)
	}

	if resp.Provinces[0].Id == 1 {
		t.Fatalf("ListProvinces() error = %v, wantErr %v", err, false)
	}
}

func Test_vnProvinceService_GetProvince(t *testing.T) {
	s := NewVNProvinceService()
	ctx := context.Background()

	req := &proto.GetProvinceRequest{Name: "provinces/1"}
	resp, err := s.GetProvince(ctx, req)
	if err != nil {
		t.Fatalf("GetProvince() error = %v, wantErr %v", err, false)
	}

	if resp.Id != 1 {
		t.Fatalf("GetProvince() error = %v, wantErr %v", err, false)
	}
}

func Test_vnProvinceService_ListDistricts(t *testing.T) {
	s := NewVNProvinceService()
	ctx := context.Background()

	req := &proto.ListDistrictsRequest{PageSize: -10}
	if _, err := s.ListDistricts(ctx, req); err == nil {
		t.Fatalf("ListDistricts() error = %v, wantErr %v", err, true)
	}

	req = &proto.ListDistrictsRequest{PageSize: 10}
	resp, err := s.ListDistricts(ctx, req)
	if err == nil {
		t.Fatalf("ListDistricts() error = %v, wantErr %v", err, true)
	}

	req = &proto.ListDistrictsRequest{PageSize: 10, Parent: "provinces/1"}
	resp, err = s.ListDistricts(ctx, req)

	if err != nil {
		t.Fatalf("ListDistricts() error = %v, wantErr %v", err, false)
	}

	if len(resp.Districts) != 10 {
		t.Fatalf("ListDistricts() error = %v, wantErr %v", err, false)
	}

	if resp.NextPageToken == "" {
		t.Fatalf("ListDistricts() error = %v, wantErr %v", err, false)
	}

	req = &proto.ListDistrictsRequest{PageSize: 10, Parent: "provinces/1", Filter: "district_code=1"}
	resp, err = s.ListDistricts(ctx, req)
	if err != nil {
		t.Fatalf("ListDistricts() error = %v, wantErr %v", err, false)
	}

	if len(resp.Districts) != 1 {
		t.Fatalf("ListDistricts() error = %v, wantErr %v", err, false)
	}

	if resp.NextPageToken != "" {
		t.Fatalf("ListDistricts() error = %v, wantErr %v", err, false)
	}

	if resp.Districts[0].Id != 1 {
		t.Fatalf("ListDistricts() error = %v, wantErr %v", err, false)
	}

	req = &proto.ListDistrictsRequest{PageSize: 10, Parent: "provinces/1", Skip: 10}
	resp, err = s.ListDistricts(ctx, req)
	if err != nil {
		t.Fatalf("ListDistricts() error = %v, wantErr %v", err, false)
	}

	if len(resp.Districts) != 10 {
		t.Fatalf("ListDistricts() error = %v, wantErr %v", err, false)
	}

	if resp.Districts[0].Id == 1 {
		t.Fatalf("ListDistricts() error = %v, wantErr %v", err, false)
	}
}

func Test_vnProvinceService_GetDistrict(t *testing.T) {
	s := NewVNProvinceService()
	ctx := context.Background()

	req := &proto.GetDistrictRequest{Name: "provinces/1/districts/1"}
	resp, err := s.GetDistrict(ctx, req)
	if err != nil {
		t.Fatalf("GetDistrict() error = %v, wantErr %v", err, false)
	}

	if resp.Id != 1 {
		t.Fatalf("GetDistrict() error = %v, wantErr %v", err, false)
	}
}

func Test_vnProvinceService_ListWards(t *testing.T) {
	s := NewVNProvinceService()
	ctx := context.Background()

	req := &proto.ListWardsRequest{PageSize: -10}
	if _, err := s.ListWards(ctx, req); err == nil {
		t.Fatalf("ListWards() error = %v, wantErr %v", err, true)
	}

	req = &proto.ListWardsRequest{PageSize: 10}
	resp, err := s.ListWards(ctx, req)
	if err == nil {
		t.Fatalf("ListWards() error = %v, wantErr %v", err, true)
	}

	req = &proto.ListWardsRequest{PageSize: 10, Parent: "provinces/1/districts/1"}
	resp, err = s.ListWards(ctx, req)

	if err != nil {
		t.Fatalf("ListWards() error = %v, wantErr %v", err, false)
	}

	if len(resp.Wards) != 10 {
		t.Fatalf("ListWards() error = %v, wantErr %v", err, false)
	}

	if resp.NextPageToken == "" {
		t.Fatalf("ListWards() error = %v, wantErr %v", err, false)
	}

	req = &proto.ListWardsRequest{PageSize: 10, Parent: "provinces/1/districts/1", Filter: "ward_code=1"}
	resp, err = s.ListWards(ctx, req)
	if err != nil {
		t.Fatalf("ListWards() error = %v, wantErr %v", err, false)
	}

	if len(resp.Wards) != 1 {
		t.Fatalf("ListWards() error = %v, wantErr %v", err, false)
	}

	if resp.NextPageToken != "" {
		t.Fatalf("ListWards() error = %v, wantErr %v", err, false)
	}

	if resp.Wards[0].Id != 1 {
		t.Fatalf("ListWards() error = %v, wantErr %v", err, false)
	}

	req = &proto.ListWardsRequest{PageSize: 10, Parent: "provinces/1/districts/1", Skip: 10}
	resp, err = s.ListWards(ctx, req)
	if err != nil {
		t.Fatalf("ListWards() error = %v, wantErr %v", err, false)
	}

	if len(resp.Wards) != 4 {
		t.Fatalf("ListWards() error = %v, wantErr %v", err, false)
	}

	if resp.Wards[0].Id == 1 {
		t.Fatalf("ListWards() error = %v, wantErr %v", err, false)
	}
}

func Test_vnProvinceService_GetWard(t *testing.T) {
	s := NewVNProvinceService()
	ctx := context.Background()

	req := &proto.GetWardRequest{Name: "provinces/1/districts/1/wards/1"}
	resp, err := s.GetWard(ctx, req)
	if err != nil {
		t.Fatalf("GetWard() error = %v, wantErr %v", err, false)
	}

	if resp.Id != 1 {
		t.Fatalf("GetWard() error = %v, wantErr %v", err, false)
	}
}
