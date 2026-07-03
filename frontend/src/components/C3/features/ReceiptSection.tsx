import { ThreadAssetViewInterface } from "../domain/thread/asset.types";

export const ReceiptSection = ({
    asset,
}: {
    asset: ThreadAssetViewInterface;
}) => (

    <div className="dp-section">

        <div className="dp-section-title">

            Receipts

        </div>

        {asset?.receipts?.map(r => (

            <div
                key={r?.id}
                className="commit-row"
            >

                <div className="commit-ts">

                    {r?.time}

                </div>

                <div className="commit-body">

                    <div>

                        <span className="commit-actor">

                            {r?.vault}

                        </span>

                        <span className="commit-action">

                            {r?.status}

                        </span>

                    </div>

                </div>

            </div>

        ))}

    </div>

);