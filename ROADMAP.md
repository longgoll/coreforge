# 🗺️ COREFORGE — KẾ HOẠCH PHÁT TRIỂN TOÀN DIỆN

> **Phiên bản:** 2.0  
> **Cập nhật lần cuối:** 27/02/2026  
> **Tác giả:** longgoll  
> **Triết lý:** _"Copy, don't install" — Code được sao chép trực tiếp vào project. Dev sở hữu 100%, không vendor lock-in._

---

## 📌 MỤC LỤC

1. [Tổng Quan Dự Án](#1-tổng-quan-dự-án)
2. [Trạng Thái Hiện Tại](#2-trạng-thái-hiện-tại)
3. [Phân Tích Đối Thủ — ServerCN](#3-phân-tích-đối-thủ--servercn)
4. [Kiến Trúc Mục Tiêu — Registry 4 Tầng](#4-kiến-trúc-mục-tiêu--registry-4-tầng)
5. [Phase 1 — Nền Tảng Vững Chắc (Tuần 1–3)](#5-phase-1--nền-tảng-vững-chắc-tuần-13)
6. [Phase 2 — Mở Rộng Component (Tuần 4–8)](#6-phase-2--mở-rộng-component-tuần-48)
7. [Phase 3 — Blueprints & DX Nâng Cao (Tuần 9–14)](#7-phase-3--blueprints--dx-nâng-cao-tuần-914)
8. [Phase 4 — Website & Ecosystem (Tuần 15–20)](#8-phase-4--website--ecosystem-tuần-1520)
9. [Phase 5 — Đẳng Cấp Thượng Thừa (Tuần 21+)](#9-phase-5--đẳng-cấp-thượng-thừa-tuần-21)
10. [Catalog Component Đầy Đủ](#10-catalog-component-đầy-đủ)
11. [Chi Tiết Kỹ Thuật](#11-chi-tiết-kỹ-thuật)
12. [Chiến Lược Testing](#12-chiến-lược-testing)
13. [Chiến Lược Community & Open Source](#13-chiến-lược-community--open-source)

---

## 1. Tổng Quan Dự Án

**CoreForge** là CLI tool viết bằng Go, giúp Backend Developer thêm các component sẵn dùng (error handler, JWT auth, logger...) vào project bằng **một câu lệnh**, hỗ trợ **3 ngôn ngữ / framework** cùng lúc:

| Stack | Language | Framework | Convention |
|-------|----------|-----------|------------|
| `nodejs_express` | Node.js | Express | camelCase, `/middlewares/`, `/services/` |
| `csharp_dotnet-webapi` | C# | .NET Web API | PascalCase, `/Middlewares/`, `/Services/` |
| `golang_gin` | Golang | Gin | snake_case, `/middlewares/`, `/services/` |

### Lợi thế cạnh tranh cốt lõi

| Feature | CoreForge | ServerCN | Yeoman | Plop |
|---------|-----------|----------|--------|------|
| Multi-language (Node, C#, Go) | ✅ | ❌ (chỉ Node) | ❌ | ❌ |
| Single binary (Go) | ✅ | ❌ (npm) | ❌ | ❌ |
| No vendor lock-in | ✅ | ✅ | ❌ | ✅ |
| Backend-focused component registry | ✅ | ✅ | ❌ | ❌ |
| Interactive TUI (charmbracelet/huh) | ✅ | ❌ | ❌ | ❌ |
| Remote registry + caching | ✅ | ✅ | ✅ | ❌ |
| Project template scaffolding | ✅ | ✅ | ✅ | ❌ |

---

## 2. Trạng Thái Hiện Tại

### ✅ Đã hoàn thành

**CLI Commands (6/6):**
- [x] `forge init` — Setup project, chọn language/framework, tạo `.forge.json`, download template `.zip`
- [x] `forge add <component>` — Download component, cài dependencies, ghi file, hiện post-install guide
- [x] `forge remove <component>` — Xóa file + cập nhật config, confirm prompt, `--force` flag
- [x] `forge list [--installed]` — Liệt kê components từ registry hoặc đã cài
- [x] `forge search <keyword>` — Tìm kiếm theo name, description, tags, category
- [x] `forge doctor` — Check environment (node, npm, dotnet, go, git)

**Infrastructure:**
- [x] Go Template Engine — `processTemplate()` inject biến vào code
- [x] Remote Registry — HTTP fetch từ GitHub raw content + `--remote` flag
- [x] Caching — `~/.forge/cache/` + ETag/Last-Modified + TTL 1h + offline fallback
- [x] Cross-compile — `build.ps1` + `Makefile` cho 6 platforms (Win/Linux/macOS × amd64/arm64)
- [x] CI/CD — GitHub Actions: `ci.yml` (test 3 OS) + `release.yml` (auto-release on tag)
- [x] Project Templates — 3 bộ `.zip` skeleton (Express, .NET Web API, Gin)

**Components (6 components × 3 stacks = 18 implementations):**
- [x] `error-handler` — Global error handling middleware
- [x] `logger` — Request/response structured logging
- [x] `jwt-auth` — JWT authentication middleware + token service
- [x] `cors` — CORS configuration middleware
- [x] `rate-limiter` — API abuse prevention (rate limiting)
- [x] `validation` — Request validation (Joi / FluentValidation / go-playground/validator)

### 🔲 Chưa làm (→ Roadmap bên dưới)

- [x] Phân tầng Registry (Foundations / Schemas / Blueprints)
- [ ] 20+ component mới
- [x] Database-aware init flow
- [x] Conflict resolution thông minh
- [x] `.env` auto-configuration  
- [x] `forge update` command
- [ ] Website registry catalog
- [ ] AST Code Injection
- [ ] Docker Compose auto-integration
- [ ] Community contribution system

---

## 3. Phân Tích Đối Thủ — ServerCN

### Điểm CoreForge MẠNH hơn ServerCN
- **3 ngôn ngữ** vs ServerCN chỉ có Node.js → "Killer feature" cho đa team
- **Go single binary** vs ServerCN phải cài qua npm/npx → Cold boot = 0, không cần runtime
- **TUI đẹp** (charmbracelet/huh) vs ServerCN dùng CLI prompts tiêu chuẩn

### Điểm ServerCN MẠNH hơn CoreForge (Cần học hỏi)

| Tính năng ServerCN | Mô tả | Mức độ ưu tiên |
|---------------------|--------|----------------|
| **Architecture-Aware** | Code sinh ra linh hoạt theo MVC hoặc Feature-Based | 🔴 P0 — Quan trọng |
| **Database-Aware** | Tự sinh code kết nối MongoDB/Postgres/MySQL phù hợp ORM | 🔴 P0 — Quan trọng |
| **Non-Destructive Updates** | Không ghi đè code user đã sửa — hỏi merge/skip/rename | 🟡 P1 — Cần thiết |
| **Auto .env Config** | Tự động thêm biến môi trường vào `.env` / `.env.example` | 🟡 P1 — Cần thiết |
| **Landing Page** | Website showcase component với live code preview | 🟢 P2 — Nâng tầm |

---

## 4. Kiến Trúc Mục Tiêu — Registry 4 Tầng

Hiện tại `manifest.json` gộp tất cả vào `components`. Mục tiêu là phân tầng rõ ràng:

```
                    ┌─────────────────────────────────────┐
                    │          COREFORGE REGISTRY          │
                    └─────────────────────────────────────┘
                                     │
              ┌──────────┬───────────┼───────────┬──────────┐
              ▼          ▼           ▼           ▼          ▼
        ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐
        │FOUNDATIONS│ │COMPONENTS│ │ SCHEMAS  │ │BLUEPRINTS│
        │ (Tầng 1) │ │ (Tầng 2) │ │ (Tầng 3) │ │ (Tầng 4) │
        └──────────┘ └──────────┘ └──────────┘ └──────────┘
             │            │            │            │
     Boilerplate     Lego blocks   DB Models    Workflows
     project setup   đơn lẻ        ORM-ready    tổ hợp
```

### Tầng 1: Foundations (Nền móng dự án)
> Dùng với `forge init` — Tải bộ skeleton hoàn chỉnh

| Foundation ID | Mô tả | Bao gồm |
|---------------|--------|---------|
| `express-mvc` | Express.js theo kiến trúc MVC | `app.js`, folder structure, `tsconfig`, `.gitignore`, `Dockerfile`, linter |
| `express-feature-based` | Express.js theo Feature-Based | Mỗi feature 1 folder chứa đủ route/controller/service |
| `dotnet-clean-architecture` | .NET Clean Architecture | `Program.cs`, layers (API/Application/Domain/Infra), `.csproj` |
| `dotnet-minimal-api` | .NET Minimal API | Lightweight, single-file endpoints |
| `gin-standard` | Gin standard layout | `main.go`, `/cmd`, `/internal`, `/pkg` theo Go project layout |
| `gin-feature-based` | Gin Feature-Based | Mỗi feature 1 package riêng |

**Triển khai:** Mỗi foundation là 1 file `.zip` trên registry. `forge init` hỏi user chọn foundation → download → extract.

### Tầng 2: Components (Thành phần đơn lẻ — đang có)
> Dùng với `forge add` — Các khối Lego nhỏ nhất, độc lập

Đây là tầng hiện tại đang hoạt động. Giữ nguyên và mở rộng thêm (xem [Catalog](#10-catalog-component-đầy-đủ)).

### Tầng 3: Schemas (Mô hình dữ liệu)
> Dùng với `forge add schema-<name>` — Đóng gói sẵn DB model

| Schema ID | Mô tả | Node.js (ORM) | C# (ORM) | Go (ORM) |
|-----------|--------|---------------|----------|----------|
| `schema-user-auth` | User model cho auth | Prisma / Mongoose | EF Core | GORM |
| `schema-product` | Product + Category | Prisma / Mongoose | EF Core | GORM |
| `schema-blog-post` | Blog post + comments | Prisma / Mongoose | EF Core | GORM |
| `schema-order` | E-commerce order | Prisma / Mongoose | EF Core | GORM |

**Triển khai:** Thêm key `schemas` vào `manifest.json` v2. Mỗi schema bao gồm:
- File model/entity
- File migration (nếu ORM hỗ trợ)
- Seed data mẫu

### Tầng 4: Blueprints (Luồng nghiệp vụ tổ hợp)
> Dùng với `forge add blueprint-<name>` — Kéo về cả bộ component + schema + router

| Blueprint ID | Kéo về | Kết quả |
|--------------|--------|---------|
| `blueprint-auth` | `jwt-auth` + `password-hashing` + `schema-user-auth` + auth controller + auth routes | Luồng Đăng Nhập / Đăng Ký hoạt động ngay |
| `blueprint-crud-api` | `validation` + `error-handler` + `response-formatter` + CRUD controller template | REST CRUD endpoint scaffold |
| `blueprint-file-upload` | `file-upload` + upload controller + S3/local config | Upload ảnh/tài liệu end-to-end |

**Triển khai:** Thêm key `blueprints` vào manifest. Mỗi blueprint tham chiếu tới danh sách components + schemas cần kéo. CLI tự động resolve dependencies.

---

## 5. Phase 1 — Nền Tảng Vững Chắc (Tuần 1–3)

> **Mục tiêu:** Refactor core, chuẩn bị cho mở rộng. Không thêm feature mới, chỉ vững nền.

### Task 1.1: Refactor manifest.json Schema v2
**Effort:** 🟡 Medium (2-3 ngày)

**Thay đổi schema:**
```jsonc
// manifest.json v2
{
  "schemaVersion": "2.0.0",         // ← MỚI: versioning schema
  "version": "1.1.0",
  "languages": { /* giữ nguyên */ },
  
  "foundations": {                   // ← MỚI: Tầng 1
    "express-mvc": {
      "description": "Express.js MVC Architecture",
      "language": "nodejs",
      "framework": "express",
      "templateUrl": "...",
      "includes": ["tsconfig", "eslint", "dockerfile", "gitignore"]
    }
  },
  
  "components": { /* giữ nguyên cấu trúc hiện tại */ },
  
  "schemas": {                       // ← MỚI: Tầng 3
    "schema-user-auth": {
      "description": "User model for authentication",
      "category": "schema",
      "tags": ["user", "auth", "database"],
      "implementations": { /* cùng format với components */ }
    }
  },
  
  "blueprints": {                    // ← MỚI: Tầng 4
    "blueprint-auth": {
      "description": "Complete auth flow (register + login + JWT)",
      "includes": ["jwt-auth", "password-hashing", "schema-user-auth"],
      "files": { /* files riêng của blueprint: controller, routes */ },
      "implementations": { /* cùng format */ }
    }
  }
}
```

**File cần sửa:**
- `internal/registry/manifest.go` — Thêm struct `Foundation`, `Schema`, `Blueprint`
- `cmd/add.go` — Resolve từ cả 3 tầng (components, schemas, blueprints)
- `cmd/list.go` — Hiển thị theo tầng (grouped output)
- `cmd/search.go` — Tìm kiếm cross-tầng

**Backward compatibility:**
- Nếu `schemaVersion` không tồn tại hoặc < 2.0 → đọc theo format v1 (chỉ có `components`)
- CLI in warning: _"Registry đang dùng schema v1. Một số tính năng mới không khả dụng."_

### Task 1.2: Database-Aware Init Flow
**Effort:** 🟡 Medium (2 ngày)

Mở rộng `forge init` để hỏi thêm:

```
┌─────────────────────────────────────────────────────┐
│  ⚡ CoreForge — Project Setup                       │
│                                                     │
│  1. Language:   [Node.js ▼]                         │
│  2. Framework:  [Express ▼]                         │
│  3. Database:   [PostgreSQL ▼]          ← MỚI      │
│  4. ORM:        [Prisma ▼]              ← MỚI      │
│  5. Source Dir:  ./src                              │
│  6. Architecture: [MVC ▼]              ← MỚI      │
└─────────────────────────────────────────────────────┘
```

**Lưu vào `.forge.json` v2:**
```json
{
  "language": "nodejs",
  "framework": "express",
  "database": "postgresql",       // ← MỚI
  "orm": "prisma",                // ← MỚI
  "architecture": "mvc",          // ← MỚI
  "sourceDir": "./src",
  "installedComponents": []
}
```

**File cần sửa:**
- `internal/config/config.go` — Thêm fields mới vào `ForgeConfig`
- `internal/tui/prompts.go` — Thêm prompts cho Database, ORM, Architecture
- `cmd/init.go` — Truyền thêm options mới

### Task 1.3: Cải thiện `.forge.json` v2 — Backward Compatible
**Effort:** 🟢 Small (1 ngày)

```go
// config.go — đọc được cả v1 và v2
type ForgeConfig struct {
    // V1 fields (giữ nguyên)
    Language            string               `json:"language"`
    Framework           string               `json:"framework"`
    SourceDir           string               `json:"sourceDir"`
    InstalledComponents []InstalledComponent  `json:"installedComponents"`
    
    // V2 fields (mới, optional)
    ConfigVersion       string `json:"configVersion,omitempty"`  // "2.0"
    Database            string `json:"database,omitempty"`       // "postgresql", "mongodb", "mysql"
    ORM                 string `json:"orm,omitempty"`            // "prisma", "mongoose", "efcore", "gorm"
    Architecture        string `json:"architecture,omitempty"`   // "mvc", "feature-based", "clean"
}
```

### Task 1.4: Unit Tests cho Core
**Effort:** 🟡 Medium (2-3 ngày)

Viết test cho những phần quan trọng nhất:

```
internal/
├── config/
│   ├── config.go
│   └── config_test.go          ← Test Load/Save/HasComponent/Migration v1→v2
├── registry/
│   ├── manifest.go
│   └── manifest_test.go        ← Test LoadManifest/ResolveComponent/SchemaParsing
└── tui/
    └── prompts.go              (TUI test riêng biệt — skip ở phase này)
```

**Checklist Phase 1:**
- [x] manifest.json v2 schema + backward compat reader
- [x] Database/ORM/Architecture prompts trong `forge init`
- [x] `.forge.json` v2 struct + auto-migration
- [x] Unit tests cho config + registry (27 tests — all pass)
- [x] Verify build + go vet xanh sau refactor

---

## 6. Phase 2 — Mở Rộng Component (Tuần 4–8)

> **Mục tiêu:** Tăng từ 6 → 18 components (×3 stacks = 54 implementations)

### Ưu tiên triển khai (theo nhóm giá trị)

#### ✅ Đợt 2A (Tuần 4-5): Components nền tảng — 4 components mới — HOÀN THÀNH

| # | Component | Mô tả | Node.js | C# | Go | Effort |
|---|-----------|--------|---------|----|----|--------|
| 1 | `env-config` | Validate biến môi trường lúc startup, crash sớm nếu thiếu | dotenv + Joi | built-in IConfiguration | godotenv + validator | ✅ Done |
| 2 | `response-formatter` | JSON response chuẩn `{success, data, message, errors}` | Utility | ActionResult wrapper | Gin response helper | ✅ Done |
| 3 | `not-found-handler` | Catch unknown routes → 404 formatted response | Middleware | Middleware | Gin NoRoute handler | ✅ Done |
| 4 | `health-check` | Endpoint `/health` + `/ready` cho K8s/ELB monitoring | Route handler | Minimal API endpoint | Gin handler | ✅ Done |

**Tổng effort đợt 2A:** ✅ Hoàn thành (4 components × 3 stacks = 12 files)

#### ✅ Đợt 2B (Tuần 6-7): Components bảo mật — 4 components mới — HOÀN THÀNH

| # | Component | Mô tả | Node.js | C# | Go | Effort |
|---|-----------|--------|---------|----|----|--------|
| 5 | `password-hashing` | Hash/verify password (Bcrypt/Argon2id) | bcryptjs | BCrypt.Net | golang.org/x/crypto/bcrypt | ✅ Done |
| 6 | `security-headers` | Anti-XSS, Clickjacking, HSTS, CSP | helmet | Custom middleware | Custom middleware | ✅ Done |
| 7 | `rbac` | Role-Based Access Control middleware | Middleware | Authorize attribute + custom | Gin middleware | ✅ Done |
| 8 | `generate-otp` | OTP generation + time-limited verification | crypto + custom | custom | crypto/rand + custom | ✅ Done |

**Tổng effort đợt 2B:** ✅ Hoàn thành (4 components × 3 stacks = 12 files)

#### ✅ Đợt 2C (Tuần 8): Components utility — 4 components mới — HOÀN THÀNH

| # | Component | Mô tả | Node.js | C# | Go | Effort |
|---|-----------|--------|---------|----|----|--------|
| 9 | `shutdown-handler` | Graceful shutdown — đóng DB connections, hoàn thành requests | process.on SIGTERM | IHostedService | os.Signal + context | ✅ Done |
| 10 | `async-handler` | Wrap async functions tự động bắt lỗi (Express-specific) | Wrapper function | Không cần (C# có async/await) | Không cần (Gin recovery) | ✅ Done |
| 11 | `swagger-docs` | Auto-generate API documentation | swagger-jsdoc + swagger-ui | Swashbuckle | swag (swaggo) | ✅ Done |
| 12 | `pagination` | Pagination utility với format chuẩn | Utility | Extension method | Utility | ✅ Done |

**Tổng effort đợt 2C:** ✅ Hoàn thành (4 components × 3 stacks = 12 files)

### Quy trình thêm mỗi component

```
Bước 1: Viết code cho 3 stacks
  └─ mock-registry/components/nodejs_express/<component>/
  └─ mock-registry/components/csharp_dotnet-webapi/<component>/
  └─ mock-registry/components/golang_gin/<component>/

Bước 2: Thêm entry vào manifest.json
  └─ components.<component>.implementations.nodejs_express
  └─ components.<component>.implementations.csharp_dotnet-webapi
  └─ components.<component>.implementations.golang_gin

Bước 3: Test cài đặt
  └─ forge add <component> (cho cả 3 stacks)
  └─ Verify file output đúng vị trí
  └─ Verify dependencies cài đúng

Bước 4: Viết postInstall instructions rõ ràng
  └─ Đoạn code user cần thêm vào file chính
  └─ Environment variables cần set
  └─ Link tài liệu tham khảo
```

**Checklist Phase 2:**
- [x] 12 components mới × 3 stacks = 36 file code mới
- [x] manifest.json cập nhật đầy đủ
- [x] Mỗi component có postInstall instructions rõ ràng
- [x] Test `forge add` + `forge remove` cho từng component mới
- [x] Cập nhật `forge list` output cho đẹp với 18 components

---

## 7. Phase 3 — Blueprints & DX Nâng Cao (Tuần 9–14)

> **Mục tiêu:** Tạo Blueprints (workflow tổ hợp), cải thiện developer experience

### Task 3.1: Conflict Resolution thông minh
**Effort:** 🟡 Medium (3 ngày)

Hiện tại `forge add` ghi đè file không hỏi. Cần xử lý:

```
╔════════════════════════════════════════════════════════╗
║  ⚠ File conflict detected!                            ║
║                                                        ║
║  File: src/middlewares/errorHandler.js                 ║
║  Status: Modified by user (differs from original)      ║
║                                                        ║
║  What would you like to do?                            ║
║                                                        ║
║  ○ Overwrite — Replace with new version                ║
║  ● Skip     — Keep your current file                  ║
║  ○ Backup   — Save current as .bak, then overwrite    ║
║  ○ Diff     — Show differences first                  ║
╚════════════════════════════════════════════════════════╝
```

**Cách implement:**
1. Khi `forge add`, check file đã tồn tại chưa
2. Nếu tồn tại, so sánh hash (MD5/SHA256) với phiên bản gốc
3. Nếu hash khác (user đã sửa) → hiện TUI prompt hỏi action
4. Lưu hash gốc vào `.forge.json` để detect thay đổi sau này

**File cần sửa:**
- `cmd/add.go` — Thêm logic check file existence + hash compare
- `internal/config/config.go` — Thêm field `fileHashes` vào InstalledComponent
- `internal/tui/prompts.go` — Thêm prompt conflict resolution

### Task 3.2: Auto .env Configuration
**Effort:** 🟢 Small (2 ngày)

Khi `forge add jwt-auth`, CLI tự động:
1. Tìm file `.env` trong project root
2. Nếu có → append thiếu biến (`JWT_SECRET`, `JWT_EXPIRES_IN`)
3. Nếu chưa có → tạo `.env` + `.env.example`
4. Sinh giá trị mặc định an toàn (crypto random cho secrets)

**Thay đổi manifest:**
```jsonc
// Thêm field mới vào implementation
{
  "envVars": [
    { "key": "JWT_SECRET", "default": "{{RANDOM_SECRET_32}}", "description": "Secret key for JWT signing" },
    { "key": "JWT_EXPIRES_IN", "default": "7d", "description": "Token expiration time" }
  ]
}
```

**File cần sửa:**
- `cmd/add.go` — Thêm step "Configure environment variables"
- `internal/registry/manifest.go` — Thêm struct `EnvVar`

### Task 3.3: Implement Blueprints
**Effort:** 🟡 Medium-Large (5-7 ngày)

**Blueprint đầu tiên:** `blueprint-auth` (Authentication flow hoàn chỉnh)

Khi user chạy: `forge add blueprint-auth`

```
⚡ CoreForge — Add Blueprint
─────────────────────────────────────────

  Blueprint: auth (Complete Authentication Flow)
  
  This blueprint will install:
    • jwt-auth         — JWT middleware + token service
    • password-hashing — Bcrypt password utilities  
    • schema-user-auth — User database model
    • auth-controller  — Register + Login endpoints
    • auth-routes      — Route configuration

  Proceed? [Y/n]

  ✓ jwt-auth installed
  ✓ password-hashing installed
  ✓ schema-user-auth installed
  ✓ auth-controller created → /controllers/authController.js
  ✓ auth-routes created → /routes/authRoutes.js

  📋 Post-install:
     1. Connect your database (see schema-user-auth docs)
     2. Add to app.js:
        const authRoutes = require('./routes/authRoutes');
        app.use('/api/auth', authRoutes);
     3. Set .env: JWT_SECRET=<auto-generated>

  ✓ Blueprint auth installed successfully!
```

**Logic:**
```go
// cmd/add.go — mở rộng
func addBlueprint(name string, cfg *config.ForgeConfig, manifest *registry.Manifest) error {
    blueprint := manifest.Blueprints[name]
    
    // Step 1: Install included components (skip if already installed)
    for _, compName := range blueprint.Includes {
        if !cfg.HasComponent(compName) {
            addComponent(compName, cfg, manifest) // reuse existing logic
        }
    }
    
    // Step 2: Install blueprint-specific files (controllers, routes)
    for _, file := range blueprint.Files {
        downloadAndWrite(file)
    }
    
    // Step 3: Auto-configure .env
    configureEnvVars(blueprint.EnvVars)
    
    return nil
}
```

### Task 3.4: `forge update` Command
**Effort:** 🟡 Medium (3 ngày)

```bash
forge update jwt-auth     # Cập nhật 1 component
forge update --all        # Cập nhật tất cả
forge update --check      # Chỉ kiểm tra, không cập nhật
```

**Logic:**
1. Load manifest mới nhất từ remote
2. So sánh version trong manifest với version trong `.forge.json`
3. Nếu có bản mới → hiện diff + hỏi user
4. Respect conflict resolution (Task 3.1)

**File cần tạo:**
- `cmd/update.go` — Command mới

### Task 3.5: Post-Install Instructions nâng cao (Dev Deps)
**Effort:** 🟢 Small (1 ngày)

Phân tách `dependencies` thành:
```jsonc
{
  "dependencies": ["jsonwebtoken"],           // production
  "devDependencies": ["@types/jsonwebtoken"], // dev only
  "installCmd": "npm install jsonwebtoken",
  "installDevCmd": "npm install -D @types/jsonwebtoken"
}
```

**Checklist Phase 3:**
- [x] Conflict resolution (Overwrite / Skip / Backup) — hash-based detection + TUI prompt
- [x] Auto .env configuration — `internal/env/dotenv.go` + auto-append + secret generation
- [x] Blueprint system + `blueprint-auth` đầu tiên — 3 stacks (Node/C#/Go)
- [x] `forge update` command — `cmd/update.go` (--all, --check)
- [x] devDependencies phân tách — `installDevCmd` đã có từ Phase 1
- [x] Unit tests cho Phase 3 (33 tests — all pass)
- [x] `go vet` + `go build` xanh

---

## 8. Phase 4 — Website & Ecosystem (Tuần 15–20)

> **Mục tiêu:** Tạo website showcase, push registry lên GitHub, marketing

### Task 4.1: Publish GitHub Registry
**Effort:** 🟡 Medium (2-3 ngày)

Tách `mock-registry/` ra repo riêng: `github.com/longgoll/coreforge-registry`

```
coreforge-registry/
├── manifest.json              # Production manifest
├── components/
│   ├── nodejs_express/
│   ├── csharp_dotnet-webapi/
│   └── golang_gin/
├── schemas/
├── blueprints/
├── foundations/
└── .github/
    └── workflows/
        └── validate.yml       # CI: validate manifest + code syntax
```

**Đổi default registry URL:**
```go
// cmd/root.go
const defaultRegistryURL = "https://raw.githubusercontent.com/longgoll/coreforge-registry/main/manifest.json"
```

### Task 4.2: Website — registry.coreforge.dev
**Effort:** 🔴 Large (2-3 tuần)

**Tech stack cho website:**
- **Framework:** Next.js 14+ (App Router)
- **Styling:** Tailwind CSS + shadcn/ui (vì CoreForge lấy cảm hứng từ shadcn, website cũng nên dùng)
- **Hosting:** Vercel (free tier)
- **Content:** Đọc trực tiếp từ `coreforge-registry` repo

**Trang chính:**
```
┌─────────────────────────────────────────────────────────────────┐
│  🏗️ CoreForge                              [Docs] [GitHub] [★] │
│                                                                  │
│  The Backend Component Registry                                  │
│  Production-ready code blocks for Node.js, C#, and Go           │
│                                                                  │
│  [Get Started]    [Browse Components]                            │
│                                                                  │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐                      │
│  │ npm i -g │  │ dotnet   │  │ go       │                      │
│  │ forge    │  │ tool     │  │ install  │                      │
│  │          │  │ forge    │  │ forge    │                      │
│  │ 18 comps │  │ 18 comps │  │ 18 comps │                      │
│  └──────────┘  └──────────┘  └──────────┘                      │
│                                                                  │
│  ── Components ──────────────────────────────────────────        │
│                                                                  │
│  [error-handler]  [logger]  [jwt-auth]  [cors]                  │
│  [rate-limiter]   [validation]  [env-config]  [rbac]            │
│  ...                                                             │
│                                                                  │
│  ── Code Preview ────────────────────────────────────────        │
│  ┌─ Tab: [Node.js] [C#] [Go] ─────────────────────────┐        │
│  │  // errorHandler.js                                  │        │
│  │  const errorHandler = (err, req, res, next) => {     │        │
│  │    const statusCode = err.statusCode || 500;          │        │
│  │    ...                                                │        │
│  │  };                                                   │        │
│  └──────────────────────────────────────────────────────┘        │
│                                                                  │
│  [Copy Code]  [forge add error-handler]                         │
└─────────────────────────────────────────────────────────────────┘
```

**Tính năng website:**
- Dark mode mặc định
- Code preview tab switching (Node/C#/Go) cho từng component
- Copy-to-clipboard cho cả code và CLI command
- Search + filter theo category/tags
- Version compatibility matrix
- SEO optimized (để Google tìm được khi dev search "express error handler boilerplate")

### Task 4.3: Version Compatibility Matrix
**Effort:** 🟢 Small (1 ngày)

Thêm vào manifest:
```jsonc
{
  "components": {
    "jwt-auth": {
      "implementations": {
        "nodejs_express": {
          "compatibility": {
            "express": ">=4.18.0",
            "node": ">=18.0.0"
          }
        }
      }
    }
  }
}
```

CLI check version khi `forge add` và cảnh báo nếu không tương thích.

**Checklist Phase 4:**
- [x] Registry repo riêng trên GitHub (`github.com/longgoll/coreforge-registry`)
- [x] CI validate manifest + syntax check (`.github/workflows/validate.yml`)
- [ ] Website cơ bản (component catalog + code preview)
- [ ] Dark mode + responsive
- [ ] Search/filter components
- [ ] Version compatibility warnings
- [ ] Domain setup (coreforge.dev hoặc tương tự)

---

## 9. Phase 5 — Đẳng Cấp Thượng Thừa (Tuần 21+)

> **Mục tiêu:** Vượt xa đối thủ. Tính năng "chưa ai làm" trong giới Backend scaffold tools.

### Task 5.1: AST Code Injection — Tự Động Nối Code
**Effort:** 🔴 XL (2-4 tuần)

**Vấn đề hiện tại:** Khi `forge add jwt-auth`, CLI in ra:
```
📋 Post-install steps:
   Add to your app.js:
     const { authMiddleware } = require('./middlewares/authMiddleware');
     app.use(authMiddleware);
```
→ User phải tự mở file và paste. Dễ sai, dễ quên.

**Giải pháp:** Dùng AST parser để tự động chèn code vào đúng vị trí:

| Stack | AST Parser | Cách inject |
|-------|-----------|-------------|
| Node.js | `github.com/nicolo-ribaudo/tc39-proposal` hoặc dùng regex đơn giản trước | Tìm `app.use(` cuối cùng → chèn dòng mới sau |
| C# | Roslyn — nhưng quá nặng cho CLI Go → dùng regex pattern matching | Tìm `app.MapControllers()` → chèn trước |
| Go | `go/ast` + `go/parser` (native Go) — dễ nhất | Parse `main.go` → thêm middleware registration |

**Chiến lược triển khai thực tế (pragmatic):**

Thay vì full AST parsing ngay, **Phase 1** dùng **Pattern-Based Injection**:
```go
// Tìm pattern đặc trưng trong file và chèn code
type InjectionPoint struct {
    File       string   // "app.js", "Program.cs", "main.go"
    Pattern    string   // regex: `app\.use\(.*\);`
    Position   string   // "after-last" | "before-first" | "after-pattern"
    Code       string   // code cần chèn
}
```

**Phase 2** mới dùng full AST cho Go (vì Go có built-in AST parser).

### Task 5.2: Docker Compose Auto-Integration
**Effort:** 🟡 Medium (3-5 ngày)

Khi `forge add` component cần external service:

```bash
forge add redis-cache
```

CLI tự động:
1. Detect `docker-compose.yml` trong project root
2. Nếu có → parse YAML → thêm Redis service block
3. Nếu chưa có → tạo `docker-compose.yml` mới

```yaml
# Tự động gắn vào docker-compose.yml
services:
  redis:                    # ← CoreForge tự thêm
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

volumes:
  redis_data:               # ← CoreForge tự thêm
```

**Components có Docker integration:**
| Component | Docker Service | Port |
|-----------|---------------|------|
| `redis-cache` | Redis 7 Alpine | 6379 |
| `schema-*` (PostgreSQL) | PostgreSQL 16 | 5432 |
| `schema-*` (MongoDB) | MongoDB 7 | 27017 |
| `email-service` | MailHog (dev) | 1025/8025 |

### Task 5.3: Live Sandbox Playground trên Website
**Effort:** 🔴 XL (4+ tuần)

**Ý tưởng:** Trên website, user bấm "Test" để thấy component hoạt động:
- `rate-limiter` → Bấm "Send Request" 10 lần → thấy HTTP 429
- `jwt-auth` → Nhập credentials → nhận JWT token → gọi protected route
- `health-check` → Xem response `{"status": "ok", "uptime": "..."}`

**Cách triển khai:**
- Dùng **WebContainer** (StackBlitz technology) để chạy Node.js trong browser
- Hoặc dùng **mock API** bằng JavaScript thuần mô phỏng behavior
- Phase đầu: chỉ cần animated demo (GIF/video), chưa cần real sandbox

### Task 5.4: Multi-Language Architecture Translation
**Effort:** 🔴 XL (ongoing)

Đảm bảo mỗi component trên 3 ngôn ngữ có:
- **Cùng function signature concept** (tên hàm tương đương)
- **Cùng response format** (`{"success": true, "data": {}, "message": "..."}`)
- **Cùng error codes** (401, 403, 404, 429...)
- **Cùng env variables** (`JWT_SECRET`, `PORT`, `DB_URL`...)

Tạo file **`architecture-spec.md`** trong registry — mô tả "contract" cho mỗi component mà tất cả language phải follow.

### Task 5.5: OAuth Social Login  
**Effort:** 🟡 Medium-Large (5-7 ngày)

Component `oauth` hỗ trợ:
- Google Login
- GitHub Login
- Facebook Login (optional)

Mỗi provider là 1 file riêng có thể chọn thêm:
```bash
forge add oauth              # Hỏi chọn providers
forge add oauth --google     # Chỉ Google
```

### Task 5.6: Plugin System — 3rd Party Registries
**Effort:** 🔴 Large (2-3 tuần)

Cho phép community tạo registry riêng:
```bash
forge registry add https://my-company.com/forge-registry/manifest.json
forge add my-company/custom-component
```

**Checklist Phase 5:**
- [ ] Pattern-Based Code Injection (pragmatic AST)
- [ ] Docker Compose auto-integration
- [ ] Website sandbox/demo (ít nhất animated demo)
- [ ] Architecture contract spec
- [ ] OAuth component
- [ ] Plugin system (nếu có community demand)

---

## 10. Catalog Component Đầy Đủ

### Tổng quan: 24 components mục tiêu

| # | Component | Category | Trạng thái | Phase |
|---|-----------|----------|------------|-------|
| 1 | `error-handler` | middleware | ✅ Done | — |
| 2 | `logger` | middleware | ✅ Done | — |
| 3 | `jwt-auth` | auth | ✅ Done | — |
| 4 | `cors` | middleware | ✅ Done | — |
| 5 | `rate-limiter` | middleware | ✅ Done | — |
| 6 | `validation` | middleware | ✅ Done | — |
| 7 | `env-config` | utility | ✅ Done | Phase 2A |
| 8 | `response-formatter` | utility | ✅ Done | Phase 2A |
| 9 | `not-found-handler` | middleware | ✅ Done | Phase 2A |
| 10 | `health-check` | utility | ✅ Done | Phase 2A |
| 11 | `password-hashing` | auth | ✅ Done | Phase 2B |
| 12 | `security-headers` | middleware | ✅ Done | Phase 2B |
| 13 | `rbac` | auth | ✅ Done | Phase 2B |
| 14 | `generate-otp` | auth | ✅ Done | Phase 2B |
| 15 | `shutdown-handler` | utility | ✅ Done | Phase 2C |
| 16 | `async-handler` | middleware | ✅ Done | Phase 2C |
| 17 | `swagger-docs` | documentation | ✅ Done | Phase 2C |
| 18 | `pagination` | utility | ✅ Done | Phase 2C |
| 19 | `file-upload` | utility | 🔲 Todo | Phase 3+ |
| 20 | `cron-job` | background | 🔲 Todo | Phase 3+ |
| 21 | `email-service` | service | 🔲 Todo | Phase 3+ |
| 22 | `redis-cache` | database | 🔲 Todo | Phase 5 |
| 23 | `database-config` | database | 🔲 Todo | Phase 5 |
| 24 | `oauth` | auth | 🔲 Todo | Phase 5 |

### Chi tiết library/framework cho mỗi stack

| Component | Node.js (Express) | C# (.NET) | Go (Gin) |
|-----------|-------------------|-----------|----------|
| `env-config` | dotenv + Joi validate | IConfiguration + DataAnnotations | godotenv + custom validate |
| `response-formatter` | Utility function | ActionResult\<ApiResponse\<T\>\> | gin.JSON wrapper |
| `not-found-handler` | Express middleware | Middleware / UseStatusCodePages | r.NoRoute() |
| `health-check` | GET /health handler | MapHealthChecks() | GET /health handler |
| `password-hashing` | bcryptjs | BCrypt.Net-Next | golang.org/x/crypto/bcrypt |
| `security-headers` | helmet | Custom middleware | Custom middleware |
| `rbac` | Custom middleware | AuthorizeAttribute + Policy | Custom middleware |
| `generate-otp` | crypto.randomInt | RNGCryptoServiceProvider | crypto/rand |
| `shutdown-handler` | process.on('SIGTERM') | IHostApplicationLifetime | os.Signal + context.Done |
| `async-handler` | Wrapper function | N/A (async/await built-in) | N/A (Gin recovery) |
| `swagger-docs` | swagger-jsdoc + swagger-ui-express | Swashbuckle.AspNetCore | swaggo/swag + gin-swagger |
| `pagination` | Utility + query parser | Extension method IQueryable | Utility + query parser |
| `file-upload` | multer + S3 SDK | IFormFile + S3 SDK | gin multipart + S3 SDK |
| `cron-job` | node-cron | Quartz.NET | robfig/cron |
| `email-service` | nodemailer / Resend SDK | MailKit / SendGrid SDK | gomail / SendGrid SDK |
| `redis-cache` | ioredis | StackExchange.Redis | go-redis |
| `database-config` | Prisma / Mongoose setup | EF Core DbContext | GORM setup |
| `oauth` | passport.js | AspNet.Security.OAuth | golang.org/x/oauth2 |

---

## 11. Chi Tiết Kỹ Thuật

### 11.1 File Structure sau khi nâng cấp

```
CoreForge/
├── main.go
├── go.mod
├── go.sum
├── Makefile
├── build.ps1
├── .gitignore
├── README.md
├── AI_CONTEXT.md
├── ROADMAP.md                       ← FILE NÀY
│
├── .github/workflows/
│   ├── ci.yml
│   └── release.yml
│
├── cmd/
│   ├── root.go
│   ├── init.go                      # Mở rộng: DB/ORM/Architecture prompts
│   ├── add.go                       # Mở rộng: Blueprint support, conflict resolution, .env config
│   ├── remove.go
│   ├── list.go                      # Mở rộng: Grouped output theo tầng
│   ├── search.go                    # Mở rộng: Cross-tầng search
│   ├── doctor.go
│   ├── update.go                    # ← MỚI: Phase 3
│   ├── template.go
│   └── zip.go
│
├── internal/
│   ├── config/
│   │   ├── config.go                # V2: thêm Database/ORM/Architecture + migration
│   │   └── config_test.go           # ← MỚI: Phase 1
│   ├── registry/
│   │   ├── manifest.go              # V2: Foundations/Schemas/Blueprints support
│   │   └── manifest_test.go         # ← MỚI: Phase 1
│   ├── tui/
│   │   └── prompts.go               # Mở rộng: DB/ORM prompts, conflict resolution
│   ├── env/                         # ← MỚI: Phase 3
│   │   └── dotenv.go                # Parse/write .env files
│   └── inject/                      # ← MỚI: Phase 5
│       └── pattern.go               # Pattern-based code injection
│
├── mock-registry/                   # DEV — sẽ tách ra repo riêng ở Phase 4
│   ├── manifest.json
│   ├── components/     (24 components × 3 stacks)
│   ├── schemas/        (MỚI)
│   ├── blueprints/     (MỚI)
│   ├── foundations/    (MỚI - thay thế templates/)
│   └── templates/      (GIỮ cho backward compat, deprecated)
│
└── test-project/                    # Testing sandbox
```

### 11.2 manifest.json v2 Migration Flow

```
CLI khởi động
  │
  ├─ Parse manifest.json
  │     │
  │     ├─ Có field "schemaVersion"?
  │     │     ├─ YES → dùng v2 parser
  │     │     └─ NO  → dùng v1 parser (backward compat)
  │     │
  │     └─ Return Manifest struct (unified)
  │
  └─ Parse .forge.json
        │
        ├─ Có field "configVersion"?
        │     ├─ YES → dùng v2 struct
        │     └─ NO  → dùng v1 struct, auto-fill v2 defaults
        │
        └─ Return ForgeConfig struct (unified)
```

### 11.3 Quy tắc naming convention

| Ngôn ngữ | Files | Folders | Functions/Methods | Classes/Types |
|-----------|-------|---------|-------------------|---------------|
| Node.js | camelCase.js | camelCase/ | camelCase() | PascalCase |
| C# | PascalCase.cs | PascalCase/ | PascalCase() | PascalCase |
| Go | snake_case.go | snake_case/ | PascalCase() (exported) / camelCase() (unexported) | PascalCase |

---

## 12. Chiến Lược Testing

### 12.1 Unit Tests (Phase 1)

```go
// config_test.go
func TestLoadV1Config(t *testing.T) { ... }
func TestLoadV2Config(t *testing.T) { ... }
func TestV1ToV2Migration(t *testing.T) { ... }
func TestHasComponent(t *testing.T) { ... }
func TestAddComponent(t *testing.T) { ... }

// manifest_test.go
func TestLoadManifestV1(t *testing.T) { ... }
func TestLoadManifestV2(t *testing.T) { ... }
func TestResolveComponent(t *testing.T) { ... }
func TestResolveBlueprint(t *testing.T) { ... }
```

### 12.2 Integration Tests (Phase 2+)

Tạo folder `tests/` với test scripts:

```bash
# tests/test_add_flow.sh
# 1. Init project
forge init --language nodejs --framework express --source ./src --non-interactive

# 2. Add component
forge add error-handler

# 3. Verify files exist
test -f ./src/middlewares/errorHandler.js || exit 1
test -f ./src/utils/AppError.js || exit 1

# 4. Verify .forge.json updated
grep -q "error-handler" .forge.json || exit 1

# 5. Remove component
forge remove error-handler --force

# 6. Verify files removed
test ! -f ./src/middlewares/errorHandler.js || exit 1

echo "✓ All tests passed"
```

### 12.3 CI Pipeline mở rộng

```yaml
# .github/workflows/ci.yml — thêm test job
jobs:
  unit-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
      - run: go test ./internal/... -v -cover

  integration-test:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, windows-latest, macos-latest]
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
      - run: go build -o forge .
      - run: |
          mkdir test-tmp && cd test-tmp
          ../forge init --non-interactive --language nodejs --framework express
          ../forge add error-handler
          ../forge list --installed
          ../forge remove error-handler --force
```

### 12.4 Component Code Validation

Mỗi khi thêm component mới, chạy syntax check:
- **Node.js:** `node --check <file>.js` hoặc `npx acorn --ecma2020 <file>.js`
- **C#:** `dotnet build` (trong folder test)
- **Go:** `go vet` + `gofmt -l`

---

## 13. Chiến Lược Community & Open Source

### 13.1 Contributing Guidelines

Tạo `CONTRIBUTING.md` hướng dẫn:

1. **Thêm component mới:**
   - Fork repo `coreforge-registry`
   - Tạo code cho **cả 3 stacks** (hoặc tối thiểu 1, ghi rõ unavailable cho stack khác)
   - Thêm entry trong `manifest.json`
   - Viết `postInstall` instructions rõ ràng
   - Submit PR

2. **Component Authoring Spec:**
   Mỗi component PHẢI có:
   - Code hoạt động ngay (copy → paste → chạy)
   - Comment giải thích ở những chỗ quan trọng
   - Không hardcode giá trị, dùng environment variables
   - Response format theo chuẩn: `{"success": bool, "data": any, "message": string}`
   - Error codes theo HTTP standard

3. **Review Checklist:**
   - [ ] Code chạy được trên stack tương ứng
   - [ ] Có postInstall instructions
   - [ ] Dependencies liệt kê đầy đủ
   - [ ] Naming convention đúng stack
   - [ ] Không conflict với component khác

### 13.2 Marketing & Visibility

| Channel | Action | Timeline |
|---------|--------|----------|
| GitHub | Viết README xuất sắc + badges + GIF demo | Phase 1 |
| DEV.to | Bài viết: "I built a shadcn/ui for Backend" | Phase 2 |
| Reddit | Post trên r/node, r/golang, r/dotnet | Phase 4 |
| YouTube | Video demo 5 phút | Phase 4 |
| Product Hunt | Launch | Phase 4 |
| X (Twitter) | Thread giới thiệu + GIF demo | Phase 4 |

---

## 📅 Timeline Tổng Hợp

```
 Tuần 1─3    Phase 1: Nền tảng vững chắc
  ┃            ├─ Manifest v2 schema
  ┃            ├─ DB/ORM/Architecture init flow
  ┃            ├─ .forge.json v2 + migration
  ┃            └─ Unit tests core
  ┃
 Tuần 4─5    Phase 2A: 4 components nền tảng
  ┃            └─ env-config, response-formatter, not-found-handler, health-check
  ┃
 Tuần 6─7    Phase 2B: 4 components bảo mật
  ┃            └─ password-hashing, security-headers, rbac, generate-otp
  ┃
 Tuần 8      Phase 2C: 4 components utility
  ┃            └─ shutdown-handler, async-handler, swagger-docs, pagination
  ┃
 Tuần 9─11   Phase 3A: DX nâng cao
  ┃            ├─ Conflict resolution
  ┃            ├─ Auto .env config
  ┃            └─ forge update command
  ┃
 Tuần 12─14  Phase 3B: Blueprints
  ┃            ├─ Blueprint system implementation
  ┃            ├─ blueprint-auth
  ┃            └─ blueprint-crud-api
  ┃
 Tuần 15─17  Phase 4A: Registry & Website
  ┃            ├─ Publish registry repo
  ┃            └─ Website v1 (catalog + code preview)
  ┃
 Tuần 18─20  Phase 4B: Polish & Launch
  ┃            ├─ README / docs / video demo
  ┃            ├─ Product Hunt launch
  ┃            └─ Blog posts / social media
  ┃
 Tuần 21+    Phase 5: Next-Level
               ├─ AST Code Injection
               ├─ Docker Compose integration
               ├─ OAuth component
               ├─ Plugin system
               └─ Live Sandbox
```

---

## ⚡ Quick Reference — Chạy Ngay Hôm Nay

Nếu bạn chỉ có **1 ngày**, hãy làm những thứ này:

1. **Viết 1 component mới đơn giản nhất:** `health-check` (×3 stacks, ~2 giờ)
2. **Thêm vào manifest.json** (~15 phút)
3. **Test `forge add health-check`** (~15 phút)
4. **Commit + push** (~5 phút)

→ Bạn đã thêm giá trị thực sự cho project mà không cần đụng vào core. 🚀
