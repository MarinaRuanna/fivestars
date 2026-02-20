# FiveStars (5tars)

**Avaliações reais de quem realmente esteve lá.**

Plataforma mobile de avaliações de estabelecimentos baseada em check-in real: só quem fez check-in no local pode avaliar, em até 5 dias, com 1 review por check-in. Isso aumenta a credibilidade e reduz fake reviews.

---

## Requisitos

- [Go](https://go.dev/) 1.23+
- [PostgreSQL](https://www.postgresql.org/) (local ou Docker)

## Variáveis de ambiente

| Variável | Descrição |
|----------|-----------|
| `DATABASE_URL` | URL de conexão Postgres (ex.: `postgres://user:pass@localhost:5432/fivestars?sslmode=disable`) |
| `PORT` | Porta da API (opcional; padrão `8080`) |

## Executar

Crie o banco (ex.: `createdb fivestars`) e defina `DATABASE_URL`. As migrations rodam automaticamente na subida da API.

```bash
export DATABASE_URL="postgres://localhost:5432/fivestars?sslmode=disable"
go run ./cmd/api
```

- **Health:** `GET http://localhost:8080/health` → `{"status":"ok"}`
- **Estabelecimentos:** `GET http://localhost:8080/establishments` → lista em JSON (CORS habilitado)

## Build

```bash
go build -o fivestars ./cmd/api
./fivestars   # com DATABASE_URL definida
```

---

## Documentação do produto

| Documento | Descrição |
|-----------|-----------|
| [docs/PRD.md](docs/PRD.md) | Product Requirements Document — conceito, funcionalidades (check-in, review, moedas, perfil, seguir amigos, código influencer), busca, gamificação, expansões |
| [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) | Arquitetura backend, modelo de dados (User, Establishment, Checkin, Review, Like, Highlight, Wallet, Follow, Referral), regras de negócio e estrutura de pastas em Go |
| [docs/EXECUTION_PLAN.md](docs/EXECUTION_PLAN.md) | Plano de execução em fases (0–8): ambiente, backend, auth, check-in, reviews, highlights, moedas, app mobile, busca e polish |

---

## Próximos passos sugeridos

1. **Migrations** — Criar schema PostgreSQL (entidades do ARCHITECTURE).
2. **Auth** — Registro, login, JWT; opcional: OAuth (Google/Apple).
3. **Core** — CRUD Establishment, fluxo Check-in → Review (janela 5 dias), Like, Highlight.
4. **Moedas** — Wallet, transações, bônus por review/highlight/like e por código influencer.
5. **Social** — Follow (seguir amigos), feed “reviews dos amigos”.
6. **Busca** — Categorias, filtros (nota, localização, mais curtidos, mais recentes).
7. **Assets** — Upload de fotos (reviews, avatar) para armazenamento S3-compatible.

---

## Estrutura do repositório (Fase 1)

```
fivestars/
├── cmd/api/
│   ├── main.go              # entrypoint da API
│   └── migrations/          # SQL (embed; roda na subida)
├── internal/
│   ├── config/              # env (DATABASE_URL, PORT, JWT_SECRET)
│   ├── domain/              # User, Establishment
│   ├── handler/             # health, establishments, CORS
│   └── repository/          # postgres pool, EstablishmentRepository
├── migrations/              # cópia das migrations (referência / CLI)
├── docs/                    # PRD, ARCHITECTURE, EXECUTION_PLAN
├── go.mod
└── README.md
```

Evoluir como produto real: priorizar MVP (check-in + review + estabelecimento + moedas), depois rede social (seguir amigos + código influencer) e gamificação completa.
