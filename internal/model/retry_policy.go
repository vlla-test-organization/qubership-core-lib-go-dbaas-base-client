package model

var (
	NoRetryPolicyForGetConnection = RetryPolicy{notRetrybleHttpCodes: []int{400: 404}}
	NoRetryPolicy                 = RetryPolicy{notRetrybleHttpCodes: []int{401, 403}}
	EmptyRetryPolicy              = RetryPolicy{notRetrybleHttpCodes: []int{}}
)

type RetryPolicy struct {
	notRetrybleHttpCodes []int
}

func (this *RetryPolicy) HasNotRetryableHttpCode(httpCode int) bool {
	for _, v := range this.notRetrybleHttpCodes {
		if v == httpCode {
			return true
		}
	}
	return false
}
