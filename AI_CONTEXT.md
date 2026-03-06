# 🧠 AI CONTEXT — Forge CLI

> **Đọc file này trước khi code bất kỳ thứ gì trong dự án.**
> File này giúp AI (hoặc developer mới) hiểu nhanh dự án đang làm gì, kiến trúc như thế nào, và đang ở giai đoạn nào.

---

## 1. Dự án là gì?

**Forge CLI** là một công cụ dòng lệnh (CLI) đa năng dành cho Backend Developer, lấy cảm hứng từ triết lý của [shadcn/ui](https://ui.shadcn.com/).

### Ý tưởng cốt lõi

> "Copy, don't install" — Code được **sao chép trực tiếp** vào project của người dùng.
> Họ sở hữu code 100%, tự do chỉnh sửa, không bị phụ thuộc (vendor lock-in).

### Forge CLI giải quyết vấn đề gì?

Backend developers thường phải copy-paste cùng một boilerplate code (error handler, JWT auth, logger...) giữa các dự án. Forge CLI biến việc đó thành **một câu lệnh duy nhất**:

```bash
forge add error-handler   # → Copy error handler vào project, tự cài dependencies
forge add jwt-auth        # → Copy JWT auth middleware + token service vào project
```

### Điểm khác biệt so với các tool khác

| Feature                       | Forge CLI | Yeoman | Plop | Hygen |
| ----------------------------- | --------- | ------ | ---- | ----- |
| Multi-language (Node, C#, Go) | ✅        | ❌     | ❌   | ❌    |
| No vendor lock-in             | ✅        | ❌     | ✅   | ✅    |
| Single binary (Go)            | ✅        | ❌     | ❌   | ❌    |
| Backend focused               | ✅        | ❌     | ❌   | ❌    |
| Remote registry               | ✅        | ✅     | ❌   | ❌    |

---

## 2. Tech Stack

| Thành phần          | Công nghệ                                                 | Vai trò                                       |
| ------------------- | --------------------------------------------------------- | --------------------------------------------- |
| **Ngôn ngữ lõi**    | Golang                                                    | Biên dịch thành single binary, cross-platform |
| **CLI Framework**   | [spf13/cobra](https://github.com/spf13/cobra)             | Quản lý commands, flags, help text            |
| **Interactive TUI** | [charmbracelet/huh](https://github.com/charmbracelet/huh) | Menu chọn ngôn ngữ/framework, confirm prompts |
| **Terminal Colors** | [fatih/color](https://github.com/fatih/color)             | Output có màu (✓ xanh, ✗ đỏ, ⚠ vàng)          |
| **Go Module**       | `github.com/longgoll/forge-cli`                           | Module path                                   |

---

## 3. Kiến trúc Hệ thống (3 thành phần)

```
┌─────────────────────────────────────────────────────────────────┐
│                        FORGE CLI (Go Binary)                     │
│                                                                  │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────────┐│
│  │ forge    │  │ forge    │  │ forge    │  │ forge            ││
│  │ init     │  │ add      │  │ list     │  │ doctor           ││
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────────┬─────────┘│
│       │              │             │                  │          │
│  ┌────▼──────────────▼─────────────▼──────────────────▼────────┐│
│  │                     INTERNAL PACKAGES                        ││
│  │  config/     → .forge.json read/write                        ││
│  │  registry/   → manifest.json parse + fetch                   ││
│  │  tui/        → Interactive prompts (huh)                     ││
│  └──────────────────────┬──────────────────────────────────────┘│
└─────────────────────────┼───────────────────────────────────────┘
                          │ HTTP / File read
                          ▼
┌─────────────────────────────────────────────────────────────────┐
│                   REMOTE REGISTRY (GitHub)                        │
│                                                                  │
│  manifest.json ──── "Cuốn thực đơn" liệt kê components          │
│  components/   ──── Source code files cho từng stack              │
│                                                                  │
│  Repo: github.com/longgoll/forge-registry (tương lai)            │
│  Hiện tại: dùng mock-registry/ local để dev/test                 │
└─────────────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────────┐
│                USER'S PROJECT (bất kỳ thư mục nào)               │
│                                                                  │
│  .forge.json ──── Config file: language, framework, installed    │
│  src/                                                            │
│   ├── middlewares/   ← forge add error-handler, logger, jwt-auth │
│   ├── services/      ← forge add jwt-auth (tokenService)        │
│   └── utils/         ← forge add error-handler (AppError)       │
└─────────────────────────────────────────────────────────────────┘
```

---

## 4. Cấu trúc Thư mục

```
CoreForge/
├── main.go                          # Entry point → gọi cmd.Execute()
├── go.mod                           # Go module (github.com/longgoll/forge-cli)
├── go.sum                           # Dependency checksums
├── forge.exe                        # Binary đã build (gitignore)
├── Makefile                         # Cross-compile (Linux/macOS)
├── build.ps1                        # Cross-compile (Windows PowerShell)
├── .gitignore
├── README.md
├── AI_CONTEXT.md                    # ← FILE NÀY
│
├── .github/workflows/               # ── CI/CD ──
│   ├── ci.yml                       # Test + build on push/PR (3 OSes)
│   └── release.yml                  # Auto-release on version tag (v*)
│
├── cmd/                             # ── COBRA COMMANDS ──
│   ├── root.go                      # Root command + ASCII banner + --remote flag
│   ├── init.go                      # forge init  — setup project → tạo .forge.json
│   ├── add.go                       # forge add   — download component → cài vào project
│   ├── remove.go                    # forge remove — xóa component (files + config)
│   ├── list.go                      # forge list  — hiện components (available / installed)
│   ├── search.go                    # forge search — tìm component theo keyword/tag
│   └── doctor.go                    # forge doctor — check environment + registry info
│
├── internal/                        # ── CORE LOGIC (không export ra ngoài) ──
│   ├── config/
│   │   └── config.go                # Struct ForgeConfig, đọc/ghi .forge.json
│   ├── registry/
│   │   └── manifest.go              # Load local/remote, ETag caching, offline fallback
│   └── tui/
│       └── prompts.go               # Interactive prompts dùng charmbracelet/huh
│
├── mock-registry/                   # ── LOCAL TEST REGISTRY ──
│   ├── manifest.json                # Danh sách languages + 6 components
│   └── components/                  # Source code templates
│       ├── nodejs_express/          # error-handler, logger, jwt-auth, cors, rate-limiter, validation
│       ├── csharp_dotnet-webapi/    # Cấu trúc tương tự nodejs
│       └── golang_gin/             # Cấu trúc theo Go convention
│
├── dist/                            # ── BUILD OUTPUT (gitignore) ──
│   ├── forge-windows-amd64.exe
│   ├── forge-linux-amd64
│   ├── forge-darwin-arm64
│   └── checksums.txt
│
└── test-project/                    # ── THƯ MỤC TEST (không commit) ──
    ├── .forge.json
    └── src/
```

---

## 5. Các Data Structure Quan Trọng

### 5.1 `.forge.json` (nằm ở root project của user)

```json
{
  "language": "nodejs",
  "framework": "express",
  "sourceDir": "./src",
  "installedComponents": [
    {
      "name": "error-handler",
      "version": "1.0.0",
      "installedAt": "2026-02-27T15:05:06+07:00"
    }
  ]
}
```

- **Ai tạo?** `forge init` (interactive) hoặc user tạo tay
- **Ai đọc?** `forge add` đọc để biết stack → match component
- **Ai cập nhật?** `forge add` thêm vào `installedComponents` sau khi cài

### 5.2 `manifest.json` (registry — "cuốn thực đơn")

```json
{
  "version": "1.0.0",
  "languages": {
    "<language_key>": {
      "name": "Display Name",
      "frameworks": {
        "<framework_key>": {
          "name": "Framework Name",
          "description": "...",
          "minVersion": "x.x.x"
        }
      }
    }
  },
  "components": {
    "<component_name>": {
      "description": "...",
      "category": "middleware | auth | database | ...",
      "tags": ["keyword1", "keyword2"],
      "implementations": {
        "<language_key>_<framework_key>": {
          "files": [
            { "url": "path/or/url/to/file.js", "target": "/dest/in/project.js" }
          ],
          "dependencies": ["package-name"],
          "installCmd": "npm install package-name",
          "postInstall": "Hướng dẫn sau khi cài",
          "requires": ["other-component"],
          "conflicts": ["conflicting-component"]
        }
      }
    }
  }
}
```

**Quy tắc key:** Stack key = `<language>_<framework>` (VD: `nodejs_express`, `csharp_dotnet-webapi`, `golang_gin`)

---

## 6. Luồng Xử Lý Chi Tiết

### `forge init`

```
User chạy "forge init"
  │
  ├→ Check .forge.json đã tồn tại? → Hỏi overwrite
  ├→ Load manifest.json (từ mock-registry/ hoặc remote URL)
  ├→ TUI: Chọn Language (Node.js / C# / Golang)
  ├→ TUI: Chọn Framework (Express / .NET Web API / Gin)
  ├→ TUI: Nhập source directory (default: ./src)
  └→ Ghi .forge.json vào thư mục hiện tại
```

### `forge add <component>`

```
User chạy "forge add error-handler"
  │
  ├→ Đọc .forge.json → lấy language + framework (VD: nodejs_express)
  ├→ Check component đã cài chưa? → Báo "already installed"
  ├→ Load manifest.json
  ├→ Tìm component "error-handler" → match implementation cho "nodejs_express"
  ├→ Với mỗi file trong implementation:
  │     ├→ Download nội dung file (local path hoặc HTTP URL)
  │     ├→ Tạo thư mục đích nếu chưa có
  │     └→ Ghi file vào sourceDir + target path
  ├→ Nếu có installCmd → chạy (VD: npm install jsonwebtoken)
  ├→ Cập nhật .forge.json → thêm vào installedComponents
  └→ Hiện postInstall instructions + required components
```

### `forge list [--installed]`

```
Không có flag  → Load manifest.json → hiện tất cả components
--installed    → Đọc .forge.json → hiện components đã cài
```

### `forge doctor`

```
Check từng tool: node, npm, dotnet, go, git
  → Có: hiện ✓ + version
  → Không: hiện ✗ + link cài đặt
```

### `forge remove <component>` (NEW)

```
User chạy "forge remove error-handler"
  │
  ├→ Đọc .forge.json → kiểm tra component đã cài chưa
  ├→ Load manifest.json → tìm files thuộc component
  ├→ Liệt kê files sẽ bị xóa
  ├→ TUI: Confirm prompt (bỏ qua nếu --force)
  ├→ Xóa từng file, cleanup empty directories
  ├→ Cập nhật .forge.json → xóa khỏi installedComponents
  └→ Cảnh báo dependencies đã cài (không tự xóa)
```

### `forge search <keyword>` (NEW)

```
User chạy "forge search auth"
  │
  ├→ Load manifest.json
  ├→ Tìm keyword trong: name, description, category, tags
  ├→ Hiện kết quả match kèm: match reason, tags, stacks
  └→ Hiện ✓/✗ nếu đang ở project có .forge.json (stack compatibility)
```

---

## 7. Trạng Thái Hiện Tại

### ✅ Đã hoàn thành

- [x] CLI skeleton (Cobra) với 6 commands: `init`, `add`, `remove`, `list`, `search`, `doctor`
- [x] Data structures: `ForgeConfig`, `Manifest`, `Component`, `Implementation`
- [x] Interactive TUI prompts (charmbracelet/huh)
- [x] `.forge.json` read/write với installed component tracking
- [x] Manifest loader (local file, relative to binary)
- [x] Component downloader (local + HTTP support)
- [x] Auto dependency install (`npm install`, `dotnet add`, `go get`)
- [x] Post-install instructions
- [x] Duplicate install protection
- [x] `forge remove` — xóa file + cập nhật config, confirm prompt, `--force` flag
- [x] `forge search` — tìm theo name, description, tags, category
- [x] All 6 components × 3 stacks = **18 implementations**:
  - `error-handler` → Express, .NET, Gin
  - `logger` → Express, .NET, Gin
  - `jwt-auth` → Express, .NET, Gin
  - `cors` → Express, .NET, Gin
  - `rate-limiter` → Express, .NET, Gin
  - `validation` → Express, .NET, Gin

### ⚠️ Cần sửa / lưu ý

- [x] ~~**manifest.json cần update URL cho golang_gin**~~ → ĐÃ SỬA

### 🔲 Chưa làm (Roadmap)

- [x] ~~**Remote Registry**~~: HTTP fetch từ GitHub raw content + `--remote` flag → ĐÃ LÀM
- [x] ~~**Caching**~~: Cache manifest ở `~/.forge/cache/` + ETag/Last-Modified + TTL 1h + offline fallback → ĐÃ LÀM
- [x] ~~**Thêm components**~~: cors, rate-limiter, validation (3 × 3 stacks = 9 implementations) → ĐÃ LÀM
- [x] ~~**Cross-compile**~~: build.ps1 + Makefile cho 6 platforms (Win/Linux/macOS × amd64/arm64) → ĐÃ LÀM
- [x] ~~**CI/CD**~~: GitHub Actions — ci.yml (test on 3 OSes) + release.yml (auto-release on tag) → ĐÃ LÀM
- [ ] **Project Template Generator**: `forge init` tải .zip skeleton project
- [ ] **Go Template Engine**: Dùng `text/template` để inject biến vào code (project name, options)
- [ ] **Thêm components tiếp**: database-config, docker-compose, env-config...
- [ ] **`forge update`**: Cập nhật component đã cài lên version mới
- [ ] **Publish GitHub Registry**: Push manifest + components lên github.com/longgoll/forge-registry

---

## 8. Quy Tắc Khi Code

### Quy tắc chung

1. **CLI không chứa hardcode code mẫu** — tất cả template code nằm trong registry
2. **Mỗi component phải có đủ implementation cho cả 3 stacks** (hoặc ghi rõ trong manifest là "unavailable")
3. **File paths trong manifest dùng relative path** với prefix `./mock-registry/...` (dev) hoặc full HTTPS URL (prod)
4. **Binary resolve paths relative to itself**, không phải CWD — xem `getExeDir()` trong `add.go` và `manifest.go`

### Khi thêm component mới

1. Tạo file source code trong `mock-registry/components/<stack>/<component>/`
2. Thêm entry trong `manifest.json` → `components.<name>.implementations`
3. Đảm bảo `target` path theo convention của từng stack:
   - **Node.js**: `/middlewares/`, `/services/`, `/utils/` (camelCase)
   - **C#**: `/Middlewares/`, `/Services/`, `/Exceptions/` (PascalCase)
   - **Golang**: `/middlewares/`, `/services/`, `/utils/` (snake_case)

### Khi thêm command mới

1. Tạo file trong `cmd/<command>.go`
2. Register với `rootCmd.AddCommand()` trong `func init()`
3. Dùng color helpers từ `root.go`: `cyan()`, `green()`, `yellow()`, `red()`, `bold()`, `dimmed()`
4. Follow output format: `⚡ Forge CLI — <Title>` + divider line

---

## 9. Cách Build & Test

```bash
# Download dependencies
go mod tidy

# Build binary (current OS)
go build -o forge.exe .
# Hoặc dùng script (có version injection):
.\build.ps1              # Windows PowerShell

# Cross-compile cho tất cả platforms
.\build.ps1 -All         # Windows PowerShell
make build-all           # Linux/macOS (Makefile)

# Test commands
./forge.exe --help
./forge.exe doctor
./forge.exe list
./forge.exe search auth

# Test full flow (tạo thư mục test riêng)
mkdir my-test && cd my-test
../forge.exe init          # Interactive setup
../forge.exe add error-handler
../forge.exe add cors
../forge.exe remove cors --force
../forge.exe list --installed

# Release (GitHub Actions — tự động khi push tag)
git tag v1.0.0
git push origin v1.0.0     # → CI builds + creates GitHub Release
```

---

## 10. Liên Kết

- **CLI Engine**: [github.com/longgoll/coreforge](https://github.com/longgoll/coreforge)
- **Registry** (tương lai): [github.com/longgoll/coreforge](https://github.com/longgoll/coreforge)
- **Cảm hứng**: [shadcn/ui](https://ui.shadcn.com/) và https://servercn.vercel.app/ — "Copy and paste. Use the CLI."
