# Plano de execução — FiveStars (5tars)

Plano em fases para ir do zero ao MVP: onde começar e passos a seguir.

---

## Visão geral

| Fase | Foco | Resultado esperado |
|------|------|--------------------|
| **0** | Ambiente e decisões | Tudo instalado e definido |
| **1** | Fundação do backend | API Go + Postgres + entidades base |
| **2** | Autenticação | Login/registro e perfis |
| **3** | Estabelecimentos e check-in | Fluxo de check-in funcionando |
| **4** | Reviews e curtidas | Publicar e curtir reviews |
| **5** | Highlights e dashboard | Estabelecimento destaca reviews |
| **6** | Moedas e gamificação | Wallet, níveis, avatar |
| **7** | App mobile | Telas principais + integração com a API |
| **8** | Busca, filtros e polish | Busca por categoria/local + ajustes finais |

---

## Fase 0 — Ambiente e decisões (1–2 dias)

**Objetivo:** Ter o ambiente pronto e as últimas decisões fechadas.

- **Go:** Confirmar versão (1.21+), IDE/config ok.
- **PostgreSQL:** Instalar localmente (ou Docker); criar banco `5tars` (ou `fivestars`).
- **Repositório:** Renomear módulo para algo como `github.com/seu-usuario/5tars-api` (opcional).
- **Mobile:** Decidir stack: Expo (React Native) ou Flutter. Instalar SDK.
- **Google Maps:** Criar projeto no Google Cloud, ativar APIs (Maps, Places ou Geocoding), gerar API key (restringir por app/bundle ID depois).

**Entregável:** Postgres rodando, Go buildando, stack mobile escolhida e instalada, API key do Google em env (não no código).

---

## Fase 1 — Fundação do backend (3–5 dias)

**Objetivo:** API HTTP em Go, conexão com Postgres, estrutura de pastas e entidades principais.

- **Estrutura:** `cmd/api` (main), `internal/domain`, `internal/handler`, `internal/repository`, `internal/service`, `internal/config`, `migrations`.
- **Config:** Variáveis de ambiente: `DATABASE_URL`, `PORT`, `JWT_SECRET`. Conexão Postgres (pgx). Endpoint `GET /health` com `SELECT 1`.
- **Migrações:** golang-migrate ou scripts SQL versionados. Tabelas iniciais: `users`, `establishments`.
- **Entidades:** User, Establishment (structs em `internal/domain`).
- **Primeiro endpoint:** `GET /establishments` retornando JSON; CORS configurado.

**Entregável:** API sobe com `go run`, Postgres com tabelas, um GET que retorna JSON, CORS ok.

---

## Fase 2 — Autenticação e usuário (3–5 dias)

**Objetivo:** Registro, login e proteção de rotas com JWT.

- Tabela `users`: id, email, password_hash, name, avatar_url, level, created_at, updated_at.
- `POST /auth/register`, `POST /auth/login` (bcrypt, JWT).
- Middleware de autenticação (Bearer token → user_id no context).
- `GET /users/me` com dados do usuário logado.

**Entregável:** Registro e login; rotas protegidas; cliente obtém token e chama `GET /users/me`.

---

## Fase 3 — Estabelecimentos e check-in (4–6 dias)

**Objetivo:** CRUD básico de estabelecimentos e fluxo de check-in com validação de proximidade.

- Estabelecimentos: `GET /establishments`, `GET /establishments/:id`, `POST /establishments` (protegido).
- Tabela `checkins`: id, user_id, establishment_id, lat, lng, checked_at, created_at.
- Regra: distância (Haversine/PostGIS) ≤ raio (ex.: 200 m); opcional: 1 check-in por dia por lugar.
- `POST /checkins` com body `{ establishment_id, lat, lng }`; `GET /checkins/me` ou `GET /users/me/checkins`.

**Entregável:** Check-in validado por distância; usuário vê seus check-ins.

---

## Fase 4 — Reviews e curtidas (4–6 dias)

**Objetivo:** Review só com check-in válido (janela 5 dias); curtir reviews.

- Tabelas: `reviews` (id, user_id, establishment_id, checkin_id, rating, text, created_at; opcional photos/tags), `likes` (user_id, review_id).
- Regras: check-in existente e dentro de 5 dias; 1 review por check-in; texto com mínimo de caracteres.
- Endpoints: `POST /reviews`, `GET /establishments/:id/reviews`, `POST /reviews/:id/like`, `DELETE /reviews/:id/like`, `GET /reviews/:id`.

**Entregável:** Publicar review com check-in recente; curtidas; lista de reviews por estabelecimento.

---

## Fase 5 — Highlights e dashboard do estabelecimento (2–4 dias)

**Objetivo:** Estabelecimento destaca reviews; noção de dashboard.

- Tabela `highlights` ou many-to-many establishment ↔ review; limite (ex.: 3–5) por estabelecimento.
- `POST` para adicionar highlight; `GET /establishments/:id` inclui highlights.
- Opcional: `GET /establishments/:id/stats` (nota média, total reviews).
- Autenticação do estabelecimento: role ou tabela establishment_users.

**Entregável:** Reviews em destaque no topo da página do estabelecimento; dashboard mínimo.

---

## Fase 6 — Moedas e gamificação (4–6 dias)

**Objetivo:** Wallet, ganho por review/like/highlight, itens de avatar, níveis.

- `wallets` (user_id, balance); `wallet_transactions` (user_id, amount, type, reference_id).
- `avatar_items`, `user_avatar_items`; `POST /store/purchase`.
- Nível e medalhas (badges); `GET /users/me` com level e badges; `GET /users/me/avatar`.

**Entregável:** Moedas creditadas/debitadas; compra de itens; nível e medalhas no perfil.

---

## Fase 7 — App mobile (1–2 semanas)

**Objetivo:** App (Expo ou Flutter) com telas principais e integração com a API.

- Setup: base URL, cliente HTTP, token em SecureStore, Authorization em requisições.
- Telas: Login/Registro; lista e detalhe de estabelecimentos; check-in (localização); formulário de review; perfil.
- Opcional: mapa com estabelecimentos próximos.
- Teste ponta a ponta: registrar → login → check-in → review.

**Entregável:** App instalável que usa a API para auth, estabelecimentos, check-in e reviews.

---

## Fase 8 — Busca, filtros e polish (3–5 dias)

**Objetivo:** Busca por categoria e localização; filtros; ajustes para lançamento.

- Backend: `GET /establishments?category=...&min_rating=...&lat=...&lng=...&radius=...`; ordenação.
- App: tela de busca com filtros; pull-to-refresh; erros e estados vazios.
- Polish: mensagens amigáveis; validações; documentar API (OpenAPI/Swagger).

**Entregável:** Busca e filtros funcionando; experiência estável para primeiro teste com usuários.

---

## Ordem sugerida (resumo)

1. Fase 0 — Ambiente e decisões  
2. Fase 1 — Backend: Go, Postgres, migrações, entidades, primeiro GET  
3. Fase 2 — Auth: registro, login, JWT, GET /users/me  
4. Fase 3 — Estabelecimentos e check-in (validação de distância)  
5. Fase 4 — Reviews e likes (5 dias, 1 review por check-in)  
6. Fase 5 — Highlights e dashboard  
7. Fase 6 — Moedas, avatar, níveis/medalhas  
8. Fase 7 — App mobile  
9. Fase 8 — Busca, filtros e polish  

**Dica:** Fase 7 pode começar em paralelo após Fase 4 (app consumindo check-in e reviews), evoluindo o app enquanto se implementam Fases 5 e 6 no backend.

---

*Documento de referência do plano de execução do MVP FiveStars.*
