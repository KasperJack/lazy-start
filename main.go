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

    //log.Println(runtimeconfig.App.Server.Host)
    //log.Println(runtimeconfig.App.Server.Port)
    //log.Println(runtimeconfig.App.Paths.ConfigDir)
    //og.Println(runtimeconfig.App.Paths.LogsDir)
  




    err := LoadServiceConfigs("services.json")
    if err != nil {
        log.Fatalf("Error loading configs: %v", err)
    }

    r := mux.NewRouter()
    r.HandleFunc("/start/{service}", startHandler).Methods("GET")





	addr := fmt.Sprintf("%s:%d", runtimeconfig.App.Server.Host, runtimeconfig.App.Server.Port)

	fmt.Printf("Listening on http://%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
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

