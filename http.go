package main

import (
	//"io"
	"fmt"
	"net/http"
	"strings"

	"github.com/FredFoonly/wingo/commands"
	"github.com/FredFoonly/wingo/logger"
	"github.com/gorilla/mux"
)

func httpAddress() string {
	// Take the http addr from the command line if possible
	if len(flagHttpAddr) > 0 {
		return strings.TrimSpace(flagHttpAddr)
	}

	// We weren't handed a path on a plate, so have to synthesize it as best we can
	return ":8080"
}

func serveHttp(httpAddr string) {
	r := mux.NewRouter()
	g := r.Methods("GET").Subrouter()

	g.HandleFunc("/Windows/Current", doGetActiveWindow)

	g.HandleFunc("/Clients", doGetAllClients)
	g.HandleFunc("/Clients/{client}/Height", doGetClientHeight)
	g.HandleFunc("/Clients/{client}/Name", doGetClientName)
	g.HandleFunc("/Clients/{client}/States", doGetClientStates)
	g.HandleFunc("/Clients/{client}/Type", doGetClientType)
	g.HandleFunc("/Clients/{client}/Width", doGetClientWidth)
	g.HandleFunc("/Clients/{client}/Workspace", doGetClientWorkspace)
	g.HandleFunc("/Clients/{client}/X", doGetClientX)
	g.HandleFunc("/Clients/{client}/Y", doGetClientY)

	g.HandleFunc("/Workspaces", doGetAllWorkspaces)
	g.HandleFunc("/Workspaces/Current/Name", doGetCurrentWorkspaceName)
	g.HandleFunc("/Workspaces/Current/Id", doGetCurrentWorkspaceId)
	g.HandleFunc("/Workspaces/Next/Name", doGetNextWorkspaceName)
	g.HandleFunc("/Workspaces/Prev/Name", doGetPrevWorkspaceName)
	g.HandleFunc("/Workspaces/{wsp}/Clients", doGetWorkspaceClients)
	g.HandleFunc("/Workspaces/{wsp}/Layout", doGetWorkspaceLayout)

	go http.ListenAndServe(httpAddr, r)
}

func doGetActiveWindow(w http.ResponseWriter, r *http.Request) {
	runCmds(w, "GetActive")
}

func doGetAllClients(w http.ResponseWriter, r *http.Request) {
	runCmds(w, "GetAllClients")
}

func doGetClientHeight(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	runCmds(w, fmt.Sprintf("GetClientHeight \"%s\"", vars["client"]))
}

func doGetClientName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	runCmds(w, fmt.Sprintf("GetClientName \"%s\"", vars["client"]))
}

func doGetClientStates(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	runCmds(w, fmt.Sprintf("GetClientStates \"%s\"", vars["client"]))
}

func doGetClientType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	runCmds(w, fmt.Sprintf("GetClientType \"%s\"", vars["client"]))
}

func doGetClientWidth(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	runCmds(w, fmt.Sprintf("GetClientWidth \"%s\"", vars["client"]))
}

func doGetClientWorkspace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	runCmds(w, fmt.Sprintf("GetClientWorkspace \"%s\"", vars["client"]))
}

func doGetClientX(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	runCmds(w, fmt.Sprintf("GetClientX \"%s\"", vars["client"]))
}

func doGetClientY(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	runCmds(w, fmt.Sprintf("GetClientY \"%s\"", vars["client"]))
}


func doGetAllWorkspaces(w http.ResponseWriter, r *http.Request) {
	runCmds(w, "GetWorkspaceList")
}

func doGetCurrentWorkspaceName(w http.ResponseWriter, r *http.Request) {
	runCmds(w, "GetWorkspace")
}

func doGetCurrentWorkspaceId(w http.ResponseWriter, r *http.Request) {
	runCmds(w, "GetWorkspaceId")
}

func doGetNextWorkspaceName(w http.ResponseWriter, r *http.Request) {
	runCmds(w, "GetWorkspaceNext")
}

func doGetPrevWorkspaceName(w http.ResponseWriter, r *http.Request) {
	runCmds(w, "GetWorkspacePrev")
}

func doGetWorkspaceClients(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	runCmds(w, fmt.Sprintf("GetClientList \"%s\"", vars["wsp"]))
}

func doGetWorkspaceLayout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	runCmds(w, fmt.Sprintf("GetLayout \"%s\"", vars["wsp"]))
}

func runCmds(w http.ResponseWriter, msg string) {
	logger.Lots.Printf("Running command from HTTP: '%s'.", msg)

	commands.Env.Verbose = true
	val, err := commands.Env.RunMany(msg)
	commands.Env.Verbose = false
	if err != nil {
		logger.Lots.Printf("ERROR running command: '%s'.", err)
		fmt.Fprintf(w, "ERROR: %s%c", err, 0)
		// set http error code as well
		return
	}

	// Fetch the return value of the command that was executed, and send
	// it back to the client. If the return value is nil, send an empty
	// response back. Otherwise, we need to type switch on all possible
	// return values.
	if val != nil {
		var retVal string
		switch v := val.(type) {
		case string:
			retVal = v
		case int:
			retVal = fmt.Sprintf("%d", v)
		case float64:
			retVal = fmt.Sprintf("%f", v)
		default:
			logger.Error.Fatalf("BUG: Unknown Gribble return type: %T", v)
		}

		fmt.Fprintf(w, "%s%c", retVal, 0)
	} else {
		fmt.Fprintf(w, "%c", 0)
	}
}
