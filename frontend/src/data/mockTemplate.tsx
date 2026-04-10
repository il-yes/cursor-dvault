

export const devsecops_incident_v1 = {
    template_id: "devsecops.incident.v1",
    record_type: "incident",
    schema_version: 1,

    incident_id: "INC-2026-0007",
    title: "Suspicious OAuth token usage detected",
    severity: "high", // low | medium | high | critical
    status: "investigating", // detected | investigating | contained | eradicated | recovered | closed

    detected_at: "2026-03-28T17:12:00.000Z",
    reported_by: "SIEM",
    environment: "production", // prod | staging | dev
    affected_systems: ["api-gateway", "auth-service", "billing-web"],

    summary: "Multiple refresh token exchanges from new ASN + impossible travel pattern.",
    customer_impact: "Potential unauthorized access to 3 customer accounts (unconfirmed).",

    indicators: {
        ip_addresses: ["203.0.113.10", "198.51.100.23"],
        domains: ["evil-example[.]com"],
        hashes: [],
        asn: ["AS12345"]
    },

    timeline: [
        { at: "2026-03-28T17:12:00.000Z", event: "Alert triggered by SIEM rule: oauth_anomaly_v3" },
        { at: "2026-03-28T17:20:00.000Z", event: "On-call acknowledged; began triage" }
    ],

    containment_actions: [
        "Revoked suspected tokens",
        "Forced password reset for impacted users",
        "Blocked IPs at WAF"
    ],

    root_cause_hypothesis: "Stolen refresh token via compromised endpoint or leaked client secret.",
    remediation_tasks: [
        { owner: "Manuel", task: "Rotate OAuth client secrets", due: "2026-03-29" },
        { owner: "Manuel", task: "Add device binding for refresh tokens", due: "2026-04-05" }
    ],

    evidence_links: [
        { label: "SIEM alert", url: "https://example.local/siem/alerts/123" },
        { label: "Auth logs", url: "https://example.local/logs/auth?query=INC-2026-0007" }
    ],

    postmortem: {
        completed: false,
        lessons_learned: "",
        followups: []
    }
}
export const legal_matter_v1 = {
    template_id: "legal.matter.v1",
    record_type: "matter",
    schema_version: 1,

    matter_id: "MAT-2026-0012",
    matter_type: "contract_review", // contract_review | dispute | compliance | employment | ip | other
    status: "open", // intake | open | on_hold | closed

    client_name: "Acme Mining Ltd",
    counterparty: "VendorCo International",
    jurisdiction: "Kenya",
    governing_law: "Kenya",
    forum: "Nairobi",

    opened_at: "2026-03-20",
    assigned_to: "Manuel",
    priority: "medium", // low | medium | high

    summary: "Review MSA + DPA for sovereign data handling and auditability requirements.",
    key_issues: [
        "Data residency and cross-border transfer clauses",
        "Audit rights and breach notification timelines",
        "Subprocessor disclosure + approval mechanism"
    ],

    key_dates: {
        signature_deadline: "2026-04-02",
        renewal_date: "2027-04-02"
    },

    documents: [
        { name: "MSA_v3.pdf", classification: "confidential" },
        { name: "DPA_draft.docx", classification: "confidential" }
    ],

    contacts: {
        client: [
            { name: "Jane Doe", role: "GC", email: "jane@example.com" }
        ],
        counterparty: [
            { name: "John Smith", role: "Legal", email: "john@vendorco.com" }
        ]
    },

    risk_assessment: {
        overall: "medium", // low | medium | high
        notes: "Audit rights acceptable; subprocessor clause needs tightening."
    },

    next_actions: [
        { owner: "Manuel", task: "Propose redlines for subprocessor clause", due: "2026-03-30" },
        { owner: "Manuel", task: "Confirm breach notification SLA (<=72h)", due: "2026-03-30" }
    ],

    billing: {
        fee_type: "fixed", // fixed | hourly | pro_bono
        amount_usd: 1500,
        invoice_status: "not_invoiced" // not_invoiced | invoiced | paid
    }
}