package main



import (
    "fmt"
    "html/template"
    "log"
    "net/http"

    "github.com/gorilla/mux"


    "lazy-start/runtimeconfig"
)






var tmpl = template.Must(template.ParseFiles("templates/index.html"))

func main() {



    runtimeconfig.LoadAppConfig()

    log.Panicln(runtimeconfig.App.Server.Host)
    log.Panicln(runtimeconfig.App.Server.Port)
    log.Panicln(runtimeconfig.App.Paths.ConfigDir)
    log.Panicln(runtimeconfig.App.Paths.LogsDir)
    return




    err := LoadServiceConfigs("services.json")
    if err != nil {
        log.Fatalf("Error loading configs: %v", err)
    }

    r := mux.NewRouter()
    r.HandleFunc("/start/{service}", startHandler).Methods("GET")


    fmt.Println("Listening on http://localhost:7712")
    log.Fatal(http.ListenAndServe("127.0.0.1:7712", r))
}











func startHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    serviceName := vars["service"]
    

   
    svc, ok := GetService(serviceName)
    if !ok {
        http.Error(w, "Service not found", http.StatusNotFound)
        return
    }



    svc.mu.Lock()
    //isRunning := svc.IsRunning  //service is running 
    isStarting := svc.IsStarting  //service is beeing started 
	isDown := svc.IsDown // service  is perma down 

    svc.mu.Unlock()

	if isDown { 
		//redicrect to static not avalabe page
		return
	}





	

    if !isStarting {
        StartServiceIfNeeded(svc)
    }

	w.Header().Set("Content-Type", "text/html")
    tmpl.Execute(w, svc)
}

