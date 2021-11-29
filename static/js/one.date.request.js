
function callback(func) {
    func()
}

function ODTask(method, url, parameters, callback) {
    var $ = layui.$
    $.ajax({
        url: url,
        type: method,
        success: function (response) {

            layer.msg("请求成功")
            if (func == undefined || func == 'nil') {
                return
            }

            callback(response)
        },
        error: function (error) {
            console.info(error)
            layer.msg("请求出错" + error.statusText)
        }
    })
}