![Logo do Quick](/Logo.png)

[![GoDoc](https://godoc.org/github.com/jeffotoni/quick?status.svg)](https://godoc.org/github.com/jeffotoni/quick) [![Github Release](https://img.shields.io/github/v/release/jeffotoni/quick?include_prereleases)](https://img.shields.io/github/v/release/jeffotoni/quick) [![CircleCI](https://dl.circleci.com/status-badge/img/gh/jeffotoni/quick/tree/master.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/jeffotoni/quick/tree/master) [![Go Report](https://goreportcard.com/badge/github.com/jeffotoni/quick)](https://goreportcard.com/badge/github.com/jeffotoni/quick) [![License](https://img.shields.io/github/license/jeffotoni/quick)](https://img.shields.io/github/license/jeffotoni/quick) ![CircleCI](https://img.shields.io/circleci/build/github/jeffotoni/quick/master) ![Coveralls](https://img.shields.io/coverallsCoverage/github/jeffotoni/quick)

# Exemplos Quick

  

### **Bem-vindo ao repositório de exemplos do Quick!**

  

O **Quick** é um gerenciador de rotas para Go bem flexível e extensível com diversas funcionalidades. O repositório do Framework Quick pode ser encontrado em [aqui](https://github.com/jeffotoni/quick).

  

Este repositório contém exemplos práticos de como utilizar a biblioteca Quick em Go, uma biblioteca de teste baseada em propriedades que permite escrever testes mais robustos e completos para sua aplicação.

  

Os exemplos apresentados aqui mostram como utilizar o Quick em diferentes tipos de testes, desde simples até mais complexos, ajudando a começar rapidamente e aprender as melhores práticas de teste.

  

O Quick é desenvolvido por **jeffotoni** e é uma excelente opção para escrever testes em Go, encontrando falhas em sua aplicação que podem não ser encontradas em testes tradicionais e aumentando a qualidade do código.

  

Sinta-se à vontade para explorar o repositório, contribuir com seus próprios exemplos e melhorias para a biblioteca Quick. Obrigado por usar Quick!

  

## Quais exemplos você encontrará no repositório?

  

* [Group](/group/)

* [Middleware](/middleware/)

* [Delete](quick.delete/)

* [Get](quick.get/)

* [Post](quick.post/)

* [Put](quick.put/)

* [Regex](quick.regex/)

* [Start](quick.start/)

  

```go

package main

  

import  "github.com/jeffotoni/quick"

  

func  main() {

q  := quick.New()

q.Get("/v1/user", func(c *quick.Ctx) error {

c.Set("Content-Type", "application/json")

return c.Status(200).SendString("Quick in action com Cors❤️!")

})

q.Listen("0.0.0.0:8080")

}
```

## 🚀 **Apoiadores do Projeto Quick** 🙏

O Projeto Quick visa desenvolver e disponibilizar softwares de qualidade para a comunidade de desenvolvedores. 💻 Para continuarmos a melhorar nossas ferramentas, contamos com o apoio de nossos patrocinadores no Patreon. 🤝

Agradecemos a todos os nossos apoiadores! 🙌 Se você também acredita em nosso trabalho e quer contribuir para o avanço da comunidade de desenvolvimento, considere apoiar o Projeto Quick em nosso Patreon [aqui](https://www.patreon.com/jeffotoni_quick)

Juntos podemos continuar a construir ferramentas incríveis! 🚀
