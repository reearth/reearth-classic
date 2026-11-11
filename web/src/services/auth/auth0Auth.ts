import { useAuth0 } from "@auth0/auth0-react";

import { e2eAccessToken, logOutFromTenant } from "@reearth/services/config";
import { clearSentinelToken } from "@reearth/services/sentinel";

import type { AuthHook } from "./authHook";

export const errorKey = "reeartherror";

export const useAuth0Auth = (): AuthHook => {
  const { isAuthenticated, error, isLoading, loginWithRedirect, logout, getAccessTokenSilently } =
    useAuth0();

  return {
    isAuthenticated: !!e2eAccessToken() || (isAuthenticated && !error),
    isLoading,
    error: error?.message ?? null,
    getAccessToken: () => getAccessTokenSilently(),
    login: () => {
      logOutFromTenant();
      return loginWithRedirect();
    },
    logout: () => {
      logOutFromTenant();
      clearSentinelToken();
      return logout({
        logoutParams: {
          returnTo: error
            ? `${window.location.origin}?${errorKey}=${encodeURIComponent(error?.message)}`
            : window.location.origin,
        },
      });
    },
  };
};
