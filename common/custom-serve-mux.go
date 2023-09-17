package common

import "net/http"

type CustomServeMux struct {
	mux                *http.ServeMux
	registeredPatterns []string
}

func NewCustomServeMux() *CustomServeMux {
	return &CustomServeMux{
		mux: http.NewServeMux(),
	}
}

func (csm *CustomServeMux) HandleFunc(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	csm.mux.HandleFunc(pattern, handler)
	csm.registeredPatterns = append(csm.registeredPatterns, pattern)
}

func (csm *CustomServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	csm.mux.ServeHTTP(w, r)
}

func (csm *CustomServeMux) GetRegisteredPatterns() []string {
	return csm.registeredPatterns
}
