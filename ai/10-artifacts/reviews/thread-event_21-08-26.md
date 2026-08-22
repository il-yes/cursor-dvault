The audit also exposed a very serious engineering smell:

_ = json.Unmarshal(payloadJson, &ref)

That should absolutely not survive in the final architecture.

An invalid event payload currently becomes a valid-looking empty event and gets sent to Cloud.

That's exactly the kind of boundary failure your Runtime Verification principles are supposed to catch.