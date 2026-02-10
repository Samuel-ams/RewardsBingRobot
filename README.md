# Automação Rewards Bing - MS Edge

Este projeto foi desenvolvido em **Go** exclusivamente para Windows e utiliza as seguintes bibliotecas e ferramentas:

## Dependências
- [Go](https://golang.org/) — linguagem de programação principal  
- [GoCV](https://gocv.io/) — binding para [OpenCV](https://opencv.org/)  
- [RobotGo](https://github.com/go-vgo/robotgo) — automação de mouse, teclado e tela  
- [NSIS](https://nsis.sourceforge.io/) — utilizado para gerar o instalador do build (não incluso no repositório, está no `.gitignore`)  
- [windres](https://www.mingw.org/wiki/Windows_Resource_Compiler) — usado para gerar ícones (`.syso`) para o executável  

## Configuração do Ambiente
1. Instale Go (versão recomendada: 1.20+).  
2. Configure o `GOPATH` e certifique-se de que o `go` está disponível no terminal.  
3. Instale as dependências do projeto:  
   ```bash
   go mod tidy
## Comandos Principais  
Para compilar e executar o projeto, siga os passos abaixo:
1. Compilar e executar o projeto
    ```bash
    go run ./cmd/
---
Para gerar build deste projeto, utilize os comandos:
1.  Gerar o arquivo .syso (ícone do executável):
    ```bash
    go run tools.go
2. Gerar o arquivo executável
    ```bash
    go build -ldflags="-s -w -H=windowsgui" -o rewardsAutomation.exe ./cmd/

