package oreillyapi

import "time"

type SearchResponse struct {
	Previous *string  `json:"previous"`
	Next     string   `json:"next"`
	Page     int      `json:"page"`
	Results  []Result `json:"results"`
	Total    int      `json:"total"`
	// Facets   Facets   `json:"facets"`
	Meta Meta `json:"meta"`
}

type Result struct {
	Id               string    `json:"id"`
	ContentType      string    `json:"content_type"`
	ArchiveId        string    `json:"archive_id"`
	Ourn             string    `json:"ourn"`
	Title            string    `json:"title"`
	Timestamp        time.Time `json:"timestamp"`
	LastModifiedTime time.Time `json:"last_modified_time"`
	DateAdded        time.Time `json:"date_added"`
	Issued           time.Time `json:"issued"`
	Isbn             string    `json:"isbn"`
	Format           string    `json:"format"` // NOTE: convert to enum ?
	ContentFormat    string    `json:"content_format"`
	VirtualPages     int       `json:"virtual_pages"`
	MinutesRequired  float64   `json:"minutes_required"`
	DurationSeconds  int       `json:"duration_seconds"`
	Language         string    `json:"language"`
	// NaturalKey          []string  `json:"natural_key"`
	Source              string   `json:"source"`
	VideoClassification []string `json:"video_classification,omitempty"`
	HasAssessment       bool     `json:"has_assessment"`
	AcademicExcluded    bool     `json:"academic_excluded"`
	Url                 string   `json:"url"`
	WebUrl              string   `json:"web_url"`
	CoverUrl            string   `json:"cover_url"`
	Description         string   `json:"description"`
	Populartity         int      `json:"popularity"`
	AverageRating       int      `json:"average_rating"`
	ReportScore         int      `json:"report_score"`
	NumberOfReviews     int      `json:"number_of_reviews"`
	NumberOfFollowers   int      `json:"number_of_followers"`
	NumberOfItems       int      `json:"number_of_items"`
	Authors             []string `json:"authors"`
	Publishers          []string `json:"publishers"`
	Topics              []string `json:"topics"`
	TopicsPayload       []Topic  `json:"topics_payload"`
}

type Topic struct {
	Uuid  string   `json:"uuid"`
	Slug  string   `json:"slug"`
	Name  string   `json:"name"`
	Score *float64 `json:"score"`
}

type Facets struct {
	FacetQueries   FacetQueries   `json:"facet_queries"`
	FacetFields    FacetFields    `json:"facet_fields"`
	FacetRanges    FacetRanges    `json:"facet_ranges"`
	FacetIntervals FacetIntervals `json:"facet_intervals"`
}

type FacetQueries struct {
	// TODO:
}
type FacetFields struct {
	// TODO:
}
type FacetRanges struct {
	// TODO:
}
type FacetIntervals struct {
	// TODO:
}

type Meta struct {
	QueryIdentifier string `json:"query_identifier"`
}
