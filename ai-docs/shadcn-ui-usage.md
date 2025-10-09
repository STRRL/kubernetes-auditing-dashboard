# shadcn/ui Usage Guide

This document describes how shadcn/ui is configured and used in the Kubernetes Auditing Dashboard project.

## Overview

shadcn/ui is a collection of re-usable components built with Radix UI and Tailwind CSS. Unlike traditional component libraries, shadcn/ui components are copied directly into your project, giving you full ownership and control over the code.

## Project Setup

### Configuration Files

The project uses the following configuration for shadcn/ui:

#### `ui/components.json`
```json
{
  "$schema": "https://ui.shadcn.com/schema.json",
  "style": "new-york",
  "rsc": false,
  "tsx": true,
  "tailwind": {
    "config": "tailwind.config.js",
    "css": "styles/globals.css",
    "baseColor": "neutral",
    "cssVariables": true,
    "prefix": ""
  },
  "iconLibrary": "lucide",
  "aliases": {
    "components": "@/components",
    "utils": "@/lib/utils",
    "ui": "@/components/ui",
    "lib": "@/lib",
    "hooks": "@/hooks"
  }
}
```

- **Style**: `new-york` - Uses the New York style variant (recommended by shadcn)
- **RSC**: `false` - Not using React Server Components
- **Icon Library**: `lucide-react` - For icons in components

#### Path Aliases

The following import aliases are configured:
- `@/components` → `ui/components`
- `@/lib` → `ui/lib`
- `@/components/ui` → `ui/components/ui`
- `@/hooks` → `ui/hooks`

These are defined in both `components.json` and `tsconfig.json`.

## Adding New Components

### Using the CLI (Recommended)

To add a new shadcn/ui component to the project:

```bash
cd ui
npx shadcn@latest add [component-name]
```

**Examples:**

```bash
npx shadcn@latest add button
npx shadcn@latest add card
npx shadcn@latest add dialog
npx shadcn@latest add dropdown-menu

npx shadcn@latest add button card table sheet
```

This command will:
1. Download the component code from the shadcn/ui registry
2. Install any required dependencies (e.g., Radix UI primitives)
3. Place the component in `ui/components/ui/[component-name].tsx`
4. Update the component with the project's configuration

### Currently Installed Components

The following shadcn/ui components are currently installed:

- **Button** (`ui/components/ui/button.tsx`) - For actions and navigation
- **Card** (`ui/components/ui/card.tsx`) - For displaying content in a card layout
- **Table** (`ui/components/ui/table.tsx`) - For displaying tabular data
- **Sheet** (`ui/components/ui/sheet.tsx`) - For side panels and drawers

## Using Components

### Import Syntax

Always import components from the `@/components/ui` path:

```typescript
import { Button } from "@/components/ui/button"
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card"
import { Table, TableHeader, TableBody, TableRow, TableCell } from "@/components/ui/table"
import { Sheet, SheetTrigger, SheetContent } from "@/components/ui/sheet"
```

### Example Usage

#### Button Component

```tsx
import { Button } from "@/components/ui/button"

export function MyComponent() {
  return (
    <div>
      <Button variant="default">Click me</Button>
      <Button variant="outline">Outline</Button>
      <Button variant="ghost">Ghost</Button>
      <Button variant="destructive">Delete</Button>
      <Button size="sm">Small</Button>
      <Button size="lg">Large</Button>
      <Button disabled>Disabled</Button>
    </div>
  )
}
```

#### Card Component

```tsx
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card"

export function StatCard() {
  return (
    <Card>
      <CardHeader>
        <CardDescription>Total Events</CardDescription>
        <CardTitle className="text-3xl">12,543</CardTitle>
      </CardHeader>
      <CardContent>
        <p>Additional content here</p>
      </CardContent>
    </Card>
  )
}
```

#### Table Component

```tsx
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from "@/components/ui/table"

export function DataTable() {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Name</TableHead>
          <TableHead>Status</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow>
          <TableCell>Item 1</TableCell>
          <TableCell>Active</TableCell>
        </TableRow>
      </TableBody>
    </Table>
  )
}
```

#### Sheet Component (Drawer)

```tsx
import { Sheet, SheetTrigger, SheetContent } from "@/components/ui/sheet"
import { Button } from "@/components/ui/button"

export function SideDrawer() {
  return (
    <Sheet>
      <SheetTrigger asChild>
        <Button>Open Menu</Button>
      </SheetTrigger>
      <SheetContent side="left">
        <nav>
          <ul>
            <li>Home</li>
            <li>Events</li>
          </ul>
        </nav>
      </SheetContent>
    </Sheet>
  )
}
```

## Utility Functions

### `cn()` Helper

The `cn()` function from `lib/utils.ts` is used to merge Tailwind CSS classes conditionally:

```typescript
import { cn } from "@/lib/utils"

<Button className={cn("my-custom-class", isActive && "bg-primary")} />
```

This function combines `clsx` and `tailwind-merge` to:
1. Handle conditional classes
2. Merge Tailwind classes correctly (avoiding conflicts)

## Styling and Theming

### CSS Variables

All component colors use CSS custom properties defined in `ui/styles/globals.css`:

```css
:root {
  --background: 0 0% 100%;
  --foreground: 0 0% 3.9%;
  --primary: 0 0% 9%;
  --primary-foreground: 0 0% 98%;
  /* ... more variables */
}
```

These variables can be customized to change the entire theme.

### Customizing Component Styles

Since components are copied into your project, you can directly modify them:

```typescript
<Button className="bg-blue-500 hover:bg-blue-600">
  Custom Blue Button
</Button>
```

Or modify the component file directly at `ui/components/ui/button.tsx`.

## Project-Specific Components

### Sidebar Component

A custom responsive sidebar component that wraps shadcn's Sheet component:

**Location**: `ui/components/Sidebar.tsx`

**Usage**:
```tsx
import { Sidebar } from "@/components/Sidebar"

export default function MyPage() {
  return (
    <Sidebar>
      <div className="p-4">
        <h1>Page Content</h1>
      </div>
    </Sidebar>
  )
}
```

Features:
- Desktop: Fixed sidebar on the left
- Mobile: Collapsible drawer using Sheet component
- Automatic route highlighting

## Dependencies

The following packages are required for shadcn/ui:

```json
{
  "@radix-ui/react-dialog": "^1.0.5",
  "@radix-ui/react-slot": "^1.0.2",
  "class-variance-authority": "^0.7.0",
  "clsx": "^2.0.0",
  "lucide-react": "^0.294.0",
  "tailwind-merge": "^2.0.0",
  "tailwindcss-animate": "^1.0.7"
}
```

These are automatically installed when you add components via the CLI.

## Best Practices

### 1. Always Use the CLI

Don't manually create component files. Always use:
```bash
npx shadcn@latest add [component-name]
```

This ensures:
- Correct dependencies are installed
- Components match the project configuration
- TypeScript types are properly set up

### 2. Import from `@/components/ui`

Always use the alias path:
```typescript
import { Button } from "@/components/ui/button"
```

Never use relative paths like `"../components/ui/button"`.

### 3. Use the `cn()` Helper

For conditional or merged classes:
```typescript
import { cn } from "@/lib/utils"

<div className={cn(
  "base-class",
  condition && "conditional-class",
  "another-class"
)} />
```

### 4. Customize Responsibly

While you can modify component code directly, document significant changes and consider:
- Will updates from shadcn/ui break your changes?
- Could this be handled with props or CSS classes instead?

### 5. Check the Official Docs

For component-specific APIs and variants, always refer to:
- https://ui.shadcn.com/docs/components/[component-name]

## Updating Components

To update a component to the latest version:

```bash
npx shadcn@latest add [component-name]
```

The CLI will ask if you want to overwrite the existing file. Review changes before accepting.

## Troubleshooting

### Import Errors

If you get import errors like "Module not found":
1. Check that the component is installed: `ls ui/components/ui/`
2. Verify path aliases in `tsconfig.json`
3. Restart your dev server

### Styling Issues

If components don't look right:
1. Verify `tailwind.config.js` includes shadcn configuration
2. Check that `styles/globals.css` has CSS variables defined
3. Ensure `tailwindcss-animate` plugin is installed

### Type Errors

If TypeScript complains:
1. Make sure component files are `.tsx` not `.jsx`
2. Check that `typescript` and `@types/react` are installed
3. Run `npm install` to ensure all dependencies are present

## Migration Notes

This project was migrated from DaisyUI to shadcn/ui. Key changes:

- **Drawer** → **Sheet**: DaisyUI's drawer replaced with Sheet component
- **Stats** → **Card**: Stats cards replaced with Card component
- **Table**: DaisyUI table replaced with shadcn Table
- **btn** → **Button**: Button classes replaced with Button component

All migrations maintain the same functionality with improved accessibility and customization.

## Resources

- **Official Docs**: https://ui.shadcn.com
- **Component Examples**: https://ui.shadcn.com/docs/components
- **Radix UI Docs**: https://www.radix-ui.com/primitives
- **Tailwind CSS**: https://tailwindcss.com/docs
