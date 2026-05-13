package farmer

import "testing"

func TestCreateFarmerPayloadValidateTrimsFields(t *testing.T) {
	p := &CreateFarmerPayload{
		Name:       "  John  ",
		NationalID: "  ABC123  ",
		Community:  "  North  ",
	}

	if err := p.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if p.Name != "John" || p.NationalID != "ABC123" || p.Community != "North" {
		t.Fatalf("expected trimmed fields, got %+v", p)
	}
}

func TestGetFarmersQueryValidateDefaults(t *testing.T) {
	q := &GetFarmersQuery{}
	if err := q.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if q.Page == nil || *q.Page != 1 {
		t.Fatalf("expected default page 1, got %#v", q.Page)
	}
	if q.Limit == nil || *q.Limit != 20 {
		t.Fatalf("expected default limit 20, got %#v", q.Limit)
	}
	if q.Sort == nil || *q.Sort != "created_at" {
		t.Fatalf("expected default sort created_at, got %#v", q.Sort)
	}
	if q.Order == nil || *q.Order != "desc" {
		t.Fatalf("expected default order desc, got %#v", q.Order)
	}
}
