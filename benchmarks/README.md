# 📊 Benchmarking de Frameworks HTTP em Go

## 📏 O que é Benchmarking?

Benchmarking em desenvolvimento de software é o processo de medir e comparar o desempenho de um software ou sistema em relação a um padrão estabelecido. No contexto deste projeto, o benchmarking é utilizado para avaliar a performance de diferentes frameworks HTTP em Go, medindo métricas como tempo de resposta, taxa de requisições, e uso de recursos sob variadas condições de carga. Este processo ajuda a identificar quais frameworks oferecem a melhor eficiência, estabilidade e escalabilidade para desenvolvimento de aplicações web.

## 🎯 Objetivo do Benchmarking

O objetivo deste projeto de benchmarking é fornecer uma análise comparativa detalhada dos principais frameworks HTTP em Go, como Beego, Echo, Fiber, Gin, Iris, e Quick. Ao simular cenários de uso realista com as ferramentas k6 e Vegeta, buscamos entender como cada framework se comporta sob pressão e qual deles pode ser mais adequado dependendo das necessidades específicas de desempenho e arquitetura de um projeto.

## 🛠️ Ferramentas de Benchmark

### 🚀 k6

k6 é uma ferramenta de teste de desempenho para APIs e serviços web. É usado para simular tráfego e medir a resposta dos serviços sob carga.

#### Comandos k6

- **Executar um teste**: `k6 run script.js`
- **Ajustar número de usuários virtuais e duração**: `k6 run --vus 10 --duration 30s script.js`

### 🎯 Vegeta

Vegeta é uma ferramenta de ataque HTTP multifacetada, utilizada para realizar testes de carga e medir o desempenho dos endpoints.

### Comandos Vegeta

##### 🚧 Em breve! Estamos trabalhando nisso!

---

## 🗂️ Estrutura de Diretórios

```plaintext
/benchmarks
├── /servers
│   ├── /beego
│   ├── /echo
│   ├── /fiber
│   ├── /gin
│   └── /iris
└── /k6
    ├── post.file.list.payload.js
    └── post.js
```
---
## ▶️ Executando os Testes
Para executar qualquer um dos servidores, navegue até o diretório específico do servidor e, em seguida, entre na subpasta correspondente à funcionalidade desejada para executar o arquivo `main.go`.

### 🔹 Serviço Beego
```bash
cd servers/beego/post.simples/
go run main.go
```
### 🔹 Serviço Fiber
```bash
cd servers/fiber/post.simples/[subpasta]
go run main.go
```
### 🔹 Serviço Echo
```bash
cd servers/echo/post.simples/[subpasta]
go run main.go
```
### 🔹 Serviço Gin
```bash
cd servers/gin/post.simples/[subpasta]
go run main.go
```
### 🔹 Serviço Iris
```bash
cd servers/iris/post.simples/[subpasta]
go run main.go
```
### 🔹 Serviço Quick
```bash
cd servers/quick/post.simples/[subpasta]
go run main.go
```

### 🧪 Testando com k6
```bash
cd k6/
k6 run post.js
```