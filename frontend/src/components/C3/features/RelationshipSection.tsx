import { ThreadAssetViewInterface } from "../domain/thread/asset.types";

export const RelationshipSection = ({
    asset,
}: {
    asset: ThreadAssetViewInterface;
}) => (

    <div className="dp-section">

        <div className="dp-section-title">

            C3 Relationship

        </div>

        <div className="c3-extended-box">

            <div className="c3-ext-header">

                <span className="c3-ext-icon">

                    ⛓

                </span>

                <span className="c3-ext-title">

                    {asset?.externalVault}

                </span>

                <span className="c3-ext-active-badge">

                    Active

                </span>

            </div>

            <div className="c3-ext-detail">

                {asset?.direction}

                {" · "}

                {asset?.role}

            </div>

        </div>

    </div>

);