let conn = null;

window.onload = function () {
    const msg = document.getElementById("msg");
    const log = document.getElementById("log");
    const form = document.getElementById("form");
    const connectBtn = document.getElementById("connect");
    const disconnectBtn = document.getElementById("disconnect");

    function appendLog(item) {
        const doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(item);
        if (doScroll) {
            log.scrollTop = log.scrollHeight - log.clientHeight;
        }
    }

    function showMessage(text, bold = false) {
        const item = document.createElement("div");
        item.innerHTML = bold ? `<b>${text}</b>` : text;
        appendLog(item);
    }

    connectBtn.onclick = function () {
        if (conn) return;
        if (!window["WebSocket"]) {
            showMessage("Your browser does not support WebSockets.", true);
            return;
        }

        conn = new WebSocket("ws://127.0.0.1:8080/api/v1/ws/connect");

        conn.onopen = function () {
            showMessage("Connected.", true);
        };

        conn.onclose = function () {
            showMessage("Connection closed.", true);
            conn = null;
        };

        conn.onmessage = function (evt) {
            const messages = evt.data.split('\n');
            for (let msg of messages) {
                const item = document.createElement("div");
                item.innerText = msg;
                appendLog(item);
            }
        };
    };

    connectBtn.onclick = function () {
        if (conn) return;

        const uid = encodeURIComponent(document.getElementById("uid").value.trim());
        const username = encodeURIComponent(document.getElementById("username").value.trim());

        // 构建查询参数
        const params = new URLSearchParams();
        if (uid) params.append("uid", uid);
        if (username) params.append("username", username);

        const wsUrl = `ws://127.0.0.1:8080/api/v1/ws/chat${params.toString() ? "?" + params.toString() : ""}`;

        if (!window["WebSocket"]) {
            showMessage("Your browser does not support WebSockets.", true);
            return;
        }

        conn = new WebSocket(wsUrl);

        conn.onopen = function () {
            showMessage(`Connected ${uid ? `(UID: ${uid})` : ''} ${username ? `(Username: ${username})` : ''}.`, true);
        };

        conn.onclose = function () {
            showMessage("Connection closed.", true);
            conn = null;
        };

        conn.onmessage = function (evt) {
            const messages = evt.data.split('\n');
            for (let msg of messages) {
                const item = document.createElement("div");
                item.innerText = msg;
                appendLog(item);
            }
        };
    };



    form.onsubmit = function () {
        if (!conn || !msg.value.trim()) return false;
        conn.send(msg.value.trim());
        msg.value = "";
        return false;
    };
};
