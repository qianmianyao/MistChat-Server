let conn = null;

window.onload = function () {
    const msg = document.getElementById("msg");
    const log = document.getElementById("log");
    const form = document.getElementById("form");
    const connectBtn = document.getElementById("connect");
    const disconnectBtn = document.getElementById("disconnect");
    const msgType = document.getElementById("msgType");
    const destination = document.getElementById("destination");

    function appendLog(item) {
        const doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(item);
        if (doScroll) {
            log.scrollTop = log.scrollHeight - log.clientHeight;
        }
    }

    function showMessage(text, bold = false, isSystemMessage = false) {
        const item = document.createElement("div");
        if (isSystemMessage) {
            item.className = "p-2 my-3 bg-gray-200 dark:bg-gray-700 rounded-md text-center";
        } else {
            item.className = "my-3";
        }
        item.innerHTML = bold ? `<b>${text}</b>` : text;
        appendLog(item);
    }

    function formatTimestamp() {
        const now = new Date();
        return `${now.getHours().toString().padStart(2, '0')}:${now.getMinutes().toString().padStart(2, '0')}:${now.getSeconds().toString().padStart(2, '0')}`;
    }

    connectBtn.onclick = function () {
        if (conn) return;

        const uid = encodeURIComponent(document.getElementById("uid").value.trim());
        const username = encodeURIComponent(document.getElementById("username").value.trim());

        // 构建查询参数
        const params = new URLSearchParams();
        if (uid) params.append("uid", uid);
        if (username) params.append("username", username);

        const wsUrl = `ws://127.0.0.1:8080/api/v1/chat/connect${params.toString() ? "?" + params.toString() : ""}`;

        if (!window["WebSocket"]) {
            showMessage("您的浏览器不支持 WebSocket 连接。", true, true);
            return;
        }

        conn = new WebSocket(wsUrl);

        conn.onopen = function () {
            showMessage(`${formatTimestamp()} 已连接 ${uid ? `(用户ID: ${decodeURIComponent(uid)})` : ''} ${username ? `(用户名: ${decodeURIComponent(username)})` : ''}.`, true, true);
        };

        conn.onclose = function () {
            showMessage(`${formatTimestamp()} 连接已关闭。`, true, true);
            conn = null;
        };

        conn.onmessage = function (evt) {
            try {
                const data = JSON.parse(evt.data);
                
                if (data.source && data.message) {
                    const source = data.source;
                    const message = data.message;
                    const time = data.timestamp ? new Date(data.timestamp) : new Date();
                    const timeStr = `${time.getHours().toString().padStart(2, '0')}:${time.getMinutes().toString().padStart(2, '0')}`;
                    
                    if (message.type === 'system') {
                        const systemMessage = `<span class="text-gray-500">[${timeStr}]</span> <span class="text-yellow-500">系统消息:</span> ${message.content.text}`;
                        showMessage(systemMessage, false, true);
                    } else {
                        const isCurrentUser = source.uid === document.getElementById("uid").value.trim();
                        
                        // 使用气泡样式
                        let bubbleHtml = '';
                        if (isCurrentUser) {
                            // 自己发送的消息 - 右侧气泡
                            bubbleHtml = `
                                <div class="flex flex-col items-end">
                                    <div class="font-medium text-gray-500 text-xs mb-1">${timeStr}</div>
                                    <div class="flex items-end">
                                        <div class="max-w-[80%] bg-green-500 text-white p-3 rounded-lg rounded-br-none shadow">
                                            ${message.content.text}
                                        </div>
                                    </div>
                                </div>
                            `;
                        } else {
                            // 收到的消息 - 左侧气泡
                            bubbleHtml = `
                                <div class="flex flex-col items-start">
                                    <div class="font-medium text-blue-500 text-sm mb-1">${source.name || source.uid}</div>
                                    <div class="flex items-end">
                                        <div class="max-w-[80%] bg-blue-500 text-white p-3 rounded-lg rounded-bl-none shadow">
                                            ${message.content.text}
                                        </div>
                                    </div>
                                    <div class="text-gray-500 text-xs mt-1">${timeStr}</div>
                                </div>
                            `;
                        }
                        
                        showMessage(bubbleHtml, false, false);
                    }
                } else {
                    showMessage(`${formatTimestamp()} 收到消息: ${evt.data}`, false, true);
                }
            } catch (e) {
                // 尝试按旧格式处理消息
                const messages = evt.data.split('\n');
                for (let msg of messages) {
                    showMessage(`${formatTimestamp()} ${msg}`, false, true);
                }
            }
        };
    };
    
    disconnectBtn.onclick = function () {
        if (!conn) return;
        conn.close();
        showMessage(`${formatTimestamp()} 正在断开连接...`, true, true);
    };

    form.onsubmit = function () {
        if (!conn || !msg.value.trim()) return false;
        
        const uid = document.getElementById("uid").value.trim();
        const username = document.getElementById("username").value.trim();
        
        // 构造消息对象
        const messageObj = {
            source: {
                uid: uid,
                name: username
            },
            message: {
                type: msgType.value,
                content: {
                    text: msg.value.trim()
                }
            }
        };
        
        // 添加接收者（如果提供）
        if (destination.value.trim()) {
            messageObj.destination = destination.value.trim();
        }
        
        // 发送消息
        conn.send(JSON.stringify(messageObj));
        
        // 在本地显示发送的消息 - 气泡样式
        const timeStr = formatTimestamp();
        const messageContent = msg.value.trim();
        
        const bubbleHtml = `
            <div class="flex flex-col items-end">
                <div class="font-medium text-gray-500 text-xs mb-1">${timeStr}</div>
                <div class="flex items-end">
                    <div class="max-w-[80%] bg-green-500 text-white p-3 rounded-lg rounded-br-none shadow">
                        ${messageContent}
                    </div>
                </div>
            </div>
        `;
        
        showMessage(bubbleHtml, false, false);
        
        msg.value = "";
        return false;
    };
};
