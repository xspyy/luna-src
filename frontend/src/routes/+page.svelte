<script>
  import Sidebarbuttons from "../lib/sidebar/sidebarbuttons.svelte";
  import Search from "../lib/search/search.svelte";
  import Swal from "sweetalert2";
  import * as monaco from "monaco-editor";
    import { pushState } from "$app/navigation";

  window.Swal = Swal;
  window.monaco = monaco
</script>

<body
  style="-webkit-touch-callout: none;
  -webkit-user-select: none;
  -khtml-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  user-drag: none;
  -webkit-user-drag: none;
  user-select: none;
  -ms-user-select: none;margin:0"
>
  <div class="toast-nx-container"></div>

  <!--<div id="loading-screen">AWDAWD
    <div class="loader"></div>
  </div>-->

  <div id="topbarandcontainer">
    <div>
      <div
        style="--wails-draggable:drag;border-bottom: solid #6c678047 1px;"
        class="topbar"
      >
        <div style="display: flex;gap:5px">
          <img
            alt="NX"
            src="ik.png"
            style="margin-left:5px;width:55px;height:55px"
          />
        </div>
        <div style="display: flex;gap:5px;margin-left:5px">
          <Sidebarbuttons
            customContent={true}
            MarginTop={false}
            iconpath="themes.svg"
            onclicked="window.go.main.App.KillActiveRoblox()"
          >
            <img style="width: 20px;" alt="right" src="closeroblox.svg" />
          </Sidebarbuttons>
          <Sidebarbuttons
            customContent={true}
            MarginTop={false}
            onclicked="openClient()"
            iconpath="themes.svg"
          >
            <img style="width: 20px;" alt="right" src="dropdown.svg" />
          </Sidebarbuttons>
        </div>
        <div style="display: flex; gap: 6px;margin-right: 5px;">
          <Sidebarbuttons
            customContent={true}
            MarginTop={false}
            iconpath="themes.svg"
            onclicked="window.go.main.App.Minimize()"
          >
            <img style="width: 20px;" alt="right" src="minimize.svg" />
          </Sidebarbuttons>
          <Sidebarbuttons
            customContent={true}
            MarginTop={false}
            iconpath="themes.svg"
            onclicked="window.go.main.App.Close()"
          >
            <img style="width: 20px;" alt="right" src="close.svg" />
          </Sidebarbuttons>
        </div>
      </div>
      <div class="sidebar" id="buttons" style="border-right: solid #6c678047 1.5px;">
        <div id="home" class="btn-pushable">
          <Sidebarbuttons iconpath="home.svg"></Sidebarbuttons>
        </div>

        <div id="script"class="btn-pushable">
          <Sidebarbuttons iconpath="editor.svg" classCustom="active"
          ></Sidebarbuttons>
        </div>

        <div id="script-hub"class="btn-pushable">
          <Sidebarbuttons iconpath="scripts.svg"></Sidebarbuttons>
        </div>

        <div id="themes"class="btn-pushable">
          <Sidebarbuttons iconpath="themes.svg"></Sidebarbuttons>
        </div>

        <div
          class="circle red"
          style="margin-top: auto;margin-bottom:25px;"
          id="updateExecutionStatus"
        ></div>
      </div>
    </div>
    <div id="main-container" style="display: flex;flex-grow: 1;">
      <div class="container" id="home-page" style="display:none;">
        <div class="sidecontainer">
          <Search></Search>

          <div class="home-container">
            <div class="top-section">
              <img alt="right" src="cloud.svg" />
              <div class="text-section">
                <span class="main-text">Vulnerability found in</span>
                <span class="sub-text">...</span>
              </div>
            </div>
            <div class="bottom-section">
              <div class="date-text">9/7/2024</div>
            </div>
          </div>

          <div class="home-container">
            <div class="top-section">
              <img alt="right" src="cloud.svg" />
              <div class="text-section">
                <span class="main-text">Vulnerability found in</span>
                <span class="sub-text">...</span>
              </div>
            </div>
            <div class="bottom-section">
              <div class="date-text">9/7/2024</div>
            </div>
          </div>

          <div class="home-container">
            <div class="top-section">
              <img alt="right" src="cloud.svg" />
              <div class="text-section">
                <span class="main-text">Vulnerability found in</span>
                <span class="sub-text">...</span>
              </div>
            </div>
            <div class="bottom-section">
              <div class="date-text">9/7/2024</div>
            </div>
          </div>
        </div>
        <div
          class="home-stats-container"
          style="overflow-x: hidden;z-index:999;"
        >
          <div class="home-stats-bun">
            <div class="home-stats-bun-item" style="width: 77px; height: 77px;">
              <img
                id="avatarUrl"
                style="width:77px;height:77px;"
                src="https://tr.rbxcdn.com/30DAY-AvatarHeadshot-768C0AD1061FB65545E5FF8D66815094-Png/150/150/AvatarHeadshot/Webp/noFilter"
              />
            </div>
            <div class="home-stats-bun-item" style="width: 100%; height: 77px;">
              <div class="home-stats-info" style="margin-top:20px;">
                <img
                  alt="right"
                  src="check.svg"
                  style="margin-left: 20px;"
                /><span style="color:white;">,</span>
                <div class="text-section">
                  <span
                    class="main-text"
                    style="color:#E2E2E2"
                    id="injectedUserName">placeholder</span
                  >
                  <div style="position:absolute;">
                    <span class="text-status">Current injection status: </span>
                    <span class="status-home" id="injectedStatus"
                      >not injected</span
                    >
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="home-stats-bun">
            <div class="home-stats-bun-item" style="width: 100%; height: 45px;">
              <div
                style="
      display: flex;
      flex-direction: row;
      gap: 100px;
      justify-content: center;
      contain: content;
      margin: auto;
      position: absolute;
      width: 95.5%;
      padding: 14px;
                "
              >
                <span class="home-label-button">On Going Banwaves</span>
                <span class="home-label-button">Update Log</span>
                <span class="home-label-button">Announcements</span>
              </div>
            </div>
          </div>

          <!-- Containers with information about ban waves. -->

          <div style="margin-bottom:2px;"></div>

          <div class="home-stats-bun">
            <div
              class="home-stats-bun-item"
              style="width: 100%; height: 77px;background:radial-gradient(100% 100% at 51.43% 100%, rgba(42, 41, 49, 0.50) 0%, rgba(42, 41, 49, 0.00) 100%), rgba(27, 26, 32, 0.85);"
            >
              <div
                class="home-stats-info"
                style="margin-top:10px;margin-left:10px"
              >
                <div class="text-section">
                  <div
                    style="content: center;display: flex;gap: 7px;padding: 1px;"
                  >
                    <span class="main-text" style="color:#E2E2E2;font-size:13px"
                      >September 23rd</span
                    >

                    <img alt="right" src="warning.svg" />

                    <span
                      class="main-text text-yellow"
                      style="color:#E2E2E2;font-size:13px">23 Banned</span
                    >
                  </div>

                  <div style="width: 250px;margin-top: -5px;">
                    <span class="text-status"
                      >This banwave started on september 23rd and is currently</span
                    > <span class="text-status-yellow">on going</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="home-stats-bun">
            <div
              class="home-stats-bun-item"
              style="width: 100%; height: 77px;background:radial-gradient(100% 100% at 51.43% 100%, rgba(42, 41, 49, 0.50) 0%, rgba(42, 41, 49, 0.00) 100%), rgba(27, 26, 32, 0.85);"
            >
              <div
                class="home-stats-info"
                style="margin-top:10px;margin-left:10px"
              >
                <div class="text-section">
                  <div
                    style="content: center;display: flex;gap: 7px;padding: 1px;"
                  >
                    <span class="main-text" style="color:#E2E2E2;font-size:13px"
                      >September 10th</span
                    >

                    <img alt="right" src="progress.svg" />

                    <span class="main-text" style="color:#E2E2E2;font-size:13px"
                      >September 23rd</span
                    >

                    <img alt="right" src="successful.svg" />
                    <span
                      class="main-text text-green"
                      style="color:#E2E2E2;font-size:13px"
                      >2.25k Banned Total</span
                    >
                  </div>

                  <div style="width: 250px;margin-top: -5px;">
                    <span class="text-status"
                      >This banwave started on september 23rd and is currently</span
                    >
                    <span class="text-status-green" style="margin-top:5px"
                      >finished</span
                    >
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="home-stats-bun">
            <div
              class="home-stats-bun-item"
              style="width: 100%; height: 77px;background:radial-gradient(100% 100% at 51.43% 100%, rgba(42, 41, 49, 0.50) 0%, rgba(42, 41, 49, 0.00) 100%), rgba(27, 26, 32, 0.85);"
            >
              <div
                class="home-stats-info"
                style="margin-top:10px;margin-left:10px"
              >
                <div class="text-section">
                  <div
                    style="content: center;display: flex;gap: 7px;padding: 1px;"
                  >
                    <span class="main-text" style="color:#E2E2E2;font-size:13px"
                      >September 10th</span
                    >

                    <img alt="right" src="progress.svg" />

                    <span class="main-text" style="color:#E2E2E2;font-size:13px"
                      >September 23rd</span
                    >

                    <img alt="right" src="successful.svg" />
                    <span
                      class="main-text text-green"
                      style="color:#E2E2E2;font-size:13px"
                      >2.25k Banned Total</span
                    >
                  </div>

                  <div style="width: 250px;margin-top: -5px;">
                    <span class="text-status"
                      >This banwave started on september 23rd and is currently</span
                    >
                    <span class="text-status-green" style="margin-top:5px"
                      >finished</span
                    >
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="container" id="home-scripts">
        <div class="sidecontainer" id="applyfiles">
          <Search></Search>
        </div>
          <div id="editor-container">
            <div id="tabs" style="position:absolute;width:calc(100% - 60px);display: -webkit-inline-box;height: 50px;flex-direction: row;overflow-x:hidden;overflow-y:hidden;white-space:nowrap;">
              <div id="add-tabs" style="position: fixed;height:50px;width:60px;right:0;border-left: 1px solid rgba(108, 103, 128, 0.25);z-index: 1;">
                <div id="tab-add" class="tab-add" onclick="addTab()"></div>
              </div>
            </div>
            <div id="editor" style="width:100%;display:flex;"></div>
            <div
              style="width: 100%;
      position: fixed;
      display: flex;
      gap: 10px;
      bottom: 0;
      left: 0;
      margin: 15px;"
            >
              <div onclick="relay.setValue('')" class="btn-pushable">
                <Sidebarbuttons
                  customContent={true}
                  MarginTop={false}
                  iconpath="themes.svg"
                >
                  CLEAR
                </Sidebarbuttons>
              </div>
  
              <div onclick="openFile()" class="btn-pushable">
                <Sidebarbuttons
                  customContent={true}
                  MarginTop={false}
                  iconpath="themes.svg"
                >
                  OPEN FILE
                </Sidebarbuttons>
              </div>
              <div id="savefile" class="btn-pushable">
                <Sidebarbuttons
                  customContent={true}
                  MarginTop={false}
                  iconpath="themes.svg"
                >
                  SAVE FILE
                </Sidebarbuttons>
              </div>
              <div onclick="execute()" class="btn-pushable">
                <Sidebarbuttons
                  customContent={true}
                  MarginTop={false}
                  iconpath="themes.svg"
                >
                  EXECUTE
                </Sidebarbuttons>
              </div>
  
              <div onclick="attach()" style="position: absolute;right: 25px;" class="btn-pushable">
                <Sidebarbuttons customContent={true} MarginTop={false}>
                  <img src="attach.svg" />
                </Sidebarbuttons>
              </div>
            </div>
        </div>
      </div>

      <div class="container" id="home-settings" style="display:none;">
        <div class="sidecontainer">
          <Search></Search>

          <div
            class="settings-container"
            style="border-bottom: 1px solid rgba(108, 103, 128, 0.25);
      background: radial-gradient(100% 100% at 51.43% 100%, rgb(42, 41, 49) 0%, rgba(42, 41, 49, 0) 100%);
      backdrop-filter: blur(15px);"
          >
            <div class="middle-section">
              <img
                alt="right"
                src="box.svg"
                style="margin-top: -7px;margin-left:-4px;margin-right: -5px;"
              />
              <div class="text-section">
                <span class="main-text">General</span>
              </div>
            </div>
          </div>

          <div class="settings-container">
            <div class="middle-section">
              <img alt="right" src="core.svg" style="margin-top: 2px;" />
              <div class="text-section">
                <span class="main-text">Core</span>
              </div>
            </div>
          </div>

          <div class="settings-container">
            <div class="middle-section">
              <img alt="right" src="ui.svg" style="margin-top: 2px;" />
              <div class="text-section">
                <span class="main-text">UI</span>
              </div>
            </div>
          </div>
        </div>
        <div class="home-stats-container" style="z-index:999">
          <div class="home-stats-bun">
            <div
              class="home-stats-bun-item"
              style="width: 100px;height: 30px;contain: content;background:none;border:none;"
            >
              <div class="top-section" style="position:static;">
                <div class="text-section">
                  <span class="main-text" style="font-size:16px">General</span>
                </div>
              </div>
            </div>
          </div>

          <div class="home-stats-bun">
            <div
              class="home-stats-bun-item"
              style="width: 100%;height: 93px;contain: content;background:none;"
            >
              <div class="top-section">
                <img alt="right" src="ui.svg" />
                <div class="text-section">
                  <span class="main-text">UI</span>
                </div>
              </div>

              <div class="middle-section" style="top:45%">
                <span class="text-status" style="width: 400px;"
                  >Yea idk what to put here but yea like skbiidi
                  truewthewtoewtewoitoiewoitowiehto
                </span>
              </div>

              <div class="right-section" style="right:15px;">
                <div class="toggle-container active">
                  <div class="toggle-button">OFF</div>
                </div>
              </div>
            </div>
          </div>

          <div class="home-stats-bun">
            <div
              class="home-stats-bun-item"
              style="width: 100%;height: 93px;contain: content;background:none;"
            >
              <div class="top-section">
                <img alt="right" src="ui.svg" />
                <div class="text-section">
                  <span class="main-text">UI</span>
                </div>
              </div>

              <div class="middle-section" style="top:45%">
                <span class="text-status" style="width: 400px;"
                  >Yea idk what to put here but yea like skbiidi
                  truewthewtoewtewoitoiewoitowiehto
                </span>
              </div>

              <div class="right-section" style="right:15px;">
                <div class="toggle-container">
                  <div class="toggle-button">ON</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="container" id="home-scripthub" style="display:none;">
        <div class="sidecontainer">
          <Search></Search>
        </div>
      </div>
    </div>
  </div>

  <script>

    const scrollContainer = document.getElementById("tabs");

    scrollContainer.addEventListener("wheel", (event) => {
      event.preventDefault();

      const scrollSpeed = 775;
      const startPosition = scrollContainer.scrollLeft;
      const targetPosition = startPosition + event.deltaY;
      const distance = Math.abs(targetPosition - startPosition);
      const duration = distance / scrollSpeed * 1000;

      let startTime = null;

      function animateScroll(currentTime) {
        if (!startTime) startTime = currentTime;
        const elapsed = currentTime - startTime;
        const progress = Math.min(elapsed / duration, 1);
        const easedProgress = easeOutQuad(progress);

        scrollContainer.scrollLeft = startPosition + (targetPosition - startPosition) * easedProgress;

        if (progress < 1) {
          requestAnimationFrame(animateScroll);
        }
      }

      requestAnimationFrame(animateScroll);
    });

    function easeOutQuad(t) {
      return t * (2 - t);
    }


    let makePopupHTML = (content, extrahtml, action, action_) => {
      const id = new Date().getTime();

      return `<div class="modal" id="${id}" style="-webkit-touch-callout: none;
  -webkit-user-select: none;
  -khtml-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  user-drag: none;
  -webkit-user-drag: none;
  user-select: none;
  -ms-user-select: none;">
                    <div class="modal-box">
                        ${content}
                        <div class="modal-content">
                            ${extrahtml}
                        </div>
                        <div class="modal-action">
                            ${action.replace("{id}", id)} 
                            ${action_.replace("{id}", id)}
                        </div>
                    </div>
                </div>`;
    };

    let closePopup = (id, checked, result_yes, result_no) => {
      const popup = document.getElementById(id);
      if (checked === "Yes") {
        result_yes();
      } else if (checked === "No") {
        result_no();
      }
      popup.remove();
    };

    let showPopup = (html) => {
      document.body.insertAdjacentHTML("beforeend", html);
    };

    let popInfo = (content, extrahtml, result_yes, result_no) => {
      const action = `<button onclick="closePopup('{id}', 'Yes', ${result_yes}, ${result_no})" class="btn">Yes</button>`;
      const action_no = `<button onclick="closePopup('{id}', 'No', ${result_yes}, ${result_no})" class="btn">No</button>`;
      const popup = makePopupHTML(content, extrahtml, action, action_no);
      showPopup(popup);
    };

    var currpid = "";

    function openFile() {
      const input = document.createElement("input");
      input.type = "file";
      input.accept = ".*";
      input.onchange = (event) => {
        const file = event.target.files[0];
        if (file) {
          const reader = new FileReader();
          reader.onload = (e) => {
            const content = e.target.result;
            let key = file.name;
            let sidec = document.getElementById("applyfiles");
            let ext = key.split(".");
            let uuid_ = uuidv4();

            var found = false;

            for (let i = 0; i < cached_files.length; i++) {
              if (cached_files[i].name == key) {
                found = true;
                break;
              }
            }

            if (!found) {
              relay.setValue(content);

              window.go.main.App.Savefrontend(key, relay.getValue());

              sidec.insertAdjacentHTML(
                "beforeEnd",
                `<div id="${uuid_}"
style="overflow: hidden;width: calc(100% - 10px); height: 50px; display: flex; gap: 7px; padding-left: 10px; align-items: center; color: #FFFFFFaa; border-bottom: solid 1.5px #2a2831;"
class=${`"${sidec.classList[1]} tab-item"`}
>
<img alt="search" src="file.svg" />
<slot/>
${ext[0]}
<label style="opacity: 50%; font-weight: 400; margin-left: -6px;"
    >
    .${ext[1]}
</label>
</div>`,
              );

              cached_files.push({
                uuid: uuid_,
                name: ext[0],
                extension: ext[1],
                value: content,
                full: key,
              });

              document.getElementById(uuid_).onclick = function () {
                //document
                //  .getElementById(currently_active_container)
                //  .classList.remove("active_container");

                for (let i = 0; i < cached_files.length; i++) {
                  if (cached_files[i].uuid == currpid) {
                    cached_files[i].value = relay.getValue();
                    //window.go.main.App.Savefrontend(cached_files[i].full, relay.getValue());
                  }
                }

                currpid = uuid_;

                for (let i = 0; i < cached_files.length; i++) {
                  if (cached_files[i].uuid == currpid) {
                    relay.setValue(cached_files[i].value);
                    CreateToast("success", "Loaded");
                    break;
                  }
                }
              };
            }
          };
          reader.readAsText(file);
        }
      };
      input.click();
    }

    let listener = (message) => {
      let circle = document.getElementById("updateExecutionStatus");
      if (circle !== undefined) {
        if (circle.classList.contains("red")) {
          return;
        }
        if (circle.classList.contains("yellow")) {
          circle.classList.remove("yellow");
          circle.classList.add("red");
        } else {
          circle.classList.add("red");
        }
      }
    };

    let cachedProcesses = [];

    let isRunning = false;
    let waitQueue = Promise.resolve();

    let robloxOpened = async (message) => {
      await waitQueue;

      let release;
      waitQueue = new Promise((resolve) => (release = resolve));

      try {
        isRunning = true;

        let found = false;
        for (let index = 0; index < cachedProcesses.length; index++) {
          if (cachedProcesses[index].Pid === message.Pid) {
            cachedProcesses[index] = message;
            found = true;
            break;
          }
        }
        if (!found) {
          let UWP = message.ExeName === "Windows10Universal.exe"; // Will be used later for visuals
          cachedProcesses.push(message);
        }
      } finally {
        isRunning = false;
        release(); // Allow the next function call to proceed
      }
    };
    let robloxClosure = (message) => {
      let updatedList = [];
      cachedProcesses.forEach(function (element, _) {
        if (element.Pid !== message.Pid) {
          updatedList.push(element);
        }
      });
      cachedProcesses = updatedList;
      console.log(cachedProcesses, "Removed..", message);
    };

    window.runtime.EventsOn("closure", listener);
    window.runtime.EventsOn("robloxUpdate", robloxOpened);
    window.runtime.EventsOn("robloxClosure", robloxClosure);

    async function openClient() {
      let htmls = ``;

      for (const element of cachedProcesses) {
        const inj = element.Injected ? "is injected" : "is not injected";
        const user = element.Username;
        const Uwp = element.ExeName !== "RobloxPlayerBeta.exe";
        htmls += `
            <div class="home-stats-bun-item" style="width: 100%; height: 77px; margin-bottom: 3px;">
                <div class="home-stats-info" style="margin-top: 20px;">
                    <img
                        style="width: 50px;height: 50px;bottom: 6.5px;left: 10px;margin-right: 10px;"
                        src="${element.Avatar}"
                    />
                    <div class="text-section">
                        <span class="main-text" style="color: #E2E2E2">${user}</span>
                        <div style="position: absolute;">
                            <span class="text-status">Current injection status is: </span>
                            <span class="status-home">${inj}</span>
                        </div>
                    </div>
                    <div class="right-section" style="right:5px;top:25%;">
                      <button class="btn" onclick="switchClient('${user}')">Use</button>
                    </div>
                </div>
            </div>`;
      }

      // Display the generated HTML using popInfo
      popInfo(
        "",
        htmls,
        function () {},
        function () {},
      );
    }

    async function switchClient(ign) {
      for (const element of cachedProcesses) {
        if (element.Username == ign) {
          let circle = document.getElementById("updateExecutionStatus");
          if (circle !== undefined) {
            if (circle.classList.contains("red")) {
              circle.classList.remove("red");
              circle.classList.add("yellow");
            }
          }
          if (element.Injected) {
            if (circle.classList.contains("red")) {
              circle.classList.remove("red");
            }
            if (circle.classList.contains("yellow")) {
              circle.classList.remove("yellow");
            }
          } else {
            circle.classList.remove("yellow");
            circle.classList.add("red");
          }
          window.go.main.App.SetClient(element.Pid);
          CreateToast("success", "Changed to " + element.Username);

          document.getElementById("injectedUserName").innerHTML =
            element.Username;
          document.getElementById("avatarUrl").src = element.Avatar;
          document.getElementById("injectedStatus").innerHTML = element.Injected
            ? "is injected"
            : "is not injected";
          return;
        }
      }
    }

    function getCurrentlyActive() {
      let active = document.getElementById("buttons").children;
      for (var i = 0; i < active.length; i++) {
        try {
          if (active[i].children[0].classList.contains("active")) {
            return active[i].id;
          }
        } catch {}
      }
    }

    function getCurrentlyActiveMenu() {
      let active = document.getElementsByClassName("container");
      for (var i = 0; i < active.length; i++) {
        try {
          if (active[i].style.display !== "none") {
            return active[i].id;
          }
        } catch {}
      }
    }

    let currently_active = getCurrentlyActive();
    let currently_active_menu = getCurrentlyActiveMenu();
    let currently_active_container = "";
    let cached_files = [];
    let anim_popin =
      "1425ms cubic-bezier(0, 0, 0, 1.26) 0s 1 normal none running nx-anim, 400ms linear 0s 1 normal none running nx-wrapper-opacity";
    let anim_fadeout =
      "300ms cubic-bezier(0, 0, 0, 1.26) 0s 1 normal none running nx-fade-out";
    let remove_listeners = new AbortController();
    let autocompleteData = {};

    function execute() {
      window.go.main.App.Execute(relay.getValue()).then((data) => {
        if (data["success"]) {
          CreateToast("success", "Succesful");
        } else {
          if (
            data["error"] != "The operation completed successfully." &&
            data["error"] !=
              "Cannot create a file when that file already exists."
          ) {
            CreateToast("error", data["error"]);
          }
        }
      });
    }

    function attach() {
      let circle = document.getElementById("updateExecutionStatus");
      if (circle !== undefined) {
        if (circle.classList.contains("red")) {
          circle.classList.remove("red");
          circle.classList.add("yellow");
        }
      }

      window.go.main.App.Attach().then((data) => {
        if (data["success"]) {
          if (circle.classList.contains("red")) {
            circle.classList.remove("red");
          }
          if (circle.classList.contains("yellow")) {
            circle.classList.remove("yellow");
          }
          CreateToast("success", "Succesfully attached into roblox!");
        } else {
          if (data["error"] != "The operation completed successfully.") {
            if (circle.classList.contains("yellow")) {
              circle.classList.remove("yellow");
              circle.classList.add("red");
            }
            CreateToast("error", data["error"]);
            return;
          }
        }
      });
    }

    function apply(self) {
      console.log(self);
      if (self.children[0].classList.contains("active")) {
        return;
      }
      try {
        document
          .getElementById(currently_active)
          .children[0].classList.remove("active");
      } catch (e) {
        console.log(e);
      }
      try {
        self.children[0].classList.add("active");
      } catch (e) {
        console.log(e);
      }
      currently_active = self.id;
    }

    document.getElementById("buttons").children[0].onclick = function () {
      apply(document.getElementById("buttons").children[0]);
      if (
        document
          .getElementById("buttons")
          .children[0].children[0].classList.contains("active")
      ) {
        document.getElementById(getCurrentlyActiveMenu()).style.display =
          "none";
        document.getElementById("home-page").style.display = "flex";
        currently_active_menu = "home-page";
      }
    };
    document.getElementById("buttons").children[1].onclick = function () {
      apply(document.getElementById("buttons").children[1]);
      if (
        document
          .getElementById("buttons")
          .children[1].children[0].classList.contains("active")
      ) {
        document.getElementById(getCurrentlyActiveMenu()).style.display =
          "none";
        document.getElementById("home-scripts").style.display = "flex";
        currently_active_menu = "home-scripts";
      }
    };
    document.getElementById("buttons").children[2].onclick = function () {
      apply(document.getElementById("buttons").children[2]);
      if (
        document
          .getElementById("buttons")
          .children[2].children[0].classList.contains("active")
      ) {
        document.getElementById(getCurrentlyActiveMenu()).style.display =
          "none";
        document.getElementById("home-scripthub").style.display = "flex";
        currently_active_menu = "home-scripthub";
      }
    };
    document.getElementById("buttons").children[3].onclick = function () {
      apply(document.getElementById("buttons").children[3]);
      if (
        document
          .getElementById("buttons")
          .children[3].children[0].classList.contains("active")
      ) {
        document.getElementById(getCurrentlyActiveMenu()).style.display =
          "none";
        document.getElementById("home-settings").style.display = "flex";
        currently_active_menu = "home-settings";
      }
    };

    function uuidv4() {
      return "10000000-1000-4000-8000-100000000000".replace(/[018]/g, (c) =>
        (
          +c ^
          (crypto.getRandomValues(new Uint8Array(1))[0] & (15 >> (+c / 4)))
        ).toString(16),
      );
    }

    function CreateToast(type, content) {
      // success, warning, error
      function uuidv4() {
        return "10000000-1000-4000-8000-100000000000".replace(/[018]/g, (c) =>
          (
            +c ^
            (crypto.getRandomValues(new Uint8Array(1))[0] & (15 >> (+c / 4)))
          ).toString(16),
        );
      }
      let uuid = uuidv4().replace("-", "");
      let html = `<div class="toast-nx">${content}<button id="new_toast_${uuid}" class="close-btn">×</button></div>`;
      document
        .getElementsByClassName("toast-nx-container")[0]
        .insertAdjacentHTML("afterBegin", html);
      let toast = document.getElementById(`new_toast_${uuid}`);
      toast.parentElement.classList.add(
        document.getElementsByClassName("toast-nx-container")[0].classList[1],
      );
      toast.classList.add(
        document.getElementsByClassName("toast-nx-container")[0].classList[1],
      );
      toast.parentElement.style.animation = anim_popin.replace(
        "1425ms",
        "1000ms",
      );
      let clickedAlr;
      setTimeout(() => {
        if (!clickedAlr) {
          toast.parentElement.classList.add("hide");
          setTimeout(() => {
            document.getElementById(`new_toast_${uuid}`).parentElement.remove();
          }, 300);
        }
      }, 2000);
      let ab = new AbortController();
      toast.addEventListener(
        "click",
        (event) => {
          clickedAlr = true;
          toast.parentElement.classList.add("hide");
          setTimeout(() => {
            document.getElementById(`new_toast_${uuid}`).parentElement.remove();
            ab.abort();
          }, 300);
        },
        {
          signal: remove_listeners.signal,
        },
      );
    }


    window.monaco.editor.defineTheme("customDarkTheme", {
        base: "vs-dark",
        inherit: true,
        rules: [
          { token: "keyword", foreground: "B8AAFF" }, // Keywords like 'if', 'else', 'function'
          { token: "constants", foreground: "7C19C9" }, // Constants (Booleans, Strings, Ints)
          { token: "global", foreground: "7C19C9" }, // Global variables
          { token: "operator", foreground: "B8AAFF" }, // Operators (e.g., '+', '==')
          { token: "variable.language.self", foreground: "7C19C9" }, // Self-references or 'this'
          { token: "comment", foreground: "608B4E", fontStyle: "italic" }, // Comments
          { token: "string", foreground: "7C19C9" }, // Strings
          { token: "number", foreground: "7C19C9" }, // Numbers (Ints, Floats, etc.)
          { token: "number.float", foreground: "7C19C9" }, // Floating-point numbers
          { token: "number.hex", foreground: "7C19C9" }, // Hexadecimal numbers
          { token: "number.binary", foreground: "7C19C9" }, // Binary numbers
          { token: "delimiter", foreground: "B8AAFF" }, // General delimiters (e.g., commas, semicolons)
          { token: "delimiter.bracket", foreground: "B8AAFF" }, // Brackets (e.g., '{', '}')
          { token: "delimiter.array", foreground: "B8AAFF" }, // Array delimiters (e.g., '[', ']')
          { token: "delimiter.parenthesis", foreground: "B8AAFF" }, // Parentheses (e.g., '(', ')')
          { token: "operator.special", foreground: "B8AAFF" }, // Special operators (e.g., assignment operator)
          { token: "variable.parameter.function", foreground: "797693" }, // Function parameters (Typography)
          { token: "entity.name.function", foreground: "B8AAFF" }, // Function names
          { token: "support.function", foreground: "B8AAFF" }, // Support functions
          { token: "support.type", foreground: "7C19C9" }, // Support types
          { token: "function.call", foreground: "B8AAFF" }, // Function calls
          { token: "class.name", foreground: "7C19C9" }, // Class names
          { token: "type.annotation", foreground: "B8AAFF" }, // Type annotations (e.g., in TypeScript)
          { token: "template.literal", foreground: "7C19C9" }, // Template literals (Strings)
          { token: "modifier", foreground: "B0C4DE" }, // Modifiers (e.g., 'public', 'private')
          { token: "annotation", foreground: "FFA500" }, // Annotations (e.g., '@Override')
          { token: "variable.other.constant", foreground: "7C19C9" }, // Other constants
          { token: "variable.property", foreground: "797693" }, // Object properties (Typography)
          { token: "variable.other.enummember.lua", foreground: "7C19C9" }, // Enum members in Lua
          { token: "support.function.library", foreground: "B8AAFF" }, // Library functions
          { token: "keyword.control.export", foreground: "B8AAFF" }, // Export keywords
          { token: "storage.type", foreground: "B8AAFF" }, // Storage types (e.g., 'class', 'struct')
          { token: "entity.name.type.alias", foreground: "7C19C9" }, // Type aliases
          {
            token: "punctuation.definition.typeparameters",
            foreground: "B8AAFF",
          }, // Type parameter definitions
          { token: "punctuation.definition.block", foreground: "B8AAFF" }, // Block definitions (e.g., '{', '}')
          { token: "punctuation.separator.table", foreground: "B8AAFF" }, // Table separators (e.g., in Lua)
          {
            token: "punctuation.definition.parameters",
            foreground: "B8AAFF",
          }, // Function parameter definitions
          { token: "variable.parameter.variadic", foreground: "B0C4DE" }, // Variadic function parameters (e.g., '...')
          {
            token: "keyword.operator.type.annotation",
            foreground: "B8AAFF",
          }, // Type annotation operators
          { token: "operator.type", foreground: "B8AAFF" }, // Type operators (e.g., 'is', 'as')
          { token: "punctuation", foreground: "B8AAFF" }, // General punctuation (e.g., '.', ',')
          { token: "error", foreground: "FF0000" }, // Errors
          { token: "", foreground: "B0C4DE" }, // Default text color
        ],
        colors: {
          "scrollbar.shadow": "#00000000",
          "widget.shadow": "#00000000",
          "editor.background": "#29292900",
          "minimap.background": "#29292900",
          "editor.foreground": "#ABB2BF",
          "editorLineNumber.foreground": "#999999",
          "editor.selectionBackground": "#264F78",
          "editorCursor.foreground": "#528BFF",
          "editor.lineHighlightBackground": "#f1f1f101",
          "editor.selectionHighlightBackground": "#264F78",
        },
      });

      window.monaco.languages.register({
        id: "lua",
        extensions: [".lua"],
        aliases: ["Lua", "lua"],
      });

      let tabs = [
        { title: "Tab 1", code: "-- auto generated (#luna)", selected: false, uuid: uuidv4()},
      ];

      window.monaco.editor.setTheme("customDarkTheme");
      const editor = window.monaco.editor.create(document.getElementById("editor"), {
        language: "lua",
        theme: "customDarkTheme",
        fontSize: 14,
        fontFamily: "'JetBrains Mono', Consolas, 'Courier New', monospace",
        folding: true,
        autoIndent: true,
        dragAndDrop: true,
        minimap: {
          enabled: false,
        },
        showFoldingControls: "always",
        smoothScrolling: true,
        stopRenderingLineAfter: 6500,
        cursorBlinking: "smooth",
        cursorSmoothCaretAnimation: "on",
        foldingHighlight: false,
        fontLigatures: true,
        formatOnPaste: true,
        showDeprecated: true,
        suggest: {
          snippetsPreventQuickSuggestions: true,
        },
        tabFocusMode: false,
        overviewRulerLanes: 0,
        automaticLayout: true,
      });
      relay = editor;
      editor.focus();

      function resizeEditor() {
        relay.layout();
      }
      window.addEventListener("resize", resizeEditor);

      window.go.main.App.GetScriptsItems().then((resp) => {
        let first = false;
        let sidec = document.getElementById("applyfiles");

        for (const [key, value] of Object.entries(resp)) {
          let ext = key.split(".");
          let uuid_ = uuidv4();
          sidec.insertAdjacentHTML(
            "beforeEnd",
            `<div id="${uuid_}"
style="overflow: hidden;width: calc(100% - 10px); height: 50px; display: flex; gap: 7px; padding-left: 10px; align-items: center; color: #FFFFFFaa; border-bottom: solid 1.5px #2a2831;"
class=${`"${sidec.classList[1]}"`}
>
<img alt="search" src="file.svg" />
<slot/>
${ext[0]}
<label style="opacity: 50%; font-weight: 400; margin-left: -6px;"
    >
    .${ext[1]}
</label>
</div>`,
          );

          cached_files.push({
            uuid: uuid_,
            name: ext[0],
            extension: ext[1],
            value: value,
            full: key,
          });

          document.getElementById(uuid_).onclick = function () {
            tabs = tabs.map((tab, i) => {
              i = i + 1;
              tab.selected = false;
              return tab;
            });

            currpid = uuid_;
            for (let i = 0; i < cached_files.length; i++) {
              if (cached_files[i].uuid == currpid) {
                relay.setValue(cached_files[i].value);
                CreateToast("success", "Loaded");
                break;
              }
            }
          };
          if (first == false) {
            currpid = uuid_;
            relay.setValue(value);
            first = true;
          }
        }
      });

    document.getElementById("savefile").onclick = function () {
        if (currpid != "") {
          for (let i = 0; i < cached_files.length; i++) {
            if (cached_files[i].uuid == currpid) {
              // logic to save the file
              cached_files[i].value = relay.getValue();
              window.go.main.App.Savefrontend(
                cached_files[i].full,
                relay.getValue(),
              );
            }
          }
        }
      };

    function selectTab(index) {
      tabs = tabs.map((tab, i) => {
         i = i + 1;
         tab.selected = i === index;
         return tab;
      });
      relay.setValue(tabs[index-1].code);
      document.getElementById(tabs[index-1].uuid).classList.add("tab-selected")
    }

    function removetab(uuid) {
      let i = tabs.length;
      let o = 0;
      let new_tabs = [];
      
      while (o < i) {
        if (tabs[o].uuid !== uuid) {
          new_tabs.push(tabs[o]);
        }
        o = o + 1;
      }
      
      tabs = new_tabs;
      
      const element = document.getElementById(uuid);
      if (element) {
        element.remove();
      }
    }

   function addTab() {
      let uuid_ = uuidv4();
      tabs.push({
         title: `Tab ${tabs.length + 1}`,
         code: `-- auto generated (#luna)`,
         selected: false,
         uuid:uuid_,
      });

      document.getElementById("tabs").innerHTML += `<div class="tab-item" id="${uuid_}" onclick="selectTab(${tabs.length})">
        <p class="main-text" style="color:#E2E2E2;font-size:14px;margin: auto;">Tab ${tabs.length}</p>
        <img onclick="removetab('${uuid_}')" style="width: 15px;margin: auto;margin-right:10px;color:red;" alt="right" src="close.svg" />
      </div>`
      
   }

   document.getElementById("tabs").innerHTML += `<div class="tab-item" id="${tabs[0].uuid}" onclick="selectTab(${tabs.length})">>
        <p class="main-text" style="color:#E2E2E2;font-size:14px;margin: auto;">Tab 1</p>
        <img onclick="removetab('${tabs[0].uuid}')" style="width: 15px;margin: auto;margin-right:10px;color:red;" alt="right" src="close.svg" />
      </div>`

    relay.setValue(tabs[0].code);
    document.getElementById(tabs[0].uuid).classList.add("tab-selected")

    relay.onDidChangeModelContent(() => {
        const newCode = relay.getValue();
        const selectedTab = tabs?.find((t) => t.selected);
        if (selectedTab) {
          selectedTab.code = newCode
        };
    });
  </script>
  <style>

    .tab-selected {
      background:black;
    }
    .modal {
      color: #666;
      font-family: Outfit;
      font-size: 16px;
      font-style: normal;
      font-weight: 700;
      line-height: normal;
      display: flex;
      position: fixed;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      background-color: rgba(0, 0, 0, 0.35); /* Semi-transparent background */
      justify-content: center;
      align-items: center;
    }

    .modal-box {
      contain: content;
      background: radial-gradient(
        100% 100% at 51.43% 100%,
        rgb(42 41 49 / 43%) 0%,
        rgb(19 18 24) 100%
      );
      padding: 20px;
      border-radius: 3px;
      max-width: 400px;
      box-shadow: 0 0 4px rgba(0, 0, 0, 0.7); /* Slight shadow */
      border: 1px solid rgba(108, 103, 128, 0.28);
      display: flex;
      flex-direction: column;
      max-height: 80vh; /* Ensure modal does not exceed viewport height */
      overflow: hidden;
      width: 100%; /* Full width */
      overflow: auto; /* Enable scroll if needed */
    }

    .modal-content {
      flex-grow: 1; /* Ensures content area grows as needed */
      flex-direction: column;
      margin-top: 5px; /* Space above the action buttons */
      display: flex;
    }

    .modal-action {
      display: flex;
      justify-content: flex-end;
      gap: 10px; /* Space between buttons */
    }

    .btn {
      font-size: 11px;
      margin-top: 5px;
      height: 32px;
      width: 50px;
      flex-shrink: 0;
      border-radius: 4px;
      border: 1px solid rgba(108, 103, 128, 0.25);
      background: rgba(42, 41, 49, 0.85);
      backdrop-filter: blur(15px);
      outline: 0px;
      color: rgba(255, 255, 255, 0.6);
      caret-color: rgba(255, 255, 255, 0.2);
    }

    .btn:hover {
      background: radial-gradient(
        100% 100% at 51.43% 100%,
        #2a2931 20%,
        rgba(42, 41, 49, 0) 100%
      ); /* Darker green on hover */
    }

    /* Container for the toggle switch */
    .toggle-container {
      width: 62px;
      height: 21px;
      border-radius: 5px;
      display: flex;
      align-items: center;
      justify-content: flex-start;
      padding: 3px;
      position: relative;
      box-shadow: 0 0 4px rgba(0, 0, 0, 0.4); /* Slight shadow */
      transition: background-color 0.3s ease;
    }

    /* The button itself */
    .toggle-button {
      color: #666;
      font-family: Outfit;
      font-size: 16px;
      font-style: normal;
      font-weight: 700;
      line-height: normal;

      width: 35px;
      height: 22px;
      border-radius: 3px;
      background: linear-gradient(180deg, #d5d5d5 0%, #f0f0f0 100%);
      display: flex;
      align-items: center; /* Vertically center the text */
      justify-content: center; /* Horizontally center the text */
      box-shadow: 0 2px 2px rgba(0, 0, 0, 0.2);
      transition: transform 0.3s ease;
      user-select: none;
    }

    /* Active state for the button */
    .toggle-container.active .toggle-button {
      transform: translateX(
        27px
      ); /* Move the button to the right when toggled */
      background: linear-gradient(180deg, #393939 0%, #2d2d2d 100%);
      color: #979797;
    }

    /* For some extra aesthetic, you can add dots */
    .toggle-container::before {
      content: "";
      position: absolute;
      width: 100%;
      height: 100%;
      background-size: 3px 3px;
      opacity: 0.5;
    }

    /* Full-screen loading background */
    #loading-screen {
      position: fixed;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      backdrop-filter: blur(10px); /* Apply blur effect */
      background: radial-gradient(
          100% 100% at 51.43% 100%,
          #2a2931 0%,
          rgba(42, 41, 49, 0) 100%
        ),
        rgba(27, 26, 32, 0.85); /* Semi-transparent background */
      display: flex;
      justify-content: center;
      align-items: center;
      flex-direction: column;
      z-index: 9999;
      opacity: 1; /* Initial opacity */
      transition: opacity 1s ease; /* Transition for opacity */
    }

    .fade-out {
      opacity: 0; /* Change to 0 when fading out */
    }

    /* Loading spinner with the specified color */
    .loader {
      border-radius: 16px;
      border: 3px solid rgba(108, 103, 128, 0.25); /* Spinner color */
      border-top: 3px solid rgba(108, 103, 128, 0.1); /* Spinner color (solid) */
      border-radius: 50%;
      width: 35px;
      height: 35px;
      animation: spin 1s linear infinite;
    }

    /* Text under the spinner */
    #loading-screen p {
      font-size: 20px;
      margin-top: 20px;
      color: rgba(108, 103, 128, 1); /* Matching text color */
    }

    /* Spinner animation */
    @keyframes spin {
      0% {
        transform: rotate(0deg);
      }
      100% {
        transform: rotate(360deg);
      }
    }

    /* Content style after loading */
    #content {
      padding: 20px;
    }

    .btn-pushable {
      border: none;
      outline: none;
      font-family: Outfit;
      border-radius: 3px;
    }

    .btn-pushable:hover {
      transition: 0.1s ease-in;
    }

    .btn-pushable:active {
      transform: translateY(2px);
      border-bottom-width: 2px;
      box-shadow: none;
    }

    input::-webkit-input-placeholder {
      color: rgba(255, 255, 255, 0.6);
      font-family: Outfit;
      font-size: 11px;
      font-style: normal;
      font-weight: 500;
      line-height: normal;
    }

    input::-moz-placeholder {
      color: rgba(255, 255, 255, 0.6);
      font-family: Outfit;
      font-size: 11px;
      font-style: normal;
      font-weight: 500;
      line-height: normal;
    }

    input:-ms-input-placeholder {
      color: rgba(255, 255, 255, 0.6);
      font-family: Outfit;
      font-size: 11px;
      font-style: normal;
      font-weight: 500;
      line-height: normal;
    }

    input:-moz-placeholder {
      color: rgba(255, 255, 255, 0.6);
      font-family: Outfit;
      font-size: 11px;
      font-style: normal;
      font-weight: 500;
      line-height: normal;
    }

    .input-button-holder {
      font-size: 11px;
      margin-top: 5px;
      height: 28px;
      flex-shrink: 0;
      border-radius: 4px;
      border: 1px solid rgba(108, 103, 128, 0.25);
      background: rgba(42, 41, 49, 0.85);
      backdrop-filter: blur(15px);
      outline: 0px;
      color: rgba(255, 255, 255, 0.6);
      caret-color: rgba(255, 255, 255, 0.2);
    }
    .text-input-holder {
      display: flex;
      width: 100%;
      height: 100%;
      flex-direction: column;
      justify-content: center;
      flex-shrink: 0;
      color: rgba(255, 255, 255, 0.8);
      text-align: center;
      font-family: Outfit;
      font-size: 11px;
      font-style: normal;
      font-weight: 700;
      line-height: normal;
    }
    .text-yellow {
      background: linear-gradient(0deg, #f3c72a 0%, #8d7318 100%);
      background-clip: text;
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
    }
    .text-green {
      background: linear-gradient(180deg, #71ce38 0%, #30a342 100%);
      background-clip: text;
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
    }

    .text-status-yellow {
      color: #f3eb2a;
      font-family: Outfit;
      font-size: 13px;
      font-style: normal;
      font-weight: 500;
      line-height: normal;
    }
    .text-status-green {
      background: linear-gradient(180deg, #70ce38 52.38%, #30a342 109.52%);
      background-clip: text;
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      font-family: Outfit;
      font-size: 13px;
      font-style: normal;
      font-weight: 500;
      line-height: normal;
    }

    .home-label-button {
      color: rgba(255, 255, 255, 0.5);
      font-family: Outfit;
      font-size: 11.5px;
      font-style: normal;
      font-weight: 500;
      line-height: normal;
    }

    .home-stats-info {
      display: flex;
      flex-direction: row;
      /* position: fixed; */
      gap: 5px;
      top: 7px;
      justify-content: center;
    }

    .text-status {
      color: rgba(255, 255, 255, 0.5);
      font-size: 11.5px;
      font-family: "Outfit", sans-serif;
      font-weight: 400;
      word-wrap: break-word;
    }

    .status-home {
      color: rgba(181, 162, 255, 0.5);
      font-size: 11.5px;
      font-family: "Outfit", sans-serif;
      font-weight: 400;
      word-wrap: break-word;
    }

    .home-stats-info img {
      position: relative;
    }

    .home-stats-bun-item {
      border-radius: 3px;
      border: 1px solid #38373c;
      background: radial-gradient(
          100% 100% at 51.43% 100%,
          #2a2931 0%,
          rgba(42, 41, 49, 0) 100%
        ),
        rgba(27, 26, 32, 0.85);
      contain: content;
    }
    .home-stats-bun {
      flex-direction: row;
      display: flex;
      gap: 10px;
    }

    .home-stats-container {
      gap: 10px;
      display: flex;
      flex-direction: column;
      position: relative;
      width: 100%;
      contain: content;
      padding: 15px;
    }

    .home-container {
      background: linear-gradient(
        0deg,
        rgba(26.86, 25.64, 31.72, 0.85) 0%,
        rgba(26.86, 25.64, 31.72, 0.85) 100%
      );
      backdrop-filter: blur(4px);
      height: 130.5px;
      position: relative;
      border-bottom: solid #6c678047 1px;
    }

    .settings-container {
      height: 42px;
      border-bottom: solid #6c678047 1px;
      justify-content: center;
      position: relative;
      contain: content;
    }

    .top-section {
      display: flex;
      flex-direction: row;
      position: fixed;
      gap: 15px;
      content: contain;
      top: 10px;
      justify-content: center;
    }

    .top-section img {
      position: relative;
      left: 10px;
    }

    .text-section {
      width: 100%;
      height: 100%;
    }

    .text-section .main-text {
      color: rgba(255, 255, 255, 0.5);
      font-size: 12px;
      font-family: Outfit, sans-serif;
      font-weight: 600;
      word-wrap: break-word;
    }

    .text-section .sub-text {
      color: rgba(255, 255, 255, 0.25);
      font-size: 12px;
      font-family: Outfit, sans-serif;
      font-weight: 600;
      word-wrap: break-word;
    }

    .bottom-section {
      display: flex;
      position: fixed;
      gap: 10px;
      content: contain;
      bottom: 15px;
      left: 10px;
      justify-content: center;
    }

    .bottom-section .date-text {
      width: 100%;
      height: 100%;
      color: rgba(255, 255, 255, 0.25);
      font-size: 12px;
      font-family: Outfit, sans-serif;
      font-weight: 600;
      word-wrap: break-word;
    }

    .active-container {
      background: linear-gradient(
        0deg,
        rgba(85, 85, 85, 0.09287464985994398) 0%,
        rgba(94, 94, 94, 0.06766456582633051) 47%,
        rgba(99, 99, 99, 0.09007352941176472) 100%
      );
    }
    .tab-item:hover {
      background: linear-gradient(
        0deg,
        rgba(85, 85, 85, 0.09287464985994398) 0%,
        rgba(94, 94, 94, 0.06766456582633051) 47%,
        rgba(99, 99, 99, 0.09007352941176472) 100%
      );
    }

    #editor-container {
      position: relative;
      width: 100%;
      contain: content;
    }

    #editor {
      position: absolute;
      top: 50px;
      bottom: 0;
      left: 0;
      right: 0;
      height: calc(100% - 115px);
      border-top: 1.5px solid rgb(42, 40, 49);
      border-bottom: 1.5px solid rgb(42, 40, 49);
    }

    .sidecontainer {
      border-right: solid #6c678047 1px;
      height: 100%;
      width: 178px;
    }

    .container {
      position: fixed;
      display: flex;
      width: calc(100% - 66px);
      height: calc(100% - 50px);
      overflow: auto;
      background-color: #151419;
      border-left: solid 1.5px #6c678040;
    }

    .topbar {
      width: 99%;
      height: 54px;
      z-index: 2;
      position: fixed;
      display: flex;
      justify-content: space-between;
      align-items: center;
      left: 0;
      top: 0;
    }

    .sidebar {
      width: 65px;
      position: fixed;
      align-items: center;
      display: flex;
      flex-direction: column;
      left: 0;
      bottom: 0;
    }

    #app {
      height: 100vh;
      text-align: center;
    }

    .toast-nx-container {
      position: fixed;
      top: 60px;
      right: 20px;
      display: flex;
      flex-direction: column;
      gap: 10px;
      z-index: 999999999999;
    }

    .toast-nx {
      -webkit-touch-callout: none;
      -webkit-user-select: none;
      -khtml-user-select: none;
      -moz-user-select: none;
      -ms-user-select: none;
      padding: 15px 20px;
      border-radius: 16px;
      color: #f0f0f0;
      font-size: 14px;
      background: linear-gradient(158deg, #0e0e0e0c, #00000009 100%);
      backdrop-filter: blur(10px);
      border-top: 1px solid rgba(68, 68, 68, 0.8);
      box-shadow:
        0 -2px 5px rgb(0 0 0 / 30%),
        0 0 15px rgb(0 0 0 / 10%);
      display: flex;
      align-items: center;
      justify-content: space-between;
      transition:
        opacity 0.3s,
        transform 0.3s;
      z-index: 999999999999;
    }

    .toast-nx .close-btn {
      background: 0 0;
      border: none;
      color: #f0f0f0;
      font-size: 16px;
      cursor: pointer;
      margin-left: 10px;
    }

    .toast-nx.hide {
      opacity: 0;
      transform: translateY(-20px);
    }

    @keyframes onlineAnimation {
      0% {
        transform: scale(1);
      }

      25% {
        transform: scale(1.1);
      }

      50% {
        transform: scale(1);
      }
    }

    .circle {
      background-color: #11ed9e;
      box-shadow: 0 0 10px #11ed9e;
      width: 12px;
      height: 12px;
      border-radius: 50%;
      animation: infinite onlineAnimation 3s;
    }

    .circle.red {
      background-color: #ff5959;
      box-shadow: 0 0 10px #ff5959;
      width: 12px;
      height: 12px;
      border-radius: 50%;
      animation: infinite onlineAnimation 3s;
    }

    .circle.yellow {
      background-color: #d9b823;
      box-shadow: 0 0 10px #d9b823;
      width: 12px;
      height: 12px;
      border-radius: 50%;
      animation: infinite onlineAnimation 3s;
    }

    .chat .messageForm {
      margin: 6px;
      display: flex;
      flex-direction: row;
      justify-content: center;
    }

    .chat .messageForm input {
      color: #fff;
      padding: 10px 12px;
      width: 100%;
      background-color: transparent;
      border: none;
      border-radius: 16px;
      outline: 0;
      box-shadow: 0 0 10px rgba(0, 0, 0, 0.3);
      background: -webkit-linear-gradient(
        top,
        rgba(32, 32, 32, 0.767),
        rgb(29, 29, 29)
      );
    }

    .chatBox {
      width: auto;
      display: flex;
      flex-direction: row;
      padding: 7px 5px 0;
      margin-bottom: 5px;
    }

    .middle-section {
      display: flex;
      flex-direction: row;
      position: fixed;
      gap: 5px;
      content: contain;
      top: 20%;
      justify-content: center;
      margin: 0 10px;
    }

    .right-section {
      top: 37%;
      right: 0px;
      position: absolute;
    }

    .chatBox .username {
      color: #fff;
      font-weight: 500;
    }

    .chatBox .pfp {
      object-fit: cover;
      height: 55px;
      width: 55px;
      border-radius: 50%;
    }

    .chatBox .content {
      display: flex;
      flex-direction: column;
      margin-left: 5px;
    }

    .chatBox .content .topHalf {
      display: flex;
      flex-direction: row;
      align-items: center;
    }

    .chat .messagesContainer {
      display: flex;
      flex-direction: column-reverse;
      overflow-y: scroll;
      overflow-x: hidden;
      justify-content: end;
      flex: auto;
      background: -webkit-linear-gradient(
        top,
        rgba(32, 32, 32, 0.767),
        rgb(29, 29, 29)
      );
      border-radius: 16px;
      box-shadow: 0 0 10px rgba(0, 0, 0, 0.3);
    }

    .chatBox .content .topHalf .username {
      margin-right: 7px;
    }

    .chatBox .content .topHalf .role {
      height: 25px;
      width: 25px;
      margin-right: 5px;
    }

    .chatBox .content .topHalf .timeStamp {
      opacity: 34%;
      color: #fff;
      font-weight: 300;
    }

    .chatBox .content .message {
      color: #b2b2b2;
      opacity: 92%;
      word-break: break-all;
    }

    .chat {
      display: flex;
      flex-direction: column;
      width: 100%;
      height: 100%;
    }

    @keyframes slideUp {
      from {
        transform: translateY(100vh);
      }

      to {
        transform: translateY(0);
      }
    }

    @keyframes nx-opacity {
      0% {
        opacity: 0;
      }

      100% {
        opacity: 1;
      }
    }

    @keyframes nx-anim {
      0% {
        transform: scaleX(2) scaleY(1.25);
      }

      30% {
        transform: scaleX(0.9) scaleY(0.95);
      }

      100% {
        transform: scaleX(1) scaleY(1);
      }
    }

    @keyframes nx-fade-out {
      0% {
        opacity: 1;
      }

      100% {
        opacity: 0;
      }
    }

    @keyframes nx-fade-in {
      0% {
        opacity: 0;
      }

      100% {
        opacity: 1;
      }
    }

    /* width */
    ::-webkit-scrollbar {
      width: 0px;
    }

    /* Track */
    ::-webkit-scrollbar-track {
      background: #2020205e;
      border-radius: 16px;
      height: 80%;
    }

    /* Handle */
    ::-webkit-scrollbar-thumb {
      background: #000000;
    }

    /* Handle on hover */
    ::-webkit-scrollbar-thumb:hover {
      background: #000000;
    }

    :root {
      background-color: #151419ce;
      font-family: Outfit;
    }

    #topbarandcontainer {
      display: flex;
      flex-direction: column;
      height: 100vh; /* Ensuring it takes up full height */
    }

    .topbar {
      width: 100%;
      height: 54px; /* Adjust as necessary */
      z-index: 2; /* Ensure it's above other elements */
      position: fixed; /* Keep it at the top */
      top: 0;
      left: 0;
    }

    #main-container {
      display: flex; /* Allowing flexbox for sidebar and containers */
      flex-grow: 1;
      margin-top: 54px; /* To offset for the top bar */
      height: calc(100% - 54px); /* Adjusting the height */
    }

    .sidebar {
      width: 65px; /* Fixed width */
      flex-shrink: 0; /* Prevent shrinking */
      position: fixed; /* Keep it visible */
      top: 54px; /* Below the top bar */
      left: 0;
      display: flex;
      flex-direction: column;
    }

    .container {
      flex: 1; /* This will allow it to take the remaining space */
      margin-left: 65px; /* Offset for sidebar */
      overflow: auto; /* Enable scrolling if content overflows */
      background-color: #151419; /* Your background */
    }
    .monaco-editor {
      width: 100%;
      height: 100%;
      outline:none;
    }

    .tab-item {
      overflow-x:hidden;
      overflow-y:hidden;
      white-space:nowrap;
      display:flex;
      contain:content;
      width:95px;
      height:35px;
      border-radius:10px;
      margin-top:7px;
      margin-left:5px;
      border:1.5px solid rgb(42, 40, 49);
      background: #00000008;
      backdrop-filter: blur(10px);
    }

    .tab-add {
      width: 25px;
      height: 25px;
      background: radial-gradient(100% 100% at 51.43% 100%, #2A2931 20%, rgba(42, 41, 49, 0.00) 100%);
      border-radius: 10px;
      margin: 12px auto auto;
      border: solid 1px rgba(108, 103, 128, 0.25);
      transition: transform 0.2s ease-in-out;
    }

    .tab-add:hover {
      transition: transform 0.2s ease-in-out;
    }
    
    .tab-add:active {
      transform: translateY(1px);
      border-bottom-width: 2px;
      box-shadow: none;
    }
  </style>
</body>
