package json

import (
	"encoding/json"
	"net/http"
)

func Write(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func Read(r *http.Request, data any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() //возращает ошибку, если что-то не то ввели в body
	/*
			ВВОД ТОЛЬКО ПО ТАКОЙ СТРУКТУРЕ:
					{
		    "customer_id": 42,
		    "items":[
		        { "productId": 1, "quantity": 1},
		        { "productId": 5, "quantity": 5}
		    ]
		}
	*/
	return decoder.Decode(data)
}
