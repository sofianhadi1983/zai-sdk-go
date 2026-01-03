// Package websearch provides types for the Web Search API.
package websearch

// SensitiveWordCheck represents sensitive word check configuration.
type SensitiveWordCheck struct {
	// Type is the sensitive word type, currently only supports "ALL"
	Type string `json:"type,omitempty"`

	// Status enables/disables sensitive word checking
	// "ENABLE" - Enable checking (default)
	// "DISABLE" - Disable checking (requires special permissions)
	Status string `json:"status,omitempty"`
}

// Sensitive word check constants
const (
	SensitiveWordTypeAll      = "ALL"
	SensitiveWordStatusEnable = "ENABLE"
	SensitiveWordStatusDisable = "DISABLE"
)

// SearchIntentResp represents search intent analysis response.
type SearchIntentResp struct {
	// Query is the search optimized query
	Query string `json:"query"`

	// Intent is the determined intent type
	Intent string `json:"intent"`

	// Keywords are the search keywords
	Keywords string `json:"keywords"`
}

// SearchResultResp represents an individual search result.
type SearchResultResp struct {
	// Title is the result title
	Title string `json:"title"`

	// Link is the result URL
	Link string `json:"link"`

	// Content is the result content/snippet
	Content string `json:"content"`

	// Icon is the website icon URL
	Icon string `json:"icon"`

	// Media is the source media
	Media string `json:"media"`

	// Refer is the reference number (e.g., "[ref_1]")
	Refer string `json:"refer"`

	// PublishDate is the publication date
	PublishDate string `json:"publish_date"`

	// Images are URLs of images associated with the result
	Images []string `json:"images,omitempty"`
}

// WebSearchResponse represents the web search response.
type WebSearchResponse struct {
	// Created is the creation timestamp (Unix time)
	Created int64 `json:"created,omitempty"`

	// RequestID is the request identifier
	RequestID string `json:"request_id,omitempty"`

	// ID is the response identifier
	ID string `json:"id,omitempty"`

	// SearchIntent contains search intent analysis
	SearchIntent *SearchIntentResp `json:"search_intent,omitempty"`

	// SearchResult contains search results
	SearchResult []SearchResultResp `json:"search_result,omitempty"`
}

// GetResults returns the search results.
func (r *WebSearchResponse) GetResults() []SearchResultResp {
	return r.SearchResult
}

// HasIntent returns true if search intent analysis is available.
func (r *WebSearchResponse) HasIntent() bool {
	return r.SearchIntent != nil
}

// WebSearchRequest represents a web search request.
type WebSearchRequest struct {
	// SearchQuery is the search query text (required)
	SearchQuery string `json:"search_query"`

	// SearchEngine specifies the search engine to use
	SearchEngine string `json:"search_engine,omitempty"`

	// RequestID is a user-provided unique identifier for the request
	// If not provided, the platform will generate one
	RequestID string `json:"request_id,omitempty"`

	// UserID is the user identifier
	UserID string `json:"user_id,omitempty"`

	// SensitiveWordCheck configures sensitive word filtering
	SensitiveWordCheck *SensitiveWordCheck `json:"sensitive_word_check,omitempty"`

	// Count is the number of search results to return
	Count int `json:"count,omitempty"`

	// SearchDomainFilter filters results by domain
	SearchDomainFilter string `json:"search_domain_filter,omitempty"`

	// SearchRecencyFilter filters results by recency
	// Options: "day", "week", "month", "year"
	SearchRecencyFilter string `json:"search_recency_filter,omitempty"`

	// ContentSize specifies the desired content size
	// Options: "small", "medium", "large"
	ContentSize string `json:"content_size,omitempty"`

	// SearchIntent enables search intent analysis
	SearchIntent bool `json:"search_intent,omitempty"`

	// IncludeImage enables image inclusion in results
	IncludeImage bool `json:"include_image,omitempty"`
}

// Recency filter constants (per Z.ai API specification)
const (
	RecencyFilterOneDay   = "oneDay"
	RecencyFilterOneWeek  = "oneWeek"
	RecencyFilterOneMonth = "oneMonth"
	RecencyFilterOneYear  = "oneYear"
	RecencyFilterNoLimit  = "noLimit"
)

// Content size constants
const (
	ContentSizeSmall  = "small"
	ContentSizeMedium = "medium"
	ContentSizeLarge  = "large"
)

// Search engine constants
const (
	// SearchEnginePrime is the Z.ai Premium Version Search Engine (required value).
	SearchEnginePrime = "search-prime"
)

// NewWebSearchRequest creates a new web search request.
// The search engine is automatically set to "search-prime" (the only supported value).
//
// Example:
//
//	req := websearch.NewWebSearchRequest("artificial intelligence trends 2024")
func NewWebSearchRequest(query string) *WebSearchRequest {
	return &WebSearchRequest{
		SearchQuery:  query,
		SearchEngine: SearchEnginePrime, // Default to required value
	}
}

// SetSearchEngine sets the search engine.
func (r *WebSearchRequest) SetSearchEngine(engine string) *WebSearchRequest {
	r.SearchEngine = engine
	return r
}

// SetRequestID sets the request ID.
func (r *WebSearchRequest) SetRequestID(id string) *WebSearchRequest {
	r.RequestID = id
	return r
}

// SetUserID sets the user ID.
func (r *WebSearchRequest) SetUserID(id string) *WebSearchRequest {
	r.UserID = id
	return r
}

// SetSensitiveWordCheck sets the sensitive word check configuration.
func (r *WebSearchRequest) SetSensitiveWordCheck(check *SensitiveWordCheck) *WebSearchRequest {
	r.SensitiveWordCheck = check
	return r
}

// SetCount sets the number of results to return.
func (r *WebSearchRequest) SetCount(count int) *WebSearchRequest {
	r.Count = count
	return r
}

// SetDomainFilter sets the domain filter.
func (r *WebSearchRequest) SetDomainFilter(domain string) *WebSearchRequest {
	r.SearchDomainFilter = domain
	return r
}

// SetRecencyFilter sets the recency filter.
func (r *WebSearchRequest) SetRecencyFilter(recency string) *WebSearchRequest {
	r.SearchRecencyFilter = recency
	return r
}

// SetContentSize sets the desired content size.
func (r *WebSearchRequest) SetContentSize(size string) *WebSearchRequest {
	r.ContentSize = size
	return r
}

// SetSearchIntent enables search intent analysis.
func (r *WebSearchRequest) SetSearchIntent(enable bool) *WebSearchRequest {
	r.SearchIntent = enable
	return r
}

// SetIncludeImage enables image inclusion.
func (r *WebSearchRequest) SetIncludeImage(include bool) *WebSearchRequest {
	r.IncludeImage = include
	return r
}
