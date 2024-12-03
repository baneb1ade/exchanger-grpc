package exchanger

import "errors"

var ErrCurrencyNotFound = errors.New("currency not found")
var ErrCurrencyRateIsZero = errors.New("currency rate is zero")
var ErrSmtWentWrong = errors.New("something went wrong")
