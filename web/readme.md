# kaiNET Web Client

Mobile-responsive web client for kaiNET encrypted terminal chat.

## Features

- **AES-256-GCM Encryption**: End-to-end encryption using Web Crypto API
- **Room-Based Isolation**: Private conversations per room
- **Real-Time Polling**: 1-second message polling for instant updates
- **Retro CRT Terminal Aesthetic**: Green-on-black terminal with scanlines and glow effects
- **Mobile Responsive**: Optimized for both desktop and mobile devices
- **Turso libSQL Backend**: Same database as CLI version

## Mobile-Responsive Features

### Desktop (>768px)
- Full CRT effect with animated scanlines
- Green terminal glow effects
- Larger ASCII banner
- Desktop-optimized layout

### Mobile (<768px)
- **Touch-Friendly Interface**: 44px minimum touch targets
- **Responsive Typography**: Scaled font sizes for readability
- **Optimized Layout**: Full-screen terminal experience
- **Virtual Keyboard Friendly**: 16px inputs prevent iOS zoom
- **Safe Area Support**: iOS notch and home indicator spacing
- **Simplified CRT Effect**: Maintains green aesthetic without heavy effects
- **Landscape Support**: Optimized layout for landscape orientation
- **Smooth Scrolling**: iOS momentum scrolling for messages

### Responsive Breakpoints
- **Desktop**: >768px - Full CRT effects
- **Tablet/Mobile**: <768px - Simplified, touch-optimized
- **Small Mobile**: <480px - Further optimized typography
- **Landscape Mobile**: Special optimizations for landscape mode

## Installation

```bash
# Install dependencies with Bun
bun install

# Or with npm
npm install
```

## Development

```bash
# Start dev server (accessible from mobile on same network)
bun run dev

# Or with npm
npm run dev
```

Access the app at:
- Desktop: http://localhost:3000
- Mobile (same network): http://[your-ip]:3000

## Build for Production

```bash
# Build static files
bun run build

# Preview production build
bun run preview
```

## Usage

1. Enter your username
2. Enter a room name
3. Click CONNECT
4. Watch the boot sequence
5. Start chatting with encrypted messages

### Commands
- Type messages and press Enter to send
- Click /BURN to wipe room history
- Click DISCONNECT to return to login

## Architecture

- **Encryption**: AES-256-GCM via Web Crypto API
- **Key Derivation**: SHA-256 hash of auth token
- **Database**: Turso libSQL (same as CLI)
- **Polling**: 1-second interval for new messages
- **Framework**: Vanilla JS with Vite bundler

## Mobile Testing

To test on mobile devices:

1. Start dev server: `bun run dev`
2. Find your local IP: `ifconfig` (macOS/Linux) or `ipconfig` (Windows)
3. On mobile, navigate to: `http://[your-ip]:3000`
4. Ensure both devices are on the same network

## Security

- Messages encrypted before storage
- AES-256-GCM with unique nonce per message
- Key derived from auth token via SHA-256
- Room-based isolation
- No plaintext message storage

## Browser Compatibility

- Chrome/Edge: Full support
- Firefox: Full support
- Safari: Full support (iOS 11+)
- Mobile browsers: Optimized for all major browsers

## File Structure

```
web/
├── index.html           # Main HTML structure
├── src/
│   ├── main.js         # Application logic & Turso integration
│   ├── crypto.js       # Web Crypto API encryption
│   └── style.css       # Mobile-responsive CRT terminal styles
├── package.json        # Bun/npm configuration
├── vite.config.js      # Vite bundler configuration
└── .gitignore         # Git ignore rules
```
