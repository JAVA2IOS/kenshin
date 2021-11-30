const FileInputDomName = 'one-date-file-input'

/* 文件上传后绑定到属性上的Url数据 */
const OneDateFileBindSrcProperty = 'one-file'

const ODToast = {

    // 成功toast
    success: function(message) {
        layer.msg(message, {icon: 1})
    },

    // 错误toast
    error: function(message) {
        layer.msg(message, {icon: 5})
    }
} 