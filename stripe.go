package main

import (
	"fmt"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/sub"
)

func isSubscribed(email string) bool {
	if email == "" {
		return false
	}
	stripe.Key = cfg.Stripe.APIKey
	params := &stripe.CustomerSearchParams{}
	fmt.Println(email)
	params.Query = *stripe.String("email:'" + email + "'")
	iter := customer.Search(params)
	for iter.Next() {
		cus := iter.Customer()

		params2 := &stripe.SubscriptionSearchParams{}
		params2.Query = *stripe.String("status:'active'")
		iter2 := sub.Search(params2)
		for iter2.Next() {
			s := iter2.Subscription()
			if s.Customer.ID == cus.ID {
				return true
			}
		}
	}
	return false
}
