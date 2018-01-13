package main

import (
	"github.com/julienschmidt/httprouter"
)

type Route struct {
	Name string
	Method string
	Pattern string
	Handle httprouter.Handle
}

type Routes []Route

//func (a *App) initializeRoutes() {
//	a.Router.GET("/", a.index)
//	a.Router.GET("/expenses", a.getExpenses)
//	a.Router.POST("/expense", a.createExpense)
//	a.Router.GET("/expense/:id", a.Get)
//	a.Router.PUT("/expense/:id", a.updateExpense)
//	a.Router.DELETE("/expense/:id", a.deleteExpense)
//}

//var routes = Routes{
//	Route{
//		"Index",
//		Index,
//	},
//	Route{
//		"Expense",
//		"POST",
//		"/expense",
//		Expense,
//	},
//	Route{
//		"Expenses",
//		"GET",
//		"/expenses",
//		GetExpenses,
//	},
//}