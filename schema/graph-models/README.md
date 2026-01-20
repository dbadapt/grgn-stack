# Graph Models

This directory contains visual graph model designs created with [Arrows.app](https://arrows.app).

## 🚀 Quick Start

1. **Open a model in Arrows.app:**
   - Go to https://arrows.app
   - Click "Open" → "From file"
   - Select one of the JSON files below

2. **Edit and save:**
   - Modify nodes, relationships, and properties
   - Export as JSON (overwrites the file)
   - Commit to git

3. **Tell Copilot:**
   - "I've updated [filename].json, please implement the changes"
   - Copilot reads the JSON and generates code

## 📁 Current Models

### ✅ Created (Starter Templates)

- **`core-model.json`** - Core User and QRCode relationships
  - Nodes: User, QRCode
  - Relationships: User -[:OWNS]-> QRCode
  - Status: Ready for import and editing

- **`auth-model.json`** - Authentication and MFA
  - Nodes: User, AuthProvider, MFAProvider
  - Relationships: AUTHENTICATED_BY, HAS_MFA
  - Supports: Google, Apple, SAML, local auth
  - Status: Ready for import and editing

- **`order-model.json`** - Order and payment flow
  - Nodes: User, Order, QRCode, PaymentProvider
  - Relationships: PLACED, CONTAINS, PAYMENT_VIA
  - Supports: Stripe, PayPal, Google Pay, etc.
  - Status: Ready for import and editing

- **`scavenger-hunt-model.json`** - Scavenger hunt application
  - Nodes: User, ScavengerHunt, Checkpoint, QRCode
  - Relationships: CREATED, PARTICIPATING_IN, HAS_CHECKPOINT, SCANNED
  - Features: Progress tracking, points system
  - Status: Ready for import and editing

### ⏳ To Be Created

- `inventory-model.json` - Inventory tracking application
- `reservation-model.json` - Reservation system
- `analytics-model.json` - Analytics and reporting nodes
- `notification-model.json` - Notification system

## 🎨 How to Use These Models

### Import into Arrows.app

1. Go to https://arrows.app
2. Click **Open** → **From file**
3. Select a `.json` file from this directory
4. Model loads with all nodes and relationships

### Edit a Model

1. Click on nodes to edit properties
2. Add/remove properties using the property panel
3. Create new relationships by dragging between nodes
4. Change colors, positions, and labels
5. **Export** → **Save to file** (overwrite the JSON)

### Create a New Model

1. Design in Arrows.app
2. Export as JSON
3. Save to this directory with descriptive name
4. Update this README with model description
5. Tell Copilot to implement it

## 🤖 Working with Copilot

### After Editing a Model

Tell Copilot:
```
I've updated [filename].json with [description of changes].
Please:
1. Read the model
2. Update DATABASE.md
3. Create/update migration files
4. Update GraphQL schema
5. Generate repository code
```

### Example Prompts

**Simple update:**
```
I added MFA support to auth-model.json. Please implement the MFAProvider node
and HAS_MFA relationship.
```

**Complex feature:**
```
I've created inventory-model.json for the inventory tracking application.
It has InventoryItem, Location, and StockMovement nodes. Please implement:
- Migration files for all new nodes/relationships
- GraphQL schema with queries for inventory
- Repository methods for tracking stock
- Resolvers with business logic
```

**Just review:**
```
I'm working on order-model.json. Can you review the current structure and
suggest improvements before I finalize it?
```

## 📐 Design Guidelines

### Node Names
- Use **PascalCase**: `ScavengerHunt`, `AuthProvider`
- Be specific: `MFAProvider` not `Provider`
- Singular form: `User` not `Users`

### Relationship Types
- Use **UPPERCASE_SNAKE_CASE**: `AUTHENTICATED_BY`, `HAS_CHECKPOINT`
- Be descriptive: `OWNS` not `HAS`
- Action-oriented: `PLACED` (past tense), `PARTICIPATING_IN` (present)

### Properties
- Include type hints: `"points": "Integer"`, `"email": "String (unique)"`
- Mark constraints: `"id": "UUID (unique)"`, `"email": "String (required)"`
- Use meaningful names: `totalAmount` not `amt`

### Colors
Use colors to group related nodes:
- 🟢 **Green (#68BC00)**: User/Account nodes
- 🔵 **Blue (#4C8EDA)**: Core business entities (QRCode)
- 🔴 **Red (#FB7E81)**: Authentication/Security
- 🟣 **Purple (#9B59B6)**: Applications (Scavenger Hunt)
- 🟠 **Orange (#E67E22)**: Supporting entities (Checkpoints)
- 🟤 **Brown (#8B5A99)**: Transactions (Orders)

## 📋 Model Review Checklist

Before asking Copilot to implement:

- [ ] All nodes have meaningful labels
- [ ] Properties include type information
- [ ] Unique/required constraints are marked
- [ ] Relationships have descriptive types
- [ ] Relationship directions make sense
- [ ] Colors help group related concepts
- [ ] JSON exports cleanly from Arrows.app
- [ ] File saved with descriptive name

## 🔄 Workflow

```
┌─────────────────┐
│  You: Design in │
│   Arrows.app    │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Export JSON    │
│  to this dir    │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Tell Copilot to │
│   implement     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Copilot:        │
│ - Updates docs  │
│ - Creates code  │
│ - Makes PRs     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Review & Test  │
└─────────────────┘
```

## 🔗 Related Documentation

- [SCHEMA-WORKFLOW.md](../../SCHEMA-WORKFLOW.md) - Complete workflow guide
- [DATABASE.md](../../DATABASE.md) - Full Neo4j schema documentation
- [GRAPHQL.md](../../GRAPHQL.md) - GraphQL implementation guide
- [Arrows.app](https://arrows.app) - Visual graph modeling tool

## 💡 Tips

- **Save often** - Export JSON frequently to avoid losing work
- **Commit changes** - Version control your models
- **Start simple** - Begin with core entities, add complexity later
- **Review together** - Show models in PRs for team feedback
- **Test in Neo4j** - Import models to Neo4j browser to visualize
- **Keep focused** - One domain per model (auth, orders, etc.)

## 🆘 Troubleshooting

### Model won't import to Arrows.app
- Ensure JSON is valid (use a JSON validator)
- Check file encoding is UTF-8
- Verify no syntax errors in the file

### Copilot doesn't understand model
- Check JSON has proper node/relationship structure
- Describe changes in plain English
- Point to similar existing implementations
- Ask Copilot to "read [filename].json first"

### Model too complex
- Break into smaller focused models
- Use separate files for different domains
- Link related models through documentation
