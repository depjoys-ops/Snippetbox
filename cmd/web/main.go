package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	_ "github.com/go-sql-driver/mysql"

	"github.com/depjoys-ops/Snippetbox/internal/models"
)

type application struct {
    logger			*slog.Logger
	snippets		*models.SnippetModel
	templateCache	map[string]*template.Template
	sessionManager	*scs.SessionManager
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:123456@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()
	

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level: slog.LevelInfo,
	}))
	

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()
	

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}


	sessionManager := scs.New()
    sessionManager.Store = mysqlstore.New(db)
    sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	app := &application{
        logger: logger,
		snippets: &models.SnippetModel{DB: db},
		templateCache: templateCache,
		sessionManager: sessionManager,
	}


	tlsConfig := &tls.Config{
        CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion: tls.VersionTLS13,
    	MaxVersion: tls.VersionTLS13,
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
		},
    }
	srv := &http.Server{
		Addr: *addr,
		Handler: app.routes(),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig: tlsConfig,
	}

	logger.Info("starting server", "addr", srv.Addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem","./tls/key.pem")
    logger.Error(err.Error())
	os.Exit(1)

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}