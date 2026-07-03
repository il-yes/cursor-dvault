import { EventSummary } from "../../domain/thread/asset.types";

export const EventTimeline = ({
    events
}: {
    events: EventSummary[]
}) => (

    <div className="event-timeline">


        {
            events?.map((event, index) => (


                <div
                    key={index}
                    className="event-item"
                >


                    <div className="event-dot">

                        {
                            event.status === 'rejected'
                                ? '!'
                                :
                                event.status === 'pending'
                                    ? '○'
                                    :
                                    '✓'
                        }

                    </div>



                    <div className="event-content">


                        <div className="event-title">

                            {event.type}

                        </div>



                        <div className="event-meta">

                            <span>
                                {event.actor}
                            </span>

                            <span>
                                {event.time}
                            </span>

                        </div>


                    </div>



                </div>


            ))
        }



    </div>

);
