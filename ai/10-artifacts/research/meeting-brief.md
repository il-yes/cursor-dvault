# XS-BIM — Meeting Brief

## 1. Objective

We are exploring a construction project intelligence and coordination operating model to determine whether there is a recurring construction coordination problem involving multiple stakeholders that creates a measurable operational consequence, and identify whether a bounded workflow could be tested as a pilot. This is an operational discovery conversation, not a software sales presentation.

---

## 2. What I need to learn

1. Where does multi-party coordination currently break down during project disruptions?
2. Which specific information delays create measurable operational consequences (e.g., crew downtime, equipment idle time, rework, rescheduling)?
3. Is there a concrete, recurring workflow on an active project that could serve as a targeted pilot?

---

## 3. Opening

> "Thanks for taking the time today. I’m currently researching how information flows across independent companies on construction projects—specifically when unexpected disruptions occur, like delivery delays, permit holds, or site changes.
>
> We’ve built an early structural concept and scenario demonstration called XS-BIM to explore how multi-party coordination could work across organizational boundaries. But the concept is still actively being validated.
>
> I'm not here to sell you software or pitch a product. What I really need is your candid, real-world experience: how your teams actually handle these situations today, where communication breaks down, and whether information latency actually impacts your operations. If our assumptions are wrong or off-target, I want to know."

---

## 4. Core questions

1. Can you walk me through a recent situation where something unexpected changed on a project and several independent organizations had to react?
2. Who needed to know immediately, and how did that information actually travel between the parties?
3. Where in that communication chain did information slow down, get lost, or become ambiguous?
4. What happened operationally on site because of that delay (e.g., idle crews, crane waiting time, rescheduled trades)?
5. When a disruption or claim happens, how do you reconstruct what happened, when it was communicated, and who authorized the decision?
6. If you could eliminate friction from one recurring multi-party interaction on your projects, which one would it be?
7. What existing workarounds or tools do your teams rely on today to bridge gaps between contractors, fabricators, and site supervisors?

---

## 5. Demo sequence

Walk through the existing XS-BIM structural steel scenario (`PRJ-001`) as an operating model demonstration:

1. **Existing project:** Display `PRJ-001` (Commercial Hub Project) with baseline schedule and multi-party structure.
2. **Material / delivery requirement:** Show structural steel girder requirement `REQ-STRUCT-001` linked to fabricator `SUP-001`.
3. **Disruption occurs:** Introduce 2-day transit delay notification `OFF-1042` triggered by DOT bridge clearance route restriction `DOC-EVID-001`.
4. **Information/evidence appears:** Show the official permit document and delay event surfacing in real time.
5. **Relevant stakeholders need to react:** Highlight affected roles—General Contractor, Site Supervisor, Transport Carrier, Structural Engineer.
6. **Stakeholders coordinate:** Open channel `CH-1042` where transport carrier, engineer, and GC exchange verified operational context.
7. **Decision is made:** Record formal decision `DEC-1042` (Re-sequence crane lift to Section B, adjust delivery window to Thursday AM).
8. **Action is executed:** Issue updated site inspection `INSP-1042` and revise trade schedule without site downtime.
9. **Event is closed:** Mark issue `ISS-1042` resolved with all cryptographic provenance hashes intact.
10. **History remains reconstructable:** Demonstrate complete historical audit trail of the event, evidence, decision, and approval.

*Live creation:* Demonstrate creating a similar live disruption event to show how the operating model captures live operational state changes.

---

## 6. What to listen for

* [ ] recurring event
* [ ] multiple stakeholders
* [ ] independent organizations
* [ ] information delay
* [ ] unclear ownership/responsibility
* [ ] waiting/downtime
* [ ] rework/rescheduling
* [ ] difficult decision-making
* [ ] evidence/traceability problem
* [ ] sensitive information
* [ ] information lost after project completion
* [ ] measurable operational consequence
* [ ] existing workaround
* [ ] willingness to test a real workflow

---

## 7. The critical discovery

> **"Can you give me a real recent example where an information or coordination problem affected what happened on the project?"**

Follow the operational discovery chain:

$$\text{REAL EVENT} \rightarrow \text{WHO} \rightarrow \text{INFORMATION} \rightarrow \text{CURRENT PROCESS} \rightarrow \text{DELAY / CONSEQUENCE} \rightarrow \text{COST OR OPERATIONAL IMPACT} \rightarrow \text{CURRENT WORKAROUND} \rightarrow \text{POSSIBLE PILOT}$$

*Do not introduce the XS-BIM solution until the current process and its friction are fully understood.*

---

## 8. Pilot discovery

A suitable pilot must possess the following characteristics:

* **Real project:** An active or upcoming construction job.
* **Real recurring workflow:** Repeatable events (e.g., off-site component delivery, concrete pours, specialized structural inspections).
* **Identifiable stakeholders:** Clear participants (e.g., GC, specific fabricator, transport company, field engineer).
* **Concrete event:** High-frequency or high-impact operational trigger.
* **Observable existing process:** Clear current method (email, phone, paper) to establish a baseline.
* **Measurable baseline:** Historical downtime hours, RFI response times, or scheduling variance.
* **Measurable outcome:** Reduced idle equipment hours, faster issue resolution, or a measurable reduction in the operational problem identified during discovery.
* **Manageable scope:** Focused on one specific workflow rather than total enterprise IT transformation.
* **Sufficient activity:** generates enough relevant events during the pilot period to allow comparison with the baseline.

---

## 9. Objections / challenges

1. **"We already have tools for this (e.g., Procore, Primavera, BIM 360)."**
   * *Response:* "That makes sense—those are standard for document storage and master scheduling. What I'm curious about is when an unexpected event happens on site today, do those tools coordinate the immediate response across your external suppliers and logistics carriers in real time, or do people switch to phone and WhatsApp?"
2. **"This is just another platform."**
   * *Response:* "That's a valid concern—no one wants another app to log into. That's why we're exploring an underlying event coordination layer rather than replacing your main management systems. How many different tools do your external partners currently have to use?"
3. **"People won't change their workflow."**
   * *Response:* "I completely agree. If a tool requires field crews to change how they work on site, it usually fails. How do your field supervisors currently report delays or accept delivery changes?"
4. **"Our projects are too different."**
   * *Response:* "Every site has unique conditions. But are there recurring logistics or coordination handoffs—like material deliveries or inspection sign-offs—that happen on almost every project?"
5. **"We already use email/WhatsApp/Teams/etc."**
   * *Response:* "WhatsApp is fast and everyone has it. The challenge we often hear is that messages get buried and decisions lack formal traceability. How do you track down what was agreed upon when a dispute arises weeks later?"
6. **"The problem isn't significant enough."**
   * *Response:* "That's exactly what I want to establish. If information delays on your projects don't cause meaningful downtime, rework, or cost, then this isn't a problem worth solving for your organization."
7. **"The data is too sensitive."**
   * *Response:* "Protecting commercial boundaries is critical. What specific information would your suppliers or sub-contractors refuse to share on a joint platform?"
8. **"This would require too many stakeholders."**
   * *Response:* "That's why we focus on bounded two- or three-party workflows. Which two organizations on your project experience the most communication friction today?"

---

## 10. Meeting success criteria

The meeting is successful if we leave with:

1. At least one concrete real-world operational problem.
2. A specific event/workflow where the problem occurs.
3. The stakeholders involved.
4. The current information flow.
5. The operational consequence.
6. A possible measurable baseline.
7. A possible pilot scenario.
8. A clear reason to continue — or a clear reason to change/discard the hypothesis.

*Do NOT define success as:* getting someone to like the demo, getting someone to agree that XS-BIM is useful, getting a commitment to buy, or proving the technology.

---

## 11. Final note to myself

> I am not there to prove that XS-BIM is right.
>
> I am there to discover whether the problem is real, where it occurs, what it costs operationally, how people deal with it today, and whether there is a measurable workflow worth testing.
