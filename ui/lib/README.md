# Shared Utilities

## Verb Colors (`verb-colors.ts`)

Unified color mapping for Kubernetes audit event verbs across the application.

### Color Mapping

| Verb | Color | Use Case |
|------|-------|----------|
| `create` | `bg-green-500` | Resource creation |
| `update` | `bg-indigo-500` | Full resource update |
| `patch` | `bg-blue-500` | Partial resource update |
| `delete` | `bg-red-500` | Resource deletion |
| `deletecollection` | `bg-red-600` | Bulk deletion |
| `get` | `bg-gray-400` | Read single resource |
| `list` | `bg-gray-500` | List resources |
| `watch` | `bg-yellow-500` | Watch resource changes |
| *other* | `bg-purple-500` | Unknown/custom verbs |

### Usage

```typescript
import { getVerbColor, VERB_COLORS } from '@/lib/verb-colors';

// Get color for a verb
const color = getVerbColor('patch'); // returns 'bg-blue-500'

// Use the color map directly
<FilterSection colorMap={VERB_COLORS} />
```

### Design Principles

- **Mutating operations** (create, update, patch, delete) use distinct colors for easy identification
- **Read operations** (get, list, watch) use gray/yellow tones to differentiate from mutations
- **Destructive operations** (delete) use red shades
- **Case insensitive** - works with both lowercase (database) and uppercase (display) verbs
