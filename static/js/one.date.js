layui.use('upload', function () {
    var $ = layui.jquery, upload = layui.upload

    upload.render({
      elem: '.one-date-file-input',
      url: '/file/upload/xlsx',
      accept: 'file',
      // acceptMime: 'application/vnd.ms-excel, application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
      exts: 'xlsx|xls|csv',
      drag: true,
      choose: function (obj) {
    
        // obj.pushFile()

        obj.preview(function (index, file, result) {
          console.log('文件名：', file.name, 'infex', index)
        })
        
      },

      before: function () {
        // this.data = {
        //   deviceName: "nice_shot"
        // }

        layer.load()
      },

      progress: function (prgress, elem, responsed, index) {
        console.log('1百分比:', prgress, '%')
      },
      done: function (response, index, uploadInstance) {
        layer.closeAll('loading')

        layer.msg('文件上传完成', response)
        var elem = this.item

        var data = response.Data

        var fileText = data['xlsx']

        elem.val(fileText)

        elem.attr(OneDateFileBindSrcProperty, data['url'])
      },
      error: function (index, uploadInstance) {

        layer.closeAll('loading')
        console.log('发生错误了')
      }
    })
  })