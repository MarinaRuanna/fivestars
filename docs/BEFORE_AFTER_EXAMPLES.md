# Before & After Code Examples

## AuthHandler Refactoring

### BEFORE (Mixed Concerns)
```go
type AuthHandler struct {
    userRepo  domain.UserRepository
    jwtSecret string
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    // 1. Parse HTTP
    var req RegisterRequest
    json.NewDecoder(r.Body).Decode(&req)
    
    // 2. Validate (business logic in handler!)
    if req.Email == "" || req.Password == "" {
        http.Error(w, "invalid input", http.StatusBadRequest)
        return
    }
    
    // 3. Check duplicate (database concern in handler!)
    existingUser, _ := h.userRepo.GetByEmail(r.Context(), req.Email)
    if existingUser != nil {
        http.Error(w, "user exists", http.StatusConflict)
        return
    }
    
    // 4. Hash password (auth concern in handler!)
    hashedPassword, _ := auth.HashPassword(req.Password)
    
    // 5. Create user (database concern in handler!)
    user := &domain.User{
        Email:    req.Email,
        Password: hashedPassword,
        Name:     req.Name,
    }
    createdUser, _ := h.userRepo.Create(r.Context(), user)
    
    // 6. Generate token (auth concern in handler!)
    token, _ := auth.GenerateJWT(createdUser.ID, h.jwtSecret)
    
    // 7. Return JSON
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "token": token,
        "user_id": createdUser.ID,
    })
}
```

**Problems:**
- ❌ Handler is too thick (200+ lines for auth)
- ❌ Business logic mixed with HTTP concerns
- ❌ Hard to test (requires DB, JWT secret, password hashing)
- ❌ Hard to reuse logic (only accessible via HTTP)
- ❌ Password hashing knowledge in handler
- ❌ Duplicate checking knowledge in handler

---

### AFTER (Clean Layers)
```go
// Layer 1: HTTP Handler - only HTTP concerns
type AuthHandler struct {
    registerUserUseCase *usecases.RegisterUserUseCase
    loginUserUseCase    *usecases.LoginUserUseCase
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    // 1. Parse HTTP request
    var req RegisterRequest
    json.NewDecoder(r.Body).Decode(&req)
    
    // 2. DELEGATE to usecase (ALL business logic)
    output, err := h.registerUserUseCase.Execute(r.Context(), usecases.RegisterUserInput{
        Email:    req.Email,
        Password: req.Password,
        Name:     req.Name,
    })
    
    // 3. Handle errors and return JSON
    if err != nil {
        // Error handling (could be validation, conflict, etc.)
        respondError(w, err)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]any{
        "token": output.Token,
        "user_id": output.UserID,
    })
}

// Layer 2: Usecase - business logic only
type RegisterUserUseCase struct {
    userRepo  domain.UserRepository  // Interface, not implementation
    jwtSecret string
}

func (uc *RegisterUserUseCase) Execute(
    ctx context.Context,
    input RegisterUserInput,
) (*RegisterUserOutput, error) {
    // 1. Validate
    if input.Email == "" || input.Password == "" || input.Name == "" {
        return nil, customerror.NewValidationError("email, password, name required")
    }
    
    if len(input.Password) < 6 {
        return nil, customerror.NewValidationError("password too short")
    }
    
    // 2. Check duplicate
    existingUser, _ := uc.userRepo.GetByEmail(ctx, input.Email)
    if existingUser != nil {
        return nil, customerror.NewConflictError("user already exists")
    }
    
    // 3. Hash password
    hashedPassword, err := auth.HashPassword(input.Password)
    if err != nil {
        return nil, err
    }
    
    // 4. Create user
    user := &domain.User{
        Email:        input.Email,
        PasswordHash: hashedPassword,
        Name:         input.Name,
    }
    err = uc.userRepo.Create(ctx, user)
    if err != nil {
        return nil, err
    }
    
    // 5. Get created user and generate token
    createdUser, _ := uc.userRepo.GetByEmail(ctx, input.Email)
    token, _ := auth.NewToken(createdUser.ID, uc.jwtSecret, 0)
    
    // 6. Return output DTO
    return &RegisterUserOutput{
        UserID: createdUser.ID,
        Token:  token,
        Name:   createdUser.Name,
    }, nil
}
```

**Benefits:**
- ✅ Handler is thin (10 lines HTTP logic)
- ✅ Usecase contains ALL business logic (reusable)
- ✅ Easy to test in isolation
- ✅ Can be used via different transports (HTTP, gRPC, CLI)
- ✅ Clear separation of concerns

---

## Dependency Injection Pattern

### BEFORE (Config-Heavy)
```go
// main.go - Config passed everywhere
func main() {
    cfg := config.Load()
    
    // repositories need config
    userRepo := repository.NewUserRepository(cfg.DatabaseURL)
    estabRepo := repository.NewEstablishmentRepository(cfg.DatabaseURL)
    
    // handlers need repositories AND config
    authHandler := controller.NewAuthHandler(userRepo, cfg.JWTSecret)
    userHandler := controller.NewUserHandler(userRepo)
    estabHandler := controller.NewEstablishmentsHandler(estabRepo)
    
    // routes need handlers
    handlers := controller.Handlers{...}
    router := controller.CreateRoutes(handlers)
    
    // server
    server := &http.Server{Addr: cfg.Port, Handler: router}
    server.ListenAndServe()
}

// No graceful shutdown
// Config scattered across handler constructors
```

---

### AFTER (Composition Root)
```go
// main.go - Simple entry point
func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    // Single line composition
    app, err := infra.BuildApp(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    // Start server
    go app.Start(ctx)
    
    // Graceful shutdown on signals
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan
    
    // Shutdown with timeout
    shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    app.Stop(shutdownCtx)
}

// runner.go - Composition Root (SINGLE POINT OF DEPENDENCY ASSEMBLY)
func BuildApp(ctx context.Context) (*App, error) {
    // 1. Config
    cfg := config.Load()
    cfg.Validate()
    
    // 2. Database
    pool := pgxpool.New(ctx, cfg.DatabaseURL)
    
    // 3. Repositories (concrete implementations)
    userRepo := repository.NewUserRepository(pool)
    estabRepo := repository.NewEstablishmentRepository(pool)
    
    // 4. Usecases (depend on repositories as INTERFACES)
    registerUserUC := usecases.NewRegisterUserUseCase(userRepo, cfg.JWTSecret)
    loginUserUC := usecases.NewLoginUserUseCase(userRepo, cfg.JWTSecret)
    getUserUC := usecases.NewGetUserUseCase(userRepo)
    listEstabUC := usecases.NewListEstablishmentsUseCase(estabRepo)
    
    // 5. Handlers (depend on usecases)
    healthHandler := controller.NewHealthHandler(pool)
    authHandler := controller.NewAuthHandler(registerUserUC, loginUserUC)
    userHandler := controller.NewUserHandler(getUserUC)
    estabHandler := controller.NewEstablishmentsHandler(listEstabUC)
    
    // 6. Routes
    handlers := inbound.Handlers{
        Health:         healthHandler,
        Auth:           authHandler,
        User:           userHandler,
        Establishments: estabHandler,
    }
    router := inbound.CreateChiRoutes(handlers)
    
    // 7. Return App
    return &App{
        Config: cfg,
        DB:     pool,
        Router: router,
    }, nil
}
```

**Benefits:**
- ✅ Single point of composition (BuildApp)
- ✅ Clear dependency order
- ✅ Easy to see entire wiring
- ✅ Easy to replace implementations (swap repositories)
- ✅ Graceful shutdown support
- ✅ Context propagation through all layers

---

## Testing Example

### BEFORE - Hard to Test
```go
func TestAuthHandlerRegister(t *testing.T) {
    // Need real database connection
    db := setupTestDB() // Complex setup
    defer db.Close()
    
    // Need to create repositories with real DB
    userRepo := repository.NewUserRepository(db)
    
    // Need JWT secret config
    jwtSecret := "test-secret"
    
    // Create handler
    handler := controller.NewAuthHandler(userRepo, jwtSecret)
    
    // Create request
    body := `{"email":"test@test.com","password":"password123","name":"Test"}`
    req := httptest.NewRequest("POST", "/auth/register", strings.NewReader(body))
    w := httptest.NewRecorder()
    
    // Call handler (tests HTTP + validation + hashing + DB + JWT all at once!)
    handler.Register(w, req)
    
    // Assertions on HTTP response only - can't test business logic separately
    assert.Equal(t, http.StatusCreated, w.Code)
}

// Problems:
// - Needs real DB connection
// - Tests HTTP + business logic together
// - Hard to test error cases
// - Slow (involves DB operations)
// - Can't isolate business logic from DB
```

---

### AFTER - Easy to Test
```go
// Test 1: Business logic in isolation (no HTTP, no DB!)
func TestRegisterUserUseCase(t *testing.T) {
    // Mock repository - super simple
    userRepoMock := &mockUserRepository{
        getByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
            return nil, nil // No user exists
        },
        createFunc: func(ctx context.Context, user *domain.User) error {
            user.ID = "123"
            return nil
        },
    }
    
    // Create usecase with mock
    uc := usecases.NewRegisterUserUseCase(userRepoMock, "secret")
    
    // Test business logic directly
    output, err := uc.Execute(context.Background(), usecases.RegisterUserInput{
        Email:    "test@test.com",
        Password: "password123",
        Name:     "Test",
    })
    
    // Assertions on business logic
    assert.NoError(t, err)
    assert.Equal(t, "123", output.UserID)
    assert.NotEmpty(t, output.Token)
}

// Test 2: HTTP handler in isolation (no DB!)
func TestAuthHandlerRegister(t *testing.T) {
    // Mock usecase
    usecaseMock := &mockRegisterUserUseCase{
        executeFunc: func(ctx context.Context, input usecases.RegisterUserInput) (*usecases.RegisterUserOutput, error) {
            return &usecases.RegisterUserOutput{
                UserID: "123",
                Token:  "token-here",
                Name:   "Test",
            }, nil
        },
    }
    
    // Create handler with mock
    handler := controller.NewAuthHandler(usecaseMock, nil)
    
    // Test HTTP handling only
    body := `{"email":"test@test.com","password":"password123","name":"Test"}`
    req := httptest.NewRequest("POST", "/auth/register", strings.NewReader(body))
    w := httptest.NewRecorder()
    
    handler.Register(w, req)
    
    // Assertions on HTTP response
    assert.Equal(t, http.StatusCreated, w.Code)
    assert.Contains(t, w.Body.String(), "token-here")
}

// Test 3: Error handling - business logic level
func TestRegisterUserUseCaseDuplicateEmail(t *testing.T) {
    // Mock with existing user
    userRepoMock := &mockUserRepository{
        getByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
            return &domain.User{ID: "existing"}, nil // User exists!
        },
    }
    
    uc := usecases.NewRegisterUserUseCase(userRepoMock, "secret")
    
    // Should return ConflictError
    _, err := uc.Execute(context.Background(), usecases.RegisterUserInput{
        Email:    "existing@test.com",
        Password: "password123",
        Name:     "Test",
    })
    
    assert.Error(t, err)
    assert.True(t, customerror.IsConflictError(err))
}

// Benefits:
// - No DB setup needed
// - Fast tests (no I/O)
// - Can test business logic separately
// - Can test error cases easily
// - Can test HTTP formatting separately
// - Tests are focused and isolated
```

---

## Summary of Improvements

| Aspect | Before | After |
|--------|--------|-------|
| **Handler Responsibility** | HTTP + validation + DB + auth | HTTP only |
| **Logic Location** | Scattered in handlers | Centralized in usecases |
| **Testability** | Difficult (needs DB + config) | Easy (mock dependencies) |
| **Reusability** | Only via HTTP | Via any interface |
| **Dependencies** | Tight coupling | Loose coupling via interfaces |
| **Deployment** | All-or-nothing | Can test logic offline |
| **Testing Speed** | Slow (involves DB) | Fast (mocks) |
