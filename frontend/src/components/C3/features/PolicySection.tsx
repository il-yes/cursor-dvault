import { ThreadAssetViewInterface } from "../domain/thread/asset.types";

export const PolicySection = ({
    asset,
}: {
    asset: ThreadAssetViewInterface;
}) => (

    <div className="dp-section">

        <div className="dp-section-title">

            Access Policy

        </div>

        <div className="policy-box">

            <div>

                READ

            </div>

            {asset?.policy?.read?.map(v => (

                <span
                    className="vb"
                    key={v}
                >
                    {v}
                </span>

            ))}

            <div>

                WRITE

            </div>

            {asset?.policy?.write?.map(v => (

                <span
                    className="vb"
                    key={v}
                >
                    {v}
                </span>

            ))}

        </div>

    </div>

);