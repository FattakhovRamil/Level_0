package server

import (
	"database/sql"
	"encoding/json"
	//"fmt"
	"net/http"

)

func serverStart(c *Cache) {
	// database, err := sql.Open("sqlite3", "./gopher.db")
	// if err != nil {
	// 	fmt.Println("Error opening database:", err)
	// 	return
	// }
	// statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY,application TEXT UNIQUE, param1 REAL, param2 TEXT, version INTEGER)") // создаем базу данных
	// statement.Exec()
	// defer statement.Close()
	// defer database.Close()
	// http.HandleFunc("/api/getstate", HundlerGet)   // GET
	// http.HandleFunc("/api/savestate", HundlerPost) // POST
	// fmt.Println("Server start listening on port 3001")
	// err = http.ListenAndServe("localhost:3001", nil)

	// if err != nil {
	// 	fmt.Println(err)
	// }
}

func HundlerGet(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		getState(w, r)
	} else {
		http.Error(w, "invalid http method", http.StatusMethodNotAllowed)
	}

}

func HundlerPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		saveState(w, r)
	} else {
		http.Error(w, "indalid http method", http.StatusMethodNotAllowed)
	}
}

func getState(w http.ResponseWriter, r *http.Request) {
	database, err := sql.Open("sqlite3", "./gopher.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer database.Close()

	decoder := json.NewDecoder(r.Body)
	var state State
	err = decoder.Decode(&state)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	rows, err := database.Query("SELECT application, param1, param2, version FROM users WHERE application = ?", state.Application)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var ser SelectInfo
	if rows.Next() {
		err := rows.Scan(&ser.Application, &ser.Param1, &ser.Param2, &ser.Version)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Record not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ser)
}

func saveState(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var req SelectInfo
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	now, check, err := checkUser(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !check {
		database, err := sql.Open("sqlite3", "./gopher.db")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer database.Close()
		statement, err := database.Prepare("INSERT INTO users (application, param1, param2, version) VALUES (?, ?, ?, 1)")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer statement.Close()

		_, err = statement.Exec(req.Application, req.Param1, req.Param2)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var ansver ReturnInfo
		ansver.Status = "User created"
		json.NewEncoder(w).Encode(ansver)
		return
	} else {
		if req.Param1 == now.Param1 && req.Param2 == now.Param2 {
			nothingUpdate(w)
			return
		}
		database, err := sql.Open("sqlite3", "./gopher.db")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer database.Close()
		statement, err := database.Prepare("UPDATE users SET param1 = ?, param2 = ?, version = ? WHERE application = ?")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer statement.Close()
		now.IncrementVersion()
		_, err = statement.Exec(req.Param1, req.Param2, now.Version, now.Application)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var ansver ReturnInfo
		ansver.Status = "Data updated, version upgraded"
		json.NewEncoder(w).Encode(ansver)
	}

}

func checkUser(req *SelectInfo) (now SelectInfo, found bool, err error) {
	database, err := sql.Open("sqlite3", "./gopher.db")
	if err != nil {
		return now, false, err
	}
	defer database.Close()

	rows, err := database.Query("SELECT application, param1, param2, version FROM users WHERE application = ?", req.Application)
	if err != nil {
		return now, false, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&now.Application, &now.Param1, &now.Param2, &now.Version)
		if err != nil {
			return now, false, err
		}
		return now, true, nil // Запись найдена
	}
	return now, false, nil // Запись не найдена
}

type State struct {
	Application string `json:"application"`
}

type SelectInfo struct {
	Application string  `json:"application"`
	Param1      float64 `json:"param1"`
	Param2      string  `json:"param2"`
	Version     int     `json:"version"`
}

func (p *SelectInfo) IncrementVersion() {
	p.Version = p.Version + 1
}

type ReturnInfo struct {
	Status string
}

func nothingUpdate(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	var ansver ReturnInfo
	ansver.Status = "The data is not updated because there is nothing to update"
	json.NewEncoder(w).Encode(ansver)
}
