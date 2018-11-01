package main

import (
	"crypto/tls"
	"log"
	"os"
	"time"

	"net/http"

	"github.com/gorilla/mux"

	libDatabox "github.com/me-box/lib-go-databox"
)

type config struct {
	cmAPIServiceStatus    libDatabox.DataSourceMetadata
	cmAPIDataSource       libDatabox.DataSourceMetadata
	cmDataDataSource      libDatabox.DataSourceMetadata
	appsDatasource        libDatabox.DataSourceMetadata
	driverDatasource      libDatabox.DataSourceMetadata
	allManifests          libDatabox.DataSourceMetadata
	cmStoreEndpoint       string
	manifestStoreEndpoint string
}

func main() {
	libDatabox.Info("Starting ....")
	//Read in the databox data passed to the app
	cmAPIDataSource, DATABOX_ZMQ_ENDPOINT_CM, _ := libDatabox.HypercatToDataSourceMetadata(os.Getenv("DATASOURCE_CM_API"))
	cmAPIServiceStatus, _, _ := libDatabox.HypercatToDataSourceMetadata(os.Getenv("CM_API_ServiceStatus"))
	cmDataDataSource, _, _ := libDatabox.HypercatToDataSourceMetadata(os.Getenv("DATASOURCE_CM_DATA"))
	appsDatasource, DATABOX_ZMQ_ENDPOINT_APP, _ := libDatabox.HypercatToDataSourceMetadata(os.Getenv("DATASOURCE_APPS"))
	driverDatasource, _, _ := libDatabox.HypercatToDataSourceMetadata(os.Getenv("DATASOURCE_DRIVERS"))
	allManifests, _, _ := libDatabox.HypercatToDataSourceMetadata(os.Getenv("DATASOURCE_ALL"))

	//create a config object to pass to handlers
	cfg := config{
		cmAPIServiceStatus:    cmAPIServiceStatus,
		cmAPIDataSource:       cmAPIDataSource,
		cmDataDataSource:      cmDataDataSource,
		appsDatasource:        appsDatasource,
		driverDatasource:      driverDatasource,
		allManifests:          allManifests,
		cmStoreEndpoint:       DATABOX_ZMQ_ENDPOINT_CM,
		manifestStoreEndpoint: DATABOX_ZMQ_ENDPOINT_APP,
	}

	//setup webserver routes
	router := mux.NewRouter()

	//HTTPS API
	router.HandleFunc("/status", statusEndpoint).Methods("GET")
	router.HandleFunc("/ui/api/appStore", getApps(&cfg)).Methods("GET")
	router.HandleFunc("/ui/api/containerStatus", containerStatus(&cfg)).Methods("GET")
	router.HandleFunc("/ui/api/containerStatus2", containerStatus2(&cfg)).Methods("GET")
	router.HandleFunc("/ui/api/dataSources", dataSources(&cfg)).Methods("GET")
	router.HandleFunc("/ui/api/drivers", getDrivers(&cfg)).Methods("POST")
	router.HandleFunc("/ui/api/manifest/{name}", getManifest(&cfg)).Methods("GET")
	router.HandleFunc("/ui/api/install", install(&cfg)).Methods("POST")
	router.HandleFunc("/ui/api/uninstall", uninstall(&cfg)).Methods("POST")
	router.HandleFunc("/ui/api/restart", restart(&cfg)).Methods("POST")
	router.HandleFunc("/ui/api/qrcode.png", qrcode(&cfg)).Methods("GET")
	router.HandleFunc("/ui/cert.pem", certPub(&cfg)).Methods("GET")
	router.HandleFunc("/ui/cert.der", certPubDer(&cfg)).Methods("GET")
	router.PathPrefix("/ui/{css|icons|js}").Handler(http.StripPrefix("/ui", http.FileServer(http.Dir("./www")))).Methods("GET")
	router.PathPrefix("/").Handler(http.HandlerFunc(serveIndex)).Methods("GET")

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
		},
	}

	srv := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
		TLSConfig:    tlsConfig,
		Handler:      router,
	}
	libDatabox.Info("Waiting for https requests ....")
	log.Fatal(srv.ListenAndServeTLS(libDatabox.GetHttpsCredentials(), libDatabox.GetHttpsCredentials()))
	libDatabox.Info("Exiting ....")
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	libDatabox.Info(r.RequestURI)
	http.ServeFile(w, r, "./www/index.html")
}
