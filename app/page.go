package app

import (
	"hirine/helpers"
	"hirine/models"
	"net/http"
	"sync"
	"text/template"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/tj/go/http/response"

	"github.com/gorilla/csrf"
)

func (app *App) IndexGet() http.HandlerFunc {
	var (
		init sync.Once
		tpl  *template.Template
		err  error
	)
	return func(w http.ResponseWriter, r *http.Request) {
		init.Do(func() {
			tpl, err = template.ParseFiles("./templates/layout/base.html", "./templates/index.html")
		})
		if err != nil {
			response.InternalServerError(w, err.Error())
			return
		}
		session, _ := app.S.Get(r, "r_u_n_a_w_a_y")
		flash := session.Flashes()
		session.Save(r, w)
		var messages []string
		_, ok := session.Values["n_0"]
		login := false
		if ok {
			login = true
		}
		if flash != nil {
			for _, f := range flash {
				fString, ok := f.(string)
				if ok {
					messages = append(messages, fString)
				}

			}
		}
		var jobs []models.Job
		err := app.DB.GetJobs(&jobs, 10, 0)
		if err != nil {
			response.InternalServerError(w, err.Error())
			return
		}
		tpl.Execute(w, map[string]interface{}{
			"login":          login,
			csrf.TemplateTag: csrf.TemplateField(r),
			"jobs":           jobs,
			"messages":       messages,
		})
	}

}

func (app *App) DashboardRegisterCompanyGet() http.HandlerFunc {
	var (
		init sync.Once
		tpl  *template.Template
		err  error
	)
	return func(w http.ResponseWriter, r *http.Request) {
		init.Do(func() {
			tpl, err = template.ParseFiles(
				"./templates/layout/base-dashboard.html", "./templates/company-register.html")
		})
		if err != nil {
			response.InternalServerError(w, err.Error())
		}
		session, _ := app.S.Get(r, "r_u_n_a_w_a_y")
		login := false
		if _, ok := session.Values["n_0"]; ok {
			login = true
		}
		flash := session.Flashes()
		session.Save(r, w)
		var messages []string
		if flash != nil {
			for _, f := range flash {
				fString, ok := f.(string)
				if ok {
					messages = append(messages, fString)
				}

			}
		}
		tpl.Execute(w, map[string]interface{}{
			"login":          login,
			csrf.TemplateTag: csrf.TemplateField(r),
			"messages":       messages,
		})
	}
}

func (app *App) RegisterCompanyPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := app.S.Get(r, "r_u_n_a_w_a_y")
		r.ParseForm()
		companyName := r.FormValue("name")
		var userID string

		// @todo
		userID, ok := session.Values["n_0"].(string)
		if !ok {
			session.AddFlash("Please login first")
			session.Save(r, w)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		if len(companyName) == 0 {
			session.AddFlash("Please enter company name")
			session.Save(r, w)
			http.Redirect(w, r, "/company-register", http.StatusFound)
			return
		}
		userObjectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			session.AddFlash("Something is wrong")
			session.Save(r, w)
			http.Redirect(w, r, "/company-register", http.StatusFound)
			return
		}
		id, err := app.DB.CreateCompany(&models.CreateCompanyPayload{
			CompanyName:  companyName,
			ProfileImage: "",
			Admin:        userObjectID,
		})

		if err != nil {
			session.AddFlash("Something is wrong")
			session.Save(r, w)
			http.Redirect(w, r, "/company-register", http.StatusFound)
			return
		}
		http.Redirect(w, r, "/dashboard/company/"+id+"/admin", http.StatusSeeOther)
		return
	}

}

func (app *App) PostJobPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := app.S.Get(r, "r_u_n_a_w_a_y")
		userID, ok := session.Values["n_0"].(string)
		if !ok {
			session.AddFlash("Please login first")
			session.Save(r, w)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		userObjectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			session.AddFlash("Something is wrong")
			session.Save(r, w)
			http.Redirect(w, r, "/dashboard/post-job", http.StatusFound)
			return
		}
		r.ParseForm()
		title := r.FormValue("title")
		category := r.FormValue("category")
		description := r.FormValue("description")
		companyID := r.FormValue("company")
		companyObjectID, err := primitive.ObjectIDFromHex(companyID)
		if err != nil {
			session.AddFlash("Something is wrong")
			session.Save(r, w)
			http.Redirect(w, r, "/dashboard/post-job", http.StatusFound)
			return
		}
		// firstSalaryString := r.FormValue("firstSalary")
		// firstSalary, err := strconv.ParseFloat(firstSalaryString, 64)
		// if err != nil {
		// 	http.Redirect(w, r, "/dashboard/post-job", http.StatusFound)
		// }
		// secondSalaryString := r.FormValue("secondSalary")
		// secondSalary, err := strconv.ParseFloat(secondSalaryString, 64)
		// paymentMethod := r.FormValue("paymentMethod")
		// currency := r.FormValue("currency")
		// rate := r.FormValue("rate")
		// startDateString := r.FormValue("startDate")
		// startDate := helpers.ParseJavascriptTimeString(startDateString)
		// endDateString := r.FormValue("endDate")
		// endDate := helpers.ParseJavascriptTimeString(endDateString)
		// startTimeString := r.FormValue("startTime")
		// startTime := helpers.ParseJavascriptTimeString(startTimeString)
		// endTimeString := r.FormValue("endTime")
		// endTime := helpers.ParseJavascriptTimeString(endTimeString)
		// reminder := r.FormValue("reminder")

		jobID, err := app.DB.CreateJob(&models.CreateJobPayload{
			Title:    title,
			Category: category,
			// FirstSalary:   firstSalary,
			// SecondSalary:  secondSalary,
			// PaymentMethod: paymentMethod,
			// Currency:      currency,
			// Rate:          rate,
			// StartDate:     startDate,
			// EndDate:       endDate,
			// StartTime:     startTime,
			// EndTime:       endTime,
			Description: description,
			// Reminder:      reminder,
			Company: companyObjectID,
			Creator: userObjectID,
		})
		if err != nil {
			http.Redirect(w, r, "/dashboard/post-job", http.StatusFound)
		}
		http.Redirect(w, r, "/dashboard/job/"+jobID, http.StatusFound)
	}
}

func (app *App) ApplyJobPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := app.S.Get(r, "r_u_n_a_w_a_y")
		userID, ok := session.Values["n_0"].(string)
		if !ok {
			session.AddFlash("Please login first")
			session.Save(r, w)
			app.L.WithField("user", userID)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		userObjectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			session.AddFlash("Something is wrong")
			session.Save(r, w)
			app.L.Debugln("cant convert user object id")
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		r.ParseForm()
		jobID := r.FormValue("jobID")
		jobObjectID, err := primitive.ObjectIDFromHex(jobID)
		if err != nil {
			session.AddFlash("Something is wrong")
			session.Save(r, w)
			app.L.Debugln("cant convert job object id")
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		var application models.Application
		err = app.DB.GetApplicationByApplicantAndJob(&application, jobObjectID, userObjectID)
		if (models.Application{}) != application {
			session.AddFlash("existed application")
			session.Save(r, w)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		if err != nil && err != mongo.ErrNoDocuments {
			session.AddFlash("something is wrong")
			session.Save(r, w)
			app.L.WithError(err).Info("")
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		_, err = app.DB.CreateApplication(userObjectID, jobObjectID, StatusApplying)
		if err != nil {
			session.AddFlash("something is wrong")
			session.Save(r, w)
			app.L.WithError(err).Info("")
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		session.AddFlash("application create")
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

const (
	StatusApplying = iota + 1
	StatusNoConsidered
	StatusInterviewing
	StatusOfferMake
)

func (app *App) ScheduleInterviewPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		referer := r.Referer()
		session, _ := app.S.Get(r, "r_u_n_a_w_a_y")
		userID, ok := session.Values["n_0"].(string)
		if !ok {
			session.AddFlash("Please login first")
			session.Save(r, w)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		userObjectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			session.AddFlash("Something is wrong")
			session.Save(r, w)
			http.Redirect(w, r, referer, http.StatusFound)
			return
		}
		r.ParseForm()
		applicationString := r.FormValue("application")
		var timeOptions []time.Time
		time1String := r.FormValue("time1")
		if time1, err := helpers.ParseDateTimeLocalString(time1String); err == nil {
			timeOptions = append(timeOptions, time1)
		}
		time2String := r.FormValue("time2")
		if time2, err := helpers.ParseDateTimeLocalString(time2String); err == nil {
			timeOptions = append(timeOptions, time2)
		}
		time3String := r.FormValue("time3")
		if time3, err := helpers.ParseDateTimeLocalString(time3String); err == nil {
			timeOptions = append(timeOptions, time3)
		}
		if len(timeOptions) == 0 {
			session.AddFlash("no time provide")
			session.Save(r, w)
			http.Redirect(w, r, referer, http.StatusFound)
			return
		}

		var application models.Application
		err = app.DB.GetApplicationByApplicationID(&application, applicationString)
		if err != nil {
			app.L.WithError(err).Debugln("GetApplicationByApplicationID")
			session.AddFlash("Something is wrong")
			session.Save(r, w)
			http.Redirect(w, r, referer, http.StatusFound)
			return
		}
		if application.Status != StatusApplying {
			session.AddFlash("You are not allowed to create interview event at this moment")
			session.Save(r, w)
			http.Redirect(w, r, referer, http.StatusFound)
			return
		}
		err = app.DB.CreateInterviewEvent(application.ApplicationID, application.Applicant, userObjectID, timeOptions)
		if err != nil {
			app.L.WithError(err).Debugln("CreateInterviewEvent")
			session.AddFlash("Something is wrong")
			session.Save(r, w)
			http.Redirect(w, r, referer, http.StatusFound)
			return
		}
		err = app.DB.UpdateApplicationStatus(application.ApplicationID, StatusInterviewing)
		if err != nil {
			app.L.WithError(err).Debugln("UpdateApplicationStatus")
			session.AddFlash("Something is wrong")
			session.Save(r, w)
			http.Redirect(w, r, referer, http.StatusFound)
			return
		}
		session.AddFlash("interview create")
		session.Save(r, w)
		http.Redirect(w, r, referer, http.StatusSeeOther)
		return
	}
}

func (app *App) ConfirmInterviewPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		referer := r.Referer()
		session, _ := app.S.Get(r, "r_u_n_a_w_a_y")
		userID, ok := session.Values["n_0"].(string)
		if !ok {
			session.AddFlash("Please login first")
			session.Save(r, w)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		_, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			session.AddFlash("Something is wrong")
			session.Save(r, w)
			http.Redirect(w, r, referer, http.StatusFound)
			return
		}
		r.ParseForm()
		eventString := r.FormValue("event")
		timeString := r.FormValue("time")
		time, err := helpers.ParseUnixString(timeString)
		if err != nil {
			app.L.WithError(err).Debugln("GetEventByEventID")
			session.AddFlash("time is required")
			session.Save(r, w)
			http.Redirect(w, r, referer, http.StatusFound)
			return
		}
		var event models.Event
		err = app.DB.GetEventByEventID(&event, eventString)
		if err != nil {
			app.L.WithError(err).Debugln("GetEventByEventID")
			session.AddFlash("Something is wrong")
			session.Save(r, w)
			http.Redirect(w, r, referer, http.StatusFound)
			return
		}
		err = app.DB.ConfirmInterviewEvent(event.EventID, time)
		if err != nil {
			app.L.WithError(err).Debugln("ConfirmInterviewEvent")
			session.AddFlash("Something is wrong")
			session.Save(r, w)
			http.Redirect(w, r, referer, http.StatusFound)
			return
		}

		err = app.DB.UpdateApplicationStatus(event.Attendee, StatusInterviewing)
		if err != nil {
			app.L.WithError(err).Debugln("UpdateApplicationStatus")
			session.AddFlash("Something is wrong")
			session.Save(r, w)
			http.Redirect(w, r, referer, http.StatusFound)
			return
		}
		session.AddFlash("interview confirmed")
		session.Save(r, w)
		http.Redirect(w, r, referer, http.StatusSeeOther)
		return
	}
}

func ScheduleWorkPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		referer := r.Referer()
		session, _ := app.S.Get(r, "r_u_n_a_w_a_y")
		userID, ok := session.Values["n_0"].(string)
		if !ok {
			session.AddFlash("Please login first")
			session.Save(r, w)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		userObjectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			session.AddFlash("Something is wrong")
			session.Save(r, w)
			http.Redirect(w, r, referer, http.StatusFound)
			return
		}
		r.ParseForm()
		// applicationString := r.FormValue("application")
		var workTimeString time.Time
		workTime, err := helpers.ParseDateTimeLocalString(workTimeString)
		if err == nil {
			session.AddFlash("no time provide")
			session.Save(r, w)
			http.Redirect(w, r, referer, http.StatusFound)
			return
		}

		err = app.DB.CreateWorkEvent(userObjectID, userObjectID, workTime)
		if err != nil {
			app.L.WithError(err).Debugln("CreateWorkEvent")
			session.AddFlash("Something is wrong")
			session.Save(r, w)
			http.Redirect(w, r, referer, http.StatusFound)
			return
		}

		session.AddFlash("work create")
		session.Save(r, w)
		http.Redirect(w, r, referer, http.StatusSeeOther)
		return
	}
}
