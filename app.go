package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"main/packages/Memory/bytecode"
	"main/packages/Memory/instance"
	"main/packages/Memory/memory"

	"github.com/lxn/win"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/sys/windows"
)

var (
	bc bytecode.Bytecode
)

func (a *App) FullScreen() {
	if runtime.WindowIsFullscreen(a.ctx) {
		runtime.WindowUnfullscreen(a.ctx)
	} else {
		runtime.WindowFullscreen(a.ctx)
	}
}

func (a *App) Minimize() {
	runtime.WindowMinimise(a.ctx)
}

func (a *App) Close() {
	runtime.Quit(a.ctx)
}

type Resp struct {
	Error   string `json:"error"`
	Success bool   `json:"success"`
}

type NXConfig struct {
	Version        string `json:"version"`
	Active         bool   `json:"active"`
	OffsetsUpdated bool   `json:"offsets_updated"`
}

var Config NXConfig

func init() {
	resp, err := http.Get("https://raw.githubusercontent.com/suffz/luna/refs/heads/main/config.json")
	if err != nil {
		// handle the error here.
	}
	if resp.StatusCode == 200 {
		data, _ := io.ReadAll(resp.Body)
		json.Unmarshal(data, &Config)
		if AcceptedVersion != Config.Version {
			Config.Active = false
			Config.OffsetsUpdated = false
		}

	}
}

type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	//hwnd := win.FindWindow(nil, syscall.StringToUTF16Ptr("NXBeta"))
	//win.SetWindowLong(hwnd, win.GWL_EXSTYLE, win.GetWindowLong(hwnd, win.GWL_EXSTYLE)|win.WS_EX_LAYERED)
	go Start()
}

var (
	user32                       = syscall.NewLazyDLL("user32.dll")
	procGetForegroundWindow      = user32.NewProc("GetForegroundWindow")
	procSetForegroundWindow      = user32.NewProc("SetForegroundWindow")
	procEnumWindows              = user32.NewProc("EnumWindows")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
)

func GetForegroundWindow() syscall.Handle {
	ret, _, _ := procGetForegroundWindow.Call()
	return syscall.Handle(ret)
}

func SetForegroundWindow(hwnd syscall.Handle) bool {
	ret, _, _ := procSetForegroundWindow.Call(uintptr(hwnd))
	return ret != 0
}

func (a *App) Version() string {
	cl := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errors.New("Found")
		},
	}
	req, _ := http.NewRequest("GET", "https://www.roblox.com/download/client?os=win", nil)
	resp, err := cl.Do(req)
	if strings.Contains(err.Error(), "Found") {
		url, _ := resp.Location()
		s := url.String()
		if strings.Contains(s, "version-") {
			version := "version-" + strings.Split(s, "-")[1]
			return version
		}
	}
	return ""
}

func (a *App) Savefrontend(name, value string) {
	os.WriteFile("scripts/"+name, []byte(value), 0644)
}

func GetHWNDFromPID(pid uint32) windows.HWND {
	var foundHwnd windows.HWND

	enumWindowsCallback := syscall.NewCallback(func(hwnd uintptr, lParam uintptr) uintptr {
		var procID uint32
		procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&procID)))

		if procID == pid {
			foundHwnd = windows.HWND(hwnd)
			return 0
		}

		return 1
	})

	procEnumWindows.Call(enumWindowsCallback, 0)

	return foundHwnd
}

func initializeScriptHook() {
	hwnd := GetHWNDFromPID(uint32(Active.Pid))
	SetForegroundWindow(syscall.Handle(hwnd))
	memory.SendMessage(win.HWND(hwnd), 0x0010, 0, 0)
}

func (a *App) Attach() Resp {

	var Data Resp

	if !Config.Active && !Config.OffsetsUpdated {
		Data.Error = "Offsets currently not updated"
		Logger.Error(Data.Error)
		Data.Success = false
		return Data
	}

	if Active == nil {
		Data.Error = "No selected roblox instances."
		Logger.Error(Data.Error)
		Data.Success = false
		return Data
	}

	if Active.Mem == nil || Active.Error {
		Data.Error = "Client " + Active.Username + " had a error, please rechoose it in the menu."
		Logger.Error(Data.Error)
		Data.Success = false
		return Data
	}

	if !memory.IsHandleValid(Active.Mem.ProcessHandle) {
		mem, err := memory.NewLuna(uint32(Active.Pid))
		if err != nil {
			Active.Error = true
			Data.Error = "Client " + Active.Username + " | memory cant be restored."
			Logger.Error(Data.Error)
			Data.Success = false
			return Data
		}
		Active.Mem = mem
		Active.Error = false
	}

	if Active.Injected {
		Data.Error = "Client " + Active.Username + " is already injected!"
		Logger.Error(Data.Error)
		Data.Success = false
		return Data
	}

	code, size := bc.Compile(getInit())
	if size == -10 {
		Data.Error = string(code)
		Logger.Error(Data.Error)
		Data.Success = false
		return Data
	}

	fakedm, _ := Active.Mem.ReadPointer(uintptr(Active.Instances.RenderView) + uintptr(Active.Offsets.DataModelHolder))
	realdm, _ := Active.Mem.ReadPointer(fakedm + uintptr(Active.Offsets.DataModel))
	game := instance.NewInstance(realdm, Active)

	if game.Address < 1000 {
		Data.Error = "Datamodel came back nil, please make a ticket."
		Logger.Error(Data.Error)
		Data.Success = false
		return Data
	}

	mem := Active.Mem

	//corePackages := game.FindFirstChild("CorePackages", false)
	CoreGui := game.FindFirstChildOfClass("CoreGui", false)
	RobloxGui := CoreGui.FindFirstChild("RobloxGui", false)
	Modules := RobloxGui.FindFirstChild("Modules", false)

	t := game.Name()

	switch t {
	case "Game", "Ugc":

		player_list := Modules.FindFirstChild("PlayerList", false)
		player_list_manager := player_list.FindFirstChild("PlayerListManager", false)

		GetJestGlobals := func(game instance.Instance) instance.Instance {

			corePackages := game.FindFirstChild("CorePackages", false)

			Workspace := corePackages.FindFirstChild("Workspace", false)
			Packages := Workspace.FindFirstChild("Packages", false)
			_Workspace := Packages.FindFirstChild("_Workspace", false)
			SMSProtocol := _Workspace.FindFirstChild("SMSProtocol", false)
			Dev := SMSProtocol.FindFirstChild("Dev", false)
			jest_globals := Dev.FindFirstChild("JestGlobals", false)
			if jest_globals.Address > 1000 {
				return jest_globals
			}

			Pack := corePackages.FindFirstChild("Packages", false)
			DEV2 := Pack.FindFirstChild("Dev", false)
			Jest2 := DEV2.FindFirstChild("JestGlobals", false)
			if Jest2.Address > 1000 {
				return Jest2
			}

			Jest3 := corePackages.FindFirstChild("JestGlobals", false)
			if Jest3.Address > 1000 {
				return Jest3
			}

			StarterPlayer := game.FindFirstChild("StarterPlayer", false)
			StarterPlayerScripts := StarterPlayer.FindFirstChild("StarterPlayerScripts", false)
			PlayerModule := StarterPlayerScripts.FindFirstChild("PlayerModule", false)
			ControlModule := PlayerModule.FindFirstChild("ControlModule", false)
			VRNavigation := ControlModule.FindFirstChild("VRNavigation", false)
			if VRNavigation.Address > 1000 {
				return VRNavigation
			}

			return instance.NewInstance(0, Active)
		}

		jest_globals := GetJestGlobals(game)

		if jest_globals.Address < 1000 {
			Data.Error = "JestGlobals unable to be found. report to @http2 (rebooting roblox can fix this sometimes)"
			Logger.Error(Data.Error)
			Data.Success = false
			return Data
		}

		jest_globals.SetModuleBypass()

		old, old_len := jest_globals.GetByteCode()

		mem.WritePointer(player_list_manager.Address+0x8, jest_globals.Address)
		jest_globals.SetBytecode(code, uint64(size))
		initializeScriptHook()
		Bridge := CoreGui.WaitForChild("Bridge", 10)

		if Bridge.Address > 1000 {

			jest_globals.SetBytecode(old, uint64(old_len))
			mem.WritePointer(player_list_manager.Address+0x8, player_list_manager.Address)

			if strings.Contains(Active.Username, "Unknown-") {
				Players := Bridge.FindFirstChild("LunaPlayers", false)
				LunaPlayers := Players.Value()
				Active.Username = fmt.Sprintf("%v", LunaPlayers)
			}

			// FOR NOW COMMENT OUT BECAUSE WE STILL USE HTTP
			/*
				// init bridge
				BridgeValue = bridge.NewBridge(mem)

				BridgeValue.Start(int(Pid), Bridge)

				Callbacks(mem)
			*/

		} else {
			Data.Error = "Unable to find bridge!"
			Logger.Error(Data.Error)
		}
	}

	InGameMenu := Modules.FindFirstChild("InGameMenu", false)
	Network := InGameMenu.FindFirstChild("Network", false)

	Url := Network.FindFirstChild("Url", false)

	Url.SetModuleBypass()
	Url.SetBytecode(code, uint64(size))

	Active.Injected = true
	Data.Error = "Succesfully Attached"
	Logger.Error(Data.Error)
	Data.Success = true

	SetForegroundWindow(syscall.Handle(win.FindWindow(nil, syscall.StringToUTF16Ptr("Luna"))))

	return Data
}

func (a *App) ReturnClients() []*instance.RobloxInstances {
	return Roblox
}

func (a *App) Get() *instance.RobloxInstances {
	return Active
}

func (a *App) SetClient(pid int64) {
	for _, R := range Roblox {
		if R.Pid == pid {
			Active = R
			break
		}
	}
}

func (a *App) GetScriptsItems() map[string]string {
	var Info map[string]string = make(map[string]string)
	filepath.WalkDir("scripts", func(path string, d fs.DirEntry, err error) error {
		if d != nil {
			if !d.IsDir() {
				data, _ := os.ReadFile(path)
				Info[d.Name()] = string(data)
			}
		}
		return nil
	})
	return Info
}

var sendA sync.Mutex

func (a *App) Execute(source string) Resp {
	sendA.Lock()
	defer sendA.Unlock()

	var Data Resp

	if source == "" {
		Data.Error = "Source is empty!"
		Logger.Error(Data.Error)
		Data.Success = false
		return Data
	}

	if !Config.Active && !Config.OffsetsUpdated {
		Data.Error = "Offsets currently not updated"
		Logger.Error(Data.Error)
		Data.Success = false
		return Data
	}

	if Active == nil {
		Data.Error = "No selected roblox instances."
		Logger.Error(Data.Error)
		Data.Success = false
		return Data
	}

	if Active.Mem == nil || Active.Error {
		Data.Error = "Client " + Active.Username + " had a error, please rechoose it in the menu."
		Logger.Error(Data.Error)
		Data.Success = false
		return Data
	}

	if !memory.IsHandleValid(Active.Mem.ProcessHandle) {
		mem, err := memory.NewLuna(uint32(Active.Pid))
		if err != nil {
			Active.Error = true
			Data.Error = "Client " + Active.Username + " | memory cant be restored."
			Logger.Error(Data.Error)
			Data.Success = false
			return Data
		}
		Active.Mem = mem
		Active.Error = false
	}

	fakedm, _ := Active.Mem.ReadPointer(uintptr(Active.Instances.RenderView) + uintptr(Active.Offsets.DataModelHolder))
	realdm, _ := Active.Mem.ReadPointer(fakedm + uintptr(Active.Offsets.DataModel))
	datamodel := instance.NewInstance(realdm, Active)

	if name := datamodel.Name(); name != "Ugc" && name != "Game" {
		Data.Error = "Client " + Active.Username + " isnt in a game, execution wont work!"
		Logger.Error(Data.Error)
		Data.Success = false
		return Data
	}

	if !Active.Injected {
		Data.Error = "Client " + Active.Username + " is not injected!"
		Logger.Error(Data.Error)
		Data.Success = false
		return Data
	}

	if datamodel.Address == 0 {
		Data.Error = "An error occured while getting the datamodel, address null."
		Logger.Error(Data.Error)
		Data.Success = false
		return Data
	}

	CoreGui := datamodel.FindFirstChildOfClass("CoreGui", false)
	Bridge := CoreGui.FindFirstChild("Bridge", false)
	ScriptHolder := Bridge.FindFirstChild("ModuleHolder", false)
	Execution := Bridge.FindFirstChild("Execution", false)

	Logger.Info(fmt.Sprintf(`
Datamodel: %v
Coregui: %v
Bridge: %v
ScriptHolder: %v
ScriptHolderValue: %v

Memory Valid: %v

Current Client: %v
Source: %v
`,
		datamodel.String(),
		CoreGui.String(),
		Bridge.String(),
		ScriptHolder.String(),
		func() string {
			Val := ScriptHolder.Value()
			if reflect.TypeOf(Val).String() == "instance.Instance" {
				H := Val.(instance.Instance)
				return H.String()
			}
			return fmt.Sprintf("%v", Val)
		}(),
		memory.IsHandleValid(Active.Mem.ProcessHandle),
		Active,
		func() string {
			if len(source) > 25 {
				return source[:24]
			}
			return source
		}(),
	))

	if ScriptHolder.ClassName() == "ObjectValue" {
		if Holder := ScriptHolder.Value(); reflect.TypeOf(Holder).String() == "instance.Instance" {
			code, size := bc.Compile(fmt.Sprintf(`return function(...) 
%v 
end`, source))
			if size == -10 {
				Data.Error = string(code)
				Logger.Error(Data.Error)
				Data.Success = false
				return Data
			}
			Script := Holder.(instance.Instance)
			Script.SetModuleBypass()
			Script.SetBytecode(code, uint64(size))
			Data.Success = true
			Data.Error = "Successful"
			Logger.Error(Data.Error)
			Execution.SetValue(true)
		}
	} else {
		Data.Error = "Unable to get ModuleHolder, are you using bloxstrap? (ScriptHolder: NIL)"
		Logger.Error(Data.Error)
		Data.Success = false
	}

	return Data

}

func (t *App) IsInjected() bool {
	return Active.Injected
}

type UserData struct {
	SearchResults       []SearchResults `json:"searchResults"`
	NextPageToken       string          `json:"nextPageToken"`
	FilteredSearchQuery any             `json:"filteredSearchQuery"`
	Vertical            string          `json:"vertical"`
	Sorts               any             `json:"sorts"`
	PaginationMethod    any             `json:"paginationMethod"`
}
type Contents struct {
	Username          string `json:"username"`
	DisplayName       string `json:"displayName"`
	PreviousUsernames any    `json:"previousUsernames"`
	HasVerifiedBadge  bool   `json:"hasVerifiedBadge"`
	ContentType       string `json:"contentType"`
	ContentID         int64  `json:"contentId"`
	DefaultLayoutData any    `json:"defaultLayoutData"`
}
type SearchResults struct {
	ContentGroupType string     `json:"contentGroupType"`
	Contents         []Contents `json:"contents"`
	TopicID          string     `json:"topicId"`
}

type AvatarResponse struct {
	Data []struct {
		ImageURL string `json:"imageUrl"`
	} `json:"data"`
}

func OmniSearch(ign string) (string, error) {
	userSearchURL := fmt.Sprintf("https://apis.roblox.com/search-api/omni-search?verticalType=user&searchQuery=%v&pageToken=&globalSessionId=1&sessionId=1", ign)
	resp, err := http.Get(userSearchURL)

	if err != nil {
		return "", err
	}
	if resp.StatusCode == 429 {
		return "", errors.New("Error: 429")
	}

	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	return string(data), nil

}

func KillProcess(pid int64) error {
	processes, err := process.Processes()
	if err != nil {
		return err
	}
	for _, p := range processes {
		if int64(p.Pid) == pid {
			return p.Kill()
		}
	}
	return fmt.Errorf("process not found")
}

func (t *App) KillRoblox() Resp {

	var Data Resp

	if Active == nil {
		Data.Error = "No active client found."
		Logger.Error(Data.Error)
		Data.Success = false
		return Data
	}

	if err := KillProcess(Active.Pid); err != nil {
		Data.Error = err.Error()
		Logger.Error(Data.Error)
		return Data
	}

	Data.Success = true

	return Data

}

func (t *App) Avatar(ign string) string {
	data, err := OmniSearch(ign)
	if err != nil && err.Error() == "Error: 429" {
		for i := 0; i < 5; i++ {
			data, err = OmniSearch(ign)
			if err == nil {
				break
			}
			time.Sleep(time.Millisecond * 250)
		}
	}
	if err != nil {
		return "https://tr.rbxcdn.com/30DAY-AvatarHeadshot-768C0AD1061FB65545E5FF8D66815094-Png/150/150/AvatarHeadshot/Webp/noFilter"
	}

	var userSearch UserData

	json.Unmarshal([]byte(data), &userSearch)

	if len(userSearch.SearchResults) == 0 {
		return "https://tr.rbxcdn.com/30DAY-AvatarHeadshot-768C0AD1061FB65545E5FF8D66815094-Png/150/150/AvatarHeadshot/Webp/noFilter"
	}

	userID := userSearch.SearchResults[0].Contents[0].ContentID
	avatarURL := fmt.Sprintf("https://thumbnails.roblox.com/v1/users/avatar-headshot?userIds=%d&size=150x150&format=Png&isCircular=false", userID)
	avatarResp, err := http.Get(avatarURL)
	if err != nil {
		return err.Error()
	}
	defer avatarResp.Body.Close()

	if avatarResp.StatusCode == 429 {
		for i := 0; i < 5; i++ {
			avatarURL := fmt.Sprintf("https://thumbnails.roblox.com/v1/users/avatar-headshot?userIds=%d&size=150x150&format=Png&isCircular=false", userID)
			avatarResp, err := http.Get(avatarURL)
			if err != nil {
				return err.Error()
			}
			defer avatarResp.Body.Close()
			if avatarResp.StatusCode == 200 {
				break
			}
			time.Sleep(time.Millisecond * 250)
		}
	}

	var avatarData AvatarResponse
	data2, _ := io.ReadAll(avatarResp.Body)
	json.Unmarshal(data2, &avatarData)
	if len(avatarData.Data) == 0 {
		return ""
	}
	return avatarData.Data[0].ImageURL
}
