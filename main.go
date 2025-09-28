package main



import (
    "fmt"
    "log"
    "net/http"

    "github.com/gorilla/mux"


    "lazy-start/runtimeconfig"
    "lazy-start/httpHandler"

)



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
    httphandler.RegisterRoutes(r)




	addr := fmt.Sprintf("%s:%d", runtimeconfig.App.Server.Host, runtimeconfig.App.Server.Port)

	fmt.Printf("Listening on http://%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}










