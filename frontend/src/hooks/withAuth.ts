import { auth } from "wailsjs/go/models";
import * as AppAPI from "../../wailsjs/go/main/App";
import { useAuthStore } from "@/store/useAuthStore";

export async function withAuth<T>(
  call: (jwtToken: string) => Promise<T>
): Promise<T> {
  const auth = useAuthStore.getState();

  if (!auth.jwtToken) {
    throw new Error("Not authenticated");
  }

  try {
    // First attempt
    return await call(auth.jwtToken);
  } catch (err: any) {
    // Only refresh on expiration
    const errorString =
      err?.message ||
      err?.error ||
      err?.toString?.() ||
      "";

    console.log("🚀 ~ withAuth ~ errorString:", errorString)

    if (!errorString.includes("expired token")) {
      throw err;
    }


    if (!auth.refreshToken || !auth.user?.id) {
      auth.clearAll();
      throw new Error("Session expired, please log in again");
    }
    const userId = auth.user!.id; // non-null assertion, safe because of prior check

    // 🔁 BACKEND RETURNS STRING (new access token)
    const newAccessToken: auth.TokenPairs = await AppAPI.RefreshToken(userId);

    // ✅ Store new access token
    auth.setJwtToken(newAccessToken.access_token);

    // ✅ Store new refresh token
    auth.setRefreshToken(newAccessToken.refresh_token);

    // 🔁 Retry with fresh token
    return await call(newAccessToken.access_token);
  }
}
