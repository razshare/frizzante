package send

import "net/http"

func ToLocation(writer http.ResponseWriter, location string) {
	writer.Header().Add("Location", location)
	writer.WriteHeader(http.StatusSeeOther)
}
