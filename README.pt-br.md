
![Logo do Quick](/quick_logo.png)


[![GoDoc](https://godoc.org/github.com/jeffotoni/quick?status.svg)](https://godoc.org/github.com/jeffotoni/quick) [![CircleCI](https://dl.circleci.com/status-badge/img/gh/jeffotoni/quick/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/jeffotoni/quick/tree/main) [![Go Report](https://goreportcard.com/badge/github.com/jeffotoni/quick)](https://goreportcard.com/badge/github.com/jeffotoni/quick) [![License](https://img.shields.io/github/license/jeffotoni/quick)](https://img.shields.io/github/license/jeffotoni/quick) ![CircleCI](https://img.shields.io/circleci/build/github/jeffotoni/quick/main)  ![Coveralls](https://img.shields.io/coverallsCoverage/github/jeffotoni/quick)  [![Build Status](https://github.com/alvarorichard/GoAnime/actions/workflows/ci.yml/badge.svg)](https://github.com/alvarorichard/GoAnime/actions) ![GitHub contributors](https://img.shields.io/github/contributors/jeffotoni/quick)
![GitHub stars](https://img.shields.io/github/last-commit/jeffotoni/quick) ![GitHub stars](https://img.shields.io/github/forks/jeffotoni/quick?style=social) ![GitHub stars](https://img.shields.io/github/stars/jeffotoni/quick)

<!-- [![Github Release](https://img.shields.io/github/v/release/jeffotoni/quick?include_prereleases)](https://img.shields.io/github/v/release/jeffotoni/quick) -->

<h2 align="center">
    <p>
         <a href="README.md">English</a> |
          <a href="README.pt-br.md">Рortuguês</a>
    </p> 
</h2>

# Quick - um roteador leve para Go! ![Quick](/quick.png)

🚀 Quick é um gerenciador de rotas **flexível e extensível** para a linguagem Go. Seu objetivo é ser **rápido e de alto desempenho**, além de ser **100% compatível com net/http**. Quick é um **projeto em constante desenvolvimento** e está aberto para **colaboração**, todos são bem-vindos a contribuir. 😍

💡 Se você é novo em programação, o Quick é uma ótima oportunidade para começar a aprender como trabalhar com o Go. Com **facilidade de uso** e recursos, você pode **criar rotas personalizadas** e expandir seu conhecimento do idioma.

👍 Espero que possam participar e desfrutar **Aproveite**! 😍

🔍 O repositório de exemplos do Framework Quick Run [Examples](https://github.com/jeffotoni/quick/tree/main/example).


# Quick em ação 💕🐧🚀😍
![Quick](quick_server.gif)

## 🎛️| Características

| Características   | Tem  | Status | Conclusão |
|--------------------------------------------------|------|--------|------------|
| 🛣️ Route Manager   | sim  | 🟢   | 100%   |
| 📁 Arquivos de servidor estático   | sim  | 🟢   | 100%   |
| 🔗 Http Client   | sim  | 🟢   | 100%   |
| 📤 Upload de arquivos (multipart/form-data)   | sim  | 🟢   | 100%   |
| 🚪 Agrupamento de rotas  | sim  | 🟢   | 100%   |
| 🛡️ Middleware   | sim  | 🟡   | 50%   |
| ⚡ Suporte HTTP/2   | sim  | 🟢   | 100%   |
| 🔄 Suporte para vinculação de dados em JSON, XML e formulários   | sim  | 🟢   | 100%   |
| 🔍 Suporte a Regex  | sim  | 🟡   | 80%   |
| 🌎 Site   | sim  | 🟡   | 90%   |
| 📚 Docs   | sim  | 🟡   | 40%   |


## 🗺️| Roteiro de Desenvolvimento

| Tarefa                                          | Progresso |
|-------------------------------------------------|-----------|
| Desenvolver MaxBodySize método Post   | 100%   |
| Desenvolver MaxBodySize método Put   | 100%   |
| Desenvolver configuração em  New(Config{}) não é obrigatório | 100%   |
| Criar função de impressão para não usar fmt demais | 100% |
| Criação de função própria para Concat String   | 100%   |
| Criação de benchmarking entre o. Stdout e fmt.Println | 100%   |
| Desenvolver suporte para o método GET   | 100%   |
| Desenvolver suporte para o método GET aceitando Query String | 100%   |
| Desenvolver suporte para o método GET aceitando Parametros | 100%   |
| Desenvolver suporte para o método GET que aceita Query String e parâmetros | 100% |
| Desenvolver suporte para o método GET que aceita expressão regular | 100% |
| Desenvolver suporte para o método POST   | 100%   |
| Desenvolver rotas método POST aceitando JSON   | 100%   |
| Desenvolver para MÉTODO POST o parse JSON   | 100%   |
| Desenvolver para as funções MÉTODO POST para acessar byte ou string de Parse | 100% |
| Desenvolver para MÉTODO PUT  | 100%   |
| Desenvolver para o MÉTODO PUT a parse JSON   | 100%   |
| Desenvolver para o MÉTODO PUT a parse JSON   | 100%   |
| Desenvolver para MÉTODO PUT funções para acessar byte ou string a partir do Parse | 100% |
| Desenvolver para MÉTODO DELETE   | 100%   |
| Desenvolver método para ListenAndServe   | 100%   |
| Desenvolver ServeHTTP suporte   | 100%   |
| Desenvolver suporte middleware    | 100%   |
| Desenvolver suporte para middleware compress   | 100%   |
| Desenvolver suporte para middleware cors   | 100%   |
| Desenvolver logger middleware suporte   | 100%   |
| Desenvolver suporte para maxbody middlewares   | 100%   |
| Desenvolver middleware suporte msgid   | 100%   |
| Desenvolver middleware suporte msguuid   | 100%   |
| Desenvolver suporte Cors   | 100%   |
| Desenvolver Client Get suporte  | 100%   |
| Desenvolver Client Post suporte   | 100%   |
| Desenvolver Client Put suporte   | 100%   |
| Desenvolver Client suporte Delete   | 100%   |

---

## 🚧| Rodmap em andamento

 | Tarefa   | Progresso |
|--------------------------------------------------|-----------|
| Desenvolver e relacionar-se com o Listen the Config   | 42%   |
| Desenvolve suporte para Uploads e Uploads Múltiplos | 100% |
| Desenvolve suporte para JWT | 10% |
| Desenvolver método para facilitar a manipulação de ResponseWriter | 80% |
| Desenvolver método para facilitar o tratamento do pedido | 80%  |
| Desenvolver padrão de teste unitário  | 90%   |



## 🚀| Rodmap para desenvolvimento


| Tarefa   | Progresso |
|---------------------------------------------------|-----------|
| Documentação Testes Exemplos PKG Go | 45% |
| Cobertura do teste go-cover | 74,6% |
| Cobertura de recurso regex, mas possibilidades | 0.% |
| Desenvolver para OPÇÕES MÉTODO | 100% |
| Desenvolver suporte para o método CONNECT [Veja mais](https://www.rfc-editor.org/rfc/rfc9110.html#name-connect) | 0.% |
| Desenvolver método para ListenAndServeTLS (http2) | 0.% |
| Desenvolver suporte de arquivos estáticos | 100% |
| Suporte ao WebSocket | 0.% |
| Suporte do limitador de taxa   | 0.%   |
| Modelos de motores   | 0.%   |
| Documentação Testes Exemplos PKG Go   | 45%   |
| Cobertura do teste go -cover   | 75,5%   |
| Cobertura de recursos Regex, mas possibilidades   | 0.%   |
| Desenvolver MÉTODO para OPTION   | 100%   |
| Desenvolver MÉTODO DE CONEXÃO [Veja mais](https://www.rfc-editor.org/rfc/rfc9110.html#name-connect)   | 0.%   |
| Desenvolver método para ListenAndServeTLS (http2) | 0.%   |
| Criar uma CLI (Interface de Linha de Comando) Quick.   | 0.%   |


## 📊| Cobertura de Testes

| Arquivo        | Cobertura | Status |
|---------------|-----------|--------|
| Ctx           | 84.1%     | 🟡     |
| Group         | 100.0%    | 🟢     |
| Http Status   | 7.8%      | 🔴     |
| Client        | 83.3%     | 🟢     |
| Mock          | 100.0%    | 🟢     |
| Concat        | 100.0%    | 🟢     |
| Log           | 0.0%      | 🔴     |
| Print         | 66.7%     | 🟡     |
| Cors          | 76.0%     | 🟡     |
| Logger        | 100.0%    | 🟢     |
| Maxbody       | 100.0%    | 🟢     |
| Quick         | 79.5%     | 🟡     |
| QuickTest     | 100.0%    | 🟢     |

---

### Exemplo rápido Quick
```go

package main

import "github.com/jeffotoni/quick"

func main() {
    q := quick.New()

    q.Get("/v1/user", func(c *quick.Ctx) error {
        c.Set("Content-Type", "application/json")
        return c.Status(200).SendString("Quick em ação ❤️!")
    })

    q.Listen("0.0.0.0:8080")
}

```

```bash

$ curl -i -XGET -H "Content-Type:application/json" \
'localhost:8080/v1/user'
HTTP/1.1 200 OK
Content-Type: application/json
Date: Wed, 22 Feb 2023 07:45:36 GMT
Content-Length: 23

Quick em ação ❤️!

```

### Obter parâmetros Quick
```go

package main

import "github.com/jeffotoni/quick"

func main() {
    q := quick.New()

    q.Get("/v1/customer/:param1/:param2", func(c *quick.Ctx) error {
        c.Set("Content-Type", "application/json")

        type my struct {
            Msg string `json:"msg"`
            Key string `json:"key"`
            Val string `json:"val"`
        }

        return c.Status(200).JSON(&my{
            Msg: "Quick ❤️",
            Key: c.Param("param1"),
            Val: c.Param("param2"),
        })
    })

    q.Listen("0.0.0.0:8080")
}

```

```bash

$ curl -i -XGET -H "Content-Type:application/json" \
'localhost:8080/v1/customer/val1/val2'
HTTP/1.1 200 OK
Content-Type: application/json
Date: Wed, 22 Feb 2023 07:45:36 GMT
Content-Length: 23

{"msg":"Quick ❤️","key":"val1","val":"val2"}

```

### Corpo da postagem json Quick
```go

package main

import "github.com/jeffotoni/quick"

type My struct {
    Name string `json:"name"`
    Year int    `json:"year"`
}

func main() {
    q := quick.New()
    q.Post("/v1/user", func(c *quick.Ctx) error {
        var my My
        err := c.BodyParser(&my)
        if err != nil {
            return c.Status(400).SendString(err.Error())
        }

        return c.Status(200).String(c.BodyString())
        // ou 
        // c.Status(200).JSON(&my)
    })

    q.Listen("0.0.0.0:8080")
}

```
### 📌 cURL

```bash

$ curl -i -XPOST -H "Content-Type:application/json" \
'localhost:8080/v1/user' \
-d '{"name":"jeffotoni", "year":1990}'
HTTP/1.1 200 OK
Date: Wed, 22 Feb 2023 08:10:06 GMT
Content-Length: 32
Content-Type: text/plain; charset=utf-8

{"name":"jeffotoni","year":1990}

```

### Carrega dados multipart/form-data

Quick fornece uma API simplificada para gerenciar uploads, permitindo que você facilmente recuperar e manipular arquivos.

### ✅ **Principais métodos e funcionalidades**:
| Método | Descrição |
|---------|-------------|
| `c.FormFile("file")` | Retorna um único arquivo carregado no formulário. |
| `c.FormFiles("files")` | Retorna uma lista de arquivos carregados (uploads múltiplos). |
| `c.FormFileLimit("10MB")` | Define um limite de upload (o padrão é 1MB ). |
| `uploadedFile.FileName()` | Retorna o nome do arquivo. |
| `uploadedFile.size()` | Retorna o tamanho do arquivo em bytes. |
| `uploadedFile.ContentType()` | Retorna o tipo MIME do arquivo. |
| `uploadedFile.Bytes()` | Retorna os bytes do arquivo. |
| `uploadedFile.Save("/path/")` | Salva o arquivo em um diretório especificado. |
| `uploadedFile.Save("/path", "your-name-file")` | Salva o arquivo com seu nome. |
| `uploadedFile.SaveAll("/path")` | Salva o arquivo em um diretório especificado. |

---

##### 📌 Comparação de recursos de upload de arquivos com outros frameworks

| Framework  | `FormFile()` | ‘FormFiles()’| Dynamic Limit | Methods (‘FileName()‘, ‘Size()‘) | ‘Save()‘, ‘SaveAll()’
|------------|-------------|--------------|---------------|---------------------------------|------------|
| **Quick**  | ✅ Sim | ✅ Sim | ✅ Sim | ✅ Sim | ✅ Sim |
| Fiber   | ✅ Sim | ✅ Sim | ❌ Não | ❌ Não (usa o FileHeader diretamente) | ✅ Sim |
| Gin   | ✅ Sim | ✅ Sim | ❌ Não | ❌ Não (usa o FileHeader diretamente) | ❌ Não |
| Echo   | ✅ Sim | ❌ Não  | ❌ Não | ❌ Não | ❌ Não |
| net/http   | ✅ Sim | ❌ Não  | ❌ Não | ❌ Não | ❌ Não |
---

### 📌 Arquivo Upload Exemplo

```go
package main

import (
    "fmt"
    "github.com/jeffotoni/quick"
)

func main() {
    // inicia Quick
    q := quick.New()

    q.Post("/upload", func(c *quick.Ctx) error {
        // definir limite de upload
        c.FormFileLimit("10MB")

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

        // Salvar o arquivo (opcional)
        // uploadedFile.Save("/tmp/uploads")

        return c.Status(200).JSONIN(uploadedFile)

    })

     q.Listen("0.0.0.0:8080")
}
```
### 📌 Exemplo de upload múltiplo

```go
package main

import (
    "fmt"
    "github.com/jeffotoni/quick"
)

func main() {
     // inicia Quick
    q := quick.New()

    q.Post("/upload-multiple", func(c *quick.Ctx) error {
        // definir limite de upload
        c.FormFileLimit("10MB")

        // recebe arquivos
        files, err := c.FormFiles("files")
        if err != nil {
            return c.Status(400).JSON(Msg{
                Msg:   "Upload error",
                Error: err.Error(),
            })
        }

         // listando todos os arquivos
        for _, file := range files {
            fmt.Println("Name:", file.FileName())
            fmt.Println("Size:", file.Size())
            fmt.Println("Type MINE:", file.ContentType())
            fmt.Println("Bytes:", file.Bytes())
        }

        // opcional
        // files.SaveAll("/my-dir/uploads")

        return c.Status(200).JSONIN(files)
    })
        
    q.Listen("0.0.0.0:8080")
}
```
### 📌 Testar com cURL

##### 🔹Carregar um único arquivo:
```bash
$ curl -X POST http://localhost:8080/upload -F "file=@example.png"
```

##### 🔹 Carregar vários arquivos:
```bash

$ curl -X POST http://localhost:8080/upload-multiple \
-F "files=@image1.jpg" -F "files=@document.pdf"
```

### Quick Post Bind json
```go

package main

import "github.com/jeffotoni/quick"

type My struct {
    Name string `json:"name"`
    Year int    `json:"year"`
}

func main() {
    q := quick.New()
    q.Post("/v2/user", func(c *quick.Ctx) error {
        var my My
        err := c.Bind(&my)
        if err != nil {
            return c.Status(400).SendString(err.Error())
        }
        return c.Status(200).JSON(&my)
    })

    q.Listen("0.0.0.0:8080")
}

```

### 📌 cURL
```bash

$ curl -i -XPOST -H "Content-Type:application/json" \
'localhost:8080/v2/user' \
-d '{"name":"Marcos", "year":1990}'
HTTP/1.1 200 OK
Date: Wed, 22 Feb 2023 08:10:06 GMT
Content-Length: 32
Content-Type: text/plain; charset=utf-8

{"name":"Marcos","year":1990}

```

### Cors

```go

package main

import (
    "github.com/jeffotoni/quick"
    "github.com/jeffotoni/quick/middleware/cors"
)

func main() {
    q := quick.New()
    q.Use(cors.New())

    q.Get("/v1/user", func(c *quick.Ctx) error {
        c.Set("Content-Type", "application/json")
        return c.Status(200).SendString("Quick em ação com Cors❤️!")
    })

    q.Listen("0.0.0.0:8080")
}

```

### quick.New(quick.Config{})
```go

package main

import "github.com/jeffotoni/quick"

func main() {
    q := quick.New(quick.Config{
        MaxBodySize: 5 * 1024 * 1024,
    })

    q.Get("/v1/user", func(c *quick.Ctx) error {
        c.Set("Content-Type", "application/json")
        return c.Status(200).SendString("Quick em ação com Cors❤️!")
    })

    q.Listen("0.0.0.0:8080")
}

```

### quick.Group()
```go
package main

import "github.com/jeffotoni/quick"

func main() {
    q := quick.New(quick.Config{
        MaxBodySize: 5 * 1024 * 1024,
    })

    v1 := q.Group("/v1")
    v1.Get("/user", func(c *quick.Ctx) error {
        return c.Status(200).SendString("[GET] [GROUP] /v1/user ok!!!")
    })
    v1.Post("/user", func(c *quick.Ctx) error {
        return c.Status(200).SendString("[POST] [GROUP] /v1/user ok!!!")
    })

    v2 := q.Group("/v2")
    v2.Get("/user", func(c *quick.Ctx) error {
        c.Set("Content-Type", "application/json")
        return c.Status(200).SendString("Quick em ação com [GET] /v2/user ❤️!")
    })

    v2.Post("/user", func(c *quick.Ctx) error {
        c.Set("Content-Type", "application/json")
        return c.Status(200).SendString("Quick em ação com [POST] /v2/user ❤️!")
    })

    q.Listen("0.0.0.0:8080")
}

```

### Quick Tests
```go

package main

import (
    "io"
    "strings"
    "testing"

    "github.com/jeffotoni/quick"
)

func TestQuickExample(t *testing.T) {

    // Aqui está uma função do manipulador
    testSuccessMockHandler := func(c *quick.Ctx) error {
        c.Set("Content-Type", "application/json")
        b, _ := io.ReadAll(c.Request.Body)
        resp := `"data":` + string(b)
        return c.Byte([]byte(resp))
    }

    q := quick.New()
    // Aqui você pode criar todas as rotas que deseja testar
    q.Post("/v1/user", testSuccessMockHandler)
    q.Post("/v1/user/:p1", testSuccessMockHandler)

    wantOutData := `"data":{"name":"jeff", "age":35}`
    reqBody := []byte(`{"name":"jeff", "age":35}`)
    reqHeaders := map[string]string{"Content-Type": "application/json"}

    data, err := q.QuickTest("POST", "/v1/user", reqHeaders, reqBody)
    if err != nil {
        t.Errorf("error: %v", err)
        return
    }

    s := strings.TrimSpace(data.BodyStr())
    if s != wantOutData {
        t.Errorf("esperado %s e %s recebido", wantOutData, s)
        return
    }

    t.Logf("\nOutputBodyString -> %v", data.BodyStr())
    t.Logf("\nStatusCode -> %d", data.StatusCode())
    t.Logf("\nOutputBody -> %v", string(data.Body())) 
    t.Logf("\nResponse -> %v", data.Response())
}

```

### quick.regex
```go

package main

import (
    "github.com/jeffotoni/quick"
    "github.com/jeffotoni/quick/middleware/msgid"
)

func main() {
    q := quick.New()

    q.Use(msgid.New())

    q.Get("/v1/user/{id:[0-9]+}", func(c *quick.Ctx) error {
        c.Set("Content-Type", "application/json")
        return c.Status(200).String("Quick ação total!!!")
    })

    q.Listen("0.0.0.0:8080")
}

```
### 🔑 Autenticação básica

Autenticação básica (Auth básico) é um mecanismo de autenticação simples definido no RFC 7617. É comumente usado para autenticação baseada em HTTP, permitindo que os clientes forneçam credenciais (nome de usuário e senha) no cabeçalho da solicitação.

**🔹 Como funciona** 
  1. O cliente codifica o nome de usuário e a senha no Base64 (username:password dXNlcm5hbWU6cGFzc3dvcmQ=).
  2. As credenciais codificadas são enviadas no cabeçalho Authorization:

```bash
Authorization: Basic dXNlcm5hbWU6cGFzc3dvcmQ=
```
   3. O servidor decodifica as credenciais e as verifica antes de conceder acesso.

#### **⚠️ Considerações de segurança**
- Não criptografado: Basic Auth apenas codifica as credenciais em Base64, mas não as criptografa.
- Use over HTTPS: Sempre use a autenticação básica com TLS/SSL (HTTPS) para evitar que as credenciais sejam expostas.
- Métodos alternativos de autenticação: para maior segurança, considere chaves OAuth2, JWT ou API.

O Basic Auth é adequado para casos de uso simples, mas para aplicações de produção são recomendados mecanismos de autenticação mais fortes. 🚀

#### Variáveis de ambiente básicas do Auth

Este exemplo configura a autenticação básica usando variáveis de ambiente para armazenar as credenciais com segurança.
as rotas abaixo são afetadas, para isolar o grupo de uso da rota para aplicar somente às rotas no grupo.

```go
package main

import (
	"log"
	"os"

	"github.com/jeffotoni/quick"
	middleware "github.com/jeffotoni/quick/middleware/basicauth"
)

// export USER=admin
// export PASSWORD=1234

var (
	User     = os.Getenv("USER")
	Password = os.Getenv("PASSORD")
)

func main() {

	q := quick.New()

	q.Use(middleware.BasicAuth(User, Password))

	q.Get("/protected", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.SendString("You have accessed a protected route!")
	})

	//Iniciar servidor
	log.Fatal(q.Listen("0.0.0.0:8080"))
}

```
---

### Autenticação básica com middleware rápido

Este exemplo usa o middleware BasicAuth fornecido pela Quick, oferecendo uma configuração de autenticação simples.

```go
package main

import (
	"log"

	"github.com/jeffotoni/quick"
	middleware "github.com/jeffotoni/quick/middleware/basicauth"
)

func main() {

	//Inicia Quick
	q := quick.New()

	// chama middleware
	q.Use(middleware.BasicAuth("admin", "1234"))

	// tudo abaixo vai usar o middleware
	q.Get("/protected", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.SendString("You have accessed a protected route!")
	})

	//Iniciar servidor
	log.Fatal(q.Listen("0.0.0.0:8080"))
}

```
### Autenticação básica com grupos de rota rápida

Este exemplo mostra como aplicar a autenticação básica a um grupo específico de rotas usando a funcionalidade de grupos do Quick.
Quando usamos o grupo podemos isolar o middleware, isso funciona para qualquer middleware em Quick.

```go

package main

import (
	"log"

	"github.com/jeffotoni/quick"
	middleware "github.com/jeffotoni/quick/middleware/basicauth"
)

func main() {

	q := quick.New()

	// usando o grupo para isolar rotas e middlewares
	gr := q.Group("/")

	// middleware BasicAuth
	gr.Use(middleware.BasicAuth("admin", "1234"))

	// route public
	q.Get("/v1/user", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.SendString("Public quick route")
	})

	// rota protegida
	gr.Get("/protected", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.SendString("You have accessed a protected route!")
	})

	//Iniciar servidor
	log.Fatal(q.Listen("0.0.0.0:8080"))
}

```
### BasicAuth Personalizado
Este exemplo mostra uma implementação personalizada de autenticação básica sem usar qualquer middleware. Ele verifica manualmente as credenciais do usuário e aplica a autenticação às rotas protegidas.

Em quick você tem permissão para fazer sua própria implementação personalizada diretamente no q.Use(..), ou seja, você será capaz de implementá-lo diretamente se desejar.

```go
package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"

	"github.com/jeffotoni/quick"
)

func main() {
	q := quick.New()

	// implementando middleware diretamente no Uso
	q.Use(func(next http.Handler) http.Handler {
		// credentials
		username := "admin"
		password := "1234"

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Verifique se começa com "Basic"
			if !strings.HasPrefix(authHeader, "Basic ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Descodificar credenciais
			payload, err := base64.StdEncoding.DecodeString(authHeader[len("Basic "):])
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			creds := strings.SplitN(string(payload), ":", 2)
			if len(creds) != 2 || creds[0] != username || creds[1] != password {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	q.Get("/protected", func(c *quick.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.SendString("You have accessed a protected route!")
	})

	//Iniciar servidor
	log.Fatal(q.Listen("0.0.0.0:8080"))

}

```
---

#### 📂 ARQUIVOS ESTÁTICOS

Um servidor de arquivos estático é uma característica fundamental em frameworks web, permitindo o serviço eficiente de conteúdo estático como HTML, CSS, JavaScript, imagens e outros ativos. É útil para hospedar aplicações front-end, fornecendo arquivos para download ou servindo recursos diretamente do backend.

🔹 Como funciona
    
1. O servidor escuta para solicitações HTTP que visam caminhos de arquivos estáticos.
2. Se um arquivo solicitado existir no diretório configurado, o servidor lê e retorna o arquivo como uma resposta.
3. Os tipos MIME são determinados automaticamente com base na extensão do ficheiro.

:zap: Características principais
- Manuseio eficiente: Serve arquivos diretamente sem processamento adicional.
- Detecção de tipo MIME: identifica automaticamente os tipos de arquivo para renderização adequada.
- Suporte a cache: pode ser configurado para melhorar o desempenho através de cabeçalhos HTTP.
- Listagem de diretório: (Opcional) Permite a navegação nos arquivos estáticos disponíveis.

:warning: Considerações de segurança
- Restringir o acesso a arquivos confidenciais (.env, .git, etc).
- Configure políticas CORS quando necessário.
- Use uma política de segurança de conteúdo (CSP) para mitigar os riscos XSS.

#### Servindo arquivos estáticos com o Quick Framework

Este exemplo configura um servidor web básico que serve arquivos estáticos, como HTML, CSS ou JavaScript.

```go
package main

import "github.com/jeffotoni/quick"

func main() {
    
     // Criar uma nova instância de servidor Quick
    q := quick.New()

    // Configuração de arquivos estáticos
    // Serve arquivos do diretório ". /static" 
    // no caminho de URL "/static".
    q.Static("/static", "./static")

     // Definição da rota
    // Define uma rota para servir o arquivo "index.html" ao acessar "/".
    q.Get("/", func(c *quick.Ctx) error {
        c.File("./static/index.html")
        return nil
    })

    // Iniciando o servidor
    // Inicia o servidor para ouvir as solicitações recebidas na porta 8080.
    q.Listen("0.0.0.0:8080")
}


```
---

#### 📁 EMBED
🔹 Como funcionam os arquivos estáticos incorporados
    
1. Os ativos estáticos são compilados diretamente no binário no momento da compilação (por exemplo, usando o pacote embed do Go).
2. O aplicativo serve esses arquivos da memória em vez de ler do disco.
3. Isso elimina dependências externas, facilitando a implantação.

:zap:  Vantagens dos arquivos incorporados
- Portabilidade: Distribuição binária única sem arquivos extras.
- Desempenho: acesso mais rápido a ativos estáticos, pois eles são armazenados na memória.
- Segurança: reduz a exposição a ataques externos do sistema de arquivos.

### Embedding Files
Ao incorporar arquivos estáticos em um executável binário, o servidor não depende de um sistema de arquivos externo para servir os ativos. Essa abordagem é útil para aplicativos autônomos, ferramentas CLI e implantações multiplataforma onde as dependências devem ser minimizadas.

Este exemplo incorpora arquivos estáticos no binário usando o pacote embed e os serve usando a estrutura Quick.

```go
package main

import (
	"embed"

	"github.com/jeffotoni/quick"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	// Inicialização do servidor
	// Cria uma nova instância do servidor Quick
	q := quick.New()

	// Configuração de arquivos estáticos
	// Define o diretório para servir arquivos estáticos usando os arquivos incorporados
	q.Static("/static", staticFiles)

	// Definição da rota
	// Define uma rota que serve o arquivo de índice HTML
	q.Get("/", func(c *quick.Ctx) error {
		c.File("./static/index.html") // Renderiza o arquivo index.html
		return nil
	})

	// Iniciando o servidor
	// Inicia o servidor na porta 8080, ouvindo em todos os endereços
	q.Listen("0.0.0.0:8080")
}

```
---
### 🌍 HTTP Client 
O pacote de cliente HTTP em Quick fornece uma maneira simples e flexível para fazer solicitações HTTP, suportando operações GET, POST, PUT e DELETE. É projetado para lidar com diferentes tipos de corpos de solicitação e analisar facilmente as respostas.

Este cliente abstrai o processamento HTTP de baixo nível e oferece:

- Funções de conveniência (Get, Post, Put, Delete) para fazer solicitações rápidas usando um cliente padrão.
- Solicitações personalizáveis com suporte para cabeçalhos, autenticação e configurações de transporte.
- Corpo de análise flexível, permitindo aos usuários enviar JSON, texto simples ou personalizado io. tipos de leitor.
- Marshaling e unmarshaling automático de JSON, simplificando a interação com APIs.

#### Exemplo de solicitação GET

```go
package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
    // Use o client padrão
	resp, err := client.Get("https://reqres.in/api/users")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("GET response:", string(resp.Body))
}
```

#### Exemplo de solicitação POST (usando uma estrutura)
```go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
	// Definir uma estrutura para enviar como JSON
	data := struct {
		user string `json:"user"`
	}{
		user: "Emma",
	}

	// POST request to ReqRes API
	resp, err := client.Post("https://reqres.in/api/users", data)
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal a resposta JSON (se aplicável)
	var result map[string]string
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		log.Fatal(err)
	}
	fmt.Println("POST response:", result)
}
```

#### Exemplo de solicitação PUT (usando uma string)
```go
package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
	// Definir uma estrutura com dados do usuário
	data := struct {
		user string `json:"name"`
	}{
		user: "Jeff",
	}

	// Converte struct para JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal("Erro ao codificar JSON:", err)
	}

	// PUT requisição para ReqRes API
	resp, err := client.Put("https://reqres.in/api/users/2", string(jsonData))
	if err != nil {
		log.Fatal("Erro ao fazer o pedido:", err)
	}

	// Imprimir o corpo de status e resposta HTTP
	fmt.Println("Código de status HTTP:", resp.StatusCode)
	fmt.Println("Corpo da resposta", string(resp.Body))
}
```

#### Exemplo de solicitação DELETE

```go
package main

import (
	"fmt"
	"log"

	"github.com/jeffotoni/quick/http/client"
)

func main() {

	// DELETE requisição para ReqRes API
	resp, err := client.Delete("https://reqres.in/api/users/2")
	if err != nil {
		log.Fatal("Erro ao fazer solicitação:", err)
	}

	// Imprimir o status HTTP para confirmar a exclusão
	fmt.Println("Código de status HTTP:", resp.StatusCode)

	// Como DELETE geralmente não retorna nenhum conteúdo, verificamos se está vazio
	if len(resp.Body) > 0 {
		fmt.Println("Corpo da resposta:", string(resp.Body))
	} else {
		fmt.Println("O corpo da resposta está vazio (esperado para 204 sem conteúdo)")
	}
}
```
---

# Qtest - Utilitário de teste HTTP para Quick

Qtest é uma função de teste HTTP **avançada** projetada para simplificar a validação de rotas dentro da estrutura **Quick***. Permite o teste perfeito de pedidos HTTP simulados usando o «httptest», suportando:

- **Métodos HTTP personalizados** («GET», «POST», «PUT», «DELETE», etc.).
- **Cabeçalhos personalizados**.
- **Parâmetros de consulta**.
- **Solicitar corpo**.
- **Cookies**.
- **Métodos de validação integrados** para códigos de status, cabeçalhos e corpos de resposta.


## 📌 Visão geral
A função `Qtest` pega uma estrutura `QuickTestOptions', contendo parâmetros de solicitação, executa a solicitação e retorna um objeto `QtestReturn', que fornece métodos para analisar e validar o resultado.

```go
func TestQTest_Options_POST(t *testing.T) {
    // iniciar Quick
    q := New()

    // Definir a rota POST
    q.Post("/v1/user/api", func(c *Ctx) error {
        c.Set("Content-Type", "application/json") // Configuração de cabeçalho simplificada 
        return c.Status(StatusOK).String(`{"message":"Success"}`)
    })

     // Configurar parâmetros de teste
    opts := QuickTestOptions{
        Method: "POST",
        URI:    "/v1/user/api",
        QueryParams: map[string]string{
            "param1": "value1",
            "param2": "value2",
        },
        Body: []byte(`{"key":"value"}`),
        Headers: map[string]string{
            "Content-Type": "application/json",
        },
        Cookies: []*http.Cookie{
            {Name: "session", Value: "abc123"},
        },
        LogDetails: true, // Habilita o registro detalhado
    }

    // Executa o teste
    result, err := q.Qtest(opts)
    if err != nil {
        t.Fatalf("Erro em Qtest: %v", err)
    }
    
    // Validações
    if err := result.AssertStatus(StatusOK); err != nil {
        t.Errorf("Afirmação de status falhou: %v", err)
    }

    if err := result.AssertHeader("Content-Type", "application/json"); err != nil {
        t.Errorf("Falha na asserção do cabeçalho: %v", err)
    }

    if err := result.AssertBodyContains("Success"); err != nil {
        t.Errorf("A afirmação do corpo falhou: %v", err)
    }
}
```

🚀 **Mais detalhes aqui [Qtest - Quick](https://github.com/jeffotoni/quick/tree/main/quickTest)**

---
# 🔄 Mecanismos de retentativa e failover no cliente HTTP rápido

O **Quick HTTP Client** agora inclui **built-in retry e suporte de failover**, permitindo solicitações HTTP mais resilientes e confiáveis. Esses recursos são essenciais para lidar com **falhas transientes***, **instabilidade de rede** e **tempo de inatividade do serviço** de forma eficiente.

## 🚀 Principais características
- **Tentativas automáticas**: Repetições de solicitações com falha baseadas em regras configuráveis.
- **Exponential Backoff**: aumenta gradualmente o atraso entre tentativas de repetição.
- **Retries Status-Based**: Retries somente em códigos de status HTTP especificados (por exemplo, '500', '502', '503').
- **Mecanismo de failover**: alterna para URLs de backup predefinidos se a solicitação primária falhar.
- **Logging Support**: Permite logs detalhados para o comportamento de repetição da depuração.

---

## 🔹 Como funciona o Retry & Failover
O mecanismo de repetição funciona **automaticamente reenviando a solicitação** se ela falhar, com opções para **limitar repetições**, **introduzir atrasos de backoff**, e **repetir somente para status específicos de resposta**. O sistema de failover garante **alta disponibilidade** redirecionando solicitações com falha para URLs alternativos.

### ✅ Opções de configuração:
| Opção   | Descrição |
|-----------------------|-------------|
| **MaxRetries**   | Define o número de tentativas de repetição. |
| **Delay**   | Especifica o atraso antes de cada nova tentativa. |
| **UseBackoff**   | Permite que o backoff exponencial aumente o atraso dinamicamente. |
| **Status**   | Lista de códigos de status HTTP que acionam uma nova tentativa. |
| **FailoverURLs**   | Lista de URLs de backup para failover em caso de falhas repetidas. |
| **EnableLog**   | Ativa o registro para tentativas de repetição de depuração. |

---

### **Retry com atraso exponencial**
Este exemplo demonstra **repetir uma solicitação** com um atraso crescente (backoff) quando encontrar erros.

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
	cClient := client.New(
		cliente.WithRetry(
			client.RetryConfig{
				MaxRetries: 3,   // Número máximo de tentativas
				Atraso:   1 * tempo. Segundo,  // Atraso inicial da repetição
				UseBackoff: true,   // Permite backoff exponencial
				Estados:   []int{500, 502, 503}, // Retries apenas nestes códigos de status HTTP
				EnableLog:  true,   // Ativa o registro para novas tentativas
		}),
	)

	resp, err := cClient.Get("http://localhost:3000/v1/resource")
	if err != nil {
		log.Fatal("GET request failed:", err)
	}
	fmt.Println("GET Response:", string(resp.Body))
}

```

### **Failover para URLs de backup**
Este exemplo muda para um URL de backup quando a solicitação primária falha.

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jeffotoni/quick/http/client"
)

func main() {
	cClient := cliente. Novo(
		client.WithRetry(client.RetryConfig{
			MaxRetries:   2,  // Tente a solicitação duas vezes antes de alternar
			Atraso:   2 * tempo. Segundo,  // Aguarde 2 segundos antes de tentar novamente
			Status:   []int{500, 502, 503}, // Acionar failover nesses erros
			FailoverURLs: []string{"http://backup1.com/resource", "https://reqres.in/api/users", "https://httpbin.org/post"}
			EnableLog: true, // Habilitar registros de repetição
		}),
	)


	resp, err := cClient.Get("http://localhost:3000/v1/resource")
	if err != nil {
		log.Fatal("Request failed:", err)
	}
	fmt.Println("Response:", string(resp.Body))
}

```
---
## 📝 Envio de formulário com PostForm no cliente HTTP Quick

O Quick HTTP Client agora inclui suporte embutido para `PostForm`, permitindo o manuseio perfeito de envios de formulários codificados por aplicação/formulário `x-www-urlencoded`. Esse recurso simplifica a interação com serviços da web e APIs que exigem dados codificados por formulário, tornando-o ideal para solicitações de autenticação, envio de dados e integrações de sistemas legados.


## 🔹Por que usar `PostForm`? 

| Característica   | Benefício |
|------------------------|---------|
| **Otimizado para formulários** | Simplifica o envio de dados codificados por formulário ('application/x-www-form-urlencoded’). |
| **Automatic Encoding**  | Converte `url. Values` em uma carga útil válida de envio de formulário. |
| **Header Management**   | Define automaticamente o tipo de conteúdo para aplicação/x-www-form-urlencoded.
| **Consistente API**   | Segue o mesmo design que `Post`, ‘Get’, ‘Put’, etc. |
| **Melhor compatibilidade** | Funciona com APIs que não aceitam cargas JSON. |

---
## 🔹 Como funciona o PostForm

O método PostForm codifica parâmetros de formulário, adiciona cabeçalhos necessários e envia uma solicitação HTTP POST para a URL especificada. Ele é projetado especificamente para APIs e serviços web que não aceitam cargas de JSON, mas exigem dados codificados por formulário.


### 🔹 **Servidor rápido com envio de formulário**
O exemplo a seguir demonstra como enviar dados codificados por formulário usando Quick PostForm:

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/jeffotoni/quick"
	"github.com/jeffotoni/quick/http/client"
)

func main() {
	q := quick.New()

	// Definir uma rota para processar POST form-data
	q.Post("/postform", func(c *quick.Ctx) error {
		form := c.FormValues()
		return c.JSON(map[string]any{
			"message": "Received form data",
			"data":    form,
		})
	})

	// Inicie o servidor em uma goroutine separada
	go func() {
		fmt.Println("Quick server running at http://localhost:3000")
		if err := q.Listen(":3000"); err != nil {
			log.Fatalf("Failed to start Quick server: %v", err)
		}
	}()

	// Criando um cliente HTTP antes de chamar o PostForm
	cClient := client.New(
		client.WithTimeout(5*time.Second), // Define um timeout de 5s
		client.WithHeaders(map[string]string{
			"Content-Type": "application/x-www-form-urlencoded", // Tipo correto para forms
		}),
	)

	// Verificar se o cliente HTTP foi inicializado corretamente
	if cClient == nil {
		log.Fatal("Erro: cliente HTTP não foi inicializado corretamente")
	}

	// Declare valores
	formData := url.Values{}
	formData.Set("username", "quick_user")
	formData.Set("password", "supersecret")

	// Enviar uma requisição POST
	resp, err := cClient.PostForm("http://localhost:3000/postform", formData)
	if err != nil {
		log.Fatalf("PostForm request with retry failed: %v", err)
	}

	// Verifique se a resposta é válida
	if resp == nil || resp.Body == nil {
		log.Fatal("Erro: resposta vazia ou inválida")
	}

	// Unmarshal a resposta JSON (se aplicável)
	var result map[string]any
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		log.Fatal(err)
	}
	fmt.Println("POST response:", result)
}

```


---

## 📚| Mais exemplos

Este diretório contém exemplos práticos do Quick Framework, um framework web rápido e leve desenvolvido em Go. Os exemplos são organizados em pastas separadas, cada uma contendo um exemplo completo de uso do framework em um aplicativo web simples. Se você tiver algum exemplo interessante de uso do Quick Framework, sinta-se à vontade para enviar uma solicitação de pull com sua contribuição. O repositório de exemplos do Quick Framework pode ser encontrado [aqui](https://github.com/jeffotoni/quick/tree/main/example).


## 🤝| Contribuições

Nós já temos vários exemplos, e já podemos testar e jogar 😁. Claro, estamos no início, ainda tem muito o que fazer. 
Sinta-se livre para fazer **PR** (em risco de ganhar uma t-shirt Go ❤️ e, claro, o reconhecimento como um profissional Go 😍 no mercado de trabalho).


## 🚀 **Quick Project Supporters** 🙏

O Quick Project visa desenvolver e fornecer software de qualidade para a comunidade de desenvolvedores. 💻 Para continuar melhorando nossas ferramentas, contamos com o apoio de nossos patrocinadores em Patreon. 🤝

Agradecemos a todos os nossos apoiadores! 🙌 Se você também acredita em nosso trabalho e quer contribuir para o avanço da comunidade de desenvolvimento, considere apoiar o Project Quick no nosso Patreon [aqui](https://www.patreon.com/jeffotoni_quick)

Juntos podemos continuar a construir ferramentas incríveis! 🚀


| Avatar | Usuário | Colaboração |
|--------|------|----------|
| <img src="https://avatars.githubusercontent.com/u/1092879?s=96&v=4" height=20> | [@jeffotoni](https://github.com/jeffotoni) | x 10 |
| <img src="https://avatars.githubusercontent.com/u/99341377?s=400&u=095679b08054e215561a4d4b08da764c2de619e6&v=4" height=20> | [@Crow3442](https://github.com/Crow3442) | x 5  |
| <img src="https://avatars.githubusercontent.com/u/70351793?v=4" height=20> | [@Guilherme-De-Marchi](https://github.com/Guilherme-De-Marchi) | x 5 |
| <img src="https://avatars.githubusercontent.com/u/59976892?v=4" height=20> | [@jaquelineabreu](https://github.com/jaquelineabreu) | x 1 |











