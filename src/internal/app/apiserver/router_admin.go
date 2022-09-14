package apiserver

import (
	"context"
	"html/template"
	"net/http"
	"strconv"
	"github.com/deuuus/bmstu_db_cw/src/internal/app/model"
	"github.com/joho/godotenv"
)

func (s *server) handleAdminStart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/admin/general/start.html")
	} 
}

func (s *server) handleAdminGeneralCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/admin/general/create.html")
	} 
}

func (s *server) authenticateAdminGeneral(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		_, ok := session.Values["system_password"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *server) handleAdminGeneralPass() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		password := r.FormValue("password")

		envs, err := godotenv.Read(".env")

		if err != nil {
			s.logger.Fatal("Error loading .env file")
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if password != envs["admin_password"] {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["system_password"] = "ok"
		if err := s.sessionStore.Save(r, w, session); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		http.ServeFile(w, r, "templates/admin/general/enter.html")
		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handleAdminGeneralCommitCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		surname := r.FormValue("surname")
        email := r.FormValue("email")
		password := r.FormValue("password")

		a := &model.Admin{
			Name: name,
			Surname: surname,
			Email: email,
			Password: password,
		}
		if err := s.store.Admin().Create(a); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		a.Sanitize()
		http.ServeFile(w, r, "templates/admin/main.html")
		s.respond(w, r, http.StatusCreated, a)
	} 
}

func (s *server) handleAdminGeneralLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "templates/admin/general/login.html")
	} 
}

func (s *server) handleAdminGeneralCommitLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        email := r.FormValue("email")
		password := r.FormValue("password")

		a, err := s.store.Admin().FindByEmail(email)
		if err != nil || !a.ComparePassword(password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["admin_id"] = a.ID
		if err := s.sessionStore.Save(r, w, session); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		http.ServeFile(w, r, "templates/admin/main.html")
		s.respond(w, r, http.StatusOK, nil)
	} 
}

func (s *server) handleAdminMainPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "templates/admin/main.html")
	} 
}

func (s *server) handleAdminShowStat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		specs, err := s.store.Specialization().GetAll()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		data := []*model.SpecStatView{}
		for i := 0; i < len(specs); i++ {
			ds, ps, err := s.store.Disease().Percent(specs[i].ID)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			for j := 0; j < len(ps); j++{
				disease, err := s.store.Disease().Find(ds[j])
				if err != nil {
					s.error(w, r, http.StatusInternalServerError, err)
					return
				}
				num := &model.Float{Num: ps[j]}
				data = append(data, &model.SpecStatView{Specialization: specs[i], Disease: disease, Percent: num})
			}
		}

		t, _ := template.ParseFiles("templates/admin/stat.html")
		t.Execute(w, data)
	} 
}

func (s *server) handleAdminShowDoctors() http.HandlerFunc {
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
		t, _ := template.ParseFiles("templates/admin/show_doctors.html")
		t.Execute(w, data)
	}
}

func (s *server) handleAdminShowPatients() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		patients, err := s.store.Patient().GetAll()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		t, _ := template.ParseFiles("templates/admin/show_patients.html")
		t.Execute(w, patients)
	}
}

func (s *server) handleAdminShowAdmins() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		admins, err := s.store.Admin().GetAll()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		t, _ := template.ParseFiles("templates/admin/show_admins.html")
		t.Execute(w, admins)
	}
}

func (s *server) handleAdminCabinet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["admin_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		a, err := s.store.Admin().Find(id.(int))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		t, _ := template.ParseFiles("templates/admin/cabinet.html")
		t.Execute(w, a)
	} 
}

func (s *server) authenticateAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["admin_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		a, err := s.store.Admin().Find(id.(int))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, a)))
	})
}

func (s *server) handleAdminCreateDoctor() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		specs, err := s.store.Specialization().GetAll()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		t, _ := template.ParseFiles("templates/doctor/create.html")
		t.Execute(w, specs)
	}
}

func (s *server) handleAdminCommitCreateDoctor() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		surname := r.FormValue("surname")
		work_since, _ := strconv.Atoi(r.FormValue("work_since"))
		spec_id, _ := strconv.Atoi(r.Form.Get("spec"))
        email := r.FormValue("email")
		password := r.FormValue("password")

		d := &model.Doctor{
			Name: name,
			Surname: surname,
			Work_since: work_since,
			Spec_id: spec_id,
			Email: email,
			Password: password,
		}
		if err := s.store.Doctor().Create(d); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		
		d.Sanitize()
		http.ServeFile(w, r, "templates/admin/success_doctor.html")
		s.respond(w, r, http.StatusCreated, d)
	}
}

func (s *server) handleAdminCreateAdmin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		specs, err := s.store.Specialization().GetAll()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		t, _ := template.ParseFiles("templates/admin/create.html")
		t.Execute(w, specs)
	}
}

func (s *server) handleAdminCommitCreateAdmin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		surname := r.FormValue("surname")
        email := r.FormValue("email")
		password := r.FormValue("password")

		a := &model.Admin{
			Name: name,
			Surname: surname,
			Email: email,
			Password: password,
		}
		if err := s.store.Admin().Create(a); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		
		a.Sanitize()
		http.ServeFile(w, r, "templates/admin/success_admin.html")
		s.respond(w, r, http.StatusCreated, a)
	}
}

func (s *server) handleAdminShowDoctorDoneVisits() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.URL.Query().Get("doctor_id"))
		
		records, err := s.store.Record().GetAllByDoctor(id)
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

        t, _ := template.ParseFiles("templates/admin/doctor_done_visits.html")
		t.Execute(w, data)
	})
}

func (s *server) handleAdminShowDoctorCabinet() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.URL.Query().Get("doctor_id"))
		
		d, err := s.store.Doctor().Find(id)
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

        t, _ := template.ParseFiles("templates/admin/doctor_cabinet.html")
		t.Execute(w, data)
	})
}