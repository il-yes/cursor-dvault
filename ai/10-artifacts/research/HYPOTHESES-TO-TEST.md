# XS-BIM — Hypotheses to Test

This document defines the core operational hypotheses to be validated or invalidated during stakeholder discovery meetings. These hypotheses focus on real construction operations, information flow, and cross-organizational friction, rather than software product features or technical claims.


---
Core discovery hypotheses
---

## H1 — Coordination friction

**Hypothesis**

There are recurring situations where several independent construction stakeholders need to coordinate information quickly, but the current process creates friction, delay, ambiguity, or unnecessary communication overhead.

**Why we believe it**

* **Documented observation:** In the XS-BIM scenario (and regional logistics research), off-normal events such as transport permit delays (`DOC-EVID-001`) or site access restrictions require rapid coordination between Apex Construction (GC), Apex Fabricators (Supplier), Transport Express (Logistics), Site Supervisor (Field), and Structural Engineer (Design). Currently, communication relies on phone calls, WhatsApp group chats, email threads, and paper logs.
* **Inference:** Because communications are distributed across unintegrated channels, message delivery is fragmented, status updates become stale instantly, and participants lack a shared real-time operational state.
* **Our own hypothesis:** Cross-organizational coordination friction may be caused less by a lack of communication channels and more by the absence of a shared, structured operational context when fast-moving disruptions occur.

**What would validate it**

The construction professional confirms specific, recurring operational events (e.g., fabricator delays, transport route changes, permit holds, site readiness mismatches) where multi-party coordination required extensive manual back-and-forth and caused operational confusion.

**What would weaken or invalidate it**

The professional reports that informal channels (WhatsApp/phone calls) handle multi-party changes smoothly without miscommunication, or that critical operational changes are managed cleanly bilaterally without involving multiple independent organizations.

**Meeting questions**

1. When an unexpected change occurs on site—such as a delivery delay or permit issue—what events create the most coordination friction?
2. Who needs to know immediately, and who currently communicates with whom across the different companies involved?
3. What happens when information arrives late, incomplete, or contradictory?

**Potential pilot signal**

An active construction phase involving at least three independent organizations (e.g., General Contractor, Steel Fabricator, Transport Carrier) where off-normal delivery or site logistics events occur on a weekly basis.

---

## H2 — Information latency has an operational consequence (Operational consequences)

**Hypothesis**

There are construction events where the speed at which reliable information reaches the right stakeholders has a measurable operational consequence.

**Why we believe it**

* **Documented observation:** In the delivery disruption scenario (`DEL-1042`), a 2-day delivery delay notice (`OFF-1042`) leads directly to crane scheduling conflicts, site crew idle time, and concrete pour rescheduling if not communicated immediately to site operations.
* **Inference:** Information latency forces site management to react after resources are already deployed, changing active operations from proactive scheduling to reactive disruption management.
* **Our own hypothesis:** Information latency may have a measurable operational consequence when project activities depend on timely coordination between stakeholders:
- downtime
- rescheduling
- rework
- missed delivery
- equipment idle time
- additional administrative work
- something else.

**What would validate it**

The construction professional identifies specific, recurring events where receiving information even a few hours or a day late directly caused waiting time, crew/equipment downtime, rework, or rescheduling.

**What would weaken or invalidate it**

The professional explains that project schedules have sufficient buffer, or site teams are flexible enough, such that information delays rarely alter what happens next operationally or generate downtime/rework.

**Meeting questions**

1. Can you walk me through a recent situation where information arrived late—such as a delivery delay or drawing change—and what happened on site as a result?
2. Which specific events cause crews or heavy equipment (like cranes) to stand idle while waiting for confirmation?
3. Is there a recurring event where faster, verified information would materially change your team's immediate next action?

**Potential pilot signal**

A project site with high daily equipment rental costs (e.g., mobile cranes) or tight sequential trade dependencies where schedule changes directly result in idle hours or rescheduling penalties.

---

## H3 — The problem crosses organizational boundaries (Cross-organizational boundary)

**Hypothesis**

Some important construction events involve several independent organizations, making coordination more difficult because information is distributed across companies, roles, tools, and communication channels.

**Why we believe it**

* **Documented observation:** Construction projects operate as temporary multi-organization networks. General contractors, specialized fabricators, third-party logistics providers, independent engineering consultants, and public inspectors all use distinct internal management systems (or no software systems at all).
* **Inference:** Information stops at company boundaries. Internal ERP/BIM updates inside a fabricator's system are not visible to the transport carrier or site contractor, creating blind spots across organizational seams.
* **Our own hypothesis:** Operational breakdown occurs predominantly at the boundaries *between* organizations rather than within any single company's internal workflow.

**What would validate it**

The professional states that the hardest coordination challenges occur when dealing with external sub-contractors, suppliers, or logistics companies who use different tools, formats, or communication habits.

**What would weaken or invalidate it**

The professional indicates that project management is highly vertically integrated, or that a single dominant software system (or main contractor mandate) effectively bridges all subcontractors and suppliers without friction.

**Meeting questions**

1. How do you currently share operational updates with external partners like fabricators, transport carriers, or engineering consultants?
2. Where do communication gaps or misunderstandings most frequently occur between your team and external organizations?
3. What happens when an external partner's internal schedule changes without immediate notification to your site team?

**Potential pilot signal**

A project environment characterized by multi-tiered subcontracting and independent logistics/fabrication suppliers where contractual boundaries obscure real-time status.

---

## H7 — A measurable pilot exists

**Hypothesis**

There is at least one recurring construction workflow where the problem is sufficiently concrete to establish a measurable baseline, intervention, and operational outcome.

**Why we believe it**

* **Documented observation:** Structural component delivery (e.g., steel girder delivery `DEL-1042`, inspection `INSP-1042`, and site crane erection) is a highly structured, repeatable workflow with clear handoffs and discrete stakeholders.
* **Inference:** Rather than attempting to transform an entire construction organization, a targeted pilot focused on a single high-impact recurring event can isolate data flow and prove operational value.
* **Our own hypothesis:** A successful discovery meeting will identify a specific, bounded workflow where baseline latency and downtime can be measured against live operational coordination.

**What would validate it**

The professional identifies a specific recurring workflow (e.g., off-site prefabricated element delivery, heavy equipment dispatch, material quality inspection sign-off) and agrees on how current baseline delays/downtime could be measured.

**What would weaken or invalidate it**

The professional asserts that construction workflows are too unique, unpredictable, or chaotic to isolate a single recurring process for structured measurement.

**Meeting questions**

1. Is there one specific, recurring workflow on your current project—like material deliveries, concrete pours, or specialized inspections—that causes disproportionate management headaches?
2. If we were to observe that specific workflow, what metric would tell us whether coordination improved (e.g., reduced crane idle time, faster RFI resolution, zero missed deliveries)?
3. What would a safe, low-risk test look like for your team on a single active project phase?

**Potential pilot signal**

The stakeholder explicitly offers to walk through an active project's upcoming material delivery or inspection schedule and identifies a measurable friction point suitable for observation.

---
Secondary hypotheses
---

## H4 — Evidence and traceability matter after the event

**Hypothesis**

There are situations where stakeholders need to reconstruct what happened, when it happened, who knew, what information was available, what decision was made, who approved it, and what action followed.

**Why we believe it**

* **Documented observation:** In the XS-BIM scenario, when DOT permit delay `DOC-EVID-001` occurs, the team creates an issue (`ISS-1042`), evaluates options, records decision `DEC-1042`, and issues site inspection `INSP-1042`. In conventional practice, such records are scattered across email attachments, WhatsApp screenshots, phone call logs, and verbal agreements.
* **Inference:** When delays, cost overruns, or quality disputes arise later, reconstructing the exact timeline of events and decisions requires forensic document gathering, leading to disputes and administrative burden.
* **Our own hypothesis:** Stakeholders may spend significant time and administrative effort reconstructing historical timelines and decision rationales after a disruption. (Then test it)

**What would validate it**

The professional describes painful experiences trying to settle delay claims, variation requests, or quality disputes where proving "who knew what and when" was difficult, time-consuming, or contested.

**What would weaken or invalidate it**

The professional reports that formal daily logs, RFI registers, and standard email archives already provide sufficient, undisputed traceability, or that post-event dispute reconstruction is rare and economically negligible.

**Meeting questions**

1. When a delay or dispute occurs on a project, how do you reconstruct what happened, when it was communicated, and who authorized the response?
2. How much time or effort does your team spend gathering evidence, emails, or site notes to justify schedule adjustments or cost variations?
3. Have you experienced situations where missing or informal records made it impossible to determine responsibility for an operational delay?

**Potential pilot signal**

A project with strict contractual delay damages or high variation claim frequency where post-event auditability directly impacts financial settlements.

---

## H5 — Controlled information sharing is relevant

**Hypothesis**

Construction projects contain information that should not simply be exposed to every participant, while still requiring collaboration between multiple stakeholders.

**Why we believe it**

* **Documented observation:** Fabricators and logistics providers are reluctant to expose full commercial cost structures, internal resource allocations, or proprietary mill certificates to all project participants, yet must share binding availability dates and delivery status.
* **Inference:** Fear of exposing sensitive commercial data or exposing liability leads companies to hoard information or delay sending updates until formal legal commitments are established.
* **Our own hypothesis:** Granular, sovereign permission control (sharing operational status without exposing private commercial data) is a prerequisite for getting independent companies to participate in live project data sharing.

**What would validate it**

The professional confirms that stakeholders withhold operational updates or resist centralized software tools specifically because they do not trust how their data will be shared, viewed, or used by other companies.

**What would weaken or invalidate it**

The professional states that project data is generally non-sensitive, that open transparency is already accepted across all project partners, or that existing document permissions fully satisfy privacy requirements.

**Meeting questions**

1. Are there operational or technical details that your team or your suppliers hesitate to share openly with all project participants?
2. How do you balance the need for transparency on site with protecting proprietary, financial, or commercial information?
3. Has concern over data ownership or confidentiality ever prevented project partners from adopting a shared tool or communication workflow?

**Potential pilot signal**

A joint-venture project or multi-contractor site where sensitive commercial boundaries currently inhibit open data sharing.

---

## H6 — Project information is lost as organizational memory

**Hypothesis**

Important operational knowledge remains scattered across emails, documents, messages, people, and project-management systems rather than becoming reusable project memory.

**Why we believe it**

* **Documented observation:** Post-project documentation is typically archived as unstructured PDF bundles or stored in isolated folder structures. When key personnel leave or a new project starts, operational lessons learned (e.g., fabricator lead times, transport permit bottlenecks) are forgotten.
* **Inference:** Organizations may repeatedly encounter similar operational problems because useful project experience is difficult to capture and reuse across projects.
* **Our own hypothesis:** Structuring project events, decisions, and outcomes into queryable provenance history enables organizations to reuse past experience to prevent recurring project failures.

**What would validate it**

The professional acknowledges that when experienced project managers leave, valuable operational knowledge is lost, or that teams frequently encounter known recurring issues on new projects without institutional context.

**What would weaken or invalidate it**

The professional demonstrates that their organization already maintains effective, structured knowledge-sharing databases, post-mortem repositories, or standardization templates that successfully capture lessons learned.

**Meeting questions**

1. What happens to the operational lessons and decision history of a project after it is completed and key team members move on?
2. How do new project teams access past operational data—such as supplier performance, permit lead times, or recurring site bottlenecks—when planning a new job?
3. How often do you see similar operational mistakes or coordination breakdowns repeated across different projects?

**Potential pilot signal**

A repeat developer or general contractor undertaking similar sequential projects (e.g., infrastructure, commercial builds) seeking to institutionalize operational history across project teams.

---
