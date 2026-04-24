package pagex

import (
	"errors"
	"sync"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// Page represents pagination parameters and state.
type Page struct {
	// pageNum is the page number, starting from 1.
	pageNum uint64
	// pageSize is the number of items per page.
	pageSize uint64
	// offset is the number of rows to skip.
	offset uint64
	// limit is the maximum number of rows to return.
	limit uint64
	// total is the total number of rows.
	total uint64
	// pages is the total number of pages.
	pages uint64
	// count indicates whether to include a count query.
	count bool
	// countColumn is the column name for the count query.
	countColumn string
	// orderBy is the ordering clause.
	orderBy string
	// startTime is the start time filter.
	startTime time.Time
	// endTime is the end time filter.
	endTime time.Time
	// totalOnce ensures total can only be set once.
	totalOnce sync.Once
}

func (p *Page) init() *Page {
	return p
}

func (p *Page) apply(opts ...Option) *Page {
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// PageNum returns the page number.
//
// Returns:
//   - uint64: the page number, starting from 1.
func (p *Page) PageNum() uint64 {
	return p.pageNum
}

// PageSize returns the page size.
//
// Returns:
//   - uint64: the number of items per page.
func (p *Page) PageSize() uint64 {
	return p.pageSize
}

// Offset returns the offset.
//
// Returns:
//   - uint64: the number of rows to skip.
func (p *Page) Offset() uint64 {
	return p.offset
}

// Limit returns the limit.
//
// Returns:
//   - uint64: the maximum number of rows to return.
func (p *Page) Limit() uint64 {
	return p.limit
}

// Total returns the total number of rows.
//
// Returns:
//   - uint64: the total number of rows.
func (p *Page) Total() uint64 {
	return p.total
}

// SetTotal sets the total number of rows and calculates the total pages.
//
// Parameters:
//   - total: the total number of rows.
func (p *Page) SetTotal(total uint64) {
	p.totalOnce.Do(func ()  {
		p.total = total
		p.pages = (total + p.pageSize - 1) / p.pageSize
	})
}

// Pages returns the total number of pages.
//
// Returns:
//   - uint64: the total number of pages.
func (p *Page) Pages() uint64 {
	return p.pages
}

// StartTime returns the start time.
//
// Returns:
//   - time.Time: the start time filter.
func (p *Page) StartTime() time.Time {
	return p.startTime
}

// EndTime returns the end time.
//
// Returns:
//   - time.Time: the end time filter.
func (p *Page) EndTime() time.Time {
	return p.endTime
}

// Count returns whether to include a count query.
//
// Returns:
//   - bool: true if count query is included.
func (p *Page) Count() bool {
	return p.count
}

// CountColumn returns the column name for the count query.
//
// Returns:
//   - string: the column name for the count query.
func (p *Page) CountColumn() string {
	return p.countColumn
}

// OrderBy returns the ordering clause.
//
// Returns:
//   - string: the ordering clause.
func (p *Page) OrderBy() string {
	return p.orderBy
}

// AsProto converts the Page to its protobuf representation.
//
// Returns:
//   - *PageProto: the protobuf representation of the page.
func (p *Page) AsProto() *PageProto {
	return &PageProto{
		PageNum:     p.pageNum,
		PageSize:    p.pageSize,
		Offset:      p.offset,
		Limit:       p.limit,
		Total:       p.total,
		Pages:       p.pages,
		Count:       p.count,
		CountColumn: p.countColumn,
		OrderBy:     p.orderBy,
		StartTime:   timestamppb.New(p.startTime),
		EndTime:     timestamppb.New(p.endTime),
	}
}

// Option is a functional option for configuring a Page.
type Option func(p *Page)

// Count returns an Option that sets whether to include a count query.
//
// Parameters:
//   - count: whether to include a count query.
//
// Returns:
//   - Option: the configured option.
func Count(count bool) Option {
	return func(p *Page) {
		p.count = count
	}
}

// CountColumn returns an Option that sets the count column name.
//
// Parameters:
//   - countColumn: the column name for the count query.
//
// Returns:
//   - Option: the configured option.
func CountColumn(countColumn string) Option {
	return func(p *Page) {
		p.countColumn = countColumn
	}
}

// OrderBy returns an Option that sets the ordering clause.
//
// Parameters:
//   - orderBy: the ordering clause.
//
// Returns:
//   - Option: the configured option.
func OrderBy(orderBy string) Option {
	return func(p *Page) {
		p.orderBy = orderBy
	}
}

// StartTime returns an Option that sets the start time filter.
//
// Parameters:
//   - t: the start time.
//
// Returns:
//   - Option: the configured option.
func StartTime(t time.Time) Option {
	return func(p *Page) {
		p.startTime = t
	}
}

// EndTime returns an Option that sets the end time filter.
//
// Parameters:
//   - t: the end time.
//
// Returns:
//   - Option: the configured option.
func EndTime(t time.Time) Option {
	return func(p *Page) {
		p.endTime = t
	}
}

// NewPage creates a new Page with the given page number and size.
//
// Parameters:
//   - pageNum: page number, starting from 1.
//   - pageSize: number of items per page.
//   - opts: optional pagination options.
//
// Returns:
//   - *Page: the created page.
//   - error: error if pageNum or pageSize is zero.
func NewPage(pageNum uint64, pageSize uint64, opts ...Option) (*Page, error) {
	if pageNum == 0 {
		return nil, errors.New("pagex: pageNum is zero")
	}
	if pageSize == 0 {
		return nil, errors.New("pagex: pageSize is zero")
	}
	// calculate offset and limit
	offset := (pageNum - 1) * pageSize
	limit := pageSize
	page := &Page{
		pageNum:  pageNum,
		pageSize: pageSize,
		offset:   offset,
		limit:    limit,
		count:    true,
	}
	p := page.apply(opts...).init()
	return p, nil
}

// FromProto creates a Page from its protobuf representation.
//
// Parameters:
//   - p: the protobuf representation.
//
// Returns:
//   - *Page: the created page.
func FromProto(p *PageProto) *Page {
	return &Page{
		pageNum:     p.PageNum,
		pageSize:    p.PageSize,
		offset:      p.Offset,
		limit:       p.Limit,
		total:       p.Total,
		pages:       p.Pages,
		count:       p.Count,
		countColumn: p.CountColumn,
		orderBy:     p.OrderBy,
		startTime:   p.StartTime.AsTime(),
		endTime:     p.EndTime.AsTime(),
	}
}
