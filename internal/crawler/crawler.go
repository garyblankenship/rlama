package crawler

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/dontizi/rlama/internal/domain"
	"golang.org/x/net/html/charset"
)

// WebCrawler manage web crawling operations
type WebCrawler struct {
	client       *http.Client
	visitedURLs  map[string]bool
	maxDepth     int
	concurrency  int
	baseURL      *url.URL
	excludePaths []string
	mutex        sync.Mutex
}

// NewWebCrawler creates a new web crawler instance
func NewWebCrawler(baseURLStr string, maxDepth, concurrency int, excludePaths []string) (*WebCrawler, error) {
	baseURL, err := url.Parse(baseURLStr)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	return &WebCrawler{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		visitedURLs:  make(map[string]bool),
		maxDepth:     maxDepth,
		concurrency:  concurrency,
		baseURL:      baseURL,
		excludePaths: excludePaths,
		mutex:        sync.Mutex{},
	}, nil
}

// isWebContent checks if a URL points to text/HTML content rather than binary files
func isWebContent(urlStr string) bool {
	// Extensions to ignore (binary files, etc.)
	excludeExtensions := []string{
		".zip", ".rar", ".tar", ".gz", ".pdf", ".doc", ".docx", 
		".xls", ".xlsx", ".ppt", ".pptx", ".exe", ".bin", 
		".dmg", ".iso", ".img", ".apk", ".ipa", ".mp3", 
		".mp4", ".avi", ".mov", ".flv", ".mkv",
	}
	
	lowerURL := strings.ToLower(urlStr)
	for _, ext := range excludeExtensions {
		if strings.HasSuffix(lowerURL, ext) {
			return false
		}
	}
	
	return true
}

// CrawlWebsite starts crawling from the base URL and returns processed documents
func (wc *WebCrawler) CrawlWebsite() ([]*domain.Document, error) {
	var documents []*domain.Document
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, wc.concurrency)
	resultChan := make(chan *domain.Document, 100)
	errorChan := make(chan error, 100)
	doneChan := make(chan struct{})

	// Start crawling with the base URL
	wg.Add(1)
	go wc.crawlURL(wc.baseURL.String(), 0, &wg, semaphore, resultChan, errorChan)

	// Collect results
	go func() {
		wg.Wait()
		close(doneChan)
	}()

	// Process results and errors
	for {
		select {
		case doc := <-resultChan:
			documents = append(documents, doc)
		case err := <-errorChan:
			fmt.Printf("Error during crawling: %v\n", err)
		case <-doneChan:
			return documents, nil
		}
	}
}

// crawlURL processes a single URL and extracts links to crawl further
func (wc *WebCrawler) crawlURL(urlStr string, depth int, wg *sync.WaitGroup, semaphore chan struct{}, resultChan chan<- *domain.Document, errorChan chan<- error) {
	defer wg.Done()
	
	// Respect concurrency limit
	semaphore <- struct{}{}
	defer func() { <-semaphore }()

	// First check if this is web content we want to process
	if !isWebContent(urlStr) {
		return // Skip non-web content
	}

	// Check if we've reached max depth
	if depth > wc.maxDepth {
		return
	}

	// Normalize URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		errorChan <- fmt.Errorf("error parsing URL %s: %w", urlStr, err)
		return
	}

	// Make URL absolute if it's relative
	if !parsedURL.IsAbs() {
		parsedURL = wc.baseURL.ResolveReference(parsedURL)
	}

	// Skip URLs from different domains
	if parsedURL.Host != wc.baseURL.Host {
		return
	}

	// Skip excluded paths
	for _, excludePath := range wc.excludePaths {
		if strings.HasPrefix(parsedURL.Path, excludePath) {
			return
		}
	}

	// Skip URL fragments
	parsedURL.Fragment = ""
	urlStr = parsedURL.String()

	// Skip if we've already visited this URL
	wc.mutex.Lock()
	if wc.visitedURLs[urlStr] {
		wc.mutex.Unlock()
		return
	}
	wc.visitedURLs[urlStr] = true
	wc.mutex.Unlock()

	// Fetch the page
	fmt.Printf("Crawling: %s\n", urlStr)
	resp, err := wc.client.Get(urlStr)
	if err != nil {
		errorChan <- fmt.Errorf("error fetching %s: %w", urlStr, err)
		return
	}
	defer resp.Body.Close()

	// Vérifier le Content-Type de la réponse
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") && 
	   !strings.Contains(contentType, "text/plain") && 
	   !strings.Contains(contentType, "application/xhtml+xml") {
		// Ignorer les types de contenu qui ne sont pas du texte ou HTML
		return
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		errorChan <- fmt.Errorf("received non-OK status code for %s: %d", urlStr, resp.StatusCode)
		return
	}

	// Handle character encoding
	utf8Reader, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		errorChan <- fmt.Errorf("error converting encoding for %s: %w", urlStr, err)
		return
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(utf8Reader)
	if err != nil {
		errorChan <- fmt.Errorf("error parsing HTML for %s: %w", urlStr, err)
		return
	}

	// Get page title
	title := doc.Find("title").Text()
	if title == "" {
		title = parsedURL.Path
	}

	// Extract page content and convert to Markdown
	markdown, err := extractContentAsMarkdown(doc)
	if err != nil {
		errorChan <- fmt.Errorf("error extracting content from %s: %w", urlStr, err)
		return
	}

	// Create a document
	document := domain.NewDocument(urlStr, markdown)
	document.Name = title
	document.ContentType = "text/markdown"
	
	resultChan <- document

	// Find links and queue them for crawling
	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists || href == "" || strings.HasPrefix(href, "#") {
			return
		}

		linkURL, err := url.Parse(href)
		if err != nil {
			return
		}

		// Handle relative URLs
		if !linkURL.IsAbs() {
			linkURL = parsedURL.ResolveReference(linkURL)
		}

		// Only follow links on the same domain and that are web content
		if linkURL.Host != wc.baseURL.Host || !isWebContent(linkURL.String()) {
			return
		}

		// Remove fragments
		linkURL.Fragment = ""
		nextURL := linkURL.String()

		// Avoid crawling the same URL twice
		wc.mutex.Lock()
		alreadyVisited := wc.visitedURLs[nextURL]
		wc.mutex.Unlock()

		if !alreadyVisited {
			wg.Add(1)
			go wc.crawlURL(nextURL, depth+1, wg, semaphore, resultChan, errorChan)
		}
	})
}

// extractContentAsMarkdown extracts main content from an HTML document and converts it to Markdown
func extractContentAsMarkdown(doc *goquery.Document) (string, error) {
	// Remove unwanted elements
	doc.Find("script, style, nav, footer, header, aside, .ads, .comments, .navigation").Remove()

	// Extract the main content (prefer main content areas)
	var contentNode *goquery.Selection
	
	// Try to find the main content area using common selectors
	contentSelectors := []string{"main", "article", ".content", "#content", ".post-content", ".article-content"}
	for _, selector := range contentSelectors {
		selection := doc.Find(selector)
		if selection.Length() > 0 {
			contentNode = selection
			break
		}
	}

	// If no specific content area found, use the body
	if contentNode == nil || contentNode.Length() == 0 {
		contentNode = doc.Find("body")
	}

	// Get the HTML content
	html, err := contentNode.Html()
	if err != nil {
		return "", err
	}

	// Convert HTML to Markdown (simplified approach)
	// In a real implementation, you would use a proper HTML to Markdown converter
	markdown := convertHTMLToMarkdown(html)
	
	return markdown, nil
}

// convertHTMLToMarkdown is a simplified HTML to Markdown converter
// In a real implementation, you would use a proper library
func convertHTMLToMarkdown(html string) string {
	// Replace common HTML tags with Markdown equivalents
	replacements := map[string]string{
		"<h1>":     "# ",
		"</h1>":    "\n\n",
		"<h2>":     "## ",
		"</h2>":    "\n\n",
		"<h3>":     "### ",
		"</h3>":    "\n\n",
		"<h4>":     "#### ",
		"</h4>":    "\n\n",
		"<h5>":     "##### ",
		"</h5>":    "\n\n",
		"<h6>":     "###### ",
		"</h6>":    "\n\n",
		"<p>":      "",
		"</p>":     "\n\n",
		"<strong>": "**",
		"</strong>":"**",
		"<b>":      "**",
		"</b>":     "**",
		"<em>":     "*",
		"</em>":    "*",
		"<i>":      "*",
		"</i>":     "*",
		"<br>":     "\n",
		"<br/>":    "\n",
		"<br />":   "\n",
		"<ul>":     "\n",
		"</ul>":    "\n",
		"<ol>":     "\n",
		"</ol>":    "\n",
		"<li>":     "- ",
		"</li>":    "\n",
		"<code>":   "`",
		"</code>":  "`",
		"<pre>":    "```\n",
		"</pre>":   "\n```\n",
	}

	result := html
	for tag, replacement := range replacements {
		result = strings.ReplaceAll(result, tag, replacement)
	}

	// Clean up multiple newlines
	for strings.Contains(result, "\n\n\n") {
		result = strings.ReplaceAll(result, "\n\n\n", "\n\n")
	}

	return result
}