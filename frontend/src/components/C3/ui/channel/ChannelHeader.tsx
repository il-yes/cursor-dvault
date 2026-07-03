import { ChannelView } from "../../domain/channel/channel.types";


type ChannelHeaderProps = {
    channel: ChannelView;
};


export const ChannelHeader = ({
    channel
}: ChannelHeaderProps) => (

    <div className="channel-header">


        <div className="ch-top">


            <div>

                <div className="ch-type">
                    CHANNEL
                </div>


                <h1>
                    {channel?.title}
                </h1>


                <p className="ch-subtitle">
                    {channel.subtitle}
                </p>

            </div>



            <div
                className={`c3-state c3-${channel.c3.status}`}
            >

                <span>
                    {
                        channel.c3.status === 'active'
                            ? '●'
                            :
                            channel.c3.status === 'linked'
                                ? '◐'
                                :
                                '○'
                    }
                </span>

                {
                    channel.c3.status === 'active'
                        ?
                        'C3 Active'
                        :
                        channel.c3.status === 'linked'
                            ?
                            'C3 Linked'
                            :
                            'Internal'
                }

            </div>


        </div>





        <div className="ch-body">



            <div className="participant-block">


                <label>
                    Participants
                </label>


                <div className="flow">


                    {
                        channel.participants.map(
                            (step, index) => (

                                <span key={index}>


                                    {
                                        index > 0 &&
                                        <span className="fa">
                                            →
                                        </span>
                                    }


                                    <span
                                        className="vb"
                                        style={{
                                            background: step.color
                                        }}
                                    >

                                        {step.label}

                                    </span>


                                </span>

                            )
                        )

                    }


                </div>


            </div>






            <div className="channel-stats">


                <div>

                    <label>
                        Assets
                    </label>

                    <strong>
                        {channel?.assets.total}
                    </strong>

                </div>



                <div>

                    <label>
                        Last Event
                    </label>

                    <strong>
                        {channel?.activity.lastEvent}
                    </strong>

                </div>



                <div>

                    <label>
                        Activity
                    </label>

                    <strong>
                        {channel?.activity.lastActivity}
                    </strong>

                </div>


            </div>




        </div>





        <div className="ch-footer">


            <div>

                <span>
                    Stellar Anchor
                </span>

                <strong className="stellar-val">
                    {channel.stellarTx}
                </strong>

            </div>



            {
                channel.c3.externalVault &&

                <div className="external-pill">

                    External vault:
                    {channel.c3.externalVault}

                </div>

            }


        </div>


    </div>

);