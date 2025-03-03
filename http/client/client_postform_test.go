package client

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/jeffotoni/quick"
)

// testPostFormHandler simulates a server that accepts form-encoded data.
func testPostFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Read form data
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Respond with the received form values
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(r.Form.Encode()))
}

// TestClient_PostForm verifies the POST request with form-encoded data.
func TestClient_PostForm(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testPostFormHandler))
	defer ts.Close()

	client := New()

	formData := url.Values{}
	formData.Set("username", "quick_user")
	formData.Set("password", "supersecret")

	resp, err := client.PostForm(ts.URL, formData)
	if err != nil {
		t.Fatalf("PostForm request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if strings.TrimSpace(string(resp.Body)) != formData.Encode() {
		t.Errorf("Expected body '%s', got '%s'", formData.Encode(), string(resp.Body))
	}
}

// testPostFormRetryHandler simulates a failing server that later succeeds.
func testPostFormRetryHandler(failCount int) http.HandlerFunc {
	attempts := 0
	return func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts <= failCount {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Retry Success"))
	}
}

// TestClient_PostForm_Retry verifies the retry logic in a PostForm request.
func TestClient_PostForm_Retry(t *testing.T) {
	ts := httptest.NewServer(testPostFormRetryHandler(2))
	defer ts.Close()

	client := New(
		WithRetry(
			3,                 // Maximum number of retries
			"500ms",           // Delay between attempts
			true,              // Use exponential backoff
			"500,502,503,504", // HTTP status for retry
			true,              // show Logger
		),
	)

	formData := url.Values{}
	formData.Set("username", "quick_user")
	formData.Set("password", "supersecret")

	resp, err := client.PostForm(ts.URL, formData)
	if err != nil {
		t.Fatalf("PostForm request with retry failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if string(resp.Body) != "Retry Success" {
		t.Errorf("Expected body 'Retry Success', got '%s'", string(resp.Body))
	}
}

// testPostFormFileUploadHandler simulates a server receiving a file upload.
func testPostFormFileUploadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10MB limit
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Retrieve file
	file, _, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read file content
	content, _ := io.ReadAll(file)
	if string(content) != "Fake file content, for us to test upload" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully"))
}

// TestClient_PostFormFileUpload verifies file upload through PostForm.
func TestClient_PostFormFileUpload(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(testPostFormFileUploadHandler))
	defer ts.Close()

	client := New()

	// Create a buffer to store multipart form data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add a text field
	_ = writer.WriteField("description", "Test File Upload")

	// Add a file field
	fileWriter, _ := writer.CreateFormFile("file", "test.txt")
	_, _ = fileWriter.Write([]byte("Fake file content, for us to test upload"))

	// Close the writer to finalize the multipart form
	writer.Close()

	// Ensure the correct Content-Type header is set
	client.Headers["Content-Type"] = writer.FormDataContentType()

	// Execute the request using createRequest
	resp, err := client.createRequest(ts.URL, http.MethodPost, &requestBody)
	if err != nil {
		t.Fatalf("PostForm file upload request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if string(resp.Body) != "File uploaded successfully" {
		t.Errorf("Expected body 'File uploaded successfully', got '%s'", string(resp.Body))
	}
}

// TestFormValue ensures that FormValue() correctly retrieves single values.
func TestFormValue(t *testing.T) {
	q := quick.New()

	q.Post("/formvalue", func(c *quick.Ctx) error {
		name := c.FormValue("name")
		return c.Status(200).SendString(name)
	})

	ts := httptest.NewServer(q)
	defer ts.Close()

	// Send form-urlencoded data
	form := url.Values{}
	form.Set("name", "Jefferson")

	resp, err := http.Post(ts.URL+"/formvalue", "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	if string(body) != "Jefferson" {
		t.Errorf("Expected 'Jefferson', got '%s'", string(body))
	}
}

// TestFormValues ensures that FormValues() correctly retrieves multiple values.
func TestFormValues(t *testing.T) {
	q := quick.New()

	q.Post("/formvalues", func(c *quick.Ctx) error {
		values := c.FormValues()
		jsonData, _ := json.Marshal(values)
		return c.Status(200).Send(jsonData)
	})

	ts := httptest.NewServer(q)
	defer ts.Close()

	form := url.Values{}
	form.Set("name", "Jefferson")
	form.Set("email", "jeff@example.com")

	resp, err := http.Post(ts.URL+"/formvalues", "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	var result map[string][]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(result["name"]) == 0 || result["name"][0] != "Jefferson" {
		t.Errorf("Expected name 'Jefferson', got '%v'", result["name"])
	}
	if len(result["email"]) == 0 || result["email"][0] != "jeff@example.com" {
		t.Errorf("Expected email 'jeff@example.com', got '%v'", result["email"])
	}
}

// TestFormValuesJSON ensures that FormValues() works correctly with JSON requests.
func TestFormValuesJSON(t *testing.T) {
	q := quick.New()

	q.Post("/formvaluesjson", func(c *quick.Ctx) error {
		values := c.FormValues()
		jsonData, _ := json.Marshal(values)
		return c.Status(200).Send(jsonData)
	})

	ts := httptest.NewServer(q)
	defer ts.Close()

	form := url.Values{}
	form.Set("name", "Jefferson")
	form.Set("email", "jeff@example.com")

	resp, err := http.Post(ts.URL+"/formvaluesjson", "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	var result map[string][]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(result["name"]) == 0 || result["name"][0] != "Jefferson" {
		t.Errorf("Expected name 'Jefferson', got '%v'", result["name"])
	}
	if len(result["email"]) == 0 || result["email"][0] != "jeff@example.com" {
		t.Errorf("Expected email 'jeff@example.com', got '%v'", result["email"])
	}
}

// TestFormValue_Empty ensures that FormValue() returns an empty string when the key is missing.
func TestFormValue_Empty(t *testing.T) {
	q := quick.New()

	q.Post("/formvalue_empty", func(c *quick.Ctx) error {
		value := c.FormValue("missing_key")
		return c.Status(200).SendString(value)
	})

	ts := httptest.NewServer(q)
	defer ts.Close()

	form := url.Values{}

	resp, err := http.Post(ts.URL+"/formvalue_empty", "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	if len(string(body)) != 0 {
		t.Errorf("Expected empty string, got '%s'", string(body))
	}
}

// TestFormValues_Empty ensures that FormValues() returns an empty map when no data is sent.
func TestFormValues_Empty(t *testing.T) {
	q := quick.New()

	q.Post("/formvalues_empty", func(c *quick.Ctx) error {
		values := c.FormValues()
		jsonData, _ := json.Marshal(values)
		return c.Status(200).Send(jsonData)
	})

	ts := httptest.NewServer(q)
	defer ts.Close()

	form := url.Values{}

	resp, err := http.Post(ts.URL+"/formvalues_empty", "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]any
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected empty map, got %v", result)
	}
}
