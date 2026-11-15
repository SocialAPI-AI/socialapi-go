package socialapi

// Hand-written: survives regeneration. Everything in socialapi.gen.go is generated.

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math"
	"net/http"
	"strconv"
	"time"
)

// DefaultWebhookTolerance is how far a delivery's timestamp may be from now.
const DefaultWebhookTolerance = 5 * time.Minute

// ErrInvalidWebhookSignature is returned when a delivery fails verification.
var ErrInvalidWebhookSignature = errors.New("socialapi: invalid webhook signature")

// VerifyWebhook checks the X-SocialAPI-Signature-V2 and X-SocialAPI-Timestamp
// headers of a webhook delivery against your endpoint secret. Pass the raw
// request body, before any JSON decoding. A tolerance of 0 uses DefaultWebhookTolerance.
func VerifyWebhook(secret string, body []byte, header http.Header, tolerance time.Duration) error {
	return verifyWebhookAt(secret, body, header, tolerance, time.Now())
}

func verifyWebhookAt(secret string, body []byte, header http.Header, tolerance time.Duration, now time.Time) error {
	if tolerance <= 0 {
		tolerance = DefaultWebhookTolerance
	}
	ts := header.Get("X-SocialAPI-Timestamp")
	sec, err := strconv.ParseInt(ts, 10, 64)
	if err != nil || math.Abs(now.Sub(time.Unix(sec, 0)).Seconds()) > tolerance.Seconds() {
		return ErrInvalidWebhookSignature
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(ts + "."))
	mac.Write(body)
	want := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(want), []byte(header.Get("X-SocialAPI-Signature-V2"))) {
		return ErrInvalidWebhookSignature
	}
	return nil
}
