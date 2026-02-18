# Hexagonal Architecture

<div class="navigation">
  <a href="index.md"><i class="material-icons">arrow_back</i> Guide Index</a>
</div>

---

## Hexagonal Architecture (Ports & Adapters)

Generated projects follow hexagonal architecture, also known as "Ports and Adapters".

**Fundamental principle**: The business domain (business logic) is at the center and depends on nothing. All dependencies point toward the domain.

```
┌─────────────────────────────────────────────────────────┐
│                   HTTP Layer (Fiber)                    │
│              adapters/handlers + middleware             │
│  • AuthHandler (register, login, refresh)               │
│  • UserHandler (CRUD operations)                        │
│  • AuthMiddleware (JWT verification)                    │
│  • ErrorHandler (centralized error handling)            │
└───────────────────────┬─────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────────┐
│              Shared Entities Layer                      │
│                   models/                               │
│  • User (entity with GORM tags)                         │
│  • RefreshToken (entity with GORM tags)                 │
│  • AuthResponse (DTO)                                   │
└──────────┬───────────────────────────┬──────────────────┘
           │                           │
           ▼                           ▼
┌──────────────────────┐  ┌──────────────────────────────┐
│   Interfaces Layer   │  │      Domain Layer            │
│   interfaces/        │  │      domain/user             │
│  • UserRepository    │  │  • UserService (logic)       │
│    (port)            │  │  • Business rules            │
└──────────┬───────────┘  └──────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────────────────────┐
│              Infrastructure Layer                       │
│         database + repository + server                  │
│  • GORM Database Connection                             │
│  • UserRepository (GORM implementation)                 │
│  • Fiber Server Configuration                           │
└─────────────────────────────────────────────────────────┘
```

### Complete Architecture Diagram (Mermaid)

The following diagram shows the complete hexagonal architecture with all components and their interactions:

```mermaid
flowchart TB
    subgraph External["Monde Externe"]
        Client["Client HTTP<br/>(Web, Mobile, API)"]
        DB[("PostgreSQL<br/>Database")]
    end

    subgraph Adapters["Adapters Layer"]
        direction TB
        subgraph Inbound["Inbound Adapters (Entree)"]
            Handlers["Handlers<br/>AuthHandler<br/>UserHandler"]
            Middleware["Middleware<br/>AuthMiddleware<br/>ErrorHandler"]
        end
        subgraph Outbound["Outbound Adapters (Sortie)"]
            RepoImpl["Repository GORM<br/>UserRepository"]
        end
    end

    subgraph Core["Core Business (Hexagone)"]
        direction TB
        Models["Models (Entites)<br/>User<br/>RefreshToken"]
        Domain["Domain Services<br/>UserService<br/>Business Logic"]
        Interfaces["Interfaces/Ports<br/>UserRepository<br/>UserService"]
        Errors["Domain Errors<br/>NotFound<br/>Validation<br/>Conflict"]
    end

    subgraph Infrastructure["Infrastructure Layer"]
        Server["Fiber Server<br/>Routes et Config"]
        DBConn["Database Connection<br/>GORM Setup"]
        Config["Configuration<br/>Environment vars"]
    end

    subgraph Packages["Packages Reutilisables (pkg/)"]
        Auth["Auth Package<br/>JWT Generation<br/>Token Parsing"]
        Logger["Logger Package<br/>Zerolog Config"]
        ConfigPkg["Config Package<br/>Env Loading"]
    end

    Client -->|"HTTP Request"| Server
    Server -->|"Route"| Handlers
    Handlers --> Middleware
    Handlers -->|"Appelle"| Domain
    Domain -->|"Utilise"| Interfaces
    Domain -->|"Utilise"| Models
    Domain -->|"Retourne"| Errors
    RepoImpl -.->|"Implemente"| Interfaces
    RepoImpl -->|"Utilise"| Models
    RepoImpl -->|"Query"| DBConn
    DBConn -->|"SQL"| DB
    Handlers -->|"Utilise"| Auth
    Server -->|"Utilise"| Config
    Domain -->|"Utilise"| Logger
```

### HTTP Request Flow (Sequence Diagram)

This diagram shows the complete path of an HTTP request through the architecture:

```mermaid
sequenceDiagram
    autonumber
    participant C as Client
    participant S as Server (Fiber)
    participant M as Middleware
    participant H as Handler
    participant SVC as Service
    participant P as Port (Interface)
    participant R as Repository
    participant DB as Database

    C->>S: POST /api/v1/auth/register
    S->>M: Route vers Handler
    M->>M: Validation (si protege)
    M->>H: Requete validee
    
    rect rgb(240, 248, 255)
        Note over H: Handler Layer
        H->>H: Parse JSON Body
        H->>H: Validate Input (validator)
    end
    
    H->>SVC: service.Register(email, password)
    
    rect rgb(255, 250, 240)
        Note over SVC: Domain Layer
        SVC->>SVC: Hash Password (bcrypt)
        SVC->>SVC: Business Validation
    end
    
    SVC->>P: repo.Create(user)
    P->>R: Appel implementation
    
    rect rgb(240, 255, 240)
        Note over R: Repository Layer
        R->>DB: INSERT INTO users...
        DB-->>R: User cree (ID)
    end
    
    R-->>SVC: User entity
    SVC-->>H: User + nil error
    H->>H: Generate JWT tokens
    H-->>C: HTTP 201 + JSON Response
```

### Dependency Inversion Principle

The core of hexagonal architecture relies on **Dependency Inversion**:

```mermaid
flowchart LR
    subgraph Traditional["Approche Traditionnelle"]
        direction TB
        T_Handler["Handler"] --> T_Service["Service"]
        T_Service --> T_Repo["Repository"]
        T_Repo --> T_DB["Database"]
    end

    subgraph Hexagonal["Architecture Hexagonale"]
        direction TB
        H_Handler["Handler"]
        H_Service["Service"]
        H_Interface["Interface<br/>(Port)"]
        H_Repo["Repository<br/>(Adapter)"]
        H_DB["Database"]
        
        H_Handler --> H_Service
        H_Service --> H_Interface
        H_Repo -.->|"implemente"| H_Interface
        H_Repo --> H_DB
    end
```

**Advantages of this approach**:

| Aspect | Without Hexagonal | With Hexagonal |
|--------|-------------------|----------------|
| **Testability** | Difficult (depends on DB) | Easy (mock interfaces) |
| **Database change** | Modifications everywhere | Only the repository |
| **Framework change** | Complete refactoring | Only the handlers |
| **Business logic** | Scattered | Centralized in the domain |

### File Structure and Responsibilities

```mermaid
flowchart TD
    subgraph CMD["cmd/"]
        Main["main.go<br/>Bootstrap fx.New()"]
    end

    subgraph Internal["internal/"]
        subgraph Models["models/"]
            User["user.go<br/>Entites GORM"]
        end
        
        subgraph Domain["domain/"]
            DErrors["errors.go<br/>Erreurs metier"]
            subgraph UserDomain["user/"]
                Service["service.go<br/>Logique metier"]
                Module["module.go<br/>fx.Module"]
            end
        end
        
        subgraph InterfacesPkg["interfaces/"]
            Repos["*_repository.go<br/>Ports (abstractions)"]
        end
        
        subgraph AdaptersPkg["adapters/"]
            subgraph HandlersPkg["handlers/"]
                AuthH["auth_handler.go"]
                UserH["user_handler.go"]
            end
            subgraph HttpPkg["http/"]
                Health["health.go"]
                Routes["routes.go<br/>Routes centralisees"]
            end
            subgraph MiddlewarePkg["middleware/"]
                AuthM["auth_middleware.go"]
                ErrorM["error_handler.go"]
            end
            subgraph RepoPkg["repository/"]
                UserRepo["user_repository.go<br/>Implementation GORM"]
            end
        end
        
        subgraph Infra["infrastructure/"]
            DBPkg["database/<br/>Connexion GORM"]
            ServerPkg["server/<br/>Config Fiber"]
        end
    end

    subgraph Pkg["pkg/"]
        AuthPkg["auth/<br/>JWT utilities"]
        ConfigPkg2["config/<br/>Env loading"]
        LoggerPkg["logger/<br/>Zerolog setup"]
    end

    Main --> Domain
    Main --> Infra
    Main --> Pkg
    HandlersPkg --> Domain
    HttpPkg --> HandlersPkg
    Domain --> InterfacesPkg
    RepoPkg -.-> InterfacesPkg
    RepoPkg --> Models
    Domain --> Models
```

**Data flow**:

1. **HTTP Request** → Handler (adapters/handlers)
2. **Handler** → Calls the Service via the interface (domain)
3. **Service** → Executes the business logic, calls the Repository via the interface
4. **Repository** → Persists to the DB (infrastructure)
5. **Return** → Flows back up to the Handler which returns the HTTP response

**Advantages**:

- **Testability**: The domain can be tested without DB or HTTP
- **Flexibility**: Easy to change DB (PostgreSQL → MySQL) or framework (Fiber → Gin)
- **Maintainability**: Clear separation of responsibilities
- **Scalability**: Add new features without breaking existing ones


---

## Navigation

**Next**: [Project Structure](structure.md)  
**Index**: [Guide Index](index.md)
