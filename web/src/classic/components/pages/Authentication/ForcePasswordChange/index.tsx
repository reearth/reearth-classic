import { useEffect } from "react";

import Loading from "@reearth/classic/components/atoms/Loading";
import { getAuthInfo } from "@reearth/services/config";

const STORAGE_KEY = "reearth_password_change_state";

const ForcePasswordChangePage: React.FC = () => {
  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    const ticket = params.get("ticket");
    const state = params.get("state");

    if (ticket && state) {
      // Coming from Auth0 Post-Login Action redirect.
      // Save state to sessionStorage so we can resume the login flow after password change.
      sessionStorage.setItem(STORAGE_KEY, state);
      // Redirect to Auth0's password change page.
      window.location.href = ticket;
      return;
    }

    // Coming back from Auth0 password change page (result_url redirect).
    const savedState = sessionStorage.getItem(STORAGE_KEY);
    if (savedState) {
      sessionStorage.removeItem(STORAGE_KEY);
      const authInfo = getAuthInfo();
      if (authInfo?.auth0Domain) {
        // Redirect back to Auth0 to complete the login flow.
        window.location.href = `https://${authInfo.auth0Domain}/continue?state=${savedState}`;
        return;
      }
    }

    // Fallback: redirect to top page.
    window.location.href = "/";
  }, []);

  return <Loading />;
};

export default ForcePasswordChangePage;
