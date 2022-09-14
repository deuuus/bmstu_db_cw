package apiserver

import (
	"context"
	"html/template"
	"net/http"
	"strconv"
	"github.com/deuuus/bmstu_db_cw/src/internal/app/model"
)

func (s *server) handlePatientMainPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "templates/patient/main.html")
	} 
}

func (s *server) handlePatientCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "templates/patient/create.html")
	} 
}

func (s *server) handlePatientCreateVisit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		doctors, err := s.store.Doctor().GetAll()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		data := []*model.DoctorView{}
		for i := 0; i < len(doctors); i++ {
			spec, _ := s.store.Specialization().Find(doctors[i].Spec_id)
			data = append(data, &model.DoctorView{Doctor: doctors[i], Specialization: spec})
		}
		t, _ := template.ParseFiles("templates/patient/create_visit.html")
		t.Execute(w, data)
	} 
}

func (s *server) handlePatientCommitCreateVisit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		doctor_id, _ := strconv.Atoi(r.FormValue("data"))

		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		patient_id, ok := session.Values["patient_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		v := &model.Visit{
			Status: "Active",
			Doctor_id: doctor_id,
			Patient_id: patient_id.(int),
		}

		if err := s.store.Visit().Create(v); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		http.ServeFile(w, r, "templates/patient/success_visit_create.html")
		s.respond(w, r, http.StatusCreated, v)
	} 
}

func (s *server) handlePatientCommitCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		surname := r.FormValue("surname")
		gender := r.Form.Get("gender")
		birth_year, _ := strconv.Atoi(r.FormValue("birth_year"))
		phone := r.FormValue("phone")
        email := r.FormValue("email")
		password := r.FormValue("password")

		p := &model.Patient{
			Name: name,
			Surname: surname,
			Gender: gender,
			Birth_year: birth_year,
			Phone: phone,
			Email: email,
			Password: password,
		}
		if err := s.store.Patient().Create(p); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["patient_id"] = p.ID
		if err := s.sessionStore.Save(r, w, session); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		
		p.Sanitize()
		http.ServeFile(w, r, "templates/patient/main.html")
		s.respond(w, r, http.StatusCreated, p)
	} 
}

func (s *server) authenticatePatient(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["patient_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		p, err := s.store.Patient().Find(id.(int))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, p)))
	})
}

func (s *server) handlePatientCabinet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["patient_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		p, err := s.store.Patient().Find(id.(int))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		t, _ := template.ParseFiles("templates/patient/cabinet.html")
		t.Execute(w, p)
	} 
}

func (s *server) handlePatientRecord() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["patient_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}
		
		records, err := s.store.Record().GetAllByPatient(id.(int))
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		data := []*model.RecordView{}
		for i := 0; i < len(records); i++ {
			visit, err := s.store.Visit().Find(records[i].Visit_id)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			doctor, err := s.store.Doctor().Find(visit.Doctor_id)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			spec, err := s.store.Specialization().Find(doctor.Spec_id)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			disease, err := s.store.Disease().Find(records[i].Disease_id)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			medicine, err := s.store.Medicine().Find(records[i].Medicine_id)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			data = append(data, &model.RecordView{Record: records[i], Visit: visit, Doctor: doctor, Specialization: spec, Disease: disease, Medicine: medicine})
		}

        t, _ := template.ParseFiles("templates/patient/record.html")
		t.Execute(w, data)
	} 
}

func (s *server) handlePatientShowActiveVisits() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["patient_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}
		
		visits, err := s.store.Visit().GetActiveVisitsByPatient(id.(int))
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		
		data := []*model.ActivePatientVisitView{}
		for i := 0; i < len(visits); i++ {
			doctor, err := s.store.Doctor().Find(visits[i].Doctor_id)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			spec, _ := s.store.Specialization().Find(doctor.Spec_id)
			data = append(data, &model.ActivePatientVisitView{Visit: visits[i], Doctor: doctor, Specialization: spec})
		}
		
		t, _ := template.ParseFiles("templates/patient/show_active_visits.html")
		t.Execute(w, data)
	})
}