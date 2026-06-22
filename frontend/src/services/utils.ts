export function formatMonthYear(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleDateString("en-US", {
    month: "long",
    year: "numeric",
  });
}

export function parseNotificationPayload(payload: unknown) {
  if (!payload) return null;

  if (typeof payload === "object") return payload;

  if (typeof payload === "string") {
    try {
      const json = atob(payload);
      return JSON.parse(json);
    } catch (err) {
      console.error("Failed to parse notification payload", err);
      return null;
    }
  }

  return null;
}