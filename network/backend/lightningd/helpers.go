package lightningd

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"github.com/tidwall/gjson"
)

func randomHex(lenght int) (string, error) {
	hex := []rune("0123456789abcdef")
	b := make([]rune, lenght)
	for i := range b {
		r, err := rand.Int(rand.Reader, big.NewInt(16))
		if err != nil {
			return "", err
		}
		b[i] = hex[r.Int64()]
	}
	return string(b), nil
}

func errorMessageFromPayStatus(status gjson.Result) string {
	var messages []string

	status.Get("pay").ForEach(func(_, pay gjson.Result) bool {
		pay.Get("attempts").ForEach(func(_, attempt gjson.Result) bool {
			messages = append(messages, attempt.Get("failure.message").String())
			return true
		})
		return true
	})

	return strings.Join(messages, " / ")
}

func routeFromPayStatus(status gjson.Result) (route *models.Route) {
	status.Get("pay").ForEach(func(_, pay gjson.Result) bool {
		pay.Get("attempts").ForEach(func(_, attempt gjson.Result) bool {
			if !attempt.Get("success").Exist() {
				return true
			}

			route = routeStruct(attempt.Get("route"))
			return false
		})
		if route != nil {
			return false
		}
		return true
	})

	return route
}

func routeStruct(route gjson.Result) *models.Route {
	nhops := route.Get("#").Int()
	firstHop := route.Get("0")
	lastHop := route.Get(fmt.Sprintf("%d", nhops-1))

	r := &models.Route{
		TimeLock: uint32(firstHop.Get("delay").Int()),
		Fee:      firstHop.Get("msatoshi").Int() - lastHop.Get("msatoshi").Int(),
		Amount:   firstHop.Get("msatoshi").Int(),
		Hops:     make(*models.Hop, nhops),
	}

	i := 0
	route.ForEach(func(_, hop gjson.Result) bool {
		r.Hops[i] = &models.Hop{
			Amount: hop.Get("amount").Int(),
			Expiry: uint32(hop.Get("delay").Int()),
		}
		i++
		return true
	})

	return r
}
