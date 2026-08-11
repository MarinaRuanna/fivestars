# Fluxo principal da API FiveStars

Esta API gira em torno de um fluxo simples e importante para o produto:

1. Registrar ou autenticar o usuario.
2. Identificar um estabelecimento.
3. Fazer check-in no local.
4. Publicar uma review vinculada ao check-in.
5. Consultar a review e interagir com ela.

## Arquivos gerados

- Collection Postman: [FiveStars.postman_collection.json](/Users/marinaruanna/dev/fivestars/docs/api/FiveStars.postman_collection.json)
- Este roteiro: [fluxo-principal.md](/Users/marinaruanna/dev/fivestars/docs/api/fluxo-principal.md)

## Variaveis usadas na collection

- `baseUrl`: por padrao `http://localhost:8080`
- `token`: JWT capturado no register/login
- `establishment_id`: preenchido ao listar ou criar estabelecimento
- `checkin_id`: preenchido ao criar check-in
- `review_id`: preenchido ao criar review
- `claim_code`: opcional para o fluxo de claim

## Ordem recomendada

### 1. Validar a API

Execute `GET /health`.

Resultado esperado:

- `200 OK`

### 2. Registrar ou logar

Execute nesta ordem:

1. `POST /auth/register`
2. `POST /auth/login`
3. `GET /users/me`

Observacoes:

- A collection salva automaticamente o `token`.
- Se o usuario ja existir, `POST /auth/register` pode retornar conflito. Nesse caso siga com `POST /auth/login`.

### 3. Selecionar um estabelecimento

Voce tem dois caminhos:

1. Usar `GET /establishments` se o ambiente ja tiver dados.
2. Usar `POST /establishments` para criar um estabelecimento de teste.

Depois disso execute:

1. `GET /establishments/:id`
2. `GET /establishments/:id/stats`

Observacoes:

- A collection tenta preencher `establishment_id` automaticamente.
- O `POST /establishments` exige JWT.

### 4. Fazer check-in

Execute:

1. `POST /checkins`
2. `GET /checkins/me`

Regras de negocio importantes:

- O check-in exige JWT.
- O usuario precisa informar `lat` e `lng`.
- O check-in depende da validacao de proximidade com o estabelecimento.
- Nao pode haver check-in repetido no mesmo dia para o mesmo usuario e estabelecimento.

### 5. Publicar a review

Execute:

1. `POST /reviews`
2. `GET /establishments/:id/reviews`
3. `GET /reviews/:id`

Regras de negocio importantes:

- A review exige JWT.
- A review precisa apontar para um `checkin_id` valido do proprio usuario.
- So pode existir 1 review por check-in.
- A review precisa ser publicada dentro da janela de 5 dias apos o check-in.

### 6. Interagir com a review

Execute:

1. `POST /reviews/:id/like`
2. `DELETE /reviews/:id/like`

Observacoes:

- Os dois endpoints exigem JWT.
- O `like` cria a interacao e o `delete` desfaz a interacao.

## Fluxo opcional do estabelecimento

Se quiser testar o lado do operador/dono do estabelecimento, a collection inclui:

1. `POST /establishments/:id/claim`
2. `POST /establishments/:id/highlights`
3. `DELETE /establishments/:id/highlights/:reviewId`

Observacoes:

- `claim` exige um `claim_code` valido.
- `highlights` exigem permissao sobre o estabelecimento.
- O highlight usa um `review_id` do mesmo estabelecimento.

## Dicas de uso

- Se quiser repetir o fluxo varias vezes, altere `email` e `establishment_slug`.
- Se o ambiente ja tiver estabelecimentos reais, prefira `GET /establishments` antes de criar novos.
- Se `POST /checkins` falhar por distancia, ajuste `lat` e `lng` para coordenadas proximas do estabelecimento salvo.
