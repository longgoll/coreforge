package registry

import (
	"encoding/json"
	"testing"
)

// ──────────────────────────────────────────────────────────────
// Test Data Fixtures
// ──────────────────────────────────────────────────────────────

const v1ManifestJSON = `{
  "version": "1.0.0",
  "languages": {
    "nodejs": {
      "name": "Node.js",
      "frameworks": {
        "express": {
          "name": "Express",
          "description": "Fast, unopinionated web framework",
          "templateUrl": "./templates/nodejs_express.zip",
          "minVersion": "4.18.0"
        }
      }
    }
  },
  "components": {
    "error-handler": {
      "description": "Global Error Handler middleware",
      "category": "middleware",
      "tags": ["error", "exception", "middleware"],
      "implementations": {
        "nodejs_express": {
          "files": [
            {"url": "./components/error-handler/errorHandler.js", "target": "/middlewares/errorHandler.js"}
          ],
          "dependencies": [],
          "postInstall": "Add to app.js"
        }
      }
    },
    "jwt-auth": {
      "description": "JWT Authentication middleware",
      "category": "auth",
      "tags": ["jwt", "auth", "authentication"],
      "implementations": {
        "nodejs_express": {
          "files": [
            {"url": "./components/jwt-auth/authMiddleware.js", "target": "/middlewares/authMiddleware.js"}
          ],
          "dependencies": ["jsonwebtoken"],
          "installCmd": "npm install jsonwebtoken",
          "requires": ["error-handler"]
        }
      }
    }
  }
}`

const v2ManifestJSON = `{
  "schemaVersion": "2.0.0",
  "version": "1.1.0",
  "languages": {
    "nodejs": {
      "name": "Node.js",
      "frameworks": {
        "express": {
          "name": "Express",
          "description": "Fast web framework",
          "minVersion": "4.18.0"
        }
      }
    },
    "golang": {
      "name": "Golang",
      "frameworks": {
        "gin": {
          "name": "Gin",
          "description": "HTTP web framework"
        }
      }
    }
  },
  "foundations": {
    "express-mvc": {
      "description": "Express.js MVC Architecture",
      "language": "nodejs",
      "framework": "express",
      "architecture": "mvc",
      "templateUrl": "./templates/nodejs_express.zip",
      "includes": ["tsconfig", "eslint"]
    },
    "gin-standard": {
      "description": "Gin Standard Layout",
      "language": "golang",
      "framework": "gin",
      "architecture": "mvc",
      "templateUrl": "./templates/golang_gin.zip"
    }
  },
  "components": {
    "error-handler": {
      "description": "Global Error Handler middleware",
      "category": "middleware",
      "tags": ["error", "exception"],
      "implementations": {
        "nodejs_express": {
          "files": [{"url": "./error-handler.js", "target": "/middlewares/errorHandler.js"}],
          "dependencies": [],
          "envVars": [
            {"key": "NODE_ENV", "default": "development", "description": "Node environment"}
          ]
        },
        "golang_gin": {
          "files": [{"url": "./error_handler.go", "target": "/middlewares/error_handler.go"}],
          "dependencies": []
        }
      }
    },
    "jwt-auth": {
      "description": "JWT Authentication",
      "category": "auth",
      "tags": ["jwt", "auth"],
      "implementations": {
        "nodejs_express": {
          "files": [{"url": "./auth.js", "target": "/middlewares/auth.js"}],
          "dependencies": ["jsonwebtoken"],
          "devDependencies": ["@types/jsonwebtoken"],
          "installCmd": "npm install jsonwebtoken",
          "installDevCmd": "npm install -D @types/jsonwebtoken",
          "envVars": [
            {"key": "JWT_SECRET", "default": "secret", "description": "JWT secret key"}
          ]
        }
      }
    }
  },
  "schemas": {
    "schema-user-auth": {
      "description": "User model for authentication",
      "category": "schema",
      "tags": ["user", "auth", "database"],
      "implementations": {
        "nodejs_express": {
          "files": [{"url": "./user.prisma", "target": "/prisma/models/User.prisma"}],
          "dependencies": ["@prisma/client"]
        }
      }
    }
  },
  "blueprints": {
    "blueprint-auth": {
      "description": "Complete auth flow (register + login + JWT)",
      "category": "blueprint",
      "tags": ["auth", "login", "register"],
      "includes": ["jwt-auth", "schema-user-auth"],
      "implementations": {
        "nodejs_express": {
          "files": [
            {"url": "./authController.js", "target": "/controllers/authController.js"},
            {"url": "./authRoutes.js", "target": "/routes/authRoutes.js"}
          ],
          "dependencies": [],
          "postInstall": "Add to app.js: app.use('/api/auth', authRoutes)"
        }
      }
    }
  }
}`

// ──────────────────────────────────────────────────────────────
// Tests: parseManifest — V1 format
// ──────────────────────────────────────────────────────────────

func TestParseManifestV1(t *testing.T) {
	manifest, err := parseManifest([]byte(v1ManifestJSON))
	if err != nil {
		t.Fatal("parseManifest failed:", err)
	}

	// Version
	if manifest.Version != "1.0.0" {
		t.Errorf("Version = %q, want %q", manifest.Version, "1.0.0")
	}

	// SchemaVersion should be empty for v1
	if manifest.SchemaVersion != "" {
		t.Errorf("SchemaVersion = %q, want empty for v1", manifest.SchemaVersion)
	}

	// IsV2 should be false
	if manifest.IsV2() {
		t.Error("IsV2() should return false for v1 manifest")
	}

	// Languages
	if len(manifest.Languages) != 1 {
		t.Fatalf("Languages count = %d, want 1", len(manifest.Languages))
	}
	nodejs := manifest.Languages["nodejs"]
	if nodejs.Name != "Node.js" {
		t.Errorf("nodejs.Name = %q, want %q", nodejs.Name, "Node.js")
	}

	// Components
	if len(manifest.Components) != 2 {
		t.Fatalf("Components count = %d, want 2", len(manifest.Components))
	}
	eh := manifest.Components["error-handler"]
	if eh.Description != "Global Error Handler middleware" {
		t.Errorf("error-handler.Description = %q", eh.Description)
	}
	if eh.Category != "middleware" {
		t.Errorf("error-handler.Category = %q, want %q", eh.Category, "middleware")
	}

	// V2 collections should be initialized (not nil) even for v1
	if manifest.Foundations == nil {
		t.Error("Foundations should be initialized (not nil)")
	}
	if manifest.Schemas == nil {
		t.Error("Schemas should be initialized (not nil)")
	}
	if manifest.Blueprints == nil {
		t.Error("Blueprints should be initialized (not nil)")
	}

	// V2 collections should be empty for v1
	if len(manifest.Foundations) != 0 {
		t.Errorf("Foundations count = %d, want 0 for v1", len(manifest.Foundations))
	}
}

// ──────────────────────────────────────────────────────────────
// Tests: parseManifest — V2 format
// ──────────────────────────────────────────────────────────────

func TestParseManifestV2(t *testing.T) {
	manifest, err := parseManifest([]byte(v2ManifestJSON))
	if err != nil {
		t.Fatal("parseManifest failed:", err)
	}

	// Schema version
	if manifest.SchemaVersion != "2.0.0" {
		t.Errorf("SchemaVersion = %q, want %q", manifest.SchemaVersion, "2.0.0")
	}
	if !manifest.IsV2() {
		t.Error("IsV2() should return true for v2 manifest")
	}

	// Version
	if manifest.Version != "1.1.0" {
		t.Errorf("Version = %q, want %q", manifest.Version, "1.1.0")
	}

	// Languages
	if len(manifest.Languages) != 2 {
		t.Errorf("Languages count = %d, want 2", len(manifest.Languages))
	}

	// Foundations
	if !manifest.HasFoundations() {
		t.Error("HasFoundations() should return true")
	}
	if manifest.GetFoundationCount() != 2 {
		t.Errorf("GetFoundationCount() = %d, want 2", manifest.GetFoundationCount())
	}
	expressMVC := manifest.Foundations["express-mvc"]
	if expressMVC.Language != "nodejs" {
		t.Errorf("express-mvc.Language = %q, want %q", expressMVC.Language, "nodejs")
	}
	if expressMVC.Architecture != "mvc" {
		t.Errorf("express-mvc.Architecture = %q, want %q", expressMVC.Architecture, "mvc")
	}
	if len(expressMVC.Includes) != 2 {
		t.Errorf("express-mvc.Includes count = %d, want 2", len(expressMVC.Includes))
	}

	// Components
	if manifest.GetComponentCount() != 2 {
		t.Errorf("GetComponentCount() = %d, want 2", manifest.GetComponentCount())
	}

	// Component with envVars (V2 feature)
	jwtAuth := manifest.Components["jwt-auth"]
	nodeImpl := jwtAuth.Implementations["nodejs_express"]
	if len(nodeImpl.EnvVars) != 1 {
		t.Fatalf("jwt-auth envVars count = %d, want 1", len(nodeImpl.EnvVars))
	}
	if nodeImpl.EnvVars[0].Key != "JWT_SECRET" {
		t.Errorf("envVars[0].Key = %q, want %q", nodeImpl.EnvVars[0].Key, "JWT_SECRET")
	}

	// DevDependencies (V2 feature)
	if len(nodeImpl.DevDependencies) != 1 {
		t.Fatalf("jwt-auth devDependencies count = %d, want 1", len(nodeImpl.DevDependencies))
	}
	if nodeImpl.DevDependencies[0] != "@types/jsonwebtoken" {
		t.Errorf("devDependencies[0] = %q", nodeImpl.DevDependencies[0])
	}
	if nodeImpl.InstallDevCmd != "npm install -D @types/jsonwebtoken" {
		t.Errorf("installDevCmd = %q", nodeImpl.InstallDevCmd)
	}

	// Schemas
	if !manifest.HasSchemas() {
		t.Error("HasSchemas() should return true")
	}
	if manifest.GetSchemaCount() != 1 {
		t.Errorf("GetSchemaCount() = %d, want 1", manifest.GetSchemaCount())
	}
	userSchema := manifest.Schemas["schema-user-auth"]
	if userSchema.Description != "User model for authentication" {
		t.Errorf("schema-user-auth.Description = %q", userSchema.Description)
	}
	if userSchema.Category != "schema" {
		t.Errorf("schema-user-auth.Category = %q, want %q", userSchema.Category, "schema")
	}

	// Blueprints
	if !manifest.HasBlueprints() {
		t.Error("HasBlueprints() should return true")
	}
	if manifest.GetBlueprintCount() != 1 {
		t.Errorf("GetBlueprintCount() = %d, want 1", manifest.GetBlueprintCount())
	}
	authBP := manifest.Blueprints["blueprint-auth"]
	if authBP.Description != "Complete auth flow (register + login + JWT)" {
		t.Errorf("blueprint-auth.Description = %q", authBP.Description)
	}
	if len(authBP.Includes) != 2 {
		t.Fatalf("blueprint-auth.Includes count = %d, want 2", len(authBP.Includes))
	}
}

// ──────────────────────────────────────────────────────────────
// Tests: ResolveItem (cross-tier lookup)
// ──────────────────────────────────────────────────────────────

func TestResolveItem(t *testing.T) {
	manifest, _ := parseManifest([]byte(v2ManifestJSON))

	tests := []struct {
		name      string
		wantTier  string
		wantFound bool
	}{
		{"error-handler", "component", true},
		{"jwt-auth", "component", true},
		{"schema-user-auth", "schema", true},
		{"blueprint-auth", "blueprint", true},
		{"nonexistent", "", false},
		{"", "", false},
	}

	for _, tt := range tests {
		tier, found := manifest.ResolveItem(tt.name)
		if found != tt.wantFound {
			t.Errorf("ResolveItem(%q) found = %v, want %v", tt.name, found, tt.wantFound)
		}
		if tier != tt.wantTier {
			t.Errorf("ResolveItem(%q) tier = %q, want %q", tt.name, tier, tt.wantTier)
		}
	}
}

// ──────────────────────────────────────────────────────────────
// Tests: GetAllItemNames
// ──────────────────────────────────────────────────────────────

func TestGetAllItemNames(t *testing.T) {
	manifest, _ := parseManifest([]byte(v2ManifestJSON))

	names := manifest.GetAllItemNames()
	// 2 components + 1 schema + 1 blueprint = 4
	if len(names) != 4 {
		t.Errorf("GetAllItemNames() count = %d, want 4", len(names))
	}

	// Verify all expected names are present
	expected := map[string]bool{
		"error-handler":    false,
		"jwt-auth":         false,
		"schema-user-auth": false,
		"blueprint-auth":   false,
	}
	for _, name := range names {
		expected[name] = true
	}
	for name, found := range expected {
		if !found {
			t.Errorf("GetAllItemNames() missing %q", name)
		}
	}
}

// ──────────────────────────────────────────────────────────────
// Tests: GetFoundationsForStack
// ──────────────────────────────────────────────────────────────

func TestGetFoundationsForStack(t *testing.T) {
	manifest, _ := parseManifest([]byte(v2ManifestJSON))

	// Match: nodejs + express
	nodeFoundations := manifest.GetFoundationsForStack("nodejs", "express")
	if len(nodeFoundations) != 1 {
		t.Errorf("GetFoundationsForStack(nodejs, express) count = %d, want 1", len(nodeFoundations))
	}
	if _, ok := nodeFoundations["express-mvc"]; !ok {
		t.Error("Expected express-mvc foundation")
	}

	// Match: golang + gin
	goFoundations := manifest.GetFoundationsForStack("golang", "gin")
	if len(goFoundations) != 1 {
		t.Errorf("GetFoundationsForStack(golang, gin) count = %d, want 1", len(goFoundations))
	}

	// No match: csharp + dotnet-webapi
	csFoundations := manifest.GetFoundationsForStack("csharp", "dotnet-webapi")
	if len(csFoundations) != 0 {
		t.Errorf("GetFoundationsForStack(csharp, dotnet-webapi) count = %d, want 0", len(csFoundations))
	}
}

// ──────────────────────────────────────────────────────────────
// Tests: Language/Framework helpers
// ──────────────────────────────────────────────────────────────

func TestGetLanguageKeys(t *testing.T) {
	manifest, _ := parseManifest([]byte(v2ManifestJSON))

	keys := manifest.GetLanguageKeys()
	if len(keys) != 2 {
		t.Errorf("GetLanguageKeys() count = %d, want 2", len(keys))
	}
}

func TestGetFrameworkKeys(t *testing.T) {
	manifest, _ := parseManifest([]byte(v2ManifestJSON))

	keys := manifest.GetFrameworkKeys("nodejs")
	if len(keys) != 1 {
		t.Fatalf("GetFrameworkKeys(nodejs) count = %d, want 1", len(keys))
	}
	if keys[0] != "express" {
		t.Errorf("GetFrameworkKeys(nodejs)[0] = %q, want %q", keys[0], "express")
	}

	// Non-existent language
	keys = manifest.GetFrameworkKeys("python")
	if keys != nil {
		t.Errorf("GetFrameworkKeys(python) should return nil, got %v", keys)
	}
}

func TestGetLanguageNames(t *testing.T) {
	manifest, _ := parseManifest([]byte(v2ManifestJSON))

	names := manifest.GetLanguageNames()
	if len(names) != 2 {
		t.Errorf("GetLanguageNames() count = %d, want 2", len(names))
	}
}

func TestGetFrameworkNames(t *testing.T) {
	manifest, _ := parseManifest([]byte(v2ManifestJSON))

	names := manifest.GetFrameworkNames("golang")
	if len(names) != 1 {
		t.Fatalf("GetFrameworkNames(golang) count = %d, want 1", len(names))
	}
	if names[0] != "Gin" {
		t.Errorf("GetFrameworkNames(golang)[0] = %q, want %q", names[0], "Gin")
	}
}

// ──────────────────────────────────────────────────────────────
// Tests: parseManifest — Error cases
// ──────────────────────────────────────────────────────────────

func TestParseManifest_InvalidJSON(t *testing.T) {
	_, err := parseManifest([]byte(`{invalid json}`))
	if err == nil {
		t.Error("parseManifest should return error for invalid JSON")
	}
}

func TestParseManifest_EmptyJSON(t *testing.T) {
	manifest, err := parseManifest([]byte(`{}`))
	if err != nil {
		t.Fatal("parseManifest should handle empty JSON:", err)
	}

	// All maps should be initialized
	if manifest.Components == nil {
		t.Error("Components should be initialized")
	}
	if manifest.Foundations == nil {
		t.Error("Foundations should be initialized")
	}
	if manifest.Schemas == nil {
		t.Error("Schemas should be initialized")
	}
	if manifest.Blueprints == nil {
		t.Error("Blueprints should be initialized")
	}
}

// ──────────────────────────────────────────────────────────────
// Tests: JSON Serialization Roundtrip
// ──────────────────────────────────────────────────────────────

func TestManifestRoundtrip(t *testing.T) {
	// Parse V2 manifest
	original, err := parseManifest([]byte(v2ManifestJSON))
	if err != nil {
		t.Fatal("parseManifest failed:", err)
	}

	// Serialize back to JSON
	data, err := json.MarshalIndent(original, "", "  ")
	if err != nil {
		t.Fatal("json.Marshal failed:", err)
	}

	// Parse again
	roundtrip, err := parseManifest(data)
	if err != nil {
		t.Fatal("parseManifest roundtrip failed:", err)
	}

	// Verify key fields match
	if roundtrip.SchemaVersion != original.SchemaVersion {
		t.Errorf("SchemaVersion mismatch: %q vs %q", roundtrip.SchemaVersion, original.SchemaVersion)
	}
	if roundtrip.Version != original.Version {
		t.Errorf("Version mismatch: %q vs %q", roundtrip.Version, original.Version)
	}
	if len(roundtrip.Components) != len(original.Components) {
		t.Errorf("Components count mismatch: %d vs %d", len(roundtrip.Components), len(original.Components))
	}
	if len(roundtrip.Foundations) != len(original.Foundations) {
		t.Errorf("Foundations count mismatch: %d vs %d", len(roundtrip.Foundations), len(original.Foundations))
	}
	if len(roundtrip.Schemas) != len(original.Schemas) {
		t.Errorf("Schemas count mismatch: %d vs %d", len(roundtrip.Schemas), len(original.Schemas))
	}
	if len(roundtrip.Blueprints) != len(original.Blueprints) {
		t.Errorf("Blueprints count mismatch: %d vs %d", len(roundtrip.Blueprints), len(original.Blueprints))
	}
}
