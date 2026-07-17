package send

import "net/http"

func ToLocation(writer http.ResponseWriter) {
	writer.Header().Add("Location", "/welcome")
	writer.WriteHeader(http.StatusSeeOther)
}
