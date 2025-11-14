# Hosting BT Phone Home Binaries

To enable curl-based installation, you need to host the binaries somewhere publicly accessible.

## Option 1: GitHub Releases (Recommended)

1. Create a GitHub repository (can be private with access token)
2. Build all platform binaries:
   ```bash
   # macOS Apple Silicon
   GOARCH=arm64 GOOS=darwin go build -o bt-phone-home-darwin-arm64

   # macOS Intel
   GOARCH=amd64 GOOS=darwin go build -o bt-phone-home-darwin-amd64

   # Linux
   GOARCH=amd64 GOOS=linux go build -o bt-phone-home-linux-amd64

   # Windows (64-bit)
   GOARCH=amd64 GOOS=windows go build -o bt-phone-home.exe

   # Windows (32-bit, optional)
   GOARCH=386 GOOS=windows go build -o bt-phone-home-windows-386.exe
   ```

3. Create a GitHub release and upload binaries
4. Get direct download URLs from the release
5. Update `install.sh` with your repository URL

## Option 2: GitHub Raw URLs (Simple)

1. Create a public GitHub repository
2. Create a `bin/` directory
3. Upload binaries to `bin/`:
   ```
   bin/
   ├── bt-phone-home-darwin-arm64
   ├── bt-phone-home-darwin-amd64
   ├── bt-phone-home-linux-amd64
   └── bt-phone-home.exe
   ```

4. Update installer scripts:

   In `install.sh`:
   ```bash
   BINARY_HOST="https://raw.githubusercontent.com/YOUR_USERNAME/bt-phone-home/main/bin"
   ```

   In `install.ps1`:
   ```powershell
   $BinaryHost = "https://raw.githubusercontent.com/YOUR_USERNAME/bt-phone-home/main/bin"
   ```

5. Upload both `install.sh` and `install.ps1` to repository root

6. Installation commands:

   **macOS/Linux:**
   ```bash
   curl -fsSL https://raw.githubusercontent.com/YOUR_USERNAME/bt-phone-home/main/install.sh | bash -s -- alice secret-room
   ```

   **Windows (PowerShell):**
   ```powershell
   irm https://raw.githubusercontent.com/YOUR_USERNAME/bt-phone-home/main/install.ps1 | iex
   ```
   Then run:
   ```powershell
   Install-BtPhoneHome -Username alice -RoomName secret-room
   ```

## Option 3: Vercel/Netlify Static Hosting

1. Create a `public/` directory:
   ```
   public/
   ├── install.sh
   └── bin/
       ├── bt-phone-home-darwin-arm64
       ├── bt-phone-home-darwin-amd64
       └── bt-phone-home-linux-amd64
   ```

2. Deploy to Vercel:
   ```bash
   vercel --prod
   ```

3. Update install script with your domain:
   ```bash
   BINARY_HOST="https://your-app.vercel.app/bin"
   ```

4. Installation command:
   ```bash
   curl -fsSL https://your-app.vercel.app/install.sh | bash -s -- alice secret-room
   ```

## Option 4: Cloudflare R2 / AWS S3

1. Upload binaries to bucket with public read access
2. Set CORS policy to allow downloads
3. Update install script with bucket URL
4. Installation works via direct S3/R2 URLs

## Option 5: Self-Hosted Simple Server

Create a simple static file server:

```python
# server.py
from http.server import SimpleHTTPServer, HTTPServer

class CORSRequestHandler(SimpleHTTPHandler):
    def end_headers(self):
        self.send_header('Access-Control-Allow-Origin', '*')
        SimpleHTTPHandler.end_headers(self)

httpd = HTTPServer(('0.0.0.0', 8000), CORSRequestHandler)
httpd.serve_forever()
```

```bash
python server.py
```

## Custom Installation Host

Users can override the default host:

```bash
export BT_PHONE_HOME_HOST=https://your-custom-host.com/bin
curl -fsSL https://your-custom-host.com/install.sh | bash -s -- alice secret-room
```

## Updating install.sh

After choosing a hosting method, update line 12 in `install.sh`:

```bash
BINARY_HOST="${BT_PHONE_HOME_HOST:-https://YOUR_ACTUAL_HOST/bin}"
```

## Security Considerations

- Use HTTPS for all binary downloads
- Consider signing binaries (codesign on macOS)
- Verify checksums in install script
- Use GitHub releases for version tracking
- Never commit credentials to public repos

## Quick Start with GitHub

1. Create public repo: `bt-phone-home`
2. Build all binaries:
   ```bash
   ./build-all.sh
   ```
3. Upload `bin/` directory with all binaries
4. Upload `install.sh` and `install.ps1` to root
5. Update both scripts with your GitHub URL
6. Share these commands:

   **macOS/Linux:**
   ```bash
   curl -fsSL https://raw.githubusercontent.com/USERNAME/bt-phone-home/main/install.sh | bash -s -- alice secret-room
   ```

   **Windows:**
   ```powershell
   irm https://raw.githubusercontent.com/USERNAME/bt-phone-home/main/install.ps1 | iex
   Install-BtPhoneHome -Username alice -RoomName secret-room
   ```

Replace `USERNAME` with your GitHub username.

## Build All Platforms

Use the included build script:

```bash
./build-all.sh
```

This creates:
- `bin/bt-phone-home-darwin-arm64` (macOS M1/M2)
- `bin/bt-phone-home-darwin-amd64` (macOS Intel)
- `bin/bt-phone-home-linux-amd64` (Linux)
- `bin/bt-phone-home.exe` (Windows 64-bit)
