import React from "react";
import '../styles/activation-row.css';


export const ActivationChannel = () => {
    return (
        <>
            {/* NEW ROW — top, highlighted, pulsing */}
            <tr className="row-new row-new-left-border">
                <td>
                    <span className="sdot s-new" />
                </td>
                <td>
                    <div className="th-line1 new-thread">
                        <span className="th-type" style={{ color: "#C8922A" }}>
                            Contract
                        </span>
                        contract-execution — Cipla_India
                        <span className="new-label">NEW</span>
                    </div>
                    <div className="first-slot-hint">
                        <span className="fsh-dot" />
                        vault_legal can commit{" "}
                        <code style={{ fontSize: 10, fontFamily: "monospace" }}>
                            contract_draft
                        </code>
                    </div>
                </td>
                <td>
                    <div className="flow">
                        <span className="vb" style={{ background: "#7C3AED" }}>
                            L
                        </span>
                        <span className="fa">→</span>
                        <span className="vb" style={{ background: "#2563EB" }}>
                            F
                        </span>
                        <span className="fa">→</span>
                        <span className="vb" style={{ background: "#444" }}>
                            Dir
                        </span>
                    </div>
                </td>
                <td>
                    <div className="pipeline">
                        <div className="pseg pseg-new" />
                        <div className="pseg pseg-new" />
                        <div className="pseg pseg-new" />
                    </div>
                    <div style={{ fontSize: 10, color: "#ccc", marginTop: 4 }}>
                        ⏳ slot 1 of 3 open
                    </div>
                </td>
                <td
                    className="ts"
                    style={{ color: "#C8922A", fontWeight: 500 }}
                >
                    just now
                </td>
                <td>
                    <span className="stellar-genesis">genesis ✦</span>
                </td>
                <td style={{ textAlign: "center" }}>
                    <span className="c3b c3-ext">⛓+</span>
                </td>
            </tr>
        </>
    );
}
