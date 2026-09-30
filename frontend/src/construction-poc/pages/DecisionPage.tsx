import React from "react";
import { useNavigate } from "react-router-dom";
import { getDecisionData } from "../data";
import { CONSTRUCTION_ROUTES } from "../constants/routes";

/**
 * Construction Decision — DEC-1042.
 *
 * Composition follows the BuildFlow Stitch screen for this view. Every value
 * rendered here comes from `getDecisionData()`, which projects the canonical
 * XS-BIM ConstructionDecision aggregate. Slots the Stitch composition shows but
 * the scenario does not model (per-option impact, detour surcharge, commercial
 * term, per-participant sign-off) render an explicit unavailable state rather
 * than a plausible-looking stand-in.
 *
 * Colour and radius notes: `primary`, `secondary`, `error`, `surface` and
 * `background` resolve to the Sovereign Vault theme in index.css, so Stitch
 * values for those roles are written as arbitrary hex. Everything else uses the
 * BuildFlow tokens already declared in tailwind.config.ts. `rounded-lg` is
 * likewise overridden to --radius there, so Stitch radii are explicit.
 */

const UNAVAILABLE = "Not recorded in scenario";

/** Explicit unavailable marker for a slot the canonical scenario does not model. */
const Unavailable: React.FC<{ label?: string }> = ({ label = UNAVAILABLE }) => (
  <span
    className="text-[11px] leading-[14px] font-medium text-outline italic"
    title="The canonical XS-BIM scenario carries no value for this field."
  >
    {label}
  </span>
);

const Fact: React.FC<{ label: string; value?: string }> = ({ label, value }) => (
  <div className="bg-surface-container-low rounded-[4px] p-2.5 text-left">
    <span className="text-[11px] leading-[14px] font-medium text-on-surface-variant block">{label}</span>
    {value ? (
      <span className="text-sm leading-5 font-semibold text-[#041627]">{value}</span>
    ) : (
      <Unavailable />
    )}
  </div>
);

export const DecisionPage: React.FC = () => {
  const navigate = useNavigate();
  const data = getDecisionData();

  const optionCount = data.options.length;
  const alternativeLabel = optionCount === 1 ? "Alternative Assessed" : "Alternatives Assessed";

  return (
    <div className="flex flex-col w-full pb-12 max-w-4xl mx-auto px-6 pt-6 space-y-4">
      {/* ------------------------------------------------------------- */}
      <nav aria-label="Breadcrumb" className="flex items-center gap-2 text-[11px] leading-[14px] text-on-surface-variant">
        <button
          onClick={() => navigate(CONSTRUCTION_ROUTES.PROJECTS)}
          className="hover:underline"
        >
          PRJ-METRO-001
        </button>
        <span>/</span>
        <button
          onClick={() => navigate(CONSTRUCTION_ROUTES.ISSUES)}
          className="hover:underline"
        >
          {data.related.issueReference}
        </button>
        <span>/</span>
        <span className="font-semibold text-[#041627]">Construction Decision {data.reference}</span>
      </nav>

      {/* ------------------------------------------------------------- */}
      {/* Status & metadata ribbon                                      */}
      {/* ------------------------------------------------------------- */}
      <div className="flex flex-wrap items-center justify-between gap-2 pt-2">
        <div className="flex items-center gap-2">
          <span className="px-2.5 py-1 rounded-full bg-[#006c49]/10 text-[#006c49] text-xs leading-4 font-semibold flex items-center gap-1.5 shadow-sm">
            <span className="w-1.5 h-1.5 rounded-full bg-[#006c49] animate-pulse" />
            {data.statusLabel}
          </span>
          <span className="text-on-surface-variant text-xs leading-4 font-semibold">{data.id}</span>
        </div>
        <span className="text-on-surface-variant text-[11px] leading-[14px]">
          Finalized {data.decidedOnLabel}
          {!data.decidedHasTime && (
            <span className="text-outline italic"> · no time recorded</span>
          )}
        </span>
      </div>

      {/* ------------------------------------------------------------- */}
      {/* Header decision card                                         */}
      {/* ------------------------------------------------------------- */}
      <section className="bg-surface-container-lowest rounded-[8px] p-5 shadow-sm flex flex-col gap-3">
        <div className="flex items-start justify-between gap-3">
          <h1 className="text-2xl font-bold leading-8 tracking-tight text-[#041627]">
            {data.subject}
          </h1>
        </div>

        <div className="flex flex-wrap items-center gap-2 pt-1">
          <div className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-[4px] bg-surface-container-high text-on-surface">
            <span className="material-symbols-outlined text-[15px] text-on-surface-variant">report_problem</span>
            <span className="text-[11px] leading-[14px] font-medium">
              Issue: <span className="font-semibold text-[#041627]">{data.related.issueId}</span>
            </span>
          </div>
          <div className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-[4px] bg-surface-container-high text-on-surface">
            <span className="material-symbols-outlined text-[15px] text-on-surface-variant">local_shipping</span>
            <span className="text-[11px] leading-[14px] font-medium">
              Delivery: <span className="font-semibold text-[#041627]">{data.related.deliveryId}</span>
            </span>
          </div>
          <div className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-[4px] bg-surface-container-high text-on-surface">
            <span className="material-symbols-outlined text-[15px] text-on-surface-variant">domain</span>
            <span className="text-[11px] leading-[14px] font-medium">
              Type: <span className="font-semibold text-[#041627]">{data.related.qualifierLabel}</span>
            </span>
          </div>
        </div>

        <p className="text-sm leading-5 text-on-surface-variant mt-1 leading-relaxed">{data.context}</p>

        {/* Context visual — no canonical image exists for this aggregate, so the
            slot is kept as a graphic treatment carrying only derived facts. */}
        <div className="relative w-full h-32 rounded-[4px] overflow-hidden mt-1 shadow-inner bg-gradient-to-br from-primary-container via-[#041627] to-[#006c49]">
          <div className="absolute inset-0 flex items-center justify-center">
            <span className="material-symbols-outlined text-[40px] text-secondary-fixed/60">alt_route</span>
          </div>
          <div className="absolute inset-0 bg-gradient-to-t from-[#041627]/80 via-[#041627]/30 to-transparent flex items-end p-3">
            <div className="flex items-center gap-2 text-on-primary">
              <span className="material-symbols-outlined text-[18px] text-secondary-fixed">alt_route</span>
              <span className="text-xs leading-4 font-semibold tracking-wide">
                {data.route.reference} detour
                {data.route.slipLabel ? ` · ${data.route.slipLabel} vs plan` : ""}
              </span>
            </div>
          </div>
        </div>
      </section>

      {/* ------------------------------------------------------------- */}
      {/* Options evaluated                                            */}
      {/* ------------------------------------------------------------- */}
      <section className="flex flex-col gap-2.5">
        <div className="flex items-center justify-between px-1">
          <h2 className="text-lg font-semibold leading-6 text-[#041627] flex items-center gap-2">
            <span className="material-symbols-outlined text-[#041627] text-[20px]">balance</span>
            Options Evaluated
          </h2>
          <span className="text-[11px] leading-[14px] font-medium text-on-surface-variant">
            {optionCount} {alternativeLabel}
          </span>
        </div>

        {data.options.map((option) => (
          <article
            key={option.title}
            className={
              option.isSelected
                ? "bg-surface-container-lowest rounded-[8px] p-4 shadow-md relative overflow-hidden"
                : "bg-surface-container-lowest rounded-[8px] p-4 shadow-sm opacity-85"
            }
          >
            {option.isSelected && <div className="absolute top-0 left-0 right-0 h-1 bg-[#006c49]" />}

            <div className={option.isSelected ? "flex items-start justify-between gap-2 pt-1" : "flex items-start justify-between gap-2"}>
              <div className="flex items-center gap-2">
                <span
                  className={
                    option.isSelected
                      ? "w-6 h-6 rounded-full bg-[#006c49] text-on-secondary flex items-center justify-center text-[12px] font-bold shrink-0"
                      : "w-6 h-6 rounded-full bg-surface-container-highest flex items-center justify-center text-on-surface-variant text-[12px] font-bold shrink-0"
                  }
                >
                  {option.position}
                </span>
                {option.isSelected ? (
                  <div>
                    <h3 className="text-lg font-semibold leading-6 text-[#041627]">{option.title}</h3>
                    <span className="text-[11px] leading-[14px] font-medium text-[#006c49]">Selected &amp; Authorized</span>
                  </div>
                ) : (
                  <h3 className="text-lg font-semibold leading-6 text-on-surface">{option.title}</h3>
                )}
              </div>

              <span
                className={
                  option.isSelected
                    ? "px-2.5 py-1 rounded-full bg-secondary-container text-on-secondary-container text-[11px] leading-[14px] font-semibold flex items-center gap-1 whitespace-nowrap"
                    : "px-2 py-0.5 rounded-full bg-error-container text-on-error-container text-[11px] leading-[14px] font-medium flex items-center gap-1 whitespace-nowrap"
                }
              >
                <span className="material-symbols-outlined text-[13px]">
                  {option.isSelected ? "check_circle" : "close"}
                </span>
                {option.statusLabel}
              </span>
            </div>

            {option.isSelected ? (
              <div className="mt-3.5 bg-surface-container-low rounded-[4px] p-3 space-y-2">
                {option.highlights.map((fact) => (
                  <div key={fact.label} className="flex items-center justify-between gap-3 text-left">
                    <span className="text-[11px] leading-[14px] font-medium text-on-surface-variant">{fact.label}</span>
                    {fact.value ? (
                      <span className="text-sm leading-5 font-bold text-[#041627]">{fact.value}</span>
                    ) : (
                      <Unavailable />
                    )}
                  </div>
                ))}
              </div>
            ) : (
              <div className="mt-3 grid grid-cols-2 gap-2 text-left">
                {option.impact.map((fact) => (
                  <Fact key={fact.label} label={fact.label} value={fact.value} />
                ))}
              </div>
            )}

            {option.rationale && (
              <p className="text-sm leading-5 text-on-surface-variant mt-2.5">{option.rationale}</p>
            )}
          </article>
        ))}
      </section>

      {/* ------------------------------------------------------------- */}
      {/* Approval & governance                                        */}
      {/* ------------------------------------------------------------- */}
      <section className="bg-surface-container-lowest rounded-[8px] p-5 shadow-sm space-y-4">
        <div className="flex items-center justify-between gap-3">
          <div>
            <h2 className="text-lg font-semibold leading-6 text-[#041627]">Governance Chain</h2>
            <p className="text-[11px] leading-[14px] font-medium text-on-surface-variant">
              Requested by {data.requestedBy} · {data.participantCount} participants consulted
            </p>
          </div>
          <span className="px-2.5 py-1 rounded-full bg-[#006c49]/10 text-[#006c49] text-[11px] leading-[14px] font-semibold flex items-center gap-1 whitespace-nowrap">
            <span className="material-symbols-outlined text-[14px]">verified_user</span>
            Decided by {data.decidedBy}
          </span>
        </div>

        <div className="relative pl-5 space-y-5 before:absolute before:left-2 before:top-2 before:bottom-2 before:w-0.5 before:bg-surface-container-highest">
          {data.participants.map((participant) => (
            <div key={participant.role} className="relative flex items-start justify-between gap-3">
              <span className="absolute -left-5 top-1 w-2.5 h-2.5 rounded-full bg-[#006c49] shadow-sm ring-4 ring-surface-container-lowest" />
              <div className="flex items-center gap-2.5 min-w-0">
                <div
                  className={
                    participant.isRequester
                      ? "w-8 h-8 rounded-full bg-primary-container text-on-primary font-bold text-xs flex items-center justify-center shrink-0"
                      : "w-8 h-8 rounded-full bg-surface-container text-[#041627] font-bold text-xs flex items-center justify-center shrink-0"
                  }
                >
                  {participant.initials}
                </div>
                <div className="min-w-0">
                  <p className="text-sm leading-5 font-semibold text-[#041627] truncate">{participant.role}</p>
                  <p className="text-[11px] leading-[14px] font-medium text-on-surface-variant">
                    {participant.isRequester
                      ? "Requester"
                      : participant.organization ?? "Organisation not modelled"}
                  </p>
                </div>
              </div>
              <div className="text-right shrink-0">
                {participant.signoffLabel ? (
                  <span className="text-[11px] leading-[14px] font-medium text-[#006c49]">{participant.signoffLabel}</span>
                ) : (
                  <Unavailable label="Sign-off not recorded" />
                )}
                {participant.signedAtLabel && (
                  <span className="text-[11px] leading-[14px] text-on-surface-variant block">{participant.signedAtLabel}</span>
                )}
              </div>
            </div>
          ))}
        </div>
      </section>

      {/* ------------------------------------------------------------- */}
      {/* Outcome summary                                              */}
      {/* ------------------------------------------------------------- */}
      <section className="rounded-[8px] bg-secondary-container/30 p-4 shadow-sm flex items-start gap-3">
        <div className="w-9 h-9 rounded-[4px] bg-[#006c49] text-on-secondary flex items-center justify-center shrink-0 mt-0.5">
          <span className="material-symbols-outlined text-[20px]">task_alt</span>
        </div>
        <div className="space-y-1">
          <h3 className="text-lg font-semibold leading-6 text-[#041627]">{data.outcome.title}</h3>
          <p className="text-sm leading-5 text-on-secondary-container">{data.outcome.body}</p>
          <p className="text-[11px] leading-[14px] font-medium text-on-secondary-container/80">
            Transport <span className="font-semibold text-[#041627]">{data.related.transportId}</span>
            {" · "}
            Arrival <span className="font-semibold text-[#041627]">{data.route.arrivalLabel}</span>
          </p>
        </div>
      </section>

      {/* ------------------------------------------------------------- */}
      {/* Action bar                                                   */}
      {/* ------------------------------------------------------------- */}
      <div className="flex flex-col gap-2.5 pt-2">
        <button
          onClick={() => navigate(CONSTRUCTION_ROUTES.DELIVERIES)}
          className="w-full h-12 bg-[#041627] text-on-primary rounded-[8px] text-lg font-semibold leading-6 flex items-center justify-center gap-2 shadow-md active:scale-[0.99] transition-transform hover:bg-[#041627]/90"
        >
          <span className="material-symbols-outlined text-[20px]">local_shipping</span>
          View Updated Delivery ({data.related.deliveryId})
        </button>
        <button
          onClick={() => navigate(CONSTRUCTION_ROUTES.PROVENANCE)}
          className="w-full h-11 bg-surface-container-lowest text-[#041627] rounded-[8px] text-sm font-semibold leading-5 flex items-center justify-center gap-2 shadow-sm active:bg-surface-container transition-colors hover:bg-surface-container"
        >
          <span className="material-symbols-outlined text-[18px] text-on-surface-variant">history_edu</span>
          View Audit Trail &amp; Provenance
        </button>
      </div>
    </div>
  );
};

export default DecisionPage;
