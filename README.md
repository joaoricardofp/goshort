# GoShort

[![Go Version](https://img.shields.io/badge/Go-1.27.1-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/license-Não%20especificada-lightgrey)](#licença)
[![Deploy on Railway](https://img.shields.io/badge/Deploy-Railway-0B0D0E?style=flat&logo=railway)](https://go-short.up.railway.app)

GoShort é uma API HTTP REST para encurtamento de URLs desenvolvida em Go (Golang). O projeto foi construído utilizando exclusivamente os recursos da biblioteca padrão do Go, com foco em idempotência semântica, segurança no acesso concorrente e validação rigorosa de entradas.

---

## Sumário

- [Sobre o Projeto](#sobre-o-projeto)
- [Funcionalidades](#funcionalidades)
- [Tecnologias Utilizadas](#tecnologias-utilizadas)
- [Arquitetura e Organização](#arquitetura-e-organização)
- [Requisitos](#requisitos)
- [Como Executar Localmente](#como-executar-localmente)
- [Endpoints da API](#endpoints-da-api)
  - [POST /shorten](#post-shorten)
  - [GET /{code}](#get-code)
- [Idempotência e Normalização](#idempotência-e-normalização)
  - [Regras de Normalização](#regras-de-normalização)
  - [Controle de Concorrência](#controle-de-concorrência)
  - [Codificação Base36](#codificação-base36)
- [Testes](#testes)
- [Deploy](#deploy)
- [Limitações e Melhorias Futuras](#limitações-e-melhorias-futuras)
- [Aprendizados](#aprendizados)
- [Licença](#licença)

---

## Sobre o Projeto

O GoShort resolve o problema de geração e redirecionamento de links compactos, priorizando a consistência dos dados e a utilização eficiente de recursos. Diferente de implementações ingênuas que geram um novo identificador a cada requisição, o GoShort implementa **idempotência real baseada em normalização semântica**: variações sintáticas de uma mesma URL (como maiúsculas no domínio, portas padrão explícitas ou âncoras/fragmentos) são convertidas para uma forma canônica antes do armazenamento.

Dessa forma, submissões equivalentes reutilizam o mesmo código curto existente, retornando o código de status HTTP adequado (`200 OK` em vez de `201 Created`).

---

## Funcionalidades

- **Encurtamento de URLs**: Gera identificadores curtos em formato alfanumérico compacto (Base36, com preenchimento mínimo de 6 caracteres).
- **Idempotência Garantida**: URLs equivalentes recebem sempre o mesmo código identificador, evitando duplicação em memória.
- **Redirecionamento HTTP**: Redirecionamento instantâneo via código `302 Found` para o endereço original cadastrado.
- **Validação Rigorosa de Entrada**:
  - Aceita apenas os esquemas `http` e `https`.
  - Rejeita hosts vazios ou malformados.
  - Valida representações de portas numéricas.
  - Suporta endereços IPv4 e IPv6.
- **Normalização de URLs**: Elimina divergências superficiais (conversão para minúsculas de scheme/host, remoção de portas padrão e descarte de fragmentos).
- **Proteção Contra Sobrecarga**: Limita o corpo das requisições a 1 MB através de `http.MaxBytesReader`.
- **Armazenamento Seguro para Concorrência**: Sincronização interna via `sync.Mutex`, assegurando consistência em cenários de requisições simultâneas.
- **Zero Dependências Externas**: Implementado integralmente com a biblioteca padrão do Go.

---

## Tecnologias Utilizadas

- **[Go](https://golang.org/) (v1.27.1)**: Linguagem de programação principal.
- **`net/http`**: Servidor HTTP e roteador nativo (`http.NewServeMux`) com suporte a métodos e parâmetros de caminho (`POST /shorten`, `GET /{code}`).
- **`sync` (`sync.Mutex`)**: Controle de concorrência e exclusão mútua no acesso às estruturas de armazenamento em memória.
- **`net/url` & `net`**: Análise sintática, validação de regras de rede e reconstrução canônica de URLs e hosts (incluindo tratamento de IPv6).
- **`encoding/json`**: Serialização e desserialização de dados das requisições e respostas da API.
- **`strconv`**: Conversão de contadores inteiros para representação em Base36 (`strconv.FormatInt`) e validação de portas numéricas.
- **`testing` & `net/http/httptest`**: Ferramentas nativas para execução de testes unitários, testes de concorrência com goroutines e simulação de requisições HTTP.

---

## Arquitetura e Organização

O projeto adota uma estrutura em pacote único (`package main`), adequada para utilitários e microsserviços enxutos:

```text
goshort/
├── go.mod             # Definição do módulo Go e versão da linguagem
├── main.go            # Ponto de entrada: leitura de porta e inicialização do servidor HTTP
├── handler.go         # Camada HTTP: definição de rotas, manipulação de requisições e respostas JSON
├── handler_test.go    # Testes dos endpoints HTTP (/shorten e /{code})
├── store.go           # Armazenamento em memória, mutex e algoritmo de geração de códigos Base36
├── store_test.go      # Testes de concorrência com goroutines e testes de codificação Base36
├── url.go             # Funções de validação e normalização canônica de URLs
└── url_test.go        # Testes de validação de regras sintáticas e normalização
```

### Responsabilidades dos Módulos

- **`main.go`**: Obtém a porta de execução através da variável de ambiente `PORT` (com fallback para `8080`), inicializa a instância de `Store`, conecta-a ao `Handler` e inicia o servidor HTTP via `http.ListenAndServe`.
- **`handler.go`**: Define a estrutura `Handler` e registra as rotas no roteador do Go (`net/http.ServeMux`). Gerencia o ciclo de vida da requisição HTTP, limitando o payload, decodificando o JSON, tratando erros e emitindo códigos de status HTTP apropriados.
- **`store.go`**: Define a estrutura `Store` com dois mapas em memória (`urlToCode` e `codeToURL`). Encapsula a lógica de geração de códigos alfanuméricos e garante exclusão mútua em operações de leitura e escrita através de um `sync.Mutex`.
- **`url.go`**: Concentra a lógica de validação (`validateURL`) e de transformação canônica (`normalizeURL`), atuando como base para a idempotência do sistema.

---

## Requisitos

- **Go**: Versão `1.27.1` ou superior instalada (conforme declarado em [`go.mod`](go.mod)).
- **Git** (opcional, para clonagem do repositório).

---

## Como Executar Localmente

### 1. Clonar o repositório

```bash
git clone https://github.com/joaoricardofp/goshort.git
cd goshort
```

### 2. Executar a aplicação

Para executar diretamente:

```bash
go run .
```

Por padrão, a aplicação escutará na porta `8080`. A saída no terminal indicará:

```text
Listening on :8080
```

### 3. Configurar uma porta personalizada (opcional)

A porta de execução pode ser customizada definindo a variável de ambiente `PORT`.

**No Linux / macOS (Bash / Zsh):**
```bash
PORT=3000 go run .
```

**No Windows (PowerShell):**
```powershell
$env:PORT="3000"; go run .
```

**No Windows (Prompt de Comando - CMD):**
```cmd
set PORT=3000 && go run .
```

---

## Endpoints da API

| Método | Rota | Descrição | Status de Sucesso |
| :--- | :--- | :--- | :--- |
| `POST` | `/shorten` | Recebe uma URL longa e retorna o código e a URL encurtada | `201 Created` / `200 OK` |
| `GET` | `/{code}` | Redireciona para o endereço original associado ao código | `302 Found` |

---

### POST /shorten

Cria um novo link encurtado ou retorna o link existente caso a URL informada já tenha sido cadastrada (mesmo sob variações semânticas). O corpo da requisição é restrito ao tamanho máximo de 1 MB.

#### Cabeçalhos da Requisição
```http
Content-Type: application/json
```

#### Corpo da Requisição
```json
{
  "url": "https://example.com/artigos/go-concorrencia"
}
```

#### Respostas

**Cenário 1: URL nova (primeira submissão)**
- **Status HTTP**: `201 Created`
- **Corpo**:
```json
{
  "code": "000001",
  "short_url": "http://localhost:8080/000001"
}
```

**Cenário 2: URL já cadastrada ou equivalente (reutilização idempotente)**
- **Status HTTP**: `200 OK`
- **Corpo**:
```json
{
  "code": "000001",
  "short_url": "http://localhost:8080/000001"
}
```

#### Respostas de Erro (`400 Bad Request`)

- **JSON malformado ou payload excedendo 1 MB**:
  ```json
  {
    "error": "Invalid JSON"
  }
  ```

- **URL com formato inválido, esquema ausente/incompatível ou porta não numérica**:
  ```json
  {
    "error": "Invalid URL"
  }
  ```

- **Falha no processamento interno da URL**:
  ```json
  {
    "error": "Internal Server Error"
  }
  ```

---

### GET /{code}

Recupera a URL de destino associada ao código informado e executa o redirecionamento.

#### Parâmetros de URL
- `code` (string): Identificador alfanumérico da URL encurtada (ex.: `000001`).

#### Respostas

**Cenário 1: Código existente**
- **Status HTTP**: `302 Found`
- **Cabeçalho**:
  ```http
  Location: https://example.com/artigos/go-concorrencia
  ```

**Cenário 2: Código não encontrado**
- **Status HTTP**: `404 Not Found`
- **Corpo**:
  ```json
  {
    "error": "Not Found"
  }
  ```

---

## Idempotência e Normalização

A idempotência no GoShort garante que múltiplas submissões de URLs que apontam para o mesmo recurso resultem no mesmo identificador, sem gerar registros duplicados.

### Regras de Normalização

Antes de qualquer consulta ou inserção no armazenamento, a função `normalizeURL` aplica as seguintes transformações:

1. **Esquema em minúsculas**: `HTTP://` e `HTTPS://` tornam-se `http://` e `https://`.
2. **Host em minúsculas**: Nomes de domínio são convertidos para letras minúsculas (ex.: `Example.COM` torna-se `example.com`).
3. **Remoção de portas padrão**:
   - `http://dominio:80/` tem a porta `:80` removida, tornando-se `http://dominio/`.
   - `https://dominio:443/` tem a porta `:443` removida, tornando-se `https://dominio/`.
   - Portas não padrão (como `:8080` ou `:8443`) são preservadas.
4. **Tratamento de endereços IPv6**: Hosts IPv6 sem porta recebem formatação entre colchetes (ex.: `[2001:db8::1]`).
5. **Caminho padrão**: Se a URL não contiver caminho explícito, é atribuída a barra raiz `/` (ex.: `http://example.com` torna-se `http://example.com/`).
6. **Remoção de fragmentos**: Âncoras de navegação iniciadas por `#` (ex.: `#secao-1`) são descartadas, pois dizem respeito apenas ao cliente/navegador e não ao recurso do servidor.
7. **Preservação de parâmetros e caminhos**: Parâmetros de consulta (`?query=valor`) e a diferenciação de maiúsculas/minúsculas no caminho após o domínio são preservados rigorosamente.

### Controle de Concorrência

Para permitir que a API seja atendida por múltiplas goroutines de forma concorrente sem corromper o estado dos dados, a estrutura `Store` sincroniza suas operações:

```go
type Store struct {
    mu        sync.Mutex
    urlToCode map[string]string
    codeToURL map[string]string
    counter   int
}
```

- Tanto a operação de inserção/consulta (`GetOrCreate`) quanto a de busca (`Lookup`) adquirem o bloqueio exclusivo (`s.mu.Lock()`) antes de manipular os mapas em memória, liberando-o via `defer s.mu.Unlock()`.
- O mapeamento bidirecional (`urlToCode` e `codeToURL`) assegura buscas eficientes tanto pelo endereço quanto pelo código encurtado com custo temporal O(1).

### Codificação Base36

O identificador curto é gerado a partir de um contador sequencial inteiro convertido para base 36 (composta por dígitos `0-9` e caracteres alfabéticos `a-z`):

- **Tamanho mínimo**: 6 caracteres. Códigos com representação inferior a 6 caracteres recebem preenchimento de zeros à esquerda (ex.: `1` torna-se `000001`, `10` torna-se `00000a`, `36` torna-se `000010`).
- Quando o contador ultrapassa o limite de 6 dígitos em base 36 (`2.176.782.335`), a função expande naturalmente para 7 dígitos (`1000000`).

---

## Testes

O projeto conta com uma suíte abrangente de testes automatizados cobrindo normalização, concorrência e integração dos manipuladores HTTP.

### Executar a suíte de testes

Para executar todos os testes do projeto:

```bash
go test ./...
```

Para executar em modo detalhado (verbose):

```bash
go test -v ./...
```

### Verificação de condições de corrida (Race Detector)

Para atestar que as rotinas de concorrência não possuem condições de corrida (*data races*):

```bash
go test -race ./...
```

### Comportamentos Cobertos

- **Validação e Normalização (`url_test.go`)**:
  - Aceitação de URLs válidas (HTTP, HTTPS, IPv4, IPv6 com e sem porta, maiúsculas no scheme).
  - Rejeição de esquemas não suportados (`ftp`, `mailto`), esquemas relativos (`//example.com`), strings vazias e hosts com espaços.
  - Rejeição de portas inválidas (não numéricas ou valores negativos).
  - Normalização canônica com remoção de portas 80/443 e descarte de fragmentos.
- **Armazenamento e Concorrência (`store_test.go`)**:
  - Teste de concorrência simulando 100 goroutines simultâneas inserindo a mesma URL para validar unicidade de mapeamento e ausência de corridas de dados.
  - Testes de limite e formatação de codificação Base36 (`encodeCode`).
- **Handlers HTTP (`handler_test.go`)**:
  - Teste de criação e idempotência: validação do status `201 Created` no primeiro envio e `200 OK` na reutilização com retorno do mesmo código.
  - Rejeição de payloads malformados ou inválidos com status `400 Bad Request`.
  - Redirecionamento funcional com status `302 Found` e conferência do cabeçalho `Location`.
  - Retorno de status `404 Not Found` para códigos inexistentes.

---

## Deploy

A aplicação está configurada para execução contínua e hospedada na plataforma **Railway**:

- **Ambiente de produção**: [https://go-short.up.railway.app](https://go-short.up.railway.app)

O servidor se adapta automaticamente a ambientes de nuvem ao ler a variável de ambiente `PORT` provida pela plataforma durante a inicialização em `main.go`.

---

## Limitações e Melhorias Futuras

O projeto foi intencionalmente concebido de forma minimalista para fins educacionais e experimentação prática. Seu estado atual reflete as seguintes características e oportunidades de evolução:

### Limitações do Estado Atual

- **Armazenamento em memória (volátil)**: Os links e mapeamentos são mantidos exclusivamente na memória RAM da aplicação. Toda reinicialização do processo resulta na perda dos dados cadastrados.
- **Escalabilidade horizontal restrita**: Por manter o estado localmente no processo, múltiplas instâncias da aplicação não compartilham a mesma base de dados.
- **Detecção de esquema atrás de proxy reverso**: A função `shortURL` determina o prefixo `https://` inspecionando diretamente `r.TLS`. Em ambientes onde a terminação TLS ocorre em um proxy reverso (como no Railway) sem conexão TLS direta com a aplicação, a URL encurtada retornada no JSON pode utilizar o prefixo `http://`.
- **Ausência de expiração (TTL)**: Não há mecanismo para expirar ou limpar URLs antigas.

### Oportunidades de Melhorias Futuras

- [ ] Implementar camada de persistência com banco de dados relacional (ex.: PostgreSQL, SQLite) ou chave-valor (ex.: Redis).
- [ ] Adicionar suporte aos cabeçalhos `X-Forwarded-Proto` e `X-Forwarded-Host` para geração precisa da URL curta quando atrás de proxies reversos.
- [ ] Adicionar expiração configurável (TTL) para URLs temporárias.
- [ ] Implementar endpoint de monitoramento de integridade da aplicação (`GET /healthz`).
- [ ] Adicionar suporte a aliases customizados definidos pelo usuário na requisição de encurtamento.

---

## Aprendizados

O desenvolvimento deste projeto proporcionou a prática de conceitos centrais da engenharia de software com Go:

- **APIs HTTP com Roteamento Moderno**: Utilização dos novos padrões de roteamento do `net/http.ServeMux` introduzidos a partir do Go 1.22, dispensando dependências de roteadores de terceiros.
- **Concorrência e Sincronização**: Emprego de `sync.Mutex` para gerenciar o acesso coordenado a estruturas compartilhadas entre múltiplas goroutines, validado com o analisador de concorrência (`go test -race`).
- **Idempotência em Serviços Web**: Tratamento semântico de requisições idempotentes e diferenciação explícita no protocolo HTTP entre a criação de novos recursos (`201 Created`) e o reaproveitamento de existentes (`200 OK`).
- **Validação e Normalização Canônica**: Aplicação de regras estritas de parsing de rede com `net/url` e `net`, contemplando esquemas de URL, portas canônicas e notação de endereçamento IPv6.
- **Testes Automatizados e Concorrentes**: Criação de testes de integração HTTP utilizando `httptest.ResponseRecorder` e testes de carga concorrente empregando `sync.WaitGroup`.
