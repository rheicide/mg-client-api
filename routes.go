package main

type Route struct {
	Methods []string
	Path    string
	Handler Handler
}

type Routes []Route

var routes = Routes{
	Route{
		[]string{"GET", "POST"},
		"/mails",
		MailIndex,
	},
	Route{
		[]string{"GET"},
		"/mails/{id}",
		MailById,
	},
}
