# FiveStars (5tars)

**Avaliações reais de quem realmente esteve lá.**

Plataforma mobile de avaliações de estabelecimentos baseada em check-in real: só quem fez check-in no local pode avaliar, em até 5 dias, com 1 review por check-in. Isso aumenta a credibilidade e reduz fake reviews.

---

## Requisitos

- [Go](https://go.dev/) 1.24+
- [PostgreSQL](https://www.postgresql.org/) (local ou Docker)

## Variáveis de ambiente

| Variável | Descrição |
|----------|-----------|
| `POSTGRES_HOST` | Host do PostgreSQL (ex.: `localhost`) |
| `POSTGRES_PORT` | Porta do PostgreSQL (ex.: `5432`) |
| `POSTGRES_USER` | Usuário do PostgreSQL |
| `POSTGRES_PASSWORD` | Senha do PostgreSQL |
| `POSTGRES_DATABASE` | Nome do banco (ex.: `fivestars`) |
| `POSTGRES_SSLMODE` | `disable` para ambiente local |
| `POSTGRES_MAXCONNS` | Tamanho máximo do pool (opcional) |
| `POSTGRES_MINCONNS` | Tamanho mínimo do pool (opcional) |
| `JWT_SECRET` | Segredo usado para assinar JWT |
| `APPPORT` | Porta HTTP da API (opcional; padrão `8080`) |

## .env local (opcional)

Para desenvolvimento local, você pode criar um arquivo `.env.local` na raiz do projeto.
Ele é carregado automaticamente no startup e não sobrescreve variáveis já definidas no ambiente.

Exemplo:

```bash
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DATABASE=fivestars
POSTGRES_SSLMODE=disable
POSTGRES_MAXCONNS=10
POSTGRES_MINCONNS=1
JWT_SECRET=change-me
APPPORT=8080
```

## Executar

Crie o banco (ex.: `createdb fivestars`), configure as variáveis e rode as migrations SQL.

```bash
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export POSTGRES_USER=postgres
export POSTGRES_PASSWORD=postgres
export POSTGRES_DATABASE=fivestars
export POSTGRES_SSLMODE=disable
export POSTGRES_MAXCONNS=10
export POSTGRES_MINCONNS=1
export JWT_SECRET="change-me"
export APPPORT=8080

# migrations (não rodam automaticamente no startup)
psql -h "$POSTGRES_HOST" -p "$POSTGRES_PORT" -U "$POSTGRES_USER" -d "$POSTGRES_DATABASE" -f migrations/000001_create_users_and_establishments.up.sql
psql -h "$POSTGRES_HOST" -p "$POSTGRES_PORT" -U "$POSTGRES_USER" -d "$POSTGRES_DATABASE" -f migrations/000002_create_checkins.up.sql

go run ./cmd/api
```

## Endpoints disponíveis (atual)

- `GET /health`
- `POST /auth/register`
- `POST /auth/login`
- `GET /users/me` (Bearer JWT)
- `GET /establishments`
- `POST /checkins` (Bearer JWT)
- `GET /checkins/me` (Bearer JWT)

Exemplo de health:

```bash
curl -s http://localhost:8080/health
```

## Build

```bash
go build -o fivestars ./cmd/api
./fivestars
```

---

## Documentação do produto

| Documento | Descrição |
|-----------|-----------|
| [docs/PRD.md](docs/PRD.md) | Product Requirements Document — conceito, funcionalidades (check-in, review, moedas, perfil, seguir amigos, código influencer), busca, gamificação, expansões |
| [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) | Arquitetura backend, modelo de dados (User, Establishment, Checkin, Review, Like, Highlight, Wallet, Follow, Referral), regras de negócio e estrutura de pastas em Go |
| [docs/EXECUTION_PLAN.md](docs/EXECUTION_PLAN.md) | Plano de execução em fases (0–8): ambiente, backend, auth, check-in, reviews, highlights, moedas, app mobile, busca e polish |

---

## Status atual do backend

Implementado:

1. Registro/login com JWT.
2. `GET /users/me`.
3. Listagem de estabelecimentos.
4. Check-in com validação de distância (100m) e bloqueio de check-in repetido no mesmo dia para usuário+estabelecimento.

Ainda não implementado no código:

1. Reviews (janela de 5 dias + 1 review por check-in).
2. Wallet/moedas/gamificação.
3. Follow/feed social.
4. Highlights.
5. Upload de fotos.

---

## Estrutura do repositório (atual)

```
fivestars/
├── cmd/api/
│   └── main.go                                  # entrypoint da API
├── internal/
│   ├── application/usecases/                    # casos de uso
│   ├── domain/                                  # entidades + contratos
│   └── infra/adapters/{inbound,outbound}/...   # HTTP + Postgres
├── migrations/                                  # SQL versionado
├── docs/                                        # PRD, arquitetura e planos
├── go.mod
└── README.md
```
