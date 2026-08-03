import { ThreadAssetViewInterface } from "../domain/thread/asset.types";

export const EventTimelineSection = ({
    asset,
}: {
    asset: ThreadAssetViewInterface;
}) => (

    <div className="dp-section">

        <div className="dp-section-title">

            Event Timeline

        </div>

        {asset?.events?.map(event => (

            <div
                key={event.id}
                className="commit-row"
            >

                <div className="commit-ts">

                    {event.time}

                </div>

                <div className="commit-body">

                    <div>

                        <span className="commit-actor">

                            {event.actor}

                        </span>

                        <span className="commit-action">

                            {event.type}

                        </span>

                    </div>

                    <div className="commit-cid">

                        Event {event.id}

                    </div>

                    <div className="commit-cid">

                        Payload {event.payloadCID}

                    </div>

                </div>

                <div className="commit-verify">

                    View →

                </div>

            </div>

        ))}

    </div>

);