function webSocketConnect() {
    var ws = new WebSocket('ws://' + window.location.host + '/ws')

    ws.onmessage = function conn(response) {

        console.log('socket: ' + response.data)

        var ob = {}

        ob.message = "html is alived !" + new Date().toLocaleTimeString()

        ws.send(JSON.stringify(ob))
    }
}