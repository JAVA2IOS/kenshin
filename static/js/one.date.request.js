const ODTask = {
    // 正常请求
    task: function(method, url, parameters, callback) {
        var $ = layui.$
        $.ajax({
            url: url,
            type: method,
            data: parameters,
            success: function (response) {
                // layer.msg("请求成功")
                if (callback == undefined || callback == 'nil') {
                    return
                }
    
                callback(JSON.parse(JSON.stringify(response)))
            },
            error: function (error) {
                // console.info(error)

                ODToast.error(error.statusText)
            }
        })
    },

    // GET请求
    GET: function(url, parameters, callback) {
        this.task('GET', url, parameters, callback)
    },

    // post请求
    POST: function(url, parameters, callback) {
        this.task('POST', url, parameters, callback)
    },
}