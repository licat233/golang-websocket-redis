window.onload = onlineview();

function onlineview() {
    let ws;
    const viewerID = Math.round(new Date() / 1000)+generateUUID();
    const webURL = getwebURL();
    const viewData = webURL+"*"+viewerID;
    ws = new WebSocket("ws://localhost:8888/ws");
    documentPrint();
    var view = document.getElementById("onlineview");
    var Timing = self.setInterval(
        function () {
        if(ws.readyState !== ws.OPEN){
            console.log("连接已中断!")
            window.clearInterval(Timing)
            return false
        }
        ws.send(viewData)
    }, 5000)
    ws.onmessage = function (evt) {
        view.innerHTML = evt.data + "&nbsp;人正在閱讀";
    }
}

function generateUUID() {
    var d = new Date().getTime();
    return 'xxxxxxxx'.replace(/[xy]/g, function (c) {
        var r = (d + Math.random() * 16) % 16 | 0;
        d = Math.floor(d / 16);
        return (c === 'x' ? r : (r & 0x3 | 0x8)).toString(16);
    });
};

function getwebURL() {
    var host = window.location.href.trim();
    -1 !== host.indexOf("#") && (host = host.substr(0, host.indexOf("#")));
    -1 !== host.indexOf("?") && (host = host.substr(0, host.indexOf("?")));
    "/" === host.substr(host.length - 1, 1) && (host = host.substr(0, host.length - 1));
    host_arr = host.split("/");
    web_url = host_arr[2];
    void 0 !== host_arr[3] && (web_url = host_arr[2] + "/" + host_arr[3]);
    web_url = String(web_url);
    return web_url
}

function documentPrint() {
    var style = "<style>#onlineview{bottom: 0px;text-align:center;color:#3A3A3A;font-size:13px}</style>";
    this.document.write(style)
    var d = document.createElement("div");
    d.id = "onlineview";
    d.style.bottom = "0px";
    document.body.appendChild(d);
}