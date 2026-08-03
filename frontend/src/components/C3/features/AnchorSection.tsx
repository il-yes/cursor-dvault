import { ThreadAssetViewInterface } from "../domain/thread/asset.types";


export const AnchorSection = ({ asset }: { asset: ThreadAssetViewInterface }) => (

    <div className="dp-section">

        <div className="dp-section-title">

            Anchors

        </div>

        <div className="stellar-ref-row">

            <div className="stellar-hash-full">

                {asset?.stellarTx}

            </div>

            <div className="copy-btn">

                Copy

            </div>

        </div>

    </div>

);