package main

import (
	"fmt"
	pb "github.com/qinxiaogit/go_mic_project/frontend_api/proto"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (fe *frontendApiServer) GetQuote(w http.ResponseWriter, r *http.Request) {
	log := r.Context().Value(ctxKeyLog{}).(logrus.FieldLogger)
	log.Debug("placing order")
}
func (fe *frontendApiServer) ShipOrder(w http.ResponseWriter, r *http.Request) {
	log := r.Context().Value(ctxKeyLog{}).(logrus.FieldLogger)
	log.Debug("placing order")
	_, err := fe.shippingService.GetQuote(r.Context(), &pb.GetQuoteRequest{
		Address: &pb.Address{StreetAddress: "aaaaa", City: "bbbb", State: "cccc", ZipCode: 123},
	})
	if err != nil {
		return
	}
	w.Write([]byte("hello world"))
}

func (fe *frontendApiServer) setCache(w http.ResponseWriter, r *http.Request) {
	log := r.Context().Value(ctxKeyLog{}).(logrus.FieldLogger)
	log.Debug("placing order")
	put, err := fe.cacheService.Put(r.Context(), &pb.PutRequest{Key: name, Value: "hello", Duration: "10h"})
	fmt.Println(err)
	if err != nil {
		return
	}
	w.Write([]byte("hello world" + put.String()))
	w.WriteHeader(http.StatusOK)
}
func (fe *frontendApiServer) getCache(w http.ResponseWriter, r *http.Request) {
	log := r.Context().Value(ctxKeyLog{}).(logrus.FieldLogger)
	log.Debug("placing order")
	get, err := fe.cacheService.Get(r.Context(), &pb.GetRequest{Key: name})
	if err != nil {
		return
	}
	w.Write([]byte("hello world" + string(get.GetValue())))
	w.WriteHeader(http.StatusOK)
}

func currentCurrency(r *http.Request) string {
	c, _ := r.Cookie(cookieCurrency)
	if c != nil {
		return c.Value
	}
	return defaultCurrencyCode
}
