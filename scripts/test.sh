#!/bin/bash
go test -timeout 30s -count=1 -run ^TestOrderEndpoint$ -v purchase-cart-service/tests