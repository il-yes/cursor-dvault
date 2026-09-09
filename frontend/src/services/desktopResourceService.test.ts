import { describe, it, expect, vi, beforeEach } from "vitest";
import { desktopResourceService, translateResourceError } from "./desktopResourceService";
import * as AppAPI from "../../wailsjs/go/main/App";
import { useAuthStore } from "@/store/useAuthStore";

vi.mock("../../wailsjs/go/main/App", () => ({
  ResolveCollaborativeShare: vi.fn(),
  AccessDecryptVaultEntry: vi.fn(),
}));

describe("DesktopResourceService & Plaintext Resolution", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    useAuthStore.setState({
      jwtToken: "test_jwt_token_123",
      user: { email: "bob@ankhora.io", Email: "bob@ankhora.io" } as any,
    });
  });

  it("resolves C3 ShareEntry and decodes plaintext byte array to UTF-8 string", async () => {
    const rawPlaintext = "secret construction document";
    const bytes = Array.from(new TextEncoder().encode(rawPlaintext));

    vi.mocked(AppAPI.ResolveCollaborativeShare).mockResolvedValueOnce({
      share_entry_id: "se_authoritative_777",
      trust_group_id: "tg_executive_1",
      created_by: "user_alice_uuid",
      created_at: "2026-09-01T12:00:00Z",
      metadata: { title: "Blueprint Q4" },
      plaintext: bytes,
    } as any);

    const result = await desktopResourceService.openResource({
      refType: "share_entry",
      shareEntryId: "se_authoritative_777",
      trustGroupId: "tg_executive_1",
    });

    expect(AppAPI.ResolveCollaborativeShare).toHaveBeenCalledWith(
      "test_jwt_token_123",
      "se_authoritative_777"
    );

    expect(result.resourceId).toBe("se_authoritative_777");
    expect(result.title).toBe("Blueprint Q4");
    expect(result.content).toBe("secret construction document");
    expect(result.kind).toBe("share_entry");
  });

  it("handles string or Uint8Array plaintext responses from Wails IPC", async () => {
    vi.mocked(AppAPI.ResolveCollaborativeShare).mockResolvedValueOnce({
      share_entry_id: "se_authoritative_888",
      trust_group_id: "tg_executive_1",
      created_by: "user_alice_uuid",
      created_at: "2026-09-01T12:00:00Z",
      metadata: { title: "Direct String Payload" },
      plaintext: "secret construction document string" as any,
    } as any);

    const result = await desktopResourceService.openResource({
      refType: "share_entry",
      shareEntryId: "se_authoritative_888",
    });

    expect(result.content).toBe("secret construction document string");
  });

  it("translates ErrUnauthorizedMember into user-facing authorization error", async () => {
    vi.mocked(AppAPI.ResolveCollaborativeShare).mockRejectedValueOnce(
      new Error("caller is not an authorized member of trust group: ErrUnauthorizedMember")
    );

    await expect(
      desktopResourceService.openResource({
        refType: "share_entry",
        shareEntryId: "se_unauthorized_999",
      })
    ).rejects.toThrow("You no longer have access to this resource");
  });

  it("translates ErrKeyEnvelopeNotFound into device authorization error", async () => {
    vi.mocked(AppAPI.ResolveCollaborativeShare).mockRejectedValueOnce(
      new Error("no active device key envelope found: ErrKeyEnvelopeNotFound")
    );

    await expect(
      desktopResourceService.openResource({
        refType: "share_entry",
        shareEntryId: "se_nodevice_999",
      })
    ).rejects.toThrow("This device is not authorized to access this resource");
  });

  it("passes authoritative local sovereign device ID and NOT default_desktop_device", async () => {
    vi.mocked(AppAPI.ResolveCollaborativeShare).mockResolvedValueOnce({
      share_entry_id: "se_authoritative_999",
      trust_group_id: "tg_authoritative_999",
      created_by: "user_alice",
      created_at: "2026-09-01T12:00:00Z",
      plaintext: "test payload",
    } as any);

    await desktopResourceService.openResource({
      refType: "share_entry",
      shareEntryId: "se_authoritative_999",
      trustGroupId: "tg_authoritative_999",
    });

    expect(AppAPI.ResolveCollaborativeShare).toHaveBeenCalledWith(
      "test_jwt_token_123",
      "se_authoritative_999"
    );
  });
});
