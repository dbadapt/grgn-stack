# Using the GRGN Stack Template

Quick guide to get started with this template.

## Option 1: GitHub Template (Recommended)

1. Click **"Use this template"** on GitHub
2. Name your repository and create it
3. Clone your new repository
4. Run the initialization script:

   **Windows:**

   ```powershell
   .\init-template.ps1
   ```

   **Linux/Mac:**

   ```bash
   chmod +x init-template.sh
   ./init-template.sh
   ```

## Option 2: Clone and Reinitialize

```bash
git clone https://github.com/YOUR_USERNAME/grgn-stack.git my-project
cd my-project
rm -rf .git  # Remove template's git history

# Initialize for your project
./init-template.sh  # or .\init-template.ps1 on Windows
```

## What the Init Script Does

1. ✅ Updates Go module paths to your repository
2. ✅ Replaces environment variable prefixes
3. ✅ Updates project name throughout documentation
4. ✅ Configures README badges with your username
5. ✅ Updates LICENSE with your name
6. ✅ Creates `.env` files from examples
7. ✅ Initializes fresh Git repository
8. ✅ Optionally removes template-specific files

## After Initialization

1. **Install dependencies:**

   ```bash
   npm install
   cd web && npm install && cd ..
   ```

2. **Configure environment:**
   - Edit `.env` with your database credentials
   - Edit `web/.env` with frontend config

3. **Start development:**

   ```bash
   docker-compose -f docker-compose.yml -f docker-compose.dev.yml up
   ```

4. **Verify everything works:**
   - Frontend: http://localhost:5173
   - Backend: http://localhost:8080/graphql
   - Neo4j: http://localhost:7474

## Next Steps

- 📖 Read [QUICK-REFERENCE.md](../development/QUICK-REFERENCE.md) for common commands
- 🏗️ Review [ARCHITECTURE.md](../architecture/ARCHITECTURE.md) for system overview
- 📊 Design your schema with [GRAPHQL.md](../development/GRAPHQL.md)
- 🧪 Set up testing with [TESTING-CI.md](../testing/TESTING-CI.md)

## Need Help?

- Check the [documentation table of contents](../_TOC.md)
- Review [TEMPLATE-SETUP.md](./TEMPLATE-SETUP.md) for detailed setup instructions
