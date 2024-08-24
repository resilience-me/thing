package commands

const (
    ClientTrustlines_SetTrustline      = 0
    ClientTrustlines_SyncTrustlineIn   = 1
    ClientTrustlines_SyncTrustlineOut  = 2
    ClientTrustlines_GetTrustlineIn    = 3
    ClientTrustlines_GetTrustlineOut   = 4
    ClientPayments_NewPaymentOut       = 5
    ClientPayments_NewPaymentIn        = 6
    ClientPayments_GetPayment          = 7

    ServerTrustlines_SetTrustline      = 127
    ServerTrustlines_GetTrustline      = 128
    ServerTrustlines_SetSyncOut        = 129
    ServerTrustlines_SetTimestamp      = 130
    ServerPayments_FindPathOut         = 131
    ServerPayments_FindPathIn          = 132
    ServerPayments_PathRecurse         = 133
)
