<div align="center">
  <h1>⚡ Forge CLI (CoreForge)</h1>
  <p><strong>Universal Backend CLI — Scaffolding for any stack, inspired by <a href="https://ui.shadcn.com/">shadcn/ui</a></strong></p>
  <p><i>Copy, don't install. You own the code.</i></p>
</div>

---

## 🎯 What is this?

**Forge CLI** is a powerful command-line tool that lets you add production-ready boilerplate code to your backend project with a single command.

Unlike traditional libraries, the code is **copied directly into your project** — you own it, you can modify it, and there is **no vendor lock-in**.

Need a JWT Auth middleware? A robust Error Handler? A Swagger docs generator?
Just `forge add` it.

## 🚀 Supported Stacks

| Language    | Framework    | Status   |
| :---------- | :----------- | :------- |
| **Node.js** | Express      | ✅ Ready |
| **C#**      | .NET Web API | ✅ Ready |
| **Golang**  | Gin          | ✅ Ready |

## 📦 How It Works

Forge CLI uses a 4-tier architecture to scale from tiny utilities to full workflows:

1. **Foundations**: Base project setups (`forge init`)
2. **Components**: Single middlewares or utilities (`forge add error-handler`)
3. **Schemas**: Database models (`forge add schema-user`)
4. **Blueprints**: Full workflows combining multiple components (`forge add blueprint-auth`)

## 🛠️ Installation

```bash
# Build from source
git clone https://github.com/longgoll/coreforge.git
cd coreforge
go build -o forge .

# Or install globally (Requires Go installed)
go install github.com/longgoll/forge-cli@latest
```

## 📖 Usage

### 1. Initialize a Project

Set up a new project or recognize an existing one. This creates a `.forge.json` file.

```bash
forge init
```

### 2. Add a Component

Browse available components and add them to your project.

```bash
# Search for components
forge search auth

# Add an error handler
forge add error-handler

# Add a full auth blueprint (JWT, Password hashing, routes)
forge add blueprint-auth
```

### 3. Manage Components

```bash
# List installed components
forge list --installed

# Remove a component
forge remove error-handler
```

### 4. Verify Environment

```bash
forge doctor
```

## 🔗 Architecture & Registry

Forge CLI pulls components from the **[Forge Registry](https://github.com/longgoll/forge-registry)**. The registry acts as the "menu" of all available code snippets, automatically keeping your CLI updated without needing to download a new binary.

## 🤝 Contributing (We are in BETA!)

CoreForge is currently in **v1.0.0-beta**. This is a **100% Free & Open Source** passion project.
There is absolutely no strict licensing, no paywalls, and no vendor lock-in!

We built this to solve our own headaches as backend developers, and we want to hear from you.

- Found a bug? 🐛 **Open an Issue**
- Have an idea for a component? 💡 **Start a Discussion**
- Want to add your own stack or blueprint? 🛠️ **Submit a PR!**

Let's build the ultimate "shadcn for backend" together!

## 📄 License

Do whatever you want with the code. It's yours now! (MIT License) © [longgoll](https://github.com/longgoll)
