package valvidaabs

import "time"

type IRatelimService interface {
	HitAndCheck(
		PRA []string, // primary router address

		// the key of specific counter (bucket)
		key string,

		// must be > 0; N is the same as N hits, with one request
		xMultiplier float64,

	) (
		result *Result,
		err error,
	)
}

type Result struct {
	// The main result
	Allow bool

	// 0% = fresh, 100% = full load
	// Can be used on client side, to determine the delay, to slow down the rate even before we hit the limit
	SaturationPerc float64

	// just a flag, for metrics
	// Can be used on client side to report _why_ the request failed - because of throttling, or because of fail-to-ban
	BannedByFailToBan bool

	// The time, during which any new requests are guaranteed to fail. This happens when the bucket overflows,
	// or the F2B goes to banned state.
	// This can be used on client side, to avoid excessive RPCs, when it can know for sure they'll fail.
	RetryAfter time.Duration

	PRANotFound bool

	// can return a textual error message, explaining "why" not allowed, e.g. to report to logs or user
	Why string

	// can be used to re-route to another PRA (redirect)
	RedirectPRA []string
}
