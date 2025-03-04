## 📂 Uploading Files in Quick ![Quick Logo](/quick.png)

**Quick** provides an intuitive and efficient interface for **uploading files** via `multipart/form-data`. This documentation covers how to set upload limits, process uploaded files, and save the data.

---

#### 📌 What is `multipart/form-data`?

**`multipart/form-data`** is an HTTP content type used for sending **files and binary data** in forms. Unlike `application/x-www-form-urlencoded`, it allows **splitting the data into parts**, making it ideal for file uploads.

##### Example HTML
```html
<form action="/upload" method="post" enctype="multipart/form-data">
    <input type="file" name="files" multiple>
    <button type="submit">Enviar</button>
</form>
```

#### 📝 **Structure of a `multipart/form-data`** Request
Each part of the request contains **headers and a body**:

```text
Content-Type: multipart/form-data; boundary=––WebKitFormBoundary

——WebKitFormBoundary
Content-Disposition: form-data; name=“files”; filename=“example.png”
Content-Type: image/png

(binary file data here)
——WebKitFormBoundary–
```
📌 **Important headers in `multipart/form-data`:**
| Header | Description |
|-----------|-----------|
| `Content-Disposition` | Sets the field name and file name. |
| `Content-Type` | Sets the MIME type of the uploaded file. |
| `Content-Length` | Indicates the total size of the request. |

---

#### 📌 How does it work in Quick?

Quick provides a simplified API for managing uploads, allowing you to easily retrieve and manipulate files.

✅ **Main Methods and Functionalities**:
| Method | Description |
|--------|-----------|
| `c.FormFile("file")` | Returns a single file uploaded in the form. |
| `c.FormFiles("files")` | Returns a list of uploaded files (multiple uploads). |
| `c.FormFileLimit("10MB")` | Sets an upload limit (default is `1MB`). |
| `uploadedFile.FileName()` | Returns the file name. |
| `uploadedFile.Size()` | Returns the file size in bytes. |
| `uploadedFile.ContentType()` | Returns the MIME type of the file. |
| `uploadedFile.Bytes()` | Returns the bytes of the file. |
| `uploadedFile.Save("/path")` | Saves the file to a specified directory. |
| `uploadedFile.Save("/path", "your-name-file")` | Saves the file with your name. |
| `uploadedFile.SaveAll("/path")` | Saves the file to a specified directory. |
---

#### 📌 File Upload Example

```go
q.Post("/upload", func(c *quick.Ctx) error {
    uploadedFile, err := c.FormFile("file")
    if err != nil {
        return c.Status(400).JSON(Msg{
            Msg: "Upload error",
            Error: err.Error(),
         })
    }

fmt.Println("Name:", uploadedFile.FileName())
fmt.Println("Size:", uploadedFile.Size())
fmt.Println("MIME Type:", uploadedFile.ContentType())
fmt.Println("Bytes:", file.Bytes())

// Save the file (optional)
// uploadedFile.Save("/tmp/uploads")
// uploadedFile.Save("/tmp/uploads", "your-name-file")

return c.Status(200).JSONIN(uploadedFile)
})
```
#### 📌 Multiple Upload Example

```go
q.Post("/upload-multiple", func(c *quick.Ctx) error {
    // set limit upload
    c.FormFileLimit("10MB")

    // recebereceiving files
    files, err := c.FormFiles("files")
    if err != nil {
        return c.Status(400).JSON(Msg{
            Msg:   "Upload error",
            Error: err.Error(),
        })
    }

    // listing all files
    for _, file := range files {
        fmt.Println("Name:", file.FileName())
        fmt.Println("Size:", file.Size())
        fmt.Println("Type MINE:", file.ContentType())
        fmt.Println("Bytes:", file.Bytes())
    }

    // optional
    // files.SaveAll("/my-dir/uploads")

    return c.Status(200).JSONIN(files)
})
```
#### 📌 Testing with cURL

##### 🔹Upload a single file:
```bash

$ curl -X POST http://localhost:8080/upload -F "file=@example.png"
```

##### 🔹 Upload multiple files:
```bash

$ curl -X POST http://localhost:8080/upload-multiple \
-F "files=@image1.jpg" -F "files=@document.pdf"
```


##### 📌 File Upload Feature Comparison with other Frameworks

| Framework  | `FormFile()` | `FormFiles()` | Dynamic Limit | File Metadata Methods (`FileName()`, `Size()`) | `Save()`, `SaveAll()` Method |
|------------|-------------|--------------|---------------|---------------------------------|------------|
| **Quick**  | ✅ Yes | ✅ Yes | ✅ Yes (`c.FormFileLimit("10MB")`) | ✅ Yes | ✅ Yes |
| Fiber      | ✅ Yes | ✅ Yes | ❌ No | ❌ No (uses `FileHeader` directly) | ✅ Yes |
| Gin        | ✅ Yes | ✅ Yes | ❌ No | ❌ No (uses `FileHeader` directly) | ❌ No |
| Echo       | ✅ Yes | ❌ No  | ❌ No | ❌ No | ❌ No |
| net/http   | ✅ Yes | ❌ No  | ❌ No | ❌ No | ❌ No |

### **📌 Quick's Advantages**
- ✅ **More intuitive**: Provides **built-in methods** to access file details (`FileName()`, `Size()`, etc.).
- ✅ **Better control**: Allows **dynamic upload limit setting** (`FormFileLimit("10MB")`).
- ✅ **Ease of use**: Includes a **`Save()` method** for direct file storage.

---

##### 🚀 You can now implement fast and efficient file uploads in Quick! 🔥

#### **📌 What I included in this README**

- ✅ **Detailed description of `multipart/form-data`**
- ✅ **Explanation of form-data objects**
- ✅ **Comparative table between Quick, Fiber, Gin, Echo and `net/http`**
- ✅ **Space for code examples in Go and tests with cURL**- 

Now you can **complete with your specific examples** where I left the spaces

##### 🚀 **If you need adjustments or improvements, just let me know!** 😃🔥