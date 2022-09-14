package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"text/template"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/gorilla/handlers"
	"github.com/deuuus/bmstu_db_cw/src/internal/app/model"
	"github.com/deuuus/bmstu_db_cw/src/internal/app/store"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

const (
	sessionName        = "gopherschool"
	ctxKeyUser  ctxKey = iota
	ctxKeyRequestID
)

type ctxKey int8

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

var m = map[error]string {
	errIncorrectEmailOrPassword: "Введен неверный пароль или адрес электронной почты. Повторите попытку.",
	errNotAuthenticated:         "Вы не выполнили вход в систему",
}

type server struct {
	router       *mux.Router
	logger       *logrus.Logger
	store        store.Store
	sessionStore sessions.Store
}

func newServer(store store.Store, sessionStore sessions.Store) *server {
	s := &server{
		router:       mux.NewRouter(),
		logger:       logrus.New(),
		store:        store,
		sessionStore: sessionStore,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))

	//general routes
	s.router.HandleFunc("/start", s.handleStart())
	s.router.HandleFunc("/login", s.handleLogin())
	s.router.HandleFunc("/after", s.handleAfter())

	//patient routes
	s.router.HandleFunc("/create_patient", s.handlePatientCreate())
	s.router.HandleFunc("/commit_create", s.handlePatientCommitCreate())

	//patient private routes
	patient_root := s.router.PathPrefix("/patient").Subrouter()
	patient_root.Use(s.authenticatePatient)
	patient_root.HandleFunc("/main", s.handlePatientMainPage())
	patient_root.HandleFunc("/cabinet", s.handlePatientCabinet())
	patient_root.HandleFunc("/create_visit", s.handlePatientCreateVisit())
	patient_root.HandleFunc("/commit_visit", s.handlePatientCommitCreateVisit())
	patient_root.HandleFunc("/record", s.handlePatientRecord())
	patient_root.HandleFunc("/show_active_visits", s.handlePatientShowActiveVisits())

	//doctor private routes
	doctor_root := s.router.PathPrefix("/doctor").Subrouter()
	doctor_root.Use(s.authenticateDoctor)
	doctor_root.HandleFunc("/cabinet", s.handleDoctorCabinet())
	doctor_root.HandleFunc("/main", s.handleDoctorMainPage())
	doctor_root.HandleFunc("/show_active_visits", s.handleDoctorShowActiveVisits())
	doctor_root.HandleFunc("/show_done_visits", s.handleDoctorShowDoneVisits())
	doctor_root.HandleFunc("/commit_visit", s.handleDoctorCommitVisit())
	doctor_root.HandleFunc("/record", s.handleDoctorRecord())

	//admin private routes
	admin_root := s.router.PathPrefix("/admin").Subrouter()
	admin_root.Use(s.authenticateAdmin)
	admin_root.HandleFunc("/main", s.handleAdminMainPage())
	admin_root.HandleFunc("/cabinet", s.handleAdminCabinet())
	admin_root.HandleFunc("/create_doctor", s.handleAdminCreateDoctor())
	admin_root.HandleFunc("/commit_doctor", s.handleAdminCommitCreateDoctor())
	admin_root.HandleFunc("/create_admin", s.handleAdminCreateAdmin())
	admin_root.HandleFunc("/commit_admin", s.handleAdminCommitCreateAdmin())
	admin_root.HandleFunc("/show_doctors", s.handleAdminShowDoctors())
	admin_root.HandleFunc("/show_patients", s.handleAdminShowPatients())
	admin_root.HandleFunc("/show_stat", s.handleAdminShowStat())
	admin_root.HandleFunc("/show_admins", s.handleAdminShowAdmins())
	admin_root.HandleFunc("/doctor_done_visits", s.handleAdminShowDoctorDoneVisits())
	admin_root.HandleFunc("/doctor_cabinet", s.handleAdminShowDoctorCabinet())

	//admin general routes
	s.router.HandleFunc("/admin_start", s.handleAdminStart())
	s.router.HandleFunc("/admin_pass", s.handleAdminGeneralPass())

	//admin general private routes
	admin_general_root := s.router.PathPrefix("/admin_general").Subrouter()
	admin_general_root.Use(s.authenticateAdminGeneral)
	admin_general_root.HandleFunc("/admin_create", s.handleAdminGeneralCreate())
	admin_general_root.HandleFunc("/commit_create", s.handleAdminGeneralCommitCreate())
	admin_general_root.HandleFunc("/admin_login", s.handleAdminGeneralLogin())
	admin_general_root.HandleFunc("/commit_login", s.handleAdminGeneralCommitLogin())
}

func (s *server) handleStart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/start.html")
	} 
}

func (s *server) handleAfter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "templates/after.html")
	} 
}

func (s *server) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")
		password := r.FormValue("password")

		user_cookie, user_id := "patient_id", 0
		p, err := s.store.Patient().FindByEmail(email)
		if err != nil || !p.ComparePassword(password) {

			a, err := s.store.Admin().FindByEmail(email)
			if err != nil || !a.ComparePassword(password) {

				d, err := s.store.Doctor().FindByEmail(email)
				if err != nil {
					s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
					return
				} else {
					user_cookie, user_id = "doctor_id", d.ID
				}
			} else {
				user_cookie, user_id = "admin_id", a.ID
			}
		} else {
			user_id = p.ID
		}
	
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values[user_cookie] = user_id
		if err := s.sessionStore.Save(r, w, session); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		
		if user_cookie == "patient_id" {
			http.ServeFile(w, r, "templates/patient/main.html")
		} else if user_cookie == "admin_id" {
			http.ServeFile(w, r, "templates/admin/main.html")	
		} else {
			http.ServeFile(w, r, "templates/doctor/main.html")
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		var level logrus.Level
		switch {
		case rw.code >= 500:
			level = logrus.ErrorLevel
		case rw.code >= 400:
			level = logrus.WarnLevel
		default:
			level = logrus.InfoLevel
		}
		logger.Logf(
			level,
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
	})
}

func (s *server) handleWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, r.Context().Value(ctxKeyUser).(*model.Patient))
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	t, _ := template.ParseFiles("templates/error.html")
	if val, ok := m[err]; ok {
		t.Execute(w, val)
	} else {
		t.Execute(w, nil)
	}
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
