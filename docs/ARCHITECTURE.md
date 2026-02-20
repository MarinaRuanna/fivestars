# FiveStars — Arquitetura e modelo de dados (MVP)

## Visão geral

Backend em **Go**, preparado para API REST (e futuramente GraphQL se fizer sentido).  
Mobile consome a API; estabelecimentos podem usar dashboard web.

---

## Stack sugerida (MVP)

| Camada | Tecnologia sugerida |
|--------|----------------------|
| Linguagem | Go 1.23+ |
| API | net/http ou Chi/Echo (REST) |
| Banco de dados | PostgreSQL |
| Cache / sessão | Redis (sessão, rate limit, filas leves) |
| Armazenamento de arquivos | S3-compatible (fotos de review, avatar) |
| Auth | JWT + refresh token; OAuth opcional (Google/Apple) |

---

## Modelo de dados (entidades principais)

### User

| Campo | Tipo | Descrição |
|-------|------|-----------|
| id | UUID | PK |
| email | string | único |
| password_hash | string | bcrypt/argon2 |
| display_name | string | nome exibido |
| avatar_config | JSONB | peças do avatar (pixel art) |
| level | int | nível gamificado |
| xp | int | experiência atual |
| influencer_code | string | único, ex: 5TARS-XXXX-XXXX |
| referred_by_id | UUID | FK User (quem indicou) — nullable |
| created_at | timestamp | |
| updated_at | timestamp | |

- Índices: `email` (unique), `influencer_code` (unique), `referred_by_id`.

### Establishment

| Campo | Tipo | Descrição |
|-------|------|-----------|
| id | UUID | PK |
| owner_id | UUID | FK User (dono do estabelecimento) |
| name | string | |
| slug | string | único, URL-friendly |
| category | enum/string | alimentação, entretenimento, etc. |
| qr_code_secret | string | para check-in por QR |
| location_lat | decimal | |
| location_lng | decimal | |
| address_text | string | endereço legível |
| created_at | timestamp | |
| updated_at | timestamp | |

- Índices: `slug` (unique), `owner_id`, `category`, (lat, lng) para busca por proximidade.

### Checkin

| Campo | Tipo | Descrição |
|-------|------|-----------|
| id | UUID | PK |
| user_id | UUID | FK User |
| establishment_id | UUID | FK Establishment |
| method | enum | geolocation, qr_code, manual_code |
| review_eligible_until | timestamp | checkin_at + 5 dias |
| created_at | timestamp | |

- Índice único: `(user_id, establishment_id, created_at)` ou regra de negócio: 1 review por (user, establishment) por check-in.
- Constraint: só pode existir 1 review por check-in (ver entidade Review).

### Review

| Campo | Tipo | Descrição |
|-------|------|-----------|
| id | UUID | PK |
| checkin_id | UUID | FK Checkin (único: 1 review por check-in) |
| user_id | UUID | FK User |
| establishment_id | UUID | FK Establishment |
| stars | int | 1–5 |
| body | text | texto obrigatório (mín. caracteres) |
| tags | string[] ou JSONB | atendimento, preço, ambiente, etc. |
| created_at | timestamp | |
| updated_at | timestamp | |

- Índices: `establishment_id`, `user_id`, `checkin_id` (unique).
- Regra: `created_at <= checkin.review_eligible_until`.

### Like

| Campo | Tipo | Descrição |
|-------|------|-----------|
| user_id | UUID | PK (composite) |
| review_id | UUID | PK (composite) |
| created_at | timestamp | |

- PK composta (user_id, review_id); índices em review_id e user_id.

### Highlight

| Campo | Tipo | Descrição |
|-------|------|-----------|
| id | UUID | PK |
| establishment_id | UUID | FK Establishment |
| review_id | UUID | FK Review (unique por establishment) ou ordem |
| position | int | ordem de exibição |
| created_at | timestamp | |

- Estabelecimento escolhe N reviews para destaque (ex.: máx. 5).

### ItemAvatar

| Campo | Tipo | Descrição |
|-------|------|-----------|
| id | UUID | PK |
| name | string | |
| type | string | cabelo, roupa, moldura, etc. |
| cost_coins | int | preço em moedas |
| is_seasonal | bool | |
| asset_url | string | sprite/asset pixel art |

### UserAvatarItem (relação N:N)

| Campo | Tipo | Descrição |
|-------|------|-----------|
| user_id | UUID | PK (composite) |
| item_id | UUID | PK (composite) |
| acquired_at | timestamp | |

### Wallet

| Campo | Tipo | Descrição |
|-------|------|-----------|
| user_id | UUID | PK |
| balance | int | moedas (>= 0) |
| updated_at | timestamp | |

- Transações em tabela separada (WalletTransaction) para histórico: ganho (review, highlight, like), gasto (compra ItemAvatar), bônus (código influencer).

### Follow (seguir amigos)

| Campo | Tipo | Descrição |
|-------|------|-----------|
| follower_id | UUID | PK (composite) — quem segue |
| followed_id | UUID | PK (composite) — quem é seguido |
| created_at | timestamp | |

- PK composta (follower_id, followed_id); constraint follower_id != followed_id.
- Índices: follower_id, followed_id (para feed “reviews dos que sigo”).

### InfluencerCode / Referral

| Campo | Tipo | Descrição |
|-------|------|-----------|
| id | UUID | PK |
| code | string | único (ex: 5TARS-MARIA-XY12) |
| user_id | UUID | FK User (dono do código) |
| created_at | timestamp | |

- User.referred_by_id aponta para User que usou o código no cadastro.
- Bônus de moedas: na primeira transação de “referral” (WalletTransaction type=referral_bonus).

---

## Regras de negócio (backend)

1. **Check-in:** validar localização ou QR/código; criar Checkin com `review_eligible_until = now() + 5 days`.
2. **Review:** permitir só se existir Checkin não “usado” (sem Review ligada) e `now() <= review_eligible_until`; mínimo de caracteres no body.
3. **Like:** 1 por (user, review); ao dar like, pode creditar moedas ao autor da review (regra de gamificação).
4. **Highlight:** só estabelecimento (owner) pode adicionar/remover; limite de N por estabelecimento.
5. **Moedas:** creditar em Wallet + WalletTransaction ao publicar review, receber highlight, receber like; debitar ao comprar ItemAvatar e ao dar bônus de referral.
6. **Seguidores:** não permitir follow em si mesmo; feed “reviews dos amigos” = reviews de usuários onde current_user segue.

---

## Estrutura de pastas sugerida (Go)

```
fivestars/
├── cmd/
│   └── api/
│       └── main.go          # entrypoint da API
├── internal/
│   ├── config/              # env, config
│   ├── domain/              # entidades e regras (puros)
│   ├── repository/          # persistência (Postgres)
│   ├── usecase/             # casos de uso (check-in, review, like, etc.)
│   ├── transport/           # HTTP handlers, middlewares (auth, CORS)
│   └── auth/                # JWT, hash de senha
├── pkg/                     # libs reutilizáveis (opcional)
├── docs/                    # PRD, ARCHITECTURE, OpenAPI
├── migrations/              # SQL migrations (e.g. golang-migrate)
├── go.mod
└── README.md
```

---

## Próximos passos de implementação

1. Definir e rodar migrations (PostgreSQL) para as entidades acima.
2. Implementar auth (registro, login, JWT).
3. CRUD básico: User (perfil), Establishment.
4. Fluxo Check-in → Review (com validação de janela de 5 dias).
5. Like, Highlight, Wallet + transações e bônus de referral.
6. Endpoints de Follow e feed “reviews dos amigos”.
7. Busca e filtros (categorias, localização, nota, ordem).
8. Upload de fotos (review, avatar) para S3-compatible.

---

*Este doc deve evoluir junto com o código e o PRD.*
