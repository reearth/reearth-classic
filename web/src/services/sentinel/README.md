# Sentinel Asset Security Service

This service integrates the [@reearth/sentinel](https://github.com/reearth/sentinel) library to provide secure, authenticated access to protected assets through Bearer token authentication.

## Overview

Sentinel uses a Service Worker to automatically intercept network requests to protected domains and add Bearer token authentication headers. This enables secure asset access without manual header management in your application code.

## Configuration

Add the following environment variables to configure Sentinel:

```bash
# Required: The base URL of your tile server/proxy
REEARTH_WEB_TILE_SERVER_BASE_URL=https://your-tile-server.example.com

# Optional: Authentication token for the tile server
REEARTH_WEB_TILE_SERVER_TOKEN=your-bearer-token-here
```

These will be loaded into the config as `tileServerBaseUrl` and `tileServerToken`.

## Initialization

Sentinel is automatically initialized in `main.tsx` after the config is loaded:

```typescript
import { initializeSentinel } from "./services/sentinel";

// In main.tsx
loadConfig().finally(async () => {
  // ... other initialization
  await initializeSentinel();
  // ... render app
});
```

## Usage

### Updating the Token

When a user authenticates or you receive a new token:

```typescript
import { setSentinelToken } from "@reearth/services/sentinel";

// Update token with expiration
await setSentinelToken(
  "your-jwt-token",
  Date.now() + 3600000 // Expires in 1 hour
);
```

### Clearing the Token

When a user logs out:

```typescript
import { clearSentinelToken } from "@reearth/services/sentinel";

await clearSentinelToken();
```

### Checking Security Status

To check if Sentinel is active and authenticated:

```typescript
import { getSentinelStatus, isSentinelInitialized } from "@reearth/services/sentinel";

// Check if initialized
if (isSentinelInitialized()) {
  const status = await getSentinelStatus();
  console.log("Authenticated:", status.isAuthenticated);
  console.log("Service Worker Active:", status.isRegistered);
  console.log("Token Expires:", new Date(status.tokenExpiresAt || 0));
}
```

## How It Works

1. **Service Worker Registration**: The service worker (`sentinel-sw.js`) is registered on app initialization
2. **Request Interception**: Network requests to configured protected domains are intercepted
3. **Token Injection**: Bearer tokens are automatically added to request headers
4. **Caching**: Responses are cached according to configured strategies

## Protected Domains

By default, Sentinel protects requests to the hostname of the configured `tileServerBaseUrl`. You can customize this in `src/services/sentinel/index.ts`:

```typescript
const options: AssetSecurityConfig = {
  proxyUrl,
  protectedDomains: [
    new URL(proxyUrl).hostname,
    // Add more domains as needed
    "storage.googleapis.com",
    "your-cdn.example.com",
  ],
  // ...
};
```

## Debugging

Enable debug logging by setting `developerMode` in your config. Sentinel will output detailed logs about:
- Service worker registration
- Token updates
- Request interceptions
- Authentication status

## Service Worker File

The service worker file is automatically copied from the `@reearth/sentinel` package to `/public/sentinel-sw.js` during setup. This file must be accessible at `/sentinel-sw.js` in your deployed application.

## API Reference

### `initializeSentinel(): Promise<void>`
Initializes the Sentinel service worker with configuration from the app config.

### `setSentinelToken(accessToken: string, expiresAt?: number): Promise<void>`
Updates the Bearer token used for authenticated requests.
- `accessToken`: The JWT or Bearer token
- `expiresAt`: Optional expiration timestamp in milliseconds (defaults to 24 hours)

### `clearSentinelToken(): Promise<void>`
Clears the stored authentication token.

### `getSentinelStatus(): Promise<AssetSecurityStatus>`
Returns the current security status including authentication state and service worker status.

### `isSentinelInitialized(): boolean`
Returns whether Sentinel has been successfully initialized.
