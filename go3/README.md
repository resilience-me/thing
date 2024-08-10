Step 1) A command to place a time lock on the trustlines is sent down the path. 

Step 2) A command to finalize the commit is sent down the path. This increases the time lock, and, adds a rule that the commit can only be aborted if it is verified that the next in line has aborted it, or never received it. (Thus if it reaches buyer, it cannot be cancelled unless buyer somehow cancels it... )

Step 3) A command to finalize the payment is sent down the path. A credit line has now formed, and the payment is complete.
