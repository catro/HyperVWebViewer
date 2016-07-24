package main

import (
    "bytes"
    "encoding/json"
    "html/template"
    "net/http"
    "os/exec"
    "regexp"
    "strconv"
)

// VM Infomation struct from powershell command 'Get-VM'.
type VM struct {
    Name           string `json:"Name"`
    State          int    `json:"State"`
    CPUUsage       int    `json:"CPUUsage"`
    MemoryAssigned int    `json:"MemoryAssigned"`
    Uptime         Uptime `json:"Uptime"`
    Notes          string `json:"Notes"`
    StateDesc      string `json:"StateDesc"`
    IsRunning      bool   `json:"IsRunning"`
    IsSuspended    bool   `json:"IsSuspended"`
    HHMMSS         string `json:"HHMMSS"`
}

// Uptime is the VM's Uptime
type Uptime struct {
    Days    int `json:"Days"`
    Hours   int `json:"Hours"`
    Minutes int `json:"Minutes"`
    Seconds int `json:"Seconds"`
}

// VMs is a slice of structs.
type VMs struct {
    vm []VM
}

func main() {
    http.HandleFunc("/public/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, r.URL.Path[1:])
    })
    http.HandleFunc("/", indexHandler)
    http.HandleFunc("/start", startVMHandler)
    http.HandleFunc("/resume", resumeVMHandler)
    http.HandleFunc("/save", saveVMHandler)
    http.HandleFunc("/suspend", suspendVMHandler)

    // set the port and start server.
    http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    // change the charaset encode of powershell to UTF-8.
    errEnc := setEncodeUtf8()
    if errEnc != nil {
        errMsg := errEnc.Error() + "\nCouldn't encode to UTF-8."
        http.Error(w, errMsg, http.StatusInternalServerError)
        return
    }

    // get the VM's information from powershell command "Get-VM".
    cmd := exec.Command("powershell", "Get-VM | Select-Object Name, State, CPUUsage, MemoryAssigned, Uptime, Notes | ConvertTo-Json")
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        errMsg := err.Error() + "\nCouldn't get the vm information."
        http.Error(w, errMsg, http.StatusInternalServerError)
        return
    }

    // decode JSON.
    dec := json.NewDecoder(bytes.NewReader(out.Bytes()))
    var vms VMs
    errDec := dec.Decode(&vms.vm)
    if errDec != nil {
        errMsg := err.Error() + "\nCouldn't decode JSON."
        http.Error(w, errMsg, http.StatusInternalServerError)
        return
    }

    // set the VM Struct's element values.
    setVMElem(&vms)

    // render a template web page.
    var indexTemplate = template.Must(template.ParseFiles("./public/index.html"))
    indexTemplate.Execute(w, vms.vm)
}

// start vm from POST using powershell command "Start-VM".
func startVMHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        vmName := r.FormValue("vmName")
        // check an unsupported charactert in the parameter
        if hasUnsupportedChar(vmName) {
            errMsg := "\nThe post parameter has an unsupported charactoer. Please check the parameter."
            http.Error(w, errMsg, http.StatusBadRequest)
            return
        }
        pscmd := "Start-VM -name " + vmName
        cmd := exec.Command("powershell", pscmd)
        var out bytes.Buffer
        cmd.Stdout = &out
        err := cmd.Run()
        if err != nil {
            errMsg := err.Error() + "\nCouldn't start the VM.\nThe memory on the host server is not enough to start the vm, probably."
            http.Error(w, errMsg, http.StatusInternalServerError)
            return
        }
        http.Redirect(w, r, "/", http.StatusSeeOther)
    } else {
        http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
    }
}

// resume vm from POST using powershell command "Resume-VM".
func resumeVMHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        vmName := r.FormValue("vmName")
        // check an unsupported charactert in the parameter
        if hasUnsupportedChar(vmName) {
            errMsg := "\nThe post parameter has an unsupported charactoer. Please check the parameter."
            http.Error(w, errMsg, http.StatusBadRequest)
            return
        }
        pscmd := "Resume-VM -name " + vmName
        cmd := exec.Command("powershell", pscmd)
        var out bytes.Buffer
        cmd.Stdout = &out
        err := cmd.Run()
        if err != nil {
            errMsg := err.Error()
            http.Error(w, errMsg, http.StatusInternalServerError)
            return
        }
        http.Redirect(w, r, "/", http.StatusSeeOther)
    } else {
        http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
    }
}

// save vm from POST using powershell command "Save-VM".
func saveVMHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        vmName := r.FormValue("vmName")
        // check an unsupported charactert in the parameter
        if hasUnsupportedChar(vmName) {
            errMsg := "\nThe post parameter has an unsupported charactoer. Please check the parameter."
            http.Error(w, errMsg, http.StatusBadRequest)
            return
        }
        pscmd := "Save-VM -name " + vmName
        cmd := exec.Command("powershell", pscmd)
        var out bytes.Buffer
        cmd.Stdout = &out
        err := cmd.Run()
        if err != nil {
            errMsg := err.Error()
            http.Error(w, errMsg, http.StatusInternalServerError)
            return
        }
        http.Redirect(w, r, "/", http.StatusSeeOther)
    } else {
        http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
    }
}

// suspend vm from POST using powershell command "Suspend-VM".
func suspendVMHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        vmName := r.FormValue("vmName")
        // check an unsupported charactert in the parameter
        if hasUnsupportedChar(vmName) {
            errMsg := "\nThe post parameter has an unsupported charactoer. Please check the parameter."
            http.Error(w, errMsg, http.StatusBadRequest)
            return
        }
        pscmd := "Suspend-VM -name " + vmName
        cmd := exec.Command("powershell", pscmd)
        var out bytes.Buffer
        cmd.Stdout = &out
        err := cmd.Run()
        if err != nil {
            errMsg := err.Error()
            http.Error(w, errMsg, http.StatusInternalServerError)
            return
        }
        http.Redirect(w, r, "/", http.StatusSeeOther)
    } else {
        http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
    }
}

// setVMElem set the vm struct's element values
// such as StateDesc, IsRunning, HHMMSS.
func setVMElem(vms *VMs) {
    for i := range vms.vm {
        vm := &vms.vm[i]
        vm.MemoryAssigned = vm.MemoryAssigned / 1024 / 1024 // MB
        vm.StateDesc, vm.IsRunning, vm.IsSuspended = parseState(vm.State)
        vm.HHMMSS = toHHMMSS(&vm.Uptime)
    }
}

// The parseState check the vm's state
// and return the stateDesc and the boolean value that the vm is runnning or not
func parseState(st int) (state string, isRun bool, isSuspended bool) {
    isSuspended = false
    switch st {
    case 0:
        state = "Unknown"
        isRun = false
    case 2:
        state = "Running"
        isRun = true
    case 3:
        state = "Stopped"
        isRun = false
    case 6:
        state = "Saving"
        isRun = false
    case 9:
        state = "Suspended"
        isRun = false
        isSuspended = true
    default:
        state = "other"
        isRun = false
    }
    return state, isRun, isSuspended
}

// setEncodeUtf8 change a powershell output's enode to utf-8
// because golang cannot read utf-8 directory.
func setEncodeUtf8() error {
    cmd := exec.Command("chcp", "65001")
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    return err
}

// toHHMMSS change to "HHMMSS" string from vm's uptime
func toHHMMSS(uptime *Uptime) string {
    var D, hh, mm, ss string
    if uptime.Days > 0 {
        D = strconv.Itoa(uptime.Days) + "."
    }
    if uptime.Hours < 10 {
        hh = "0" + strconv.Itoa(uptime.Hours)
    } else {
        hh = strconv.Itoa(uptime.Hours)
    }
    if uptime.Minutes < 10 {
        mm = "0" + strconv.Itoa(uptime.Minutes)
    } else {
        mm = strconv.Itoa(uptime.Minutes)
    }
    if uptime.Seconds < 10 {
        ss = "0" + strconv.Itoa(uptime.Seconds)
    } else {
        ss = strconv.Itoa(uptime.Seconds)
    }
    return D + hh + ":" + mm + ":" + ss
}

// hasUnsupportedChar check the parameters that contain unsupported characters
func hasUnsupportedChar(str string) bool {
    reg := `[^0-9a-zA-Z-_\.\s]` // add some characters that contain in the vm names
    return regexp.MustCompile(reg).Match([]byte(str))
}
