package apiserver

import (
	"context"
	"html/template"
	"net/http"
	"strconv"
	"github.com/deuuus/bmstu_db_cw/src/internal/app/model"
)

func (s *server) handleDoctorCommitVisit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		visit_id, _ := strconv.Atoi(r.FormValue("visit_id"))
		disease_id, _ := strconv.Atoi(r.FormValue("disease_id"))
		medicine_id, _ := strconv.Atoi(r.FormValue("medicine_id"))

		if err := s.store.Visit().CommitVisit(visit_id); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		rec := &model.Record{
			Visit_id: visit_id,
			Disease_id: disease_id,
			Medicine_id: medicine_id,
		}

		if err := s.store.Record().Create(rec); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
        http.ServeFile(w, r, "templates/doctor/success_commit_visit.html")
	}
}

func (s *server) handleDoctorRecord() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.URL.Query().Get("patient_id"))
		
		records, err := s.store.Record().GetAllByPatient(id)
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

        t, _ := template.ParseFiles("templates/doctor/record.html")
		t.Execute(w, data)
	} 
}

func (s *server) handleDoctorMainPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "templates/doctor/main.html")
	} 
}

func (s *server) handleDoctorShowActiveVisits() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["doctor_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}
		
		visits, err := s.store.Visit().GetActiveVisitsByDoctor(id.(int))
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		
		patients := []*model.Patient{}
		vps := []model.VisitView{}
		for i := 0; i < len(visits); i++ {
			patient, err := s.store.Patient().Find(visits[i].Patient_id)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			patients = append(patients, patient)
			vps = append(vps, model.VisitView{Visit: visits[i], Patient: patient})
		}

		doctor, err := s.store.Doctor().Find(id.(int))
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		diseases, _ := s.store.Disease().GetBySpecialization(doctor.Spec_id)
		medicines, _ := s.store.Medicine().GetAll()

		data := model.ActiveVisitView{VPs: vps, Diseases: diseases, Medicines: medicines}
		
		t, _ := template.ParseFiles("templates/doctor/show_active_visits.html")
		t.Execute(w, data)
	})
}

func (s *server) handleDoctorShowDoneVisits() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/*id, _ := strconv.Atoi(r.URL.Query().Get("patient_id"))*/

		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["doctor_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}
		
		records, err := s.store.Record().GetAllByDoctor(id.(int))
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		data := []*model.RecordDoctorView{}
		for i := 0; i < len(records); i++ {
			visit, err := s.store.Visit().Find(records[i].Visit_id)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			patient, err := s.store.Patient().Find(visit.Patient_id)
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
			data = append(data, &model.RecordDoctorView{Record: records[i], Visit: visit, Patient: patient, Disease: disease, Medicine: medicine})
		}

        t, _ := template.ParseFiles("templates/doctor/show_done_visits.html")
		t.Execute(w, data)
	})
}

func (s *server) authenticateDoctor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["doctor_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		d, err := s.store.Doctor().Find(id.(int))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, d)))
	})
}

func (s *server) handleDoctorCabinet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["doctor_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		d, err := s.store.Doctor().Find(id.(int))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		sp, err := s.store.Specialization().Find(d.Spec_id)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		data := model.DoctorView{
			Doctor: d,
			Specialization: sp,
		}
		
		t, _ := template.ParseFiles("templates/doctor/cabinet.html")
		t.Execute(w, data)
	} 
}