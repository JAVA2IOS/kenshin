const FileInputDomName = 'one-date-file-input'

/* 文件上传后绑定到属性上的Url数据 */
const OneDateFileBindSrcProperty = 'one-file'

const ODToast = {

    // 成功toast
    success: function () {

        layer.closeAll('loading')

        var i, s, numargs = arguments.length;

        s = numargs;

        for (i = 0; i < numargs; i++) {
            s += arguments[i];
        }

        layer.msg(s, { icon: 1 })

        console.info(s)
    },

    // 错误toast
    error: function () {

        layer.closeAll('loading')

        var i, s, numargs = arguments.length;

        s = numargs;

        for (i = 0; i < numargs; i++) {
            s += arguments[i];
        }

        layer.msg(s, { icon: 5 })

        console.info(s)
    }
}