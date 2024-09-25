package main

import (
	"context"
	"fmt"
	"github.com/LuisRiveraBan/gocourse_user/internal"
	"github.com/LuisRiveraBan/gocourse_user/pkg/bootstrap"
	"github.com/LuisRiveraBan/gocourse_user/pkg/handler"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"time"
)

func main() {

	// Connect to the database
	_ = godotenv.Load()

	l := bootstrap.InitLogger()

	db, err := bootstrap.ConnectToDatabase()
	if err != nil {
		fmt.Println("Error connecting to the database", err)
		os.Exit(1)
	}

	pagLimDef := os.Getenv("PAGINATOR_LIMIT_DEFAULT")

	if pagLimDef == "" {
		l.Fatal("REQUESTER VARIABLE")
	}
	// Define the user endpoints and handlers.
	ctx := context.Background()
	userRepo := user.NewRepository(l, db)
	userSrv := user.NewService(l, userRepo)
	//userEnd := user.MakeEndpoints(userSrv, user.Config{LimPageDef: pagLimDef})
	h := handler.NewUserHTTPServer(ctx, user.MakeEndpoints(userSrv, user.Config{LimPageDef: pagLimDef}))

	port := os.Getenv("PORT")

	address := fmt.Sprintf("127.0.0.1:%s", port)

	//Manejo de tiempo agotado en escucha y respuesta
	srv := &http.Server{
		Handler: accesControl(h),
		Addr:    address,
		//Tiempo de escritura
		WriteTimeout: 15 * time.Second,
		//Tiempo de lectura
		ReadTimeout: 15 * time.Second,
	}

	errCh := make(chan error)
	go func() {
		l.Println("listen in ", address)
		errCh <- srv.ListenAndServe()
	}()

	err = <-errCh
	if err != nil {
		l.Fatal("ListenAndServe: ", err)
	}

}

func accesControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}
