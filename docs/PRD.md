# FiveStars (5tars) — Product Requirements Document (PRD)

**Tagline:** *"Avaliações reais de quem realmente esteve lá."*

**Versão:** MVP 1.0  
**Última atualização:** Fevereiro 2025

---

## 1. Conceito principal

O **FiveStars (5tars)** é uma plataforma mobile de avaliações de estabelecimentos baseada em **check-in real**.

### Regras de credibilidade

| Regra | Descrição |
|-------|-----------|
| **Check-in obrigatório** | Só pode fazer review quem fez check-in no local pelo menos 1 vez |
| **Janela de 5 dias** | O review deve ser enviado em até 5 dias após o check-in |
| **1 check-in = 1 review** | Cada check-in gera direito a exatamente 1 review naquele estabelecimento |

**Objetivo:** Aumentar credibilidade e reduzir fake reviews.

---

## 2. Plataforma do usuário

### 2.1 Check-in

- **Métodos de check-in:**
  - Geolocalização (proximidade ao estabelecimento)
  - QR code do estabelecimento
  - Código único do local
- **Efeito:** Gera direito a **1 review** e inicia **contagem regressiva de 5 dias**.

### 2.2 Review

- **Nota:** 1 a 5 estrelas
- **Texto:** Obrigatório, com mínimo de caracteres
- **Opcional:** Fotos, tags (ex.: atendimento, preço, ambiente, limpeza)
- **Interação:** Pode receber curtidas de outros usuários

### 2.3 Sistema de moedas

**Ganho de moedas:**

- Publicar review
- Receber HighLights (review em destaque)
- Ter review marcada como útil (curtidas)

**Uso de moedas:**

- Itens de avatar (pixel art)
- Itens exclusivos / sazonais
- Molduras especiais de perfil

### 2.4 Perfil do usuário

- Avatar customizável (pixel art)
- Nível do usuário
- Total de reviews
- Reviews em destaque (HighLights)
- Medalhas / conquistas
- **Rede social (MVP+):** lista de amigos que segue e seguidores

### 2.5 Seguir amigos (nova funcionalidade)

- Usuário pode **seguir** outros usuários e ser **seguido**.
- Feed ou seção “Reviews dos amigos” para ver avaliações de quem segue.
- Aumenta descoberta de lugares e engajamento.
- Possível ranking entre amigos (ex.: quem mais avaliou na semana).

### 2.6 Código de influencer (nova funcionalidade)

- **Código único** por usuário (ex.: `5TARS-MARIA-XY12`) compartilhável.
- Novos usuários que se cadastrarem com o código podem:
  - Gerar bônus de moedas para o indicado e/ou para o influencer (ex.: ambos ganham X moedas).
- Estabelecimentos ou campanhas podem ter **códigos de parceiro/influencer** para rastrear indicações e possíveis benefícios (ex.: cupom, selo “Indicado por X”).
- Métricas opcionais no perfil: “Indicações” ou “Amigos que entraram pelo seu código”.

---

## 3. Plataforma do estabelecimento

### 3.1 Dashboard

- Nota média
- Total de reviews
- Avaliações recentes
- **Responder** a reviews
- **Denunciar** review inadequado
- **Selecionar HighLights** (reviews positivas para destaque)

### 3.2 HighLights

- Reviews positivas com muitas curtidas
- Exibição:
  - No topo da página do estabelecimento
  - Destaque visual especial
- Autor do review pode ganhar moedas / reconhecimento

---

## 4. Sistema de busca

### Categorias

- Alimentação
- Entretenimento
- Esportes
- Cultura
- Comércio
- Serviços
- Educação

### Filtros

- Nota mínima
- Localização
- Mais curtidos
- Mais recentes

---

## 5. Identidade visual

- **Estilo:** Pixel art, interface amigável, cores vibrantes mas suaves, UI gamificada
- **Referências:** Stardew Valley, Animal Crossing: New Horizons, Habbo Hotel
- **Elementos:** Estabelecimentos como mini construções pixeladas, avatar estilo RPG 2D, ícones em 16-bit

---

## 6. Gamificação

- **Sistema de níveis** (progressão por atividade)
- **Medalhas (exemplos):**
  - Explorador — 10 check-ins
  - Crítico Profissional — 50 reviews
  - Lendário — 500 curtidas recebidas
- **Ranking semanal** (global e/ou entre amigos)

---

## 7. Diferencial competitivo

| vs. Google Reviews / TripAdvisor | FiveStars |
|-----------------------------------|-----------|
| Qualquer um pode avaliar          | Só quem esteve no local (check-in) |
| Sem prazo                         | Janela de 5 dias |
| Pouca gamificação                 | Níveis, moedas, avatar, medalhas |
| Perfil básico                     | Avatar customizável + HighLights comunitário |

---

## 8. Expansões futuras (pós-MVP)

- Modo “Missões” (ex.: visite 3 cafeterias esta semana)
- Selo “Verificado pelo 5tars”
- Integração com mapas
- Versão web
- Sistema de reserva
- Eventos locais

---

## 9. Resumo das entidades (MVP)

- **User** — usuário (com suporte a rede social e código influencer)
- **Establishment** — estabelecimento
- **Checkin** — check-in (local, usuário, data; gera direito a 1 review em 5 dias)
- **Review** — avaliação (estrelas, texto, fotos, tags, curtidas)
- **Like** — curtida em review
- **Highlight** — review em destaque (escolhida pelo estabelecimento)
- **ItemAvatar** — itens de customização
- **Wallet** — carteira de moedas do usuário
- **Follow** *(novo)* — relação “seguir” entre usuários
- **InfluencerCode** *(novo)* — código de indicação e vínculo indicador/indicado

---

*Documento vivo: atualizar conforme decisões de produto e priorização do backlog.*
