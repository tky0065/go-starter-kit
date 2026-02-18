# Tutorial: Build a Complete Blog API with create-go-starter

Step-by-step guide to building a Blog API with `create-go-starter`, from installation to deployment.

## Tutorial Structure

This tutorial is divided into 4 progressive parts:

### Part 1: Installation and Configuration (15 min)
<i class="material-icons info">developer_board</i>
**Goal**: Install the CLI, generate the project, configure it and test it

**You will learn**:
- Install `create-go-starter`
- Generate a new project
- Configure PostgreSQL and JWT
- Test the generated authentication API

[Start Part 1](01-setup.md)

---

### Part 2: Create Your First Domain (30 min)
<i class="material-icons success">architecture</i>
**Goal**: Implement the Posts domain (blog articles)

**You will learn**:
- Create a GORM entity
- Implement a business service
- Create a repository with GORM
- Apply hexagonal architecture

[Start Part 2](02-first-domain.md)

---

### Part 3: Expose the HTTP API (30 min)
<i class="material-icons warning">http</i>
**Goal**: Create CRUD endpoints for Posts

**You will learn**:
- Create an HTTP handler with Fiber
- Register routes
- Integrate with fx (DI)
- Test the API with curl

[Start Part 3](03-api-integration.md)

---

### Part 4: Testing and Deployment (20 min)
<i class="material-icons">rocket_launch</i>
**Goal**: Add tests and deploy with Docker

**You will learn**:
- Write unit tests
- Test with mocks
- Create a Docker image
- Deploy with docker-compose

[Start Part 4](04-testing-deployment.md)

---

## Prerequisites

### Required Software

- **Go 1.25+** - [Download](https://golang.org/dl/)
- **PostgreSQL** or **Docker** - For the database
- **curl** or **Postman** - To test the API
- Code editor (VS Code, GoLand, etc.)

### Recommended Knowledge

- Go basics (structs, interfaces, error handling)
- REST API concepts
- Familiarity with SQL/PostgreSQL (basic)

No need to be an expert! This tutorial explains each step in detail.

## Estimated Total Time

**95 minutes** (~1h30) to complete the entire tutorial.

You can take breaks between parts - each part has a clear checkpoint.

## Navigation

- [Part 1: Installation and Configuration](01-setup.md)
- [Part 2: Create Your First Domain](02-first-domain.md)
- [Part 3: Expose the HTTP API](03-api-integration.md)
- [Part 4: Testing and Deployment](04-testing-deployment.md)
