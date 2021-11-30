layui.use(['upload'], function () {
  var $ = layui.jquery, upload = layui.upload

  upload.render({
    elem: '.one-date-file-input',
    url: '/file/upload/xlsx',
    accept: 'file',
    // acceptMime: 'application/vnd.ms-excel, application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
    exts: 'xlsx|xls|csv',
    drag: true,
    choose: function (obj) {
    },

    before: function () {
      layer.load()
    },

    progress: function (prgress, elem, responsed, index) {
      // console.log('1百分比:', prgress, '%')
    },
    done: function (response, index, uploadInstance) {
      layer.closeAll('loading')

      ODToast.success('文件上传完成')

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

layui.use('jquery', function () {
  var $ = layui.$
  $(function () {
    // 开启webSocket
    webSocketConnect()
  })
})

//JS 
layui.use(['element', 'layer', 'util', 'jquery'], function () {
  var element = layui.element
    , layer = layui.layer
    , util = layui.util
    , $ = layui.$;

  //头部事件
  util.event('lay-header-event', {
    //左侧菜单事件
    menuLeft: function (othis) {
      layer.msg('展开左侧菜单的操作', { icon: 0 });
    }
    , menuRight: function () {
      layer.open({
        type: 1
        , content: '<div style="padding: 15px;">处理右侧面板的操作</div>'
        , area: ['260px', '100%']
        , offset: 'rt' //右上角
        , anim: 5
        , shadeClose: true
      });
    }
  });

  util.event('one-date-button', {

    // 文件上传
    upload: function (d) {

      var files = []

      $("input[name='xlsxFile']").each(function () {

        var attrUrl = $(this).attr(OneDateFileBindSrcProperty)

        var domId = $(this).attr('id')

        if (domId == undefined || domId == 'nil') {
          ODToast.error("当前元素为空")
          return false
        }

        if (attrUrl == undefined || attrUrl == 'nil') {
          ODToast.error("文件为空")
          return false
        }

        files[$(this).attr('id')] = attrUrl
      })

      if (files.length == 0) {
        return
      }

      console.info('数据: ' + files)

      ODTask.GET('/file/xlsx/' + $(this).attr('one-date-p'), files, function (response) {
        
        if (response.Code != 200) {
          // alert("错误: " + response.Message)
          ODToast.error('错误:' + response.Message)
          return
        }

        console.info('data:' + response.Message)
      })
    }
  })
})
