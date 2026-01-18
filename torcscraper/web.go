package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

const (
	adminUser = "admin"
	adminPass = "admin"
)

func giris(r *http.Request) bool {
	c, err := r.Cookie("auth")
	return err == nil && c.Value == "1"
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !giris(r) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

func giriskontrol(w http.ResponseWriter, r *http.Request) {
	if giris(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		user := r.FormValue("user")
		pass := r.FormValue("pass")

		if user == adminUser && pass == adminPass {
			http.SetCookie(w, &http.Cookie{
				Name:     "auth",
				Value:    "1",
				Path:     "/",
				HttpOnly: true,
				MaxAge:   3600,
			})
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		} else {
			hatamesaji(w, "Geçersiz kullanıcı adı veya şifre!")
			return
		}
	}

	hatamesaji(w, "")
}

func hatamesaji(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `
		<h2>Login</h2>
		<p style="color:red">%s</p>
		<form method="POST">
			User: <input name="user" required><br><br>
			Pass: <input name="pass" type="password" required><br><br>
			<button type="submit">Giriş Yap</button>
		</form>
	`, msg)
}

func anasayfa(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
    SELECT id, title, source_title, category, criticality, created_at 
    FROM findings 
    ORDER BY created_at DESC
`)
	if err != nil {
		http.Error(w, "Veritabanı hatası: "+err.Error(), 500)
		return
	}
	defer rows.Close()

	type Row struct {
		ID          int
		Title       string
		SourceTitle string
		Category    string
		Critical    string
		CreatedAt   string
	}

	var list []Row
	for rows.Next() {
		var rr Row
		err := rows.Scan(
			&rr.ID,
			&rr.Title,
			&rr.SourceTitle,
			&rr.Category,
			&rr.Critical,
			&rr.CreatedAt,
		)
		if err != nil {
			continue
		}
		list = append(list, rr)
	}

	tpl := `
	<h1>Tor Scraper Panel</h1>
	<a href="/logout">Güvenli Çıkış</a><br><br>
	<table border="1" cellpadding="5">
		<tr>
			<th>Kaynak Adı</th>
			<th>Başlık</th>
			<th>Kategori</th>
			<th>Kritiklik</th>
			<th>Tarih</th>
		</tr>
		{{range .}}
		<tr>
			<td><a href="/detail?id={{.ID}}">{{.Title}}</a></td>
<td>{{.SourceTitle}}</td>
<td>{{.Category}}</td>
<td>{{.Critical}}</td>
<td>{{.CreatedAt}}</td>
		</tr>
		{{end}}
	</table>
	<p><i>Siber Tehdit İstihbarat Paneli</i></p>
	`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	template.Must(template.New("d").Parse(tpl)).Execute(w, list)
}

func detay(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	var title, sourceTitle, content, crit string
	err := db.QueryRow(`
		SELECT title, source_title, content, criticality
		FROM findings WHERE id=$1
	`, id).Scan(&title, &sourceTitle, &content, &crit)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	tpl := `
	<h2>{{.Title}}</h2>
	<p><b>Kaynak:</b> {{.SourceTitle}}</p>
	<p><b>Mevcut Kritiklik:</b> {{.Crit}}</p>
	<form method="POST" action="/update?id={{.ID}}">
		<select name="crit">
			<option value="low">Low</option>
			<option value="medium">Medium</option>
			<option value="high">High</option>
		</select>
		<button type="submit">Güncelle</button>
	</form>
	<hr>
	<pre style="background:#eee; padding:10px;">{{.Content}}</pre>
	<br>
	<a href="/">Anasayfaya Dön</a>
	`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	template.Must(template.New("x").Parse(tpl)).Execute(w, map[string]any{
		"Title":       title,
		"SourceTitle": sourceTitle,
		"Content":     content,
		"Crit":        crit,
		"ID":          id,
	})
}

func updateCriticality(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	newCrit := r.FormValue("crit")

	_, err := db.Exec(`UPDATE findings SET criticality=$1 WHERE id=$2`, newCrit, id)
	if err != nil {
		http.Error(w, "Güncelleme hatası", 500)
		return
	}
	http.Redirect(w, r, "/detail?id="+strconv.Itoa(id), http.StatusSeeOther)
}

func cikis(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "auth",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func StartWebServer() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost port=5432 user=postgres password=sifreniz dbname=postgres sslmode=disable"
	}

	var err error
	db, err = sql.Open("postgres", "host=findings-postgres port=5432 user=postgres password=sifreniz dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("Veritabanına ulaşılamıyor: ", err)
	}

	http.HandleFunc("/login", giriskontrol)
	http.HandleFunc("/logout", cikis)
	http.HandleFunc("/", auth(anasayfa))
	http.HandleFunc("/detail", auth(detay))
	http.HandleFunc("/update", auth(updateCriticality))

	log.Println("Panel http://localhost:8080 adresinde çalışıyor...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
