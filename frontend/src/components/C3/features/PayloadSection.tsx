import { ThreadAssetViewInterface } from "../domain/thread/asset.types";

export const PayloadSection = ({
    asset,
}: {
    asset: ThreadAssetViewInterface;
}) => (

    <div className="dp-section">

        <div className="dp-section-title">

            Payloads

        </div>

        {asset?.payloads?.map(payload => (

            <div
                className="commit-row"
                key={payload.cid}
            >

                <div className="commit-body">

                    <div>

                        {payload.name}

                    </div>

                    <div className="commit-cid">

                        {payload.cid}

                    </div>

                </div>

                <div className="commit-verify">

                    Download

                </div>

            </div>

        ))}

    </div>

);