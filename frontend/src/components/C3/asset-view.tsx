export type AssetStatus =
    | 'draft'
    | 'pending'
    | 'approved'
    | 'signed'
    | 'rejected';


export type AssetView = {

    id: string;

    channelId: string;

    type:
    | 'contract'
    | 'invoice'
    | 'report'
    | 'credential'
    | 'attestation';


    title: string;

    subtitle: string;


    status: AssetStatus;


    createdAt: string;


    lastEvent: {
        type: string;
        at: string;
    };


    participants: string[];


    stellarTx: string;


    c3Visibility:
    | 'private'
    | 'shared'
    | 'external';


};


export const AssetCard = ({
    asset
}: {
    asset: AssetView
}) => (

    <div className="asset-card">


        <div className="asset-header">


            <div>

                <span className="asset-type">

                    {asset.type}

                </span>


                <h3>
                    {asset.title}
                </h3>


                <p>
                    {asset.subtitle}
                </p>

            </div>



            <span
                className={`asset-status ${asset.status}`}
            >

                {asset.status}

            </span>


        </div>



        <div className="asset-meta">


            <div>

                <label>
                    Last Event
                </label>

                <p>
                    {asset.lastEvent.type}
                </p>

            </div>



            <div>

                <label>
                    Activity
                </label>

                <p>
                    {asset.lastEvent.at}
                </p>

            </div>


        </div>



        <div className="asset-footer">


            <div className="participants">

                {
                    asset.participants.map(p =>
                        <span key={p}>
                            {p}
                        </span>
                    )
                }

            </div>


            <div>

                <span className="stellar-val">
                    {asset.stellarTx}
                </span>

            </div>


        </div>



        <div className="asset-actions">

            <button>
                Open
            </button>


            <button>
                Timeline
            </button>


        </div>


    </div>

);