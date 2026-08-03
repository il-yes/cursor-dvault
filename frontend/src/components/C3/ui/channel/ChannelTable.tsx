import { ChannelView } from "../../domain/channel/channel.types"
import { ThreadAssetViewInterface } from "../../domain/thread/asset.types";

export const ChannelTable = ({
    channel,
    onOpenAsset
}: {
    channel: ChannelView,
    onOpenAsset: (asset:ThreadAssetViewInterface) => void,
}) => {

    
    return (
        <div className="ledger-area">
            <div className="ledger-topbar">

                <div className="ledger-title">
                    CHANNEL
                </div>

                <div className="ledger-controls">

                    <button className="ctrl-btn">
                        Assets {channel?.assets?.total}
                    </button>

                    <button className="ctrl-btn">
                        Active ▾
                    </button>

                </div>

            </div>
            <table className="ledger-table">

                <thead>

                    <tr>

                        <th>
                            Channel
                        </th>

                        <th>
                            Thread Assets
                        </th>

                        <th>
                            Events
                        </th>

                        <th>
                            Policy
                        </th>

                        <th>
                            Stellar
                        </th>

                    </tr>

                </thead>


                <tbody>


                    <tr>


                        <td>

                            <div className="th-line1">

                                <span className="sdot s-active" />

                                {channel?.title}

                            </div>


                            <div className="th-line2">

                                {channel?.subtitle}

                            </div>


                            <div className="flow">

                                {
                                    channel?.participants?.length !== 0 &&  channel?.participants?.map(
                                        (p, i) => (

                                            <span key={i}>

                                                {i > 0 &&
                                                    <span className="fa">
                                                        →
                                                    </span>
                                                }


                                                <span
                                                    className="vb"
                                                    style={{
                                                        background: p.color
                                                    }}
                                                >
                                                    {p.label}
                                                </span>


                                            </span>

                                        )
                                    )

                                }

                            </div>

                        </td>



                        <td >

                            {
                                channel?.assets !== null && channel?.assets !== undefined ? channel?.assets?.items?.map(asset => (

                                    <div
                                        className="asset-box cursor-pointer"
                                        key={asset.id}

                                        onClick={() => {
                                            onOpenAsset(asset);
                                        }}
                                    >
                                        <div className="th-line1">
                                            {asset?.title}
                                        </div>

                                        <div className="th-line2">
                                            {asset?.lastEvent}
                                        </div>



                                        <div className="pipeline">

                                            <div className="pseg pseg-done" />
                                            <div className="pseg pseg-done" />
                                            <div className="pseg pseg-wait" />

                                        </div>

                                    </div>

                                )) : (
                                    <div>
                                        No assets
                                    </div>
                                )
                            }


                        </td>



                        <td>


                            {
                                channel?.activity?.events?.length > 0 && channel?.activity?.events?.map(
                                    (e, i) => (

                                        <div
                                            className="event"
                                            key={i}
                                        >

                                            <span>
                                                {e.time}
                                            </span>

                                            &nbsp;

                                            {e.actor}

                                            :

                                            {e.type}


                                        </div>

                                    )
                                )

                            }


                        </td>



                        <td>

                            <div className="policy-box">


                                <div>

                                    READ

                                </div>


                                {
                                    channel?.policy?.read?.length > 0 && channel?.policy?.read?.map(x => (

                                        <span
                                            className="vb"
                                            key={x}
                                        >
                                            {x}
                                        </span>

                                    ))
                                }


                                <div>

                                    WRITE

                                </div>


                                {
                                    channel?.policy?.write?.length > 0 && channel?.policy?.write?.map(x => (

                                        <span
                                            className="vb"
                                            key={x}
                                        >
                                            {x}
                                        </span>

                                    ))
                                }


                            </div>


                        </td>



                        <td>


                            <span className="stellar-val">

                                {channel?.stellarTx}

                            </span>


                            <button className="c3b c3-linked">

                                ⛓✓

                            </button>


                        </td>



                    </tr>


                </tbody>

            </table>

            <div className="ledger-footer">
                Channel ID <span className="footer-highlight">{channel?.id} </span>

                <span className="mx-2">·</span>

                {channel?.assets.total} assets

                <span className="mx-2">·</span>

                {channel?.activity.events.length} events

                <span className="mx-2">·</span>

                C3 {channel?.c3.status}

                <span className="mx-2">·</span>

                Last activity: {channel?.activity.lastActivity}

            </div>
        </div>
    )
}