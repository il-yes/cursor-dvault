import { ChannelRow } from "../../domain/channel/channel.types";

export const ChannelRowView = ({
    row,
    onClick
}: {
    row: ChannelRow,
    onClick: (id: string) => void;
}) => (

    <tr onClick={() => { onClick(row.id) }}>

        <td>
            <span
                className={`sdot s-${row.status}`}
            />
        </td>


        <td style={{
            cursor: 'pointer'
        }}>

            <div className={`th-line1 ${row.id === 'budget' ? 'new-thread' : ''}`}>
                <span className="th-type">
                    {row.type}
                </span>

                {row.title}

            </div>


            <div className="th-line2">

                {row.subtitle}

            </div>

        </td>


        <td>

            <div className="flow">

                {
                    row.participants.map(
                        (step, index) => (
                            <span key={index}>

                                {index > 0 &&
                                    <span className="fa">
                                        →
                                    </span>
                                }


                                <span
                                    className="vb"
                                    style={{
                                        background: step.color
                                    }}
                                >
                                    {step.label}
                                </span>


                            </span>
                        )
                    )
                }

            </div>


        </td>


        <td>

            <div className="asset-box">

                <span>
                    {row.assetCount}
                </span>

                assets

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


            <div className="row-hover-actions">

                <button className="rha-btn" onClick={(e) => {

                    e.stopPropagation();

                    onClick(row.id);

                }}>
                    Open
                </button>


                <button className="rha-btn">
                    Export
                </button>


            </div>

        </td>



        <td style={{ textAlign: 'center' }}>


            <button
                className={
                    `c3b ${row.c3Status === 'internal'
                        ? 'c3-internal'
                        :
                        row.c3Status === 'linked'
                            ? 'c3-linked'
                            :
                            'c3-active'
                    }`
                }

                title={row.c3Status}

            >

                {row.c3Label}

            </button>


        </td>


    </tr>

);



export const LedgerRowView = ({
    row,
    onClick
}: {
    row: ChannelRow;
    onClick: (id: string) => void;
}) => (

    <tr
        className="ledger-row"
        onClick={() => onClick(row.id)}
    >


        <td>
            <span className={`sdot s-${row.status}`} />
        </td>


        <td>

            <div className="th-line1">

                <span className="th-type">
                    {row.type}
                </span>

                {row.title}

            </div>


            <div className="th-line2">
                {row.subtitle}
            </div>

        </td>



        <td>

            <div className="flow">

                {
                    row.participants.map((step, idx) => (

                        <span key={idx}>

                            {idx > 0 &&
                                <span className="fa">→</span>
                            }


                            <span
                                className="vb"
                                style={{
                                    background: step.color
                                }}
                            >
                                {step.label}
                            </span>

                        </span>

                    ))
                }

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

                className={`c3b ${row.c3Status
                    }`}

                onClick={(e) => {

                    e.stopPropagation();

                }}

            >
                {row.c3Label}

            </button>

        </td>


    </tr>

);

export const DisputeRow = () => (
    <tr className="row-sel row-new-left-border">
        <td>
            <span className="sdot s-dispute" />
        </td>

        <td>
            <div>
                <span className="th-type">Contract</span>
                Contract — Supplier X <span className="new-label">Disputed</span>
            </div>
            <div className="th-line2" style={{ marginTop: 3 }}>
                vendor: Supplier X · value: $240,000
            </div>
        </td>

        <td>
            <div className="flow">
                <span className="vb" style={{ background: '#7C3AED' }}>
                    L
                </span>
                <span className="fa">→</span>
                <span className="vb" style={{ background: '#2563EB' }}>
                    F
                </span>
                <span className="fa">→</span>
                <span className="vb" style={{ background: '#444444' }}>
                    Dir
                </span>
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
            <button className="c3b c3-active">⛓+</button>
        </td>
    </tr>
);