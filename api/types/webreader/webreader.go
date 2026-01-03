// Package webreader provides types for the Web Reader API.
package webreader

// Request represents a request to read a web page.
type Request struct {
	// URL is the target page URL to read (required).
	URL string `json:"url"`

	// RequestID is a unique request task ID (6-64 chars, optional).
	RequestID string `json:"request_id,omitempty"`

	// UserID is a unique end-user ID (6-128 chars, optional).
	UserID string `json:"user_id,omitempty"`

	// Timeout is the request timeout in seconds (optional).
	Timeout string `json:"timeout,omitempty"`

	// NoCache disables cache when true (optional).
	NoCache bool `json:"no_cache,omitempty"`

	// ReturnFormat specifies the return format, e.g. "markdown" or "text" (optional).
	ReturnFormat string `json:"return_format,omitempty"`

	// RetainImages keeps images in output when true (optional).
	RetainImages bool `json:"retain_images,omitempty"`

	// NoGFM disables GitHub Flavored Markdown when true (optional).
	NoGFM bool `json:"no_gfm,omitempty"`

	// KeepImgDataURL keeps image data URLs when true (optional).
	KeepImgDataURL bool `json:"keep_img_data_url,omitempty"`

	// WithImagesSummary includes images summary when true (optional).
	WithImagesSummary bool `json:"with_images_summary,omitempty"`

	// WithLinksSummary includes links summary when true (optional).
	WithLinksSummary bool `json:"with_links_summary,omitempty"`
}

// NewRequest creates a new web reader request.
func NewRequest(url string) *Request {
	return &Request{
		URL: url,
	}
}

// SetRequestID sets the request ID.
func (r *Request) SetRequestID(requestID string) *Request {
	r.RequestID = requestID
	return r
}

// SetUserID sets the user ID.
func (r *Request) SetUserID(userID string) *Request {
	r.UserID = userID
	return r
}

// SetTimeout sets the request timeout in seconds.
func (r *Request) SetTimeout(timeout string) *Request {
	r.Timeout = timeout
	return r
}

// SetNoCache disables cache.
func (r *Request) SetNoCache(noCache bool) *Request {
	r.NoCache = noCache
	return r
}

// SetReturnFormat sets the return format (e.g., "markdown" or "text").
func (r *Request) SetReturnFormat(format string) *Request {
	r.ReturnFormat = format
	return r
}

// SetRetainImages sets whether to keep images in output.
func (r *Request) SetRetainImages(retain bool) *Request {
	r.RetainImages = retain
	return r
}

// SetNoGFM disables GitHub Flavored Markdown.
func (r *Request) SetNoGFM(noGFM bool) *Request {
	r.NoGFM = noGFM
	return r
}

// SetKeepImgDataURL sets whether to keep image data URLs.
func (r *Request) SetKeepImgDataURL(keep bool) *Request {
	r.KeepImgDataURL = keep
	return r
}

// SetWithImagesSummary sets whether to include images summary.
func (r *Request) SetWithImagesSummary(withSummary bool) *Request {
	r.WithImagesSummary = withSummary
	return r
}

// SetWithLinksSummary sets whether to include links summary.
func (r *Request) SetWithLinksSummary(withSummary bool) *Request {
	r.WithLinksSummary = withSummary
	return r
}

// ReaderData contains the extracted web page data.
type ReaderData struct {
	// Images is a map of image URLs.
	Images map[string]string `json:"images,omitempty"`

	// Links is a map of links.
	Links map[string]string `json:"links,omitempty"`

	// Title is the page title.
	Title string `json:"title,omitempty"`

	// Description is the page description.
	Description string `json:"description,omitempty"`

	// URL is the page URL.
	URL string `json:"url,omitempty"`

	// Content is the extracted content.
	Content string `json:"content,omitempty"`

	// PublishedTime is the publication time.
	PublishedTime string `json:"publishedTime,omitempty"`

	// Metadata contains additional metadata.
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// External contains external data.
	External map[string]interface{} `json:"external,omitempty"`
}

// GetContent returns the extracted content.
func (d *ReaderData) GetContent() string {
	return d.Content
}

// GetTitle returns the page title.
func (d *ReaderData) GetTitle() string {
	return d.Title
}

// GetDescription returns the page description.
func (d *ReaderData) GetDescription() string {
	return d.Description
}

// GetImages returns the images map.
func (d *ReaderData) GetImages() map[string]string {
	if d.Images == nil {
		return make(map[string]string)
	}
	return d.Images
}

// GetLinks returns the links map.
func (d *ReaderData) GetLinks() map[string]string {
	if d.Links == nil {
		return make(map[string]string)
	}
	return d.Links
}

// HasContent returns true if the reader data has content.
func (d *ReaderData) HasContent() bool {
	return d.Content != ""
}

// Response represents the response from a web reader operation.
type Response struct {
	// ReaderResult contains the extracted web page data.
	ReaderResult *ReaderData `json:"reader_result,omitempty"`
}

// GetResult returns the reader result data.
func (r *Response) GetResult() *ReaderData {
	return r.ReaderResult
}

// HasResult returns true if the response has reader result data.
func (r *Response) HasResult() bool {
	return r.ReaderResult != nil
}

// GetContent returns the extracted content from the result.
func (r *Response) GetContent() string {
	if r.ReaderResult != nil {
		return r.ReaderResult.GetContent()
	}
	return ""
}

// GetTitle returns the page title from the result.
func (r *Response) GetTitle() string {
	if r.ReaderResult != nil {
		return r.ReaderResult.GetTitle()
	}
	return ""
}
