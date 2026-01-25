# Schema Design Quick Reference

> **TL;DR:** You design visually in Arrows.app, save JSON files, tell Copilot what you did, and Copilot generates all the code.

---

## ⚡ Quick Workflow

```
1. Open https://arrows.app
2. Import existing model or create new
3. Design nodes and relationships
4. Export → Save JSON to services/{domain}/{app}/model/
5. Tell Copilot: "I've updated X model, implement it"
6. Copilot generates code across all layers
```

---

## 🎯 Your Tools

### Visual Design: Arrows.app

- **URL:** https://arrows.app
- **Import:** Open → From file → select `.json` from `services/{domain}/{app}/model/`
- **Export:** Save → Download JSON → overwrite file

### Starter Models (Ready to Import!)

- ✅ `services/core/shared/model/core-model.json` - User & basic entities
- ✅ `services/core/auth/model/auth-model.json` - Auth providers & MFA

> 💡 **Tip:** Add your own models to your app's `model/` directory as you design them!

---

## 🤖 Copilot's Tools

When you tell Copilot to implement:

1. **Reads** your JSON model from `services/{domain}/{app}/model/`
2. **Updates** DATABASE.md documentation
3. **Generates** GraphQL schema in `services/{domain}/{app}/model/*.graphql`
4. **Creates** migration files in `services/{domain}/{app}/migrations/` (or root `migrations/` for core)
5. **Implements** repositories in `services/{domain}/{app}/controller/generated/`
6. **Writes** resolvers in `services/{domain}/{app}/controller/resolver.go`

---

## 💬 Example Conversations

### Simple Change

**You:**

```
I added a "verified" property to the User node in core-model.json.
Please add it to the schema.
```

**Copilot will:**

- Read core-model.json
- Add `verified: Boolean` to GraphQL User type
- Create migration to add property to Neo4j
- Update repository methods

### New Feature

**You:**

```
I've designed a complete inventory system in inventory-model.json.
It has InventoryItem, Location, and StockMovement nodes with
relationships for tracking item locations and movements.

Please implement:
1. Full GraphQL schema with queries and mutations
2. All database migrations
3. Repository layer
4. Resolvers with business logic
```

**Copilot will:**

- Read inventory-model.json
- Create comprehensive implementation
- Generate tests
- Update documentation

### Just Review

**You:**

```
I'm working on payment-flow-model.json but not sure if the
relationships are correct. Can you review the structure?
```

**Copilot will:**

- Read the model
- Provide feedback
- Suggest improvements
- Highlight potential issues

---

## 🎨 Design Patterns

### Node (Entity)

```json
{
  "caption": "User",
  "labels": ["User"],
  "properties": {
    "id": "UUID (unique)",
    "email": "String (unique, required)",
    "name": "String (optional)"
  }
}
```

### Relationship (Connection)

```json
{
  "type": "OWNS",
  "properties": {
    "purchasedAt": "DateTime"
  },
  "fromId": "n0",
  "toId": "n1"
}
```

---

## ✅ Design Checklist

Quick checklist before telling Copilot to implement:

- [ ] Node names are PascalCase
- [ ] Relationship types are UPPERCASE_SNAKE_CASE
- [ ] Properties have type annotations
- [ ] Unique/required fields are marked
- [ ] Relationships have clear directions
- [ ] JSON exports without errors
- [ ] File saved with descriptive name

---

## 🔍 Common Patterns

### One-to-Many

```
(User)-[:OWNS]->(QRCode)
User can own multiple QR codes
```

### Many-to-Many

```
(User)-[:PARTICIPATING_IN]->(ScavengerHunt)
Users can join multiple hunts, hunts have multiple users
```

### Hierarchical

```
(ScavengerHunt)-[:HAS_CHECKPOINT {order: Integer}]->(Checkpoint)
Hunt has ordered checkpoints
```

### Tracking/Progress

```
(User)-[:SCANNED {scannedAt: DateTime, points: Integer}]->(Checkpoint)
User progress through checkpoints
```

---

## 📂 File Structure

```
services/
├── {domain}/
│   └── {app}/
│       ├── model/               ← Visual models & GraphQL schemas
│       │   ├── {app}-model.json
│       │   └── types.graphql
│       ├── view/                ← React components & UI (web/)
│       ├── migrations/          ← Schema changes (Copilot)
│       └── controller/          ← API logic & repositories (Copilot)
```

---

## 🚀 Next Steps

### 1. Try It Now!

```bash
# Open Arrows.app in browser
start https://arrows.app

# Import core model
# Open → From file → services/core/shared/model/core-model.json
```

### 2. Make a Simple Edit

- Add a property to User node
- Change a color
- Export JSON back to file

### 3. Tell Copilot

```
I added a "phoneNumber" property to User in core-model.json.
Please update the schema.
```

### 4. Review Generated Code

- Check schema/schema.graphql
- Look at migration file
- Test with `npm run generate`

---

## 📚 Full Documentation

- **[SCHEMA-WORKFLOW.md](./SCHEMA-WORKFLOW.md)** - Complete workflow guide
- **[schema/graph-models/README.md](./schema/graph-models/README.md)** - Model directory guide
- **[DATABASE.md](./DATABASE.md)** - Full Neo4j schema docs
- **[GRAPHQL.md](./GRAPHQL.md)** - GraphQL setup guide

---

## 🆘 Need Help?

### Model Issues

```
"Arrows.app won't import my JSON"
→ Check JSON is valid, UTF-8 encoded
```

### Copilot Confusion

```
"Copilot doesn't understand my model"
→ Describe changes in plain English
→ Ask "read [filename].json first"
```

### Schema Sync

```
"GraphQL and Neo4j are out of sync"
→ Ask "reconcile schemas by reading DATABASE.md and schema.graphql"
```

---

## 🎯 Pro Tips

1. **Save Often** - Export JSON frequently from Arrows.app
2. **Commit Early** - Version control your visual models
3. **Start Simple** - Core entities first, complexity later
4. **Use Colors** - Group related nodes visually
5. **Be Specific** - Good: "AUTHENTICATED_BY", Bad: "HAS"
6. **Test Locally** - Import models to Neo4j browser
7. **Review Together** - Share models in PRs

---

## ⚡ Power User Commands

### View Schema in Neo4j

```cypher
CALL db.schema.visualization()
```

### Generate GraphQL Code

```bash
npm run generate
```

### Run Migration

```bash
./bin/migrate
```

### Test Everything

```bash
npm run test:ci
```

---

**Remember:** You design, Copilot codes. Focus on the **what** (visual models), let Copilot handle the **how** (implementation).
