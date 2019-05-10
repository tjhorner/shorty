package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
)

type apiResult struct {
	Error  *string     `json:"error"`
	Result interface{} `json:"result"`
}

func apiSuccess(result interface{}) apiResult {
	return apiResult{Result: result}
}

func apiError(err error) apiResult {
	es := err.Error()
	return apiResult{Error: &es}
}

type api struct {
	ctx *shortyContext
}

func (a *api) route(router *httprouter.Router) {
	router.GET("/", a.getRoot)
	router.GET("/:suffix", a.getSuffix)
	router.POST("/api/shorten", a.postShorten)
}

func (a *api) internalError(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	enc := json.NewEncoder(w)
	enc.Encode(apiError(errors.New("internal server error")))
}

func (a *api) badRequest(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	enc := json.NewEncoder(w)
	enc.Encode(apiError(errors.New("bad request")))
}

func (a *api) getRoot(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	http.Redirect(w, r, a.ctx.Config.DefaultRedirect, 302)
}

func (a *api) getSuffix(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var link Link

	err := a.ctx.DB.Where(&Link{Suffix: params.ByName("suffix")}).First(&link).Error
	if err != nil {
		http.Redirect(w, r, a.ctx.Config.DefaultRedirect, 302)
		return
	}

	http.Redirect(w, r, link.URL, 301)
}

func (a *api) postShorten(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		a.internalError(w, r)
		return
	}

	suffix := r.Form.Get("suffix")
	if suffix == "" {
		suffix = randString(6)
	}

	longURL := r.Form.Get("url")
	if longURL == "" {
		a.badRequest(w, r)
		return
	}

	parsedURL, err := url.ParseRequestURI(longURL)
	if err != nil {
		log.Println(err)
		a.badRequest(w, r)
		return
	}

	link := Link{
		Suffix:    suffix,
		URL:       parsedURL.String(),
		CreatorIP: remoteAddr(r),
	}

	err = a.ctx.DB.Create(&link).Error
	if err != nil {
		log.Println(err)
		a.badRequest(w, r)
		return
	}

	enc := json.NewEncoder(w)
	enc.Encode(apiSuccess(link))
}
