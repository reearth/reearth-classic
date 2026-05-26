import {
  registerAssetSecurity,
  unregisterAssetSecurity,
  updateToken,
  clearToken,
  getSecurityStatus,
  type AssetSecurityConfig,
} from "@reearth/sentinel";

import { config } from "../config";

let initialized = false;

/**
 * Wait for the service worker to take control of the page
 * On hard reload, the SW may be active but not controlling requests
 * This ensures requests can be intercepted before we set tokens
 */
async function waitForServiceWorkerControl(timeoutMs = 5000): Promise<void> {
  console.log("[Sentinel] Waiting for service worker to control page...");

  if (navigator.serviceWorker?.controller) {
    console.log("[Sentinel] Service worker already controlling the page");
    return;
  }

  const registration = await navigator.serviceWorker?.getRegistration();
  if (registration?.active) {
    registration.active.postMessage({ type: "CLAIM_CLIENTS" });
    console.log("[Sentinel] Sent CLAIM_CLIENTS message to service worker");
  }

  return new Promise<void>(resolve => {
    const onControllerChange = () => {
      clearTimeout(timeout);
      console.log("[Sentinel] Service worker gained control");
      resolve();
    };

    const timeout = setTimeout(() => {
      navigator.serviceWorker?.removeEventListener("controllerchange", onControllerChange);
      if (navigator.serviceWorker?.controller) {
        console.log("[Sentinel] Service worker control detected");
        resolve();
      } else {
        console.warn("[Sentinel] Service worker did not gain control after", timeoutMs, "ms");
        resolve();
      }
    }, timeoutMs);

    navigator.serviceWorker?.addEventListener("controllerchange", onControllerChange, {
      once: true,
    });
  });
}

/**
 * Initialize the Sentinel asset security service worker
 * This sets up automatic Bearer token authentication for protected asset domains
 */
export async function initializeSentinel(): Promise<void> {
  if (initialized) return;

  const cfg = config();
  const proxyUrl = cfg?.tileServerBaseUrl;
  const token = cfg?.tileServerToken;

  if (!proxyUrl) return;

  try {
    const options: AssetSecurityConfig = {
      proxyUrl,
      protectedDomains: [
        // Add protected domains that require authentication
        // These will be intercepted and have Bearer token added
        new URL(proxyUrl).hostname,
      ],
      namespace: "reearth-classic",
      serviceWorkerPath: "/sentinel-sw.js",
      scope: "/",
      debug: false,
      // Asset patterns for different tile types
      assetPatterns: {
        rasterTiles: /\/(tiles|imagery)\/[^/]+\/\d+\/\d+\/\d+/,
      },
      // Token configuration with caching settings
      tokenConfig: {
        memoryCacheTTL: 5 * 60 * 1000, // 5 minutes
        refreshThreshold: 60 * 1000, // 1 minute
      },
      // Cache strategies for different asset types
      cacheStrategies: {
        tiles: "network-first",
        images: "cache-first",
      },
      onTokenExpired: async () => {
        // Re-fetch config to get fresh token
        try {
          const response = await fetch("/reearth_config.json");
          if (response.ok) {
            const freshConfig = await response.json();
            if (freshConfig?.tileServerToken) {
              await updateSentinelToken(freshConfig.tileServerToken);
              return;
            }
          }
        } catch (error) {
          console.error("Sentinel: Failed to refresh config", error);
        }
        // Fallback: clear token if refresh failed
        await clearSentinelToken();
      },
    };

    const result = await registerAssetSecurity(options);

    if (!result.success) {
      console.error("Sentinel: Failed to register", result.error);
      return;
    }

    // Wait for service worker to take control before setting token
    // This prevents 401 errors on hard reload
    await waitForServiceWorkerControl();

    initialized = true;

    // If token is provided in config, set it immediately
    if (token) {
      await updateSentinelToken(token);
    }
  } catch (error) {
    console.error("Sentinel: Failed to initialize", error);
  }
}

/**
 * Update the Bearer token used for authenticated asset requests
 * @param accessToken - The JWT or Bearer token to use
 * @param expiresAt - Optional expiration timestamp in milliseconds
 */
export async function updateSentinelToken(accessToken: string, expiresAt?: number): Promise<void> {
  if (!initialized) return;

  try {
    await updateToken({
      accessToken,
      expiresAt: expiresAt || Date.now() + 24 * 60 * 60 * 1000, // Default to 24 hours
    });
  } catch (error) {
    console.error("Sentinel: Failed to update token", error);
  }
}

/**
 * Clear the stored authentication token
 * Call this on logout or when the token is no longer valid
 */
export async function clearSentinelToken(): Promise<void> {
  if (!initialized) return;

  try {
    await clearToken();
  } catch (error) {
    console.error("Sentinel: Failed to clear token", error);
  }
}

/**
 * Get the current security status
 * Returns information about authentication state and service worker status
 */
export async function getSentinelStatus() {
  if (!initialized) {
    return {
      isAuthenticated: false,
      isRegistered: false,
      tokenExpiresAt: 0,
    };
  }

  try {
    return await getSecurityStatus();
  } catch (error) {
    console.error("Sentinel: Failed to get security status", error);
    return {
      isAuthenticated: false,
      isRegistered: false,
      tokenExpiresAt: 0,
    };
  }
}

/**
 * Unregister the Sentinel service worker
 * Call this when cleaning up or resetting the application
 */
export async function unregisterSentinel(): Promise<void> {
  if (!initialized) return;

  try {
    await unregisterAssetSecurity();
    initialized = false;
  } catch (error) {
    console.error("Sentinel: Failed to unregister", error);
  }
}

/**
 * Check if Sentinel has been initialized
 */
export function isSentinelInitialized(): boolean {
  return initialized;
}
