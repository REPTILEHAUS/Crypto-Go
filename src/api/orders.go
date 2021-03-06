package api

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/mitsukomegumi/Crypto-Go/src/orders"
	"github.com/mitsukomegumi/Crypto-Go/src/pairs"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// SetupOrderRoutes - setup necessary routes for accout database
func SetupOrderRoutes(router *fasthttprouter.Router, db *mgo.Database) (*fasthttprouter.Router, error) {
	_, pErr := setOrderPosts(router, db)

	if pErr != nil {
		return router, pErr
	}

	_, err := setGeneralRoutes(router, db)

	if err != nil {
		return router, err
	}

	_, err = setOrderUpdates(router, db)

	if err != nil {
		return router, err
	}

	_, err = setOrderGets(router, db)

	if err != nil {
		return router, err
	}

	_, err = setOrderDeletes(router, db)

	if err != nil {
		return router, err
	}

	_, err = setOrderFills(router, db)

	if err != nil {
		return router, err
	}

	return router, nil
}

func setOrderGets(initRouter *fasthttprouter.Router, db *mgo.Database) (*fasthttprouter.Router, error) {
	req, err := NewRequestServer("?pair?OrderID", "/api/orders/order", "GET", db, db, "?OrderID")
	if err != nil {
		return nil, err
	}

	router, err := req.AttemptToServeRequestsWithRouter(initRouter)

	if err != nil {
		return nil, err
	}

	return router, nil
}

func setOrderPosts(initRouter *fasthttprouter.Router, db *mgo.Database) (*fasthttprouter.Router, error) {
	postReq, rErr := NewRequestServer("?pair?ordertype?orderamount?fillprice?username?password", "/api/orders", "POST", nil, db, "?pair?ordertype?orderamount?fillprice?username?password")

	if rErr != nil {
		return nil, rErr
	}

	router, pErr := postReq.AttemptToServeRequestsWithRouter(initRouter)

	if pErr != nil {
		return nil, rErr
	}

	return router, nil
}

func setOrderUpdates(initRouter *fasthttprouter.Router, db *mgo.Database) (*fasthttprouter.Router, error) {
	updateReq, err := NewRequestServer("?pair?OrderID?username?password?updatedfill?updatedamount", "/api/orders/update", "POST", nil, db, "?pair?OrderID?username?password?updatedfill?updatedamount")

	if err != nil {
		return initRouter, err
	}

	_, err = updateReq.AttemptToServeRequestsWithRouter(initRouter)

	if err != nil {
		return initRouter, err
	}

	return initRouter, nil
}

func setOrderDeletes(initRouter *fasthttprouter.Router, db *mgo.Database) (*fasthttprouter.Router, error) {
	delReq, rErr := NewRequestServer("?pair?OrderID?username?password", "/api/orders", "DELETE", nil, db, "?pair?OrderID?username?password")

	if rErr != nil {
		return nil, rErr
	}

	router, pErr := delReq.AttemptToServeRequestsWithRouter(initRouter)

	if pErr != nil {
		panic(rErr)
	}

	return router, nil
}

func setOrderFills(initRouter *fasthttprouter.Router, db *mgo.Database) (*fasthttprouter.Router, error) {
	postReq, err := NewRequestServer("?pair?OrderID", "/api/orders/fill", "POST", nil, db, "?pair?OrderID")

	if err != nil {
		return nil, err
	}

	_, err = postReq.AttemptToServeRequestsWithRouter(initRouter)

	if err != nil {
		return nil, err
	}

	return initRouter, nil
}

func setGeneralRoutes(initRouter *fasthttprouter.Router, db *mgo.Database) (*fasthttprouter.Router, error) {
	getReq, _ := NewRequestServer("?pair", "/api/orders?pair", "GET", nil, db, "?pair")
	initRouter.GET("/api/orders", getReq.HandleGETCollection)

	return initRouter, nil
}

func addOrder(database *mgo.Database, order *orders.Order) error {

	_, err := findOrder(database, order.OrderID, *order.OrderPair)

	if err != nil {
		c := (*database).C(order.OrderPair.StartingSymbol + "-" + order.OrderPair.EndingSymbol)

		iErr := c.Insert(order)

		if iErr != nil {
			return iErr
		}

		return nil
	}
	return nil
}

func removeOrder(database *mgo.Database, order *orders.Order) error {
	c := database.C(order.OrderPair.StartingSymbol + "-" + order.OrderPair.EndingSymbol)
	orderRef, _ := findOrder(database, order.OrderID, *order.OrderPair)

	err := c.Remove(orderRef)

	if err != nil {
		return err
	}

	return nil
}

func findOrder(database *mgo.Database, id string, pair pairs.Pair) (*orders.Order, error) {
	c := database.C(pair.StartingSymbol + "-" + pair.EndingSymbol)

	result := orders.Order{}

	err := c.Find(bson.M{"orderid": id}).One(&result)
	if err != nil {
		return &result, err
	}

	return &result, nil
}
