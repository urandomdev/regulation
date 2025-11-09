# Kash Out Frontend

React + TypeScript frontend for the Kash Out personal finance app ("Who Are You To Spend?").

## Tech Stack

- **React 19** - UI framework
- **TypeScript** - Type safety
- **Vite** - Build tool and dev server
- **React Router 7** - Client-side routing
- **@urandomdev/regulation** - Generated SDK for API communication

## Getting Started

### Prerequisites

- Node.js 20+
- pnpm 9+

### Development

From the monorepo root:

```bash
# Install dependencies (run from root)
pnpm install

# Start dev server
pnpm dev:frontend

# Or from this directory
cd frontend
pnpm dev
```

The app will be available at `http://localhost:5173`

### Environment Variables

Copy `.env.example` to `.env.development`:

```bash
cp .env.example .env.development
```

Available variables:

- `VITE_API_URL` - API endpoint URL (default: `http://localhost:8080`)

## Project Structure

```
frontend/
├── src/
│   ├── lib/
│   │   └── api.ts          # SDK client configuration
│   ├── pages/
│   │   ├── Home.tsx        # Landing page
│   │   ├── Login.tsx       # Login page
│   │   ├── Signup.tsx      # Sign up page
│   │   ├── Dashboard.tsx   # User dashboard
│   │   └── NotFound.tsx    # 404 page
│   ├── components/         # Reusable components
│   ├── App.tsx             # Main app with routes
│   ├── main.tsx            # Entry point
│   └── index.css           # Global styles
├── public/                 # Static assets
└── package.json
```

## Available Scripts

- `pnpm dev` - Start development server
- `pnpm build` - Build for production
- `pnpm preview` - Preview production build
- `pnpm lint` - Run ESLint

## Using the SDK

The SDK client is pre-configured in `src/lib/api.ts`:

```typescript
import { api } from './lib/api'

// Login example
const [user, error] = await api.account.login({
  email: 'user@example.com',
  password: 'password'
})

// List rules example
const [rules, error] = await api.rule.listRules()
```

The SDK uses the Result pattern returning `[data, null]` on success or `[null, error]` on failure.

## Routes

- `/` - Home page
- `/login` - User login
- `/signup` - User registration
- `/dashboard` - User dashboard (requires auth)
- `*` - 404 Not Found

## Development Notes

- The app uses React Router v7 with declarative routing
- API calls are made using the generated TypeScript SDK
- Hot Module Replacement (HMR) is enabled for fast development
- TypeScript strict mode is enabled for type safety

## Building for Production

```bash
# From monorepo root
pnpm build:frontend

# Or from this directory
pnpm build
```

The built files will be in the `dist/` directory.
