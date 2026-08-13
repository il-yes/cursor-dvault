# Frontend C3 Architectural Rules

## Purpose

These rules define how the desktop frontend must interact with the C3 collaboration layer.

The C3 frontend is a consumer of the C3 collaboration substrate. It must not turn C3 into a sector-specific business application.

---

# 1. C3 Is the Collaboration Substrate

> **C3 is the collaboration substrate; it is not the sector business domain.**

C3 provides generic collaboration capabilities that can be used by multiple business domains.

C3 concerns itself with collaboration primitives such as:

- Workspace
- Channel
- Thread
- ThreadEvent
- Asset / PayloadRef
- ShareEntry
- TrustGroup
- Federation

These concepts must remain domain-neutral.

Examples of domains that may use C3:

- construction
- pharmaceutical
- engineering
- supply chain
- banking
- manufacturing
- logistics
- other operational domains

The C3 frontend must not contain the semantics of those domains.

---

# 2. No Sector-Specific Logic Inside C3

Do not introduce sector-specific concepts into:

```text
frontend/src/components/C3/


The desktop is a Wails application, so the frontend should not behave like a browser client calling the Cloud Backend REST API directly.

The correct boundary is:

Wails React frontend
        │
        ▼
      AppAPI
   (generated Wails bindings)
        │
        ▼
   Go App methods
        │
        ▼
   Cloud Backend API

Not:

React
  │
  └── fetch()
        │
        ▼
POST /c3/...

Your existing pattern is the canonical one:

export const uploadAttachementToIPFSWithEncryption = async (
    jwtToken: string,
    fileData: number[],
    password: string
): Promise<string> => {
    const response =
        await AppAPI.UploadAttachmentToIPFSWithEncryption(
            jwtToken,
            fileData,
            password
        );

    return response;
};

So Thread must follow the same architecture, and I would explicitly tell the engineer that this is a hard constraint.