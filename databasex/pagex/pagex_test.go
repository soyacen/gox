package pagex

import (
	"context"
	"testing"
	"time"
)

func TestNewPage(t *testing.T) {
	testCases := []struct {
		name     string
		pageNum  uint64
		pageSize uint64
		opts     []Option
		wantErr  bool
	}{
		{
			name:     "valid page",
			pageNum:  1,
			pageSize: 10,
			opts:     nil,
			wantErr:  false,
		},
		{
			name:     "page num zero",
			pageNum:  0,
			pageSize: 10,
			opts:     nil,
			wantErr:  true,
		},
		{
			name:     "page size zero",
			pageNum:  1,
			pageSize: 0,
			opts:     nil,
			wantErr:  true,
		},
		{
			name:     "with options",
			pageNum:  2,
			pageSize: 20,
			opts: []Option{
				Count(false),
				CountColumn("id"),
				OrderBy("created_at DESC"),
				StartTime(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				EndTime(time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)),
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			page, err := NewPage(tc.pageNum, tc.pageSize, tc.opts...)
			if (err != nil) != tc.wantErr {
				t.Errorf("NewPage() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if err != nil {
				return
			}
			if page.PageNum() != tc.pageNum {
				t.Errorf("PageNum() = %v, want %v", page.PageNum(), tc.pageNum)
			}
			if page.PageSize() != tc.pageSize {
				t.Errorf("PageSize() = %v, want %v", page.PageSize(), tc.pageSize)
			}
			expectedOffset := (tc.pageNum - 1) * tc.pageSize
			if page.Offset() != expectedOffset {
				t.Errorf("Offset() = %v, want %v", page.Offset(), expectedOffset)
			}
			if page.Limit() != tc.pageSize {
				t.Errorf("Limit() = %v, want %v", page.Limit(), tc.pageSize)
			}
		})
	}
}

func TestPage_SetTotal(t *testing.T) {
	page, err := NewPage(1, 10)
	if err != nil {
		t.Fatalf("NewPage() error = %v", err)
	}

	if page.Total() != 0 {
		t.Errorf("Total() = %v, want 0 before SetTotal", page.Total())
	}
	if page.Pages() != 0 {
		t.Errorf("Pages() = %v, want 0 before SetTotal", page.Pages())
	}

	page.SetTotal(105)

	if page.Total() != 105 {
		t.Errorf("Total() = %v, want 105", page.Total())
	}
	if page.Pages() != 11 {
		t.Errorf("Pages() = %v, want 11", page.Pages())
	}

	// SetTotal again should not change
	page.SetTotal(200)

	if page.Total() != 105 {
		t.Errorf("Total() = %v, want 105 (should not change)", page.Total())
	}
	if page.Pages() != 11 {
		t.Errorf("Pages() = %v, want 11 (should not change)", page.Pages())
	}
}

func TestPage_Options(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	page, err := NewPage(1, 10,
		Count(false),
		CountColumn("id"),
		OrderBy("created_at DESC"),
		StartTime(start),
		EndTime(end),
	)
	if err != nil {
		t.Fatalf("NewPage() error = %v", err)
	}

	if page.Count() != false {
		t.Errorf("Count() = %v, want false", page.Count())
	}
	if page.CountColumn() != "id" {
		t.Errorf("CountColumn() = %v, want 'id'", page.CountColumn())
	}
	if page.OrderBy() != "created_at DESC" {
		t.Errorf("OrderBy() = %v, want 'created_at DESC'", page.OrderBy())
	}
	if !page.StartTime().Equal(start) {
		t.Errorf("StartTime() = %v, want %v", page.StartTime(), start)
	}
	if !page.EndTime().Equal(end) {
		t.Errorf("EndTime() = %v, want %v", page.EndTime(), end)
	}
}

func TestPage_AsProto_FromProto(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	original, err := NewPage(2, 20,
		Count(true),
		CountColumn("id"),
		OrderBy("created_at DESC"),
		StartTime(start),
		EndTime(end),
	)
	if err != nil {
		t.Fatalf("NewPage() error = %v", err)
	}
	original.SetTotal(105)

	proto := original.AsProto()
	if proto == nil {
		t.Fatal("AsProto() returned nil")
	}

	restored := FromProto(proto)
	if restored == nil {
		t.Fatal("FromProto() returned nil")
	}

	if restored.PageNum() != original.PageNum() {
		t.Errorf("PageNum() = %v, want %v", restored.PageNum(), original.PageNum())
	}
	if restored.PageSize() != original.PageSize() {
		t.Errorf("PageSize() = %v, want %v", restored.PageSize(), original.PageSize())
	}
	if restored.Offset() != original.Offset() {
		t.Errorf("Offset() = %v, want %v", restored.Offset(), original.Offset())
	}
	if restored.Limit() != original.Limit() {
		t.Errorf("Limit() = %v, want %v", restored.Limit(), original.Limit())
	}
	if restored.Total() != original.Total() {
		t.Errorf("Total() = %v, want %v", restored.Total(), original.Total())
	}
	if restored.Pages() != original.Pages() {
		t.Errorf("Pages() = %v, want %v", restored.Pages(), original.Pages())
	}
	if restored.Count() != original.Count() {
		t.Errorf("Count() = %v, want %v", restored.Count(), original.Count())
	}
	if restored.CountColumn() != original.CountColumn() {
		t.Errorf("CountColumn() = %v, want %v", restored.CountColumn(), original.CountColumn())
	}
	if restored.OrderBy() != original.OrderBy() {
		t.Errorf("OrderBy() = %v, want %v", restored.OrderBy(), original.OrderBy())
	}
	if !restored.StartTime().Equal(original.StartTime()) {
		t.Errorf("StartTime() = %v, want %v", restored.StartTime(), original.StartTime())
	}
	if !restored.EndTime().Equal(original.EndTime()) {
		t.Errorf("EndTime() = %v, want %v", restored.EndTime(), original.EndTime())
	}
}

func TestContext(t *testing.T) {
	ctx := context.Background()

	// FromContext on empty context
	page, ok := FromContext(ctx)
	if ok {
		t.Error("FromContext() ok = true, want false")
	}
	if page != nil {
		t.Error("FromContext() page != nil, want nil")
	}

	// NewContext
	ctxWithPage, err := NewContext(ctx, 1, 10)
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}

	// FromContext on context with page
	page, ok = FromContext(ctxWithPage)
	if !ok {
		t.Error("FromContext() ok = false, want true")
	}
	if page == nil {
		t.Error("FromContext() page = nil, want non-nil")
	}
	if page.PageNum() != 1 {
		t.Errorf("PageNum() = %v, want 1", page.PageNum())
	}

	// WithoutPage
	ctxWithoutPage := WithoutPage(ctxWithPage)
	page, ok = FromContext(ctxWithoutPage)
	if ok {
		t.Error("FromContext() ok = true, want false after WithoutPage")
	}
	if page != nil {
		t.Error("FromContext() page != nil, want nil after WithoutPage")
	}
}

func TestNewContext_Error(t *testing.T) {
	ctx := context.Background()

	// Should fail with pageNum 0
	_, err := NewContext(ctx, 0, 10)
	if err == nil {
		t.Error("NewContext() error = nil, want error")
	}
}
