package service

import (
	"staff-service/internal/data/pg"
	address "staff-service/internal/service/handlers/address"
	person "staff-service/internal/service/handlers/person"
	position "staff-service/internal/service/handlers/position"
	staff "staff-service/internal/service/handlers/staff"
	"staff-service/internal/service/helpers"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
)

func (s *service) router() chi.Router {
	r := chi.NewRouter()
	log := s.log.WithFields(map[string]interface{}{
		"service": "customer-service-api",
	})

	r.Use(
		ape.RecoverMiddleware(log),
		ape.LoganMiddleware(log),
		ape.CtxMiddleware(
			helpers.CtxLog(log),
			helpers.CtxAddressesQ(pg.NewAddressesQ(s.db)),
			helpers.CtxPersonsQ(pg.NewPersonsQ(s.db)),
			helpers.CtxPositionsQ(pg.NewPositionsQ(s.db)),
			helpers.CtxStaffQ(pg.NewStaffQ(s.db)),
		),
	)
	r.Route("/integrations/staff-service", func(r chi.Router) {
		r.Route("/addresses", func(r chi.Router) {
			r.Post("/", address.CreateAddress)
			r.Get("/", address.GetAddressList)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", address.GetAddress)
				r.Put("/", address.UpdateAddress)
				r.Delete("/", address.DeleteAddress)
			})
		})
		r.Route("/persons", func(r chi.Router) {
			r.Post("/", person.CreatePerson)
			r.Get("/", person.GetPersonList)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", person.GetPerson)
				r.Put("/", person.UpdatePerson)
				r.Delete("/", person.DeletePerson)
			})
		})
		r.Route("/positions", func(r chi.Router) {
			r.Post("/", position.CreatePosition)
			r.Get("/", position.GetPositionList)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", position.GetPosition)
				r.Put("/", position.UpdatePosition)
				r.Delete("/", position.DeletePosition)
			})
		})
		r.Route("/staff", func(r chi.Router) {
			r.Post("/", staff.CreateStaff)
			r.Get("/", staff.GetStaffList)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", staff.GetStaff)
				r.Put("/", staff.UpdateStaff)
				r.Delete("/", staff.DeleteStaff)
			})
		})
	})

	return r
}
