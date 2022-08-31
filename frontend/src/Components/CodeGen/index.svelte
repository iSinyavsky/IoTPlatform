<style>
    main{
            margin: 0 !important;
    }
</style>
<script>
    import "./style.css"
    import * as lrFetch from '../../common/lrFetch'
    import {afterUpdate} from "svelte";
    import {profile} from "../../store";
    import {openModal} from "../../common";
    import CodeComponent from "../CodeComponent"
    const _ = require("lodash");
    const controllers = require("./controllers.json")
    let controller = controllers["nodeMCU"];
    let modules = require("./modules.json");
    let modulesChosed = [];
    let pins = controller.pins;
    let modulePins = 0;

    let svgNode;
    let nullActiveLink = {
        from: null,
        module: null,
        pKey: null,
        mKey: null
    }
    let activeLink = {
        from: null,
        module: null,
        pKey: null,
        mKey: null
    }


    function addLineToSVG (pin, modulePin){
        let svg = document.querySelector("svg");

        svg.setAttribute("width", "100%")
        svg.setAttribute("height", "100%")
        let x1,y1,x2,y2;
        let pinElement = document.querySelector(".pin[data-id='"+pin+"']")
        let modulePinElement = document.querySelector(".module-pin[data-id='"+modulePin+"']")

        let dx = 250;
        x1 = pinElement.getBoundingClientRect().x;
        y1 = pinElement.getBoundingClientRect().y;

        x2 = modulePinElement.getBoundingClientRect().x;
        y2 = modulePinElement.getBoundingClientRect().y;

        let line = document.createElementNS("http://www.w3.org/2000/svg","line");
        line.setAttribute("x1", x1+20)
        line.setAttribute("y1", y1+10)
        line.setAttribute("x2", x2+1)
        line.setAttribute("y2", y2+2)
        line.setAttribute("stroke-width", 3)
        line.setAttribute("stroke", "red")
        line.setAttribute("data-type", "fix");
        svg.appendChild(line);
    }
    // <svg width='1920' height='1080'><line x1='50' y1='50' x2='350' y2='350' stroke='black'/></svg>
    afterUpdate(()=>{
        if(document.querySelectorAll("line").length > 0) {
            document.querySelectorAll("line").forEach(el => {
                el.remove();
            })
        }

        modulesChosed.forEach((module, i) => {
            module.pins.forEach((pin,j) => {
                if (pin !== "") {
                    let modulePin = i+""+j;
                    console.log(modulePin);
                    console.log("aa");
                    addLineToSVG(pin, modulePin);
                }
            })
        })

    })

        document.onmousemove = (event) => {
        if (activeLink.from !== null) {

                if (!document.querySelector("#svg line[data-type='temp']")) {
                    let line = document.createElementNS("http://www.w3.org/2000/svg", "line");
                    line.setAttribute("data-type", "temp")
                    line.setAttribute("stroke-width", 3)
                    line.setAttribute("stroke", "red");

                    let modulePinElement = document.querySelector(".module-pin[data-id='" + activeLink.from + "']")
                    let x1 = modulePinElement.getBoundingClientRect().x;
                    let y1 = modulePinElement.getBoundingClientRect().y;
                    line.setAttribute("x1", x1)
                    line.setAttribute("y1", y1)
                    svgNode.appendChild(line);
                }

                let line = document.querySelector("#svg line[data-type='temp']");

                line.setAttribute("x2", event.screenX - 7)
                line.setAttribute("y2", event.screenY - 77)
            svgNode.appendChild(line);
            };
        }

    const activateLink = (module, linkPin, pin, mKey, pinNum) => {

        if (pin !== ""){
            console.log(pin);
            return
        }
        console.log("activate", pin, linkPin);
        $:{
        activeLink.from = linkPin;
        activeLink.module = module;
        activeLink.pKey = pinNum;
        activeLink.mKey = mKey;
        }


        console.log(activeLink)
    }

    const selectPin = (pin) => {
        console.log("select", pin, activeLink.from, activeLink.mKey, activeLink.pKey);
        if (activeLink.mKey !== null) {
            console.log("blaa", pin)
            console.log(activeLink.mKey, "", activeLink.pKey)
            console.log("aa", modulesChosed[activeLink.mKey]);
            modulesChosed[activeLink.mKey].pins[activeLink.pKey] = pin;
            console.log(modulesChosed[activeLink.mKey].pins[activeLink.pKey])
            //modulesChosed = modulesChosed;
        }
        console.log(modulesChosed)
       //    activeLink = _.cloneDeep(nullActiveLink)
    }

    const addModule = (module) => {
        $: {
            modulesChosed = [...modulesChosed, _.cloneDeep(module)];
        }
    }

        document.querySelector("body").onclick = (e) => {
        if (e.target.classList[0] !== "module-pin" && e.target.classList[0] !== "   pin") {
            activeLink = _.cloneDeep(nullActiveLink);
            let line = document.querySelector("#svg line[data-type='temp']");
            if (line){
                line.remove();
            }
        }
    }
    let title = "";
    let ssid = "";
    let pass = "";
    let mqttToken = ""
    profile.subscribe(value=>{
        mqttToken = value.mqttToken;
    })
    const generate = () => {
        let request = {
            title: title,
            ssid: ssid,
            pass: pass,
            modulesChosed: modulesChosed,
            mqttToken: mqttToken
        }

        lrFetch.postRaw("/api/codeGen", request).then(response => {
            response.blob().then(blob => {
                var url = window.URL.createObjectURL(blob);
                var a = document.createElement('a');
                a.href = url;
                a.download = "code.txt";
                document.body.appendChild(a); // we need to append the element to the dom -> otherwise it will not work in firefox
                a.click();
                a.remove();  //afterwards we remove the element again
            })
            console.log(response)
        })
    }

    const removeModule = (key) => {
        modulesChosed = modulesChosed.filter((m, i) => {
            return i !== key;
        })

    }

    const dataTypes = {
        "bool": "on/off",
        "number": "1,2,3...n"
    }

    const getAvailable = (pin, al) => {

        if (!al.module) return "";
        console.log(pin.type, al.module.type)
        if (pin.type !== al.module.type) {
            return "bad"
        } else {
            return "good"
        }
        return "aa "
    }

</script>

<svg id="svg" style="position: fixed; left: 0; top: 0; z-index: 10000" width="100%" height="100%" xmlns="http://www.w3.org/2000/svg" version="1.1" bind:this={svgNode}>
    <line x1="815.96875" y1="214" x2="1264.78125" y2="177.4375" stroke-width="3" stroke="green" data-type="fix"></line>
</svg>
<div class="codeGen" style="display: flex; flex-direction: column; box-sizing: border-box;">
    <div style="display: flex; height: 100%">
        <div class="settings">
            <div>Настройка устройства</div>
            <h2>NodeMCU (ESP8266)</h2>
            <div class="label"> Название</div>
            <div><input bind:value={title}></div>

            <div class="label">WiFi SSID</div>
            <div><input bind:value={ssid}></div>
            <div class="label">WiFi пароль</div>
            <div><input bind:value={pass}></div>
            <div><button on:click={()=>generate()}>Сгенерировать код</button></div>

        </div>
        <div style="display: flex; margin-top: 20px; width: 100%; position: relative">
            <div style="position: absolute; left: 10px; top: 10px;">
                <h1>Генератор кода</h1>
                <div style="color: #666; margin-top: -20px;font-size: 12px">
                    Для демонстрации работы кодогенератора представлена возможность собрать устройство на плате NodeMCU<br>Выберите модуль, нажмите на его связь и соедените с пином на контроллере.
                <br>Послее добавления всех модулей, введите настройки wi-fi и сгенерируйте код. Тип: тип данных, распознаваемые в прошывке. </div>
            </div>
            <div class="controller">
                <img alt="" src="nodemcu.png">
                {#each Object.keys(pins) as pin}

                    <div class={"pin "+(activeLink ? getAvailable(pins[pin], activeLink) : " ")} on:click={()=> getAvailable(pins[pin], activeLink) === "good" ? selectPin(pins[pin].name) : ()=>{}} data-id={pins[pin].name} style={`left: ${pins[pin].x}px; top: ${pins[pin].y}px`}>
                        {pins[pin].name}
                    </div>
                {/each}
            </div>

            <div class="modules-added">
                {#each modulesChosed as module, mKey}
                    <div class="module" style={`min-height: ${12+module.pins.length*12}px; line-height: ${12+module.pins.length*12}px`}>
                        <div class="module-pins">
                            {#each module.pins as pin, pKey}
                                <div on:click={()=>activateLink(module, mKey+""+pKey, pin, mKey, pKey)} class="module-pin" data-id={mKey+""+pKey}></div>
                            {/each}
                            <div style="font-size: 6px; margin-top: 3px; line-height: 10px">{module.pins.length} pins</div>
                        </div>
                        {module.name}
                        <div><strong>Тип: </strong>{module.data} ({dataTypes[module.data]})</div>
                    <div style="position: absolute; right: -20px; top: 0; font-size: 20px; color: #fff; width: 28px; line-height: 28px; text-align: center; background: orangered; border-radius: 50%;" on:click={()=>removeModule(mKey)}>&times;</div>
                    </div>
                {/each}
            </div>
        </div>
    </div>
    <div class="available-modules">
        <div style="font-size: 12px">Доступные модули</div>
        <div style="display: flex">
        {#each modules as module, mKey}
            <div class="module" on:click={()=>addModule(module)} style={`height: ${12+module.pins.length*12}px; line-height: ${12+module.pins.length*12}px; margin-right: 40px`}>
                <div class="module-pins">
                    {#each module.pins as pin, pKey}
                        <div class="module-pin"></div>
                    {/each}
                </div>
                {module.name}

            </div>
        {/each}
        </div>
    </div>
</div>