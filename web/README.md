# kainet Web Terminal

browser-based version of the classified communications terminal.

## features

- retro CRT terminal aesthetic with scanlines and glow effects
- same AES-256-GCM encryption as CLI version
- works in any modern browser
- no installation required
- real-time message polling
- military-grade interface

## development

```bash
# install dependencies
bun install

# run dev server
bun run dev

# build for production
bun run build

# preview production build
bun run preview
```

## usage

1. open the app in browser
2. enter your username when prompted
3. type messages and press enter
4. use `/burn` to wipe history

## deployment

deploy the `dist` folder to any static hosting:

- vercel: `vercel dist`
- netlify: drag and drop `dist` folder
- github pages: push `dist` to gh-pages branch
- cloudflare pages: connect repo and set build dir to `dist`

## security

- end-to-end encryption in browser using Web Crypto API
- credentials embedded in build (regenerate for production)
- messages encrypted before leaving browser
- database only sees encrypted ciphertext

## tech stack

- vite - build tool
- vanilla js - no framework overhead
- @libsql/client - turso database client
- web crypto api - AES-256-GCM encryption
- custom CSS - retro CRT terminal effects
