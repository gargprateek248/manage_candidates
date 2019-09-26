package controllers

import "github.com/gargprateek248/manage_candidates/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	//Candidates routes
	s.Router.HandleFunc("/candidates", middlewares.SetMiddlewareJSON(s.CreateCandidate)).Methods("POST")
	s.Router.HandleFunc("/candidates", middlewares.SetMiddlewareJSON(s.GetCandidates)).Methods("GET")
	s.Router.HandleFunc("/candidates/{id}", middlewares.SetMiddlewareJSON(s.GetCandidate)).Methods("GET")
	s.Router.HandleFunc("/candidates/{id}", middlewares.SetMiddlewareJSON(s.UpdateCandidate)).Methods("PUT")
	s.Router.HandleFunc("/candidates/{id}", middlewares.SetMiddlewareJSON(s.DeleteCandidate)).Methods("DELETE")
}
