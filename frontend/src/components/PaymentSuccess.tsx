import React, { useEffect } from 'react';
import * as AppAPI from "../../wailsjs/go/main/App";

const PaymentSuccess = () => {
    useEffect(() => {
        const params = new URLSearchParams(window.location.search);
        const subId = params.get('subscription_id');
        if (subId) {
            AppAPI.NotifyPaymentSuccess(subId);
        }
    }, []);



    return (
        <div>
            <h1>    Payment Success</h1>
        </div>
    );
};

export default PaymentSuccess;