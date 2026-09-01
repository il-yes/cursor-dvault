import { useNavigate } from "react-router-dom";
import { ChannelRow } from "../../domain/channel/channel.types";
import { toRowStatusView } from "../../domain/channel/channel.mapper";
import * as ROUTES from "@/constants/routes";

export const LedgerRowView = ({
    row,
    onClick
}: {
    row: ChannelRow;
    onClick: (id: string) => void;
}) => {
    const navigate = useNavigate();
    const statusView = toRowStatusView(row.status);

    return (
        <tr
            className="ledger-row"
            onClick={() => onClick(row.id)}
        >
            <td>
                <span className={`sdot ${statusView.dotClass}`} />
                <span className={`c3-status-badge ${statusView.dotClass}`}>
                    {statusView.label}
                </span>
            </td>

            <td>
                <div className="th-line1">
                    <span className="th-type">
                        {row.type}
                    </span>
                    <span>{row.title}</span>
                </div>

                <div className="th-line2">
                    {row.subtitle}
                </div>
            </td>

            <td>
                <div className="flow">
                    {row.participants.map((step, idx) => (
                        <span key={idx}>
                            {idx > 0 && <span className="fa">→</span>}
                            <span
                                className="vb"
                                style={{ background: step.color }}
                            >
                                {step.label}
                            </span>
                        </span>
                    ))}
                </div>
            </td>

            <td>
                <div>
                    <b>{row.assetCount} </b> assets
                </div>
                <div className="event">
                    {row.lastEvent}
                </div>
            </td>

            <td className="ts">
                {row.lastActivity}
            </td>

            <td>
                <span className="stellar-val">
                    {row.stellarTx}
                </span>
                <button
                    className="rha-btn"
                    onClick={(e) => {
                        e.stopPropagation();
                        onClick(row.id);
                    }}
                >
                    Open
                </button>
            </td>

            <td>
                <button
                    className={`c3b ${row.c3Status}`}
                    onClick={(e) => {
                        e.stopPropagation();
                        navigate(ROUTES.C3_CONFIG);
                    }}
                    title="Configure C3 Channel"
                >
                    {row.c3Label} ⚙️
                </button>
            </td>
        </tr>
    );
};

export const DisputeRow = () => {
    const navigate = useNavigate();
    return (
        <tr className="row-sel row-new-left-border">
            <td>
                <span className="sdot s-dispute" />
            </td>

            <td>
                <div>
                    <span className="th-type">Contract</span>
                    Contract — Supplier X <span className="new-label">Disputed</span>
                    <button
                        className="c3-row-config-btn"
                        title="C3 Configuration"
                        onClick={(e) => {
                            e.stopPropagation();
                            navigate(ROUTES.C3_CONFIG);
                        }}
                        style={{
                            marginLeft: "8px",
                            background: "#fafafa",
                            border: "1px solid #e0e0e0",
                            borderRadius: "4px",
                            padding: "2px 6px",
                            fontSize: "12px",
                            cursor: "pointer",
                            display: "inline-flex",
                            alignItems: "center",
                            gap: "3px",
                            color: "#555"
                        }}
                    >
                        ⚙️
                    </button>
                </div>
                <div className="th-line2" style={{ marginTop: 3 }}>
                    vendor: Supplier X · value: $240,000
                </div>
            </td>

            <td>
                <div className="flow">
                    <span className="vb" style={{ background: '#7C3AED' }}>L</span>
                    <span className="fa">→</span>
                    <span className="vb" style={{ background: '#2563EB' }}>F</span>
                    <span className="fa">→</span>
                    <span className="vb" style={{ background: '#444444' }}>Dir</span>
                </div>
            </td>

            <td>
                <div className="pipeline">
                    <div className="pseg pseg-done" />
                    <div className="pseg pseg-reject" />
                    <div className="pseg pseg-wait" />
                </div>
            </td>

            <td className="ts" style={{ color: '#DC2626' }}>
                today
            </td>

            <td>
                <span className="stellar-val">tx_f2a9…</span>
            </td>

            <td style={{ textAlign: 'center' }}>
                <button
                    className="c3b c3-active"
                    onClick={(e) => {
                        e.stopPropagation();
                        navigate(ROUTES.C3_CONFIG);
                    }}
                    title="Configure C3 Channel"
                >
                    ⛓+ ⚙️
                </button>
            </td>
        </tr>
    );
};

