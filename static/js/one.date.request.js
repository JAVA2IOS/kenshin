
function request(method, url, parameters, func) {
    var $ = layui.jquery
    $.ajax({
        url: url,
        type: method,
        success: function (response) {

            layui.msg("请求成功")
            if (func == 'undefined' || func == nil) {
                return
            }



            func(response)()
        },
        error: function (error) {
            layui.msg("请求出错" + error)
        }
    })
}