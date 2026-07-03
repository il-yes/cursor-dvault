import { Fragment } from "react/jsx-runtime";
import { ThreadAssetViewInterface } from "../domain/thread/asset.types";

export const LifecycleSection = ({
    asset
}: {
    asset: ThreadAssetViewInterface;
}) => (

    <div className="dp-section">

        <div className="dp-section-title">

            Lifecycle

        </div>

        {asset?.lifecycle?.map((step, i) => (

            <Fragment key={i}>

                <div className="pipeline-step">

                    <div
                        className="ps-icon"
                        style={{
                            background: step?.color
                        }}
                    >
                        {step?.icon}
                    </div>

                    <div className="ps-content">

                        <div className="ps-label">
                            {step?.label}
                        </div>

                        <div className="ps-sublabel">
                            {step?.actor}
                        </div>

                    </div>

                    <div className={`ps-check ${step?.state}`}>
                        {step?.state === "done" ? "✓" : "○"}
                    </div>

                </div>

                {i !== asset?.lifecycle?.length - 1 &&
                    <div className="ps-connector" />
                }

            </Fragment>

        ))}

    </div>

);