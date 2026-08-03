import { AssetCard } from "./asset-view";
import { AssetSummary } from "./domain/thread/asset.types";

export type FlowStep = {
    label: string;
    color: string;
};
// Used inside Channel View

export type ChannelView = {

    id: string;


    title: string;

    subtitle: string;


    status:
    | 'active'
    | 'pending'
    | 'revoked';


    participants: FlowStep[];


    assets: {
        total: number;

        items: AssetSummary[];
    };


    activity: {
        lastEvent: string;
        lastActivity: string;

        events: EventSummary[];
    };


    policy: {
        read: string[];
        write: string[];
    };


    stellarTx: string;


    c3: {
        status:
        | 'internal'
        | 'linked'
        | 'active';

        externalVault?: string;
    };
};

export const channel: ChannelView = {

    id: 'contract-supplier',


    title:
        'Contract Execution Channel',


    subtitle:
        'Supplier X · $240,000 agreement',


    status: 'active',


    participants: [
        {
            label: 'L',
            color: '#7C3AED'
        },
        {
            label: 'F',
            color: '#2563EB'
        },
        {
            label: 'DIR',
            color: '#444'
        },
        {
            label: 'SUP',
            color: '#0891B2'
        }
    ],


    assets: {

        total: 3,


        items: [

            {
                id: 'msa',

                title:
                    'Master Service Agreement',

                type: 'contract',

                status: 'signed',

                lastEvent:
                    'contract.countersigned'
            },

            {
                id: 'pricing',

                title:
                    'Pricing Schedule',

                type: 'attachment',

                status: 'approved',

                lastEvent:
                    'finance.approved'
            }

        ]

    },


    activity: {


        lastEvent:
            'contract.countersigned',


        lastActivity:
            '2 hours ago',


        events: [

            {
                type: 'contract.created',
                actor: 'Legal',
                time: '12 Jun'
            },

            {
                type: 'finance.approved',
                actor: 'Finance',
                time: '12 Jun'
            },

            {
                type: 'contract.countersigned',
                actor: 'Supplier X',
                time: 'Today'
            }

        ]

    },



    policy: {

        read: [
            'Supplier X'
        ],

        write: [
            'Internal'
        ]

    },


    stellarTx:
        'tx_f2a9e3…',


    c3: {

        status: 'active',

        externalVault:
            'Supplier X'

    }

};

export const ChannelView = ({
    channel
}: {
    channel: ChannelView
}) => (

    <div className="channel-view">


        <ChannelHeader
            channel={channel}
        />



        <section>

            <h2>
                Assets ({channel?.assets.total})
            </h2>


            <div className="asset-grid">


                {
                    channel?.assets.items.map(asset =>

                        <AssetSummaryCard
                            key={asset.id}
                            asset={asset}
                        />

                    )
                }


            </div>


        </section>




        <section>


            <h2>
                Activity
            </h2>


            <EventTimeline
                events={channel?.activity.events}
            />


        </section>



        <section>


            <h2>
                Policy
            </h2>


            <PolicyCard
                policy={channel?.policy}
            />


        </section>


    </div>

)

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

export type EventSummary = {

    type: string;

    actor: string;

    time: string;

    status?:
    'done'
    |
    'pending'
    |
    'rejected';
};


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

type Policy = {

    read: string[];

    write: string[];

    append?: string[];

    expiresAt?: string;

};



export const PolicyCard = ({
    policy
}: {
    policy: Policy
}) => (

    <div className="policy-card">


        <div className="policy-header">

            <h3>
                Access Policy
            </h3>

            <span>
                Active
            </span>

        </div>




        <div className="policy-row">


            <div>

                <label>
                    READ
                </label>


                <div className="policy-users">


                    {
                        policy?.read.map(
                            (user) => (

                                <span key={user}>
                                    {user}
                                </span>

                            )

                        )
                    }


                </div>


            </div>


        </div>





        <div className="policy-row">


            <div>

                <label>
                    WRITE
                </label>


                <div className="policy-users">


                    {
                        policy.write.map(
                            (user) => (

                                <span key={user}>
                                    {user}
                                </span>

                            )

                        )
                    }


                </div>


            </div>


        </div>





        {
            policy.append &&
            <div className="policy-row">

                <label>
                    APPEND
                </label>


                <div className="policy-users">

                    {
                        policy.append.map(
                            (user) => (

                                <span key={user}>
                                    {user}
                                </span>

                            )

                        )
                    }

                </div>


            </div>
        }



        {
            policy.expiresAt &&

            <div className="policy-expiry">

                Expires:
                {policy.expiresAt}

            </div>

        }



    </div>

);

const AssetSummaryCard = ({
    asset
}: {
    asset: AssetSummary
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

            </div>



            <span
                className={`asset-status ${asset.status}`}
            >

                {asset.status}

            </span>


        </div>




        <div className="asset-event">

            Last event:

            <strong>
                {asset.lastEvent}
            </strong>

        </div>




        <div className="asset-actions">

            <button>
                Open
            </button>


        </div>


    </div>

);