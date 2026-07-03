import { AccessPolicy } from "../../domain/policy/policy.types";

export const PolicyCard = ({
    policy
}: {
    policy: AccessPolicy
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